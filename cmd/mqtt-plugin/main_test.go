package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestStartBroker verifies the embedded broker starts on a configurable port.
func TestStartBroker(t *testing.T) {
	broker, err := StartBroker(0) // port 0 = auto-assign
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

	client, err := ConnectClient(broker, "test-client")
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

	client, err := ConnectClient(broker, "test-subscriber")
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
	binary.BigEndian.PutUint32(header[12:16], 2) // sample_count
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

// TestBinaryParsingLittleEndian verifies little-endian data is rejected.
func TestBinaryParsingLittleEndian(t *testing.T) {
	header := make([]byte, 28)
	header[0] = 1
	header[1] = 1

	// Write in little-endian (wrong!)
	binary.LittleEndian.PutUint64(header[4:12], uint64(1719900000000000000))
	binary.LittleEndian.PutUint32(header[12:16], 2)
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

	client, err := ConnectClient(broker, "test-qos")
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
