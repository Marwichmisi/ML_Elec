// Package main implements the MQTT bridge plugin for ML-Elec.
// It runs an embedded mochi-mqtt broker, subscribes to ESP32 topics,
// validates incoming sensor data, and publishes to the NATS bus.
package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-hclog"
	natsclient "github.com/nats-io/nats.go"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"ml-elec/internal/config"
	"ml-elec/internal/storage"
	sdk "ml-elec/pkg/sdk/v1"
)

// App holds all MQTT plugin components.
type App struct {
	Config *config.Config
	Broker *mqtt.Server
	Client pahomqtt.Client
}

// StartBroker creates and starts an embedded mochi-mqtt broker on the given port.
// Use port 0 to let the OS assign a free port.
func StartBroker(port int) (*mqtt.Server, error) {
	server := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
		return nil, fmt.Errorf("adding auth hook: %w", err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "tcp1",
		Address: addr,
	})
	if err := server.AddListener(tcp); err != nil {
		return nil, fmt.Errorf("adding TCP listener: %w", err)
	}

	go func() {
		if err := server.Serve(); err != nil {
			slog.Error("broker serve error", "error", err)
		}
	}()

	// Wait for the listener to be ready
	time.Sleep(50 * time.Millisecond)

	return server, nil
}

// BrokerAddr returns the TCP address of the broker's listener.
func BrokerAddr(server *mqtt.Server) (string, error) {
	l, ok := server.Listeners.Get("tcp1")
	if !ok {
		return "", fmt.Errorf("no TCP listener found")
	}
	return l.Address(), nil
}

// ConnectClient creates a paho MQTT client connected to the embedded broker
// with LWT (Last Will and Testament) for offline status.
func ConnectClient(brokerAddr string, clientID string) (pahomqtt.Client, error) {
	opts := pahomqtt.NewClientOptions().
		AddBroker("tcp://" + brokerAddr).
		SetClientID(clientID).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetWill("sys/mqtt-plugin/status", "offline", 1, true)

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", token.Error())
	}

	// Publish online status
	client.Publish("sys/mqtt-plugin/status", 1, true, "online")

	return client, nil
}

// ConnectClientWithBackoff creates a paho MQTT client with configurable reconnect backoff.
func ConnectClientWithBackoff(brokerAddr string, clientID string, backoffCfg config.ReconnectBackoffConfig) (pahomqtt.Client, error) {
	initialInterval, err := time.ParseDuration(backoffCfg.InitialInterval)
	if err != nil {
		initialInterval = 1 * time.Second
	}
	maxInterval, err := time.ParseDuration(backoffCfg.MaxInterval)
	if err != nil {
		maxInterval = 30 * time.Second
	}

	opts := pahomqtt.NewClientOptions().
		AddBroker("tcp://" + brokerAddr).
		SetClientID(clientID).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetConnectRetryInterval(initialInterval).
		SetMaxReconnectInterval(maxInterval).
		SetWill("sys/mqtt-plugin/status", "offline", 1, true)

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", token.Error())
	}

	// Publish online status
	client.Publish("sys/mqtt-plugin/status", 1, true, "online")

	return client, nil
}

// ConnectToNATS connects to the core's embedded NATS server with exponential backoff.
// This does NOT create a new NATS server — it connects to the existing one.
func ConnectToNATS(host string, port int, backoffCfg config.ReconnectBackoffConfig) (*natsclient.Conn, error) {
	url := fmt.Sprintf("nats://%s:%d", host, port)
	initialInterval, err := time.ParseDuration(backoffCfg.InitialInterval)
	if err != nil {
		initialInterval = 1 * time.Second
	}
	conn, err := natsclient.Connect(url,
		natsclient.MaxReconnects(-1),
		natsclient.ReconnectWait(initialInterval),
	)
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS at %s: %w", url, err)
	}
	return conn, nil
}

// MessageHandler is called when a subscribed message is received.
type MessageHandler func(topic string, payload []byte)

// buildMessageHandler creates a message handler that validates, publishes to NATS, and persists to SQLite.
func buildMessageHandler(natsConn *natsclient.Conn, store *storage.Store, cfg *config.Config) MessageHandler {
	validationCfg := buildValidationCfg(cfg.Validation)
	return func(topic string, payload []byte) {
		// Validate the message
		if err := ValidateMessage(payload, validationCfg); err != nil {
			slog.Error("message validation failed", "topic", topic, "error", err, "payload_len", len(payload))
			return
		}

		// Parse JSON or binary payload
		reading, jsonErr := ParseJSONPayload(payload)
		if jsonErr != nil {
			// Try binary
			binResult, binErr := ParseBinaryPayload(payload)
			if binErr != nil {
				slog.Error("failed to parse payload", "topic", topic, "json_err", jsonErr, "bin_err", binErr)
				return
			}
			// Binary: publish raw to NATS
			natsSubject := "sensor." + topic
			if err := natsConn.Publish(natsSubject, payload); err != nil {
				slog.Error("NATS publish failed", "topic", topic, "subject", natsSubject, "error", err)
			}
			// Persist binary reading (one entry per sample count metadata)
			ts := time.Unix(0, binResult.Header.Timestamp)
			if err := store.InsertSensor(context.Background(), topic, float64(binResult.Header.SampleCount), ts); err != nil {
				slog.Error("storage insert failed", "topic", topic, "error", err)
			}
			return
		}

		// JSON: publish each sensor value to NATS
		natsSubject := "sensor." + topic
		if err := natsConn.Publish(natsSubject, payload); err != nil {
			slog.Error("NATS publish failed", "topic", topic, "subject", natsSubject, "error", err)
		}

		// Persist each sensor value to SQLite
		ts := time.Unix(reading.Timestamp, 0)
		for sensorName, value := range reading.Values {
			if err := store.InsertSensor(context.Background(), sensorName, value, ts); err != nil {
				slog.Error("storage insert failed", "sensor", sensorName, "error", err)
			}
		}
	}
}

// buildValidationCfg converts config.ValidationConfig to the local ValidationConfig type.
func buildValidationCfg(cfg config.ValidationConfig) ValidationConfig {
	maxFutureDrift, err := time.ParseDuration(cfg.Timestamp.MaxFutureDrift)
	if err != nil {
		maxFutureDrift = 5 * time.Second
	}
	maxPastDrift, err := time.ParseDuration(cfg.Timestamp.MaxPastDrift)
	if err != nil {
		maxPastDrift = 24 * time.Hour
	}

	ranges := make(map[string]RangeConfig, len(cfg.Ranges))
	for sensor, r := range cfg.Ranges {
		ranges[sensor] = RangeConfig{Min: r.Min, Max: r.Max}
	}

	return ValidationConfig{
		Ranges: ranges,
		Timestamp: TimestampConfig{
			MaxFutureDrift: maxFutureDrift,
			MaxPastDrift:   maxPastDrift,
		},
		Health: HealthConfig{
			MaxPayloadSize: cfg.Health.MaxPayloadSize,
			MinPayloadSize: cfg.Health.MinPayloadSize,
		},
	}
}

// SubscribeAndBridge subscribes to the given topic filter and calls the handler
// for each received message.
func SubscribeAndBridge(client pahomqtt.Client, topicFilter string, handler MessageHandler) error {
	token := client.Subscribe(topicFilter, 1, func(_ pahomqtt.Client, msg pahomqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("subscribe to %s: %w", topicFilter, token.Error())
	}
	return nil
}

// SubscribeAndBridgeWithQoS subscribes with a specific QoS level.
func SubscribeAndBridgeWithQoS(client pahomqtt.Client, topicFilter string, qos byte, handler MessageHandler) error {
	token := client.Subscribe(topicFilter, qos, func(_ pahomqtt.Client, msg pahomqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("subscribe to %s: %w", topicFilter, token.Error())
	}
	return nil
}

// resolveQoS returns the QoS level for a topic, using per-topic override if present.
func resolveQoS(topic string, defaultQoS byte, qosPerTopic map[string]byte) byte {
	if qos, ok := qosPerTopic[topic]; ok {
		return qos
	}
	return defaultQoS
}

// SensorReading represents a parsed JSON sensor reading per D-13.
type SensorReading struct {
	Timestamp int64
	Values    map[string]float64
}

// ParseJSONPayload parses a JSON MQTT payload per D-13 format.
func ParseJSONPayload(data []byte) (*SensorReading, error) {
	var raw struct {
		Timestamp int64              `json:"ts"`
		Values    map[string]float64 `json:"values"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if raw.Values == nil {
		return nil, fmt.Errorf("missing required 'values' field")
	}

	return &SensorReading{
		Timestamp: raw.Timestamp,
		Values:    raw.Values,
	}, nil
}

// BinaryHeaderSize is the size of the binary MQTT payload header (28 bytes per D-12).
const BinaryHeaderSize = 28

// BinaryHeader represents the fixed header for binary MQTT payloads.
type BinaryHeader struct {
	Version     byte
	SensorType  byte
	Encoding    byte
	Reserved1   byte
	Timestamp   int64
	SampleCount uint32
	SampleRate  uint32
	Reserved2   [8]byte
}

// BinaryParsingResult holds the result of parsing a binary MQTT payload.
type BinaryParsingResult struct {
	Header  BinaryHeader
	Samples []int16
}

// ParseBinaryPayload parses a binary MQTT payload with big-endian 28-byte header.
func ParseBinaryPayload(data []byte) (*BinaryParsingResult, error) {
	if len(data) < BinaryHeaderSize {
		return nil, fmt.Errorf("payload too short: %d < %d", len(data), BinaryHeaderSize)
	}

	if data[0] != 1 {
		return nil, fmt.Errorf("unsupported version: %d", data[0])
	}

	header := BinaryHeader{
		Version:     data[0],
		SensorType:  data[1],
		Encoding:    data[2],
		Reserved1:   data[3],
		Timestamp:   int64(binary.BigEndian.Uint64(data[4:12])),
		SampleCount: binary.BigEndian.Uint32(data[12:16]),
		SampleRate:  binary.BigEndian.Uint32(data[16:20]),
	}
	copy(header.Reserved2[:], data[20:28])

	// Parse samples based on encoding (int16 = 2 bytes per sample)
	payload := data[BinaryHeaderSize:]
	expectedLen := int(header.SampleCount) * 2
	if len(payload) < expectedLen {
		return nil, fmt.Errorf("payload too short for %d samples: %d < %d",
			header.SampleCount, len(payload), expectedLen)
	}

	samples := make([]int16, header.SampleCount)
	for i := uint32(0); i < header.SampleCount; i++ {
		offset := i * 2
		samples[i] = int16(binary.BigEndian.Uint16(payload[offset : offset+2]))
	}

	return &BinaryParsingResult{
		Header:  header,
		Samples: samples,
	}, nil
}

// ValidateRange checks if a sensor value is within the configured min/max range.
func ValidateRange(sensor string, value float64, ranges map[string]RangeConfig) error {
	cfg, ok := ranges[sensor]
	if !ok {
		return nil // no range configured = pass
	}
	if value < cfg.Min || value > cfg.Max {
		return fmt.Errorf("value %.2f out of range [%.2f, %.2f] for sensor %s", value, cfg.Min, cfg.Max, sensor)
	}
	return nil
}

// RangeConfig defines min/max bounds for a sensor type.
type RangeConfig struct {
	Min float64
	Max float64
}

// TimestampConfig defines acceptable timestamp drift bounds.
type TimestampConfig struct {
	MaxFutureDrift time.Duration
	MaxPastDrift   time.Duration
}

// ValidateTimestamp checks if a timestamp is within acceptable drift bounds.
func ValidateTimestamp(ts time.Time, cfg TimestampConfig) error {
	now := time.Now()
	if ts.After(now.Add(cfg.MaxFutureDrift)) {
		return fmt.Errorf("timestamp %v is in the future (drift > %v)", ts, cfg.MaxFutureDrift)
	}
	if ts.Before(now.Add(-cfg.MaxPastDrift)) {
		return fmt.Errorf("timestamp %v is too old (drift > %v)", ts, cfg.MaxPastDrift)
	}
	return nil
}

// HealthConfig defines payload health validation bounds.
type HealthConfig struct {
	MaxPayloadSize int
	MinPayloadSize int
}

// ValidateQuality checks payload size and completeness.
func ValidateQuality(payload []byte, cfg HealthConfig) error {
	if len(payload) < cfg.MinPayloadSize {
		return fmt.Errorf("payload too small: %d < %d", len(payload), cfg.MinPayloadSize)
	}
	if len(payload) > cfg.MaxPayloadSize {
		return fmt.Errorf("payload too large: %d > %d", len(payload), cfg.MaxPayloadSize)
	}
	return nil
}

// ValidationConfig holds 3-level data validation configuration.
type ValidationConfig struct {
	Ranges    map[string]RangeConfig
	Timestamp TimestampConfig
	Health    HealthConfig
}

// ValidateMessage runs three-level validation on an MQTT message payload.
func ValidateMessage(payload []byte, cfg ValidationConfig) error {
	// Level 3: Quality check (payload size)
	if err := ValidateQuality(payload, cfg.Health); err != nil {
		return fmt.Errorf("quality validation: %w", err)
	}

	// Try JSON parsing first
	reading, jsonErr := ParseJSONPayload(payload)
	if jsonErr == nil {
		// Level 2: Timestamp validation
		ts := time.Unix(reading.Timestamp, 0)
		if err := ValidateTimestamp(ts, cfg.Timestamp); err != nil {
			return fmt.Errorf("timestamp validation: %w", err)
		}

		// Level 1: Range validation
		for sensor, value := range reading.Values {
			if err := ValidateRange(sensor, value, cfg.Ranges); err != nil {
				return fmt.Errorf("range validation: %w", err)
			}
		}
		return nil
	}

	// Try binary parsing
	binResult, binErr := ParseBinaryPayload(payload)
	if binErr == nil {
		ts := time.Unix(0, binResult.Header.Timestamp)
		if err := ValidateTimestamp(ts, cfg.Timestamp); err != nil {
			return fmt.Errorf("timestamp validation: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unable to parse payload as JSON or binary")
}

// MQTTPlugin implements the PluginLifecycle and SensorCollector interfaces.
type MQTTPlugin struct {
	config   *config.Config
	broker   *mqtt.Server
	client   pahomqtt.Client
	natsConn *natsclient.Conn // NATS connection to core's embedded server
	store    *storage.Store   // SQLite storage for persistence
}

// Init initializes the plugin with the provided configuration.
func (p *MQTTPlugin) Init(_ context.Context, cfgMap map[string]string) error {
	slog.Info("mqtt plugin init", "config", cfgMap)
	p.config = config.DefaultConfig()
	return nil
}

// Start starts the MQTT broker and client.
func (p *MQTTPlugin) Start(_ context.Context) error {
	broker, err := StartBroker(p.config.MQTT.Port)
	if err != nil {
		return fmt.Errorf("starting broker: %w", err)
	}
	p.broker = broker

	addr, err := BrokerAddr(broker)
	if err != nil {
		return fmt.Errorf("getting broker address: %w", err)
	}

	client, err := ConnectClient(addr, p.config.MQTT.ClientID)
	if err != nil {
		return fmt.Errorf("connecting client: %w", err)
	}
	p.client = client

	// Connect to core's embedded NATS server
	natsConn, err := ConnectToNATS(p.config.MQTT.NATSHost, p.config.MQTT.NATSPort, p.config.MQTT.ReconnectBackoff)
	if err != nil {
		return fmt.Errorf("connecting to NATS: %w", err)
	}
	p.natsConn = natsConn

	// Initialize SQLite storage
	store, err := storage.New(&p.config.Storage)
	if err != nil {
		return fmt.Errorf("opening storage: %w", err)
	}
	p.store = store

	slog.Info("mqtt plugin started",
		"port", p.config.MQTT.Port,
		"topic", p.config.MQTT.Topics.Subscribe,
		"nats_host", p.config.MQTT.NATSHost,
		"nats_port", p.config.MQTT.NATSPort,
	)
	return nil
}

// Stop gracefully stops the MQTT plugin.
func (p *MQTTPlugin) Stop(_ context.Context) error {
	if p.natsConn != nil {
		p.natsConn.Close()
	}
	if p.store != nil {
		p.store.Close()
	}
	if p.client != nil {
		p.client.Disconnect(250)
	}
	if p.broker != nil {
		p.broker.Close()
	}
	slog.Info("mqtt plugin stopped")
	return nil
}

// Collect sends a sensor reading to the core for processing.
func (p *MQTTPlugin) Collect(_ context.Context, req *sdk.CollectRequest) (*sdk.CollectResponse, error) {
	slog.Info("sensor collect",
		"sensor_id", req.SensorId,
		"topic", req.Topic,
		"format", req.PayloadFormat,
	)
	return &sdk.CollectResponse{Accepted: true}, nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	plugin := &MQTTPlugin{config: cfg}

	// Start broker and client
	broker, err := StartBroker(cfg.MQTT.Port)
	if err != nil {
		slog.Error("failed to start broker", "error", err)
		os.Exit(1)
	}
	plugin.broker = broker

	addr, err := BrokerAddr(broker)
	if err != nil {
		slog.Error("failed to get broker address", "error", err)
		os.Exit(1)
	}

	// Connect MQTT client with exponential backoff
	client, err := ConnectClientWithBackoff(addr, cfg.MQTT.ClientID, cfg.MQTT.ReconnectBackoff)
	if err != nil {
		slog.Error("failed to connect client", "error", err)
		os.Exit(1)
	}
	plugin.client = client

	// Connect to core's embedded NATS server
	natsConn, err := ConnectToNATS(cfg.MQTT.NATSHost, cfg.MQTT.NATSPort, cfg.MQTT.ReconnectBackoff)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	plugin.natsConn = natsConn

	// Initialize SQLite storage
	store, err := storage.New(&cfg.Storage)
	if err != nil {
		slog.Error("failed to open storage", "error", err)
		os.Exit(1)
	}
	plugin.store = store

	// Build full pipeline handler: validate → NATS → SQLite
	handler := buildMessageHandler(natsConn, store, cfg)

	// Subscribe to configured topic with resolved QoS
	qos := resolveQoS(cfg.MQTT.Topics.Subscribe, cfg.MQTT.QoS, cfg.MQTT.QoSPerTopic)
	err = SubscribeAndBridgeWithQoS(client, cfg.MQTT.Topics.Subscribe, qos, handler)
	if err != nil {
		slog.Error("failed to subscribe", "error", err)
		os.Exit(1)
	}

	slog.Info("mqtt plugin running",
		"port", cfg.MQTT.Port,
		"topic", cfg.MQTT.Topics.Subscribe,
		"qos", qos,
	)

	// gRPC plugin serve (when launched by plugin manager)
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "mqtt-plugin",
		Level:  hclog.Info,
		Output: hclog.DefaultOutput,
	})

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: goplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ML_ELEC_PLUGIN",
			MagicCookieValue: "ml-elec-v1",
		},
		Plugins: map[string]goplugin.Plugin{
			"sensor": &sdk.GRPCPlugin{
				Impl:        plugin,
				CollectImpl: plugin,
			},
		},
		Logger: logger,
	})

	// Wait for shutdown
	<-ctx.Done()
	slog.Info("shutdown signal received")

	// LIFO shutdown
	if plugin.client != nil {
		plugin.client.Disconnect(250)
	}
	if plugin.broker != nil {
		plugin.broker.Close()
	}

	slog.Info("shutdown complete")
}


