package nats

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsclient "github.com/nats-io/nats.go"
)

// Server wraps the NATS server for lifecycle management.
type Server struct {
	srv   *server.Server
	ready bool
}

// New creates a new NATS server with the given port (-1 for random).
func New(port int) (*Server, error) {
	opts := &server.Options{
		Host:       "127.0.0.1",
		Port:       port,
		MaxConn:    1024,
		MaxPayload: 1048576,
	}

	s, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("creating nats server: %w", err)
	}

	return &Server{srv: s}, nil
}

// Start starts the NATS server and waits for it to be ready.
func (s *Server) Start() error {
	s.srv.ConfigureLogger()
	s.srv.Start()

	if !s.srv.ReadyForConnections(10 * time.Second) {
		return errors.New("nats server not ready")
	}

	s.ready = true
	return nil
}

// Ready returns true if the server is ready for connections.
func (s *Server) Ready() bool {
	return s.ready
}

// Client returns a NATS client connected to the embedded server.
func (s *Server) Client() (*natsclient.Conn, error) {
	return natsclient.Connect(s.srv.ClientURL())
}

// Shutdown stops the server cleanly.
func (s *Server) Shutdown() {
	s.ready = false
	s.srv.Shutdown()
	s.srv.WaitForShutdown()
}
