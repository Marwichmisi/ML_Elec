package nats

import (
	"testing"
	"time"

	"ml-elec/internal/config"
)

func TestNewCreatesServer(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer srv.Shutdown()
}

func TestStartMakesServerReady(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer srv.Shutdown()

	err = srv.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	if !srv.Ready() {
		t.Error("server should be ready after Start()")
	}
}

func TestClientReturnsWorkingConnection(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer srv.Shutdown()

	err = srv.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	conn, err := srv.Client()
	if err != nil {
		t.Fatalf("Client() failed: %v", err)
	}
	defer conn.Close()

	if !conn.IsConnected() {
		t.Error("client should be connected")
	}
}

func TestPublishSubscribe(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer srv.Shutdown()

	err = srv.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	conn, err := srv.Client()
	if err != nil {
		t.Fatalf("Client() failed: %v", err)
	}
	defer conn.Close()

	// Subscribe
	sub, err := conn.SubscribeSync("test.subject")
	if err != nil {
		t.Fatalf("SubscribeSync() failed: %v", err)
	}
	defer func() { _ = sub.Unsubscribe() }()

	// Publish
	err = conn.Publish("test.subject", []byte("hello"))
	if err != nil {
		t.Fatalf("Publish() failed: %v", err)
	}

	// Wait for message
	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		t.Fatalf("NextMsg() failed: %v", err)
	}

	if string(msg.Data) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(msg.Data))
	}
}

func TestShutdownStopsServerCleanly(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	err = srv.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	srv.Shutdown()

	if srv.Ready() {
		t.Error("server should not be ready after Shutdown()")
	}
}

func TestConcurrentPubSub(t *testing.T) {
	t.Parallel()
	srv, err := New(&config.NATSConfig{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer srv.Shutdown()

	err = srv.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	conn, err := srv.Client()
	if err != nil {
		t.Fatalf("Client() failed: %v", err)
	}
	defer conn.Close()

	// Subscribe
	sub, err := conn.SubscribeSync("concurrent.test")
	if err != nil {
		t.Fatalf("SubscribeSync() failed: %v", err)
	}
	defer func() { _ = sub.Unsubscribe() }()

	// Publish from multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			msg := []byte("message")
			err := conn.Publish("concurrent.test", msg)
			if err != nil {
				t.Errorf("Publish() failed: %v", err)
			}
		}(i)
	}

	// Wait for all publishes
	for i := 0; i < 10; i++ {
		<-done
	}

	// Receive messages
	received := 0
	for received < 10 {
		_, err := sub.NextMsg(5 * time.Second)
		if err != nil {
			t.Fatalf("NextMsg() failed after receiving %d messages: %v", received, err)
		}
		received++
	}

	if received != 10 {
		t.Errorf("expected 10 messages, got %d", received)
	}
}
