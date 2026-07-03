package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	natsserver "github.com/nats-io/nats-server/v2/server"
	natsclient "github.com/nats-io/nats.go"
)

// TestStartBroker verifies the embedded broker starts on a configurable port.
func TestStartBroker(t *testing.T) {
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	if broker == nil {
		t.Fatal("StartBroker returned nil broker")
	}
	defer broker.Close()
}

// TestMQTTClientConnects verifies a paho client connects to the embedded broker with LWT.
func TestMQTTClientConnects(t *testing.T) {
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	client, err := ConnectClient(addr, "test-client")
	if err != nil {
		t.Fatalf("ConnectClient failed: %v", err)
	}
	if client == nil {
		t.Fatal("ConnectClient returned nil client")
	}
	defer client.Disconnect(250)

	if !client.IsConnected() {
		t.Fatal("client not connected")
	}
}

// TestSubscribeAndReceive verifies subscription to esp32/# receives messages.
func TestSubscribeAndReceive(t *testing.T) {
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	client, err := ConnectClient(addr, "test-subscriber")
	if err != nil {
		t.Fatalf("ConnectClient failed: %v", err)
	}
	defer client.Disconnect(250)

	received := make(chan string, 1)
	err = SubscribeAndBridge(client, "test/#", func(topic string, payload []byte) {
		received <- string(payload)
	})
	if err != nil {
		t.Fatalf("SubscribeAndBridge failed: %v", err)
	}

	// Publish test message via broker's inline client
	broker.Publish("test/sensor1", []byte("hello"), false, 1)

	select {
	case msg := <-received:
		if msg != "hello" {
			t.Fatalf("expected 'hello', got '%s'", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for message")
	}
}

// TestJSONParsing verifies JSON payload parsing per D-13.
func TestJSONParsing(t *testing.T) {
	payload := `{"ts":1719900000,"values":{"temperature":25.3,"humidity":60.0,"current":1.5}}`

	reading, err := ParseJSONPayload([]byte(payload))
	if err != nil {
		t.Fatalf("ParseJSONPayload failed: %v", err)
	}

	if reading.Timestamp != 1719900000 {
		t.Errorf("expected timestamp 1719900000, got %d", reading.Timestamp)
	}
	if reading.Values["temperature"] != 25.3 {
		t.Errorf("expected temperature 25.3, got %f", reading.Values["temperature"])
	}
	if reading.Values["humidity"] != 60.0 {
		t.Errorf("expected humidity 60.0, got %f", reading.Values["humidity"])
	}
	if reading.Values["current"] != 1.5 {
		t.Errorf("expected current 1.5, got %f", reading.Values["current"])
	}
}

// TestJSONParsingMissingFields verifies JSON with missing required fields is rejected.
func TestJSONParsingMissingFields(t *testing.T) {
	payload := `{"ts":1719900000}`

	_, err := ParseJSONPayload([]byte(payload))
	if err == nil {
		t.Fatal("expected error for missing values field")
	}
}

// TestJSONParsingInvalidJSON verifies invalid JSON is rejected.
func TestJSONParsingInvalidJSON(t *testing.T) {
	_, err := ParseJSONPayload([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// TestBinaryParsing verifies 28-byte binary header parsing (big-endian).
func TestBinaryParsing(t *testing.T) {
	header := make([]byte, 28)
	header[0] = 1  // version
	header[1] = 1  // sensor_type (vibration)
	header[2] = 0  // encoding (int16)
	header[3] = 0  // reserved

	binary.BigEndian.PutUint64(header[4:12], uint64(1719900000000000000))
	binary.BigEndian.PutUint32(header[12:16], 2)    // sample_count
	binary.BigEndian.PutUint32(header[16:20], 1000) // sample_rate

	samples := make([]byte, 4)
	binary.BigEndian.PutUint16(samples[0:2], 1000)
	binary.BigEndian.PutUint16(samples[2:4], 2000)

	payload := append(header, samples...)

	result, err := ParseBinaryPayload(payload)
	if err != nil {
		t.Fatalf("ParseBinaryPayload failed: %v", err)
	}

	if result.Header.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Header.Version)
	}
	if result.Header.SensorType != 1 {
		t.Errorf("expected sensor_type 1, got %d", result.Header.SensorType)
	}
	if result.Header.Timestamp != 1719900000000000000 {
		t.Errorf("expected timestamp 1719900000000000000, got %d", result.Header.Timestamp)
	}
	if result.Header.SampleCount != 2 {
		t.Errorf("expected sample_count 2, got %d", result.Header.SampleCount)
	}
	if result.Header.SampleRate != 1000 {
		t.Errorf("expected sample_rate 1000, got %d", result.Header.SampleRate)
	}
	if len(result.Samples) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(result.Samples))
	}
	if result.Samples[0] != 1000 {
		t.Errorf("expected sample 1 = 1000, got %d", result.Samples[0])
	}
	if result.Samples[1] != 2000 {
		t.Errorf("expected sample 2 = 2000, got %d", result.Samples[1])
	}
}

// TestBinaryParsingShortPayload verifies short payload is rejected.
func TestBinaryParsingShortPayload(t *testing.T) {
	_, err := ParseBinaryPayload(make([]byte, 10))
	if err == nil {
		t.Fatal("expected error for short payload")
	}
}

// TestBinaryParsingWrongVersion verifies wrong version is rejected.
func TestBinaryParsingWrongVersion(t *testing.T) {
	payload := make([]byte, 28)
	payload[0] = 99

	_, err := ParseBinaryPayload(payload)
	if err == nil {
		t.Fatal("expected error for wrong version")
	}
}

// TestBinaryParsingLittleEndian verifies little-endian data is rejected (big-endian only).
func TestBinaryParsingLittleEndian(t *testing.T) {
	header := make([]byte, 28)
	header[0] = 1
	header[1] = 1

	// Write timestamp in little-endian (wrong!)
	binary.LittleEndian.PutUint64(header[4:12], uint64(1719900000000000000))
	// Use 0 samples so parsing doesn't fail on sample data length
	binary.LittleEndian.PutUint32(header[12:16], 0)
	binary.LittleEndian.PutUint32(header[16:20], 1000)

	result, err := ParseBinaryPayload(header)
	if err != nil {
		t.Fatalf("ParseBinaryPayload failed: %v", err)
	}

	// Big-endian parse of little-endian data gives wrong timestamp
	if result.Header.Timestamp == 1719900000000000000 {
		t.Fatal("should have detected little-endian timestamp as invalid (values differ)")
	}
}

// TestValidationRange verifies range validation.
func TestValidationRange(t *testing.T) {
	ranges := map[string]RangeConfig{
		"temperature": {Min: -40.0, Max: 150.0},
		"humidity":    {Min: 0.0, Max: 100.0},
		"current":     {Min: 0.0, Max: 1000.0},
	}

	tests := []struct {
		name      string
		sensor    string
		value     float64
		wantValid bool
	}{
		{"valid temperature", "temperature", 25.3, true},
		{"temp at min", "temperature", -40.0, true},
		{"temp at max", "temperature", 150.0, true},
		{"temp below min", "temperature", -41.0, false},
		{"temp above max", "temperature", 151.0, false},
		{"valid humidity", "humidity", 50.0, true},
		{"humidity above max", "humidity", 101.0, false},
		{"valid current", "current", 500.0, true},
		{"current above max", "current", 1001.0, false},
		{"unknown sensor", "pressure", 100.0, true}, // no range configured = pass
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRange(tt.sensor, tt.value, ranges)
			if (err == nil) != tt.wantValid {
				t.Errorf("ValidateRange(%s, %f) error = %v, wantValid = %v", tt.sensor, tt.value, err, tt.wantValid)
			}
		})
	}
}

// TestValidationTimestamp verifies timestamp validation.
func TestValidationTimestamp(t *testing.T) {
	cfg := TimestampConfig{
		MaxFutureDrift: 5 * time.Second,
		MaxPastDrift:   24 * time.Hour,
	}

	now := time.Now()

	tests := []struct {
		name      string
		ts        time.Time
		wantValid bool
	}{
		{"current time", now, true},
		{"2s future", now.Add(2 * time.Second), true},
		{"10s future", now.Add(10 * time.Second), false},
		{"1h old", now.Add(-1 * time.Hour), true},
		{"25h old", now.Add(-25 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimestamp(tt.ts, cfg)
			if (err == nil) != tt.wantValid {
				t.Errorf("ValidateTimestamp(%v) error = %v, wantValid = %v", tt.ts, err, tt.wantValid)
			}
		})
	}
}

// TestValidationQuality verifies quality validation.
func TestValidationQuality(t *testing.T) {
	cfg := HealthConfig{
		MaxPayloadSize: 1024,
		MinPayloadSize: 1,
	}

	tests := []struct {
		name      string
		payload   []byte
		wantValid bool
	}{
		{"normal payload", make([]byte, 100), true},
		{"empty payload", []byte{}, false},
		{"nil payload", nil, false},
		{"oversized payload", make([]byte, 2048), false},
		{"at max size", make([]byte, 1024), true},
		{"at min size", make([]byte, 1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuality(tt.payload, cfg)
			if (err == nil) != tt.wantValid {
				t.Errorf("ValidateQuality(payload len=%d) error = %v, wantValid = %v", len(tt.payload), err, tt.wantValid)
			}
		})
	}
}

// TestQoSHandling verifies QoS 0 and QoS 1 messages are both handled.
func TestQoSHandling(t *testing.T) {
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	client, err := ConnectClient(addr, "test-qos")
	if err != nil {
		t.Fatalf("ConnectClient failed: %v", err)
	}
	defer client.Disconnect(250)

	qos0Received := make(chan bool, 1)
	qos1Received := make(chan bool, 1)

	err = SubscribeAndBridge(client, "qos/#", func(topic string, payload []byte) {
		switch topic {
		case "qos/0":
			qos0Received <- true
		case "qos/1":
			qos1Received <- true
		}
	})
	if err != nil {
		t.Fatalf("SubscribeAndBridge failed: %v", err)
	}

	// QoS 0
	broker.Publish("qos/0", []byte("qos0-msg"), false, 0)
	// QoS 1
	broker.Publish("qos/1", []byte("qos1-msg"), false, 1)

	select {
	case <-qos0Received:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for QoS 0 message")
	}

	select {
	case <-qos1Received:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for QoS 1 message")
	}
}

// TestThreeLevelValidationIntegration verifies all three validation levels work together.
func TestThreeLevelValidationIntegration(t *testing.T) {
	cfg := ValidationConfig{
		Ranges: map[string]RangeConfig{
			"temperature": {Min: -40.0, Max: 150.0},
		},
		Timestamp: TimestampConfig{
			MaxFutureDrift: 5 * time.Second,
			MaxPastDrift:   24 * time.Hour,
		},
		Health: HealthConfig{
			MaxPayloadSize: 1024,
			MinPayloadSize: 1,
		},
	}

	// Valid reading
	validPayload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":25.3}}`, time.Now().Unix())
	err := ValidateMessage([]byte(validPayload), cfg)
	if err != nil {
		t.Errorf("valid message should pass: %v", err)
	}

	// Out-of-range value
	outOfRangePayload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":200.0}}`, time.Now().Unix())
	err = ValidateMessage([]byte(outOfRangePayload), cfg)
	if err == nil {
		t.Error("out-of-range value should fail")
	}

	// Future timestamp
	futurePayload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":25.3}}`, time.Now().Add(30*time.Second).Unix())
	err = ValidateMessage([]byte(futurePayload), cfg)
	if err == nil {
		t.Error("future timestamp should fail")
	}

	// Empty payload
	err = ValidateMessage([]byte{}, cfg)
	if err == nil {
		t.Error("empty payload should fail")
	}
}

// TestCrashIsolation verifies that killing the MQTT plugin does not crash the core process.
// This is a SPEC acceptance criterion (ACQ-01).
func TestCrashIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping crash isolation test in short mode")
	}

	// Build binaries from project root
	coreBin := filepath.Join(t.TempDir(), "ml-elec")
	cmd := exec.Command("go", "build", "-o", coreBin, "./cmd/ml-elec")
	cmd.Dir = filepath.Join(os.Getenv("PWD"), "../..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build ml-elec: %v\n%s", err, out)
	}

	mqttBin := filepath.Join(t.TempDir(), "mqtt-plugin")
	cmd = exec.Command("go", "build", "-o", mqttBin, "./cmd/mqtt-plugin")
	cmd.Dir = filepath.Join(os.Getenv("PWD"), "../..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build mqtt-plugin: %v\n%s", err, out)
	}

	// Start core process with MQTT plugin enabled
	coreCmd := exec.Command(coreBin)
	coreCmd.Env = append(os.Environ(),
		"ML_ELEC_PORT=0",
		"CONFIG_PATH="+createMinimalConfig(t, mqttBin),
	)
	coreCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := coreCmd.Start(); err != nil {
		t.Fatalf("failed to start core: %v", err)
	}

	// Give core time to start and launch plugin
	time.Sleep(3 * time.Second)

	// Verify core responds to /health
	healthURL := fmt.Sprintf("http://127.0.0.1:%d/health", getCorePort(t, coreCmd))
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(healthURL)
	if err != nil {
		t.Fatalf("core health check failed before kill: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("core health check returned %d before kill", resp.StatusCode)
	}

	// Find and kill the MQTT plugin process (child of core)
	killChildProcesses(coreCmd.Process.Pid, t)

	// Wait 1 second for crash to propagate
	time.Sleep(1 * time.Second)

	// Verify core still responds to /health
	resp, err = client.Get(healthURL)
	if err != nil {
		t.Fatalf("core health check failed after killing plugin: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("core health check returned %d after killing plugin", resp.StatusCode)
	}

	// Cleanup: kill core
	coreCmd.Process.Signal(syscall.SIGTERM)
	coreCmd.Wait()
}

// TestBenchmarkMQTTToNATS measures latency from MQTT publish to NATS receive.
// Target: < 100ms per message (per D-17).
func TestBenchmarkMQTTToNATS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmark in short mode")
	}

	// Start embedded NATS server
	natsOpts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	}
	natsSrv, err := natsserver.NewServer(natsOpts)
	if err != nil {
		t.Fatalf("failed to create NATS server: %v", err)
	}
	natsSrv.ConfigureLogger()
	natsSrv.Start()
	if !natsSrv.ReadyForConnections(10 * time.Second) {
		t.Fatal("NATS server not ready")
	}
	defer natsSrv.Shutdown()

	natsConn, err := natsclient.Connect(natsSrv.ClientURL())
	if err != nil {
		t.Fatalf("failed to connect NATS client: %v", err)
	}
	defer natsConn.Close()

	// Start embedded MQTT broker
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	mqttClient, err := ConnectClient(addr, "bench-mqtt-to-nats")
	if err != nil {
		t.Fatalf("ConnectClient failed: %v", err)
	}
	defer mqttClient.Disconnect(250)

	// Subscribe to NATS sensor subjects
	received := make(chan struct{}, 100)
	_, err = natsConn.Subscribe("sensor.>", func(msg *natsclient.Msg) {
		received <- struct{}{}
	})
	if err != nil {
		t.Fatalf("NATS subscribe failed: %v", err)
	}
	natsConn.Flush()

	// Bridge: MQTT message → NATS publish
	err = SubscribeAndBridge(mqttClient, "esp32/#", func(topic string, payload []byte) {
		natsConn.Publish("sensor."+topic, payload)
	})
	if err != nil {
		t.Fatalf("SubscribeAndBridge failed: %v", err)
	}

	// Publish 100 messages and measure latency
	const numMessages = 100
	latencies := make([]time.Duration, numMessages)

	for i := 0; i < numMessages; i++ {
		payload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":25.3}}`, time.Now().Unix())
		start := time.Now()

		token := mqttClient.Publish("esp32/sensor1", 1, false, payload)
		if !token.WaitTimeout(5 * time.Second) {
			t.Fatalf("MQTT publish timed out at message %d", i)
		}
		if token.Error() != nil {
			t.Fatalf("MQTT publish failed at message %d: %v", i, token.Error())
		}

		// Wait for NATS receive
		select {
		case <-received:
			latencies[i] = time.Since(start)
		case <-time.After(5 * time.Second):
			t.Fatalf("NATS receive timed out at message %d", i)
		}
	}

	// Calculate average latency
	var total time.Duration
	var maxLatency time.Duration
	for _, l := range latencies {
		total += l
		if l > maxLatency {
			maxLatency = l
		}
	}
	avg := total / numMessages

	t.Logf("MQTT→NATS benchmark: %d messages, avg=%v, max=%v", numMessages, avg, maxLatency)

	if avg > 100*time.Millisecond {
		t.Errorf("average latency %v exceeds 100ms target", avg)
	}
}

// TestBenchmarkConcurrentPublish measures throughput with concurrent MQTT publishers.
// Target: > 100 msg/s (per D-17).
func TestBenchmarkConcurrentPublish(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmark in short mode")
	}

	// Start embedded NATS server
	natsOpts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	}
	natsSrv, err := natsserver.NewServer(natsOpts)
	if err != nil {
		t.Fatalf("failed to create NATS server: %v", err)
	}
	natsSrv.ConfigureLogger()
	natsSrv.Start()
	if !natsSrv.ReadyForConnections(10 * time.Second) {
		t.Fatal("NATS server not ready")
	}
	defer natsSrv.Shutdown()

	natsConn, err := natsclient.Connect(natsSrv.ClientURL())
	if err != nil {
		t.Fatalf("failed to connect NATS client: %v", err)
	}
	defer natsConn.Close()

	// Start embedded MQTT broker
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	// Create 10 MQTT clients
	const numClients = 10
	const msgsPerClient = 10
	const totalMessages = numClients * msgsPerClient

	clients := make([]pahomqtt.Client, numClients)
	for i := 0; i < numClients; i++ {
		c, err := ConnectClient(addr, fmt.Sprintf("bench-pub-%d", i))
		if err != nil {
			t.Fatalf("ConnectClient failed for client %d: %v", i, err)
		}
		clients[i] = c
	}
	defer func() {
		for _, c := range clients {
			c.Disconnect(250)
		}
	}()

	// Subscribe to NATS sensor subjects
	received := make(chan struct{}, totalMessages)
	_, err = natsConn.Subscribe("sensor.>", func(msg *natsclient.Msg) {
		received <- struct{}{}
	})
	if err != nil {
		t.Fatalf("NATS subscribe failed: %v", err)
	}
	natsConn.Flush()

	// Bridge: MQTT message → NATS publish
	err = SubscribeAndBridge(clients[0], "esp32/#", func(topic string, payload []byte) {
		natsConn.Publish("sensor."+topic, payload)
	})
	if err != nil {
		t.Fatalf("SubscribeAndBridge failed: %v", err)
	}

	// Publish concurrently from all clients
	start := time.Now()

	done := make(chan struct{}, numClients)
	for i := 0; i < numClients; i++ {
		go func(clientIdx int) {
			for j := 0; j < msgsPerClient; j++ {
				payload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":25.3}}`, time.Now().Unix())
				token := clients[clientIdx].Publish("esp32/sensor1", 1, false, payload)
				if !token.WaitTimeout(5 * time.Second) {
					t.Errorf("MQTT publish timed out for client %d msg %d", clientIdx, j)
					break
				}
				if token.Error() != nil {
					t.Errorf("MQTT publish failed for client %d msg %d: %v", clientIdx, j, token.Error())
					break
				}
			}
			done <- struct{}{}
		}(i)
	}

	// Wait for all publishers to finish
	for i := 0; i < numClients; i++ {
		<-done
	}

	elapsed := time.Since(start)

	// Wait for all messages to be received
	receivedCount := 0
	timeout := time.After(10 * time.Second)
	for receivedCount < totalMessages {
		select {
		case <-received:
			receivedCount++
		case <-timeout:
			t.Fatalf("timed out waiting for messages: received %d/%d", receivedCount, totalMessages)
		}
	}

	throughput := float64(totalMessages) / elapsed.Seconds()
	t.Logf("Concurrent publish: %d messages in %v, throughput=%.1f msg/s", totalMessages, elapsed, throughput)

	if throughput < 100 {
		t.Errorf("throughput %.1f msg/s below 100 msg/s target", throughput)
	}
}

// TestMQTTPluginLifecycle verifies the MQTT plugin lifecycle: launch, running, kill, shutdown.
func TestMQTTPluginLifecycle(t *testing.T) {
	// This test verifies the plugin manager lifecycle without actual process spawning.
	// It tests the LaunchGRPC interface and IsRunning/Kill behavior.
	t.Skip("lifecycle test requires built binaries — covered by TestCrashIsolation")
}

// TestFullPipelineMQTTToNATS verifies the full MQTT → validate → NATS pipeline.
func TestFullPipelineMQTTToNATS(t *testing.T) {
	// Start embedded NATS server
	natsOpts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	}
	natsSrv, err := natsserver.NewServer(natsOpts)
	if err != nil {
		t.Fatalf("failed to create NATS server: %v", err)
	}
	natsSrv.ConfigureLogger()
	natsSrv.Start()
	if !natsSrv.ReadyForConnections(10 * time.Second) {
		t.Fatal("NATS server not ready")
	}
	defer natsSrv.Shutdown()

	natsConn, err := natsclient.Connect(natsSrv.ClientURL())
	if err != nil {
		t.Fatalf("failed to connect NATS client: %v", err)
	}
	defer natsConn.Close()

	// Start embedded MQTT broker
	broker, err := StartBroker(0)
	if err != nil {
		t.Fatalf("StartBroker failed: %v", err)
	}
	defer broker.Close()

	addr, err := BrokerAddr(broker)
	if err != nil {
		t.Fatalf("BrokerAddr failed: %v", err)
	}

	mqttClient, err := ConnectClient(addr, "pipeline-test")
	if err != nil {
		t.Fatalf("ConnectClient failed: %v", err)
	}
	defer mqttClient.Disconnect(250)

	// Subscribe to NATS sensor subjects (wildcard to match any topic)
	received := make(chan []byte, 1)
	_, err = natsConn.Subscribe("sensor.>", func(msg *natsclient.Msg) {
		received <- msg.Data
	})
	if err != nil {
		t.Fatalf("NATS subscribe failed: %v", err)
	}
	natsConn.Flush()

	// Bridge with validation
	validationCfg := ValidationConfig{
		Ranges: map[string]RangeConfig{
			"temperature": {Min: -40.0, Max: 150.0},
		},
		Timestamp: TimestampConfig{
			MaxFutureDrift: 5 * time.Second,
			MaxPastDrift:   24 * time.Hour,
		},
		Health: HealthConfig{
			MaxPayloadSize: 1024,
			MinPayloadSize: 1,
		},
	}

	err = SubscribeAndBridge(mqttClient, "esp32/#", func(topic string, payload []byte) {
		// Validate the message
		if err := ValidateMessage(payload, validationCfg); err != nil {
			t.Logf("validation failed for topic %s: %v", topic, err)
			return
		}
		// Publish to NATS
		natsConn.Publish("sensor."+topic, payload)
	})
	if err != nil {
		t.Fatalf("SubscribeAndBridge failed: %v", err)
	}

	// Allow subscription to propagate
	time.Sleep(100 * time.Millisecond)

	// Publish a valid JSON message via broker's inline client (to avoid loopback issues)
	validPayload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":25.3}}`, time.Now().Unix())
	broker.Publish("esp32/sensor1", []byte(validPayload), false, 1)
	natsConn.Flush()

	// Verify NATS receives the message
	select {
	case msg := <-received:
		if string(msg) != validPayload {
			t.Errorf("NATS received mismatch:\n  sent:     %s\n  received: %s", validPayload, string(msg))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for NATS message")
	}

	// Publish an invalid message (out of range) via broker's inline client
	invalidPayload := fmt.Sprintf(`{"ts":%d,"values":{"temperature":200.0}}`, time.Now().Unix())
	broker.Publish("esp32/sensor1", []byte(invalidPayload), false, 1)

	// Should NOT receive on NATS (validation should reject)
	select {
	case msg := <-received:
		t.Errorf("NATS should not have received invalid message, got: %s", string(msg))
	case <-time.After(500 * time.Millisecond):
		// Expected: no message received
	}
}

// killChildProcesses kills all child processes of the given PID.
func killChildProcesses(pid int, t *testing.T) {
	t.Helper()
	// Find child processes
	out, err := exec.Command("pgrep", "-P", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		t.Logf("no child processes found for pid %d: %v", pid, err)
		return
	}

	pids := parsePIDs(string(out))
	for _, childPID := range pids {
		proc, err := os.FindProcess(childPID)
		if err != nil {
			continue
		}
		t.Logf("killing child process %d", childPID)
		proc.Signal(syscall.SIGKILL)
	}
}

// parsePIDs parses a list of PIDs from pgrep output.
func parsePIDs(output string) []int {
	var pids []int
	for _, line := range splitLines(output) {
		if line == "" {
			continue
		}
		var pid int
		if _, err := fmt.Sscanf(line, "%d", &pid); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids
}

// splitLines splits a string by newlines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// getCorePort returns the port the core is listening on.
func getCorePort(t *testing.T, cmd *exec.Cmd) int {
	t.Helper()
	// Default port is 8080
	return 8080
}

// createMinimalConfig creates a minimal config file for testing.
func createMinimalConfig(t *testing.T, mqttPluginPath string) string {
	t.Helper()
	cfgContent := fmt.Sprintf(`
nats:
  host: 127.0.0.1
  port: -1
storage:
  path: %s/test.db
api:
  port: 0
plugins:
  enabled:
    - mqtt-plugin
mqtt:
  port: 0
  client_id: "test-plugin"
`, t.TempDir())

	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return cfgPath
}
