package v1

import (
	"context"
	"net"
	"testing"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// testLifecycleImpl is a test implementation of PluginLifecycle.
type testLifecycleImpl struct {
	initConfig map[string]string
	initErr    error
}

func (t *testLifecycleImpl) Init(_ context.Context, config map[string]string) error {
	t.initConfig = config
	return t.initErr
}

func (t *testLifecycleImpl) Start(_ context.Context) error {
	return nil
}

func (t *testLifecycleImpl) Stop(_ context.Context) error {
	return nil
}

// testCollectorImpl is a test implementation of SensorCollector.
type testCollectorImpl struct{}

func (t *testCollectorImpl) Collect(_ context.Context, req *CollectRequest) (*CollectResponse, error) {
	return &CollectResponse{
		Accepted:       true,
		RejectionReason: "",
	}, nil
}

// setupGRPCServer creates an in-memory gRPC server with the GRPCPlugin registered.
func setupGRPCServer(t *testing.T, lifecycle PluginLifecycle, collector SensorCollector) *bufconn.Listener {
	t.Helper()

	lis := bufconn.Listen(bufSize)

	srv := grpc.NewServer()
	RegisterPluginLifecycleServer(srv, &grpcServer{
		impl:       lifecycle,
		collectImpl: collector,
	})
	RegisterSensorCollectorServer(srv, &grpcServer{
		impl:       lifecycle,
		collectImpl: collector,
	})

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("server exited: %v", err)
		}
	}()

	t.Cleanup(func() {
		srv.Stop()
		lis.Close()
	})

	return lis
}

// setupGRPCClient creates a gRPC client connected to the bufconn listener.
func setupGRPCClient(t *testing.T, lis *bufconn.Listener) (*grpc.ClientConn, PluginLifecycleClient, SensorCollectorClient) {
	t.Helper()

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	return conn, NewPluginLifecycleClient(conn), NewSensorCollectorClient(conn)
}

func TestGRPCPlugin_ImplementsInterface(t *testing.T) {
	// Compile-time check: GRPCPlugin must implement goplugin.GRPCPlugin
	var _ goplugin.GRPCPlugin = (*GRPCPlugin)(nil)
}

func TestGRPCPlugin_GRPCServer_RegistersBothServices(t *testing.T) {
	lifecycle := &testLifecycleImpl{}
	collector := &testCollectorImpl{}
	p := &GRPCPlugin{
		Impl:        lifecycle,
		CollectImpl: collector,
	}

	srv := grpc.NewServer()
	err := p.GRPCServer(nil, srv)
	if err != nil {
		t.Fatalf("GRPCServer returned error: %v", err)
	}

	// Verify services are registered by checking the service info
	info := srv.GetServiceInfo()
	if _, ok := info["sdk.v1.PluginLifecycle"]; !ok {
		t.Error("PluginLifecycle service not registered")
	}
	if _, ok := info["sdk.v1.SensorCollector"]; !ok {
		t.Error("SensorCollector service not registered")
	}
}

func TestGRPCPlugin_GRPCClient_ReturnsInterfaces(t *testing.T) {
	p := &GRPCPlugin{}

	// Create a bufconn-based server for the client to connect to
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	lifecycle := &testLifecycleImpl{}
	collector := &testCollectorImpl{}
	RegisterPluginLifecycleServer(srv, &grpcServer{impl: lifecycle, collectImpl: collector})
	RegisterSensorCollectorServer(srv, &grpcServer{impl: lifecycle, collectImpl: collector})
	go srv.Serve(lis)
	t.Cleanup(func() { srv.Stop() })

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	client, err := p.GRPCClient(context.Background(), nil, conn)
	if err != nil {
		t.Fatalf("GRPCClient returned error: %v", err)
	}

	// Verify the returned client satisfies both interfaces
	if _, ok := client.(PluginLifecycle); !ok {
		t.Error("GRPCClient result does not implement PluginLifecycle")
	}
	if _, ok := client.(SensorCollector); !ok {
		t.Error("GRPCClient result does not implement SensorCollector")
	}
}

func TestLifecycle_Init_Start_Stop(t *testing.T) {
	lifecycle := &testLifecycleImpl{}
	collector := &testCollectorImpl{}

	lis := setupGRPCServer(t, lifecycle, collector)
	_, lifecycleClient, _ := setupGRPCClient(t, lis)

	ctx := context.Background()

	// Init with config
	config := map[string]string{"broker": "tcp://localhost:1883", "topic": "sensors/+"}
	_, err := lifecycleClient.Init(ctx, &InitRequest{Config: config})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Verify config was passed through
	if len(lifecycle.initConfig) != 2 {
		t.Errorf("expected 2 config entries, got %d", len(lifecycle.initConfig))
	}
	if lifecycle.initConfig["broker"] != "tcp://localhost:1883" {
		t.Errorf("expected broker=tcp://localhost:1883, got %q", lifecycle.initConfig["broker"])
	}

	// Start
	_, err = lifecycleClient.Start(ctx, &StartRequest{})
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Stop
	_, err = lifecycleClient.Stop(ctx, &StopRequest{})
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestCollector_Collect(t *testing.T) {
	lifecycle := &testLifecycleImpl{}
	collector := &testCollectorImpl{}

	lis := setupGRPCServer(t, lifecycle, collector)
	_, _, collectorClient := setupGRPCClient(t, lis)

	resp, err := collectorClient.Collect(context.Background(), &CollectRequest{
		SensorId:      "temp-001",
		Topic:         "factory/site1/line1/motor1/temperature",
		Payload:       []byte(`{"ts":1234567890,"values":{"temperature":65.2}}`),
		PayloadFormat: "json",
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}
	if !resp.Accepted {
		t.Errorf("expected accepted=true, got accepted=false, reason=%q", resp.RejectionReason)
	}
}

func TestLifecycle_Init_Error_Propagation(t *testing.T) {
	lifecycle := &testLifecycleImpl{initErr: &testError{msg: "config invalid"}}
	collector := &testCollectorImpl{}

	lis := setupGRPCServer(t, lifecycle, collector)
	_, lifecycleClient, _ := setupGRPCClient(t, lis)

	_, err := lifecycleClient.Init(context.Background(), &InitRequest{
		Config: map[string]string{},
	})
	if err == nil {
		t.Fatal("expected error from Init, got nil")
	}
	if err.Error() != "config invalid" {
		t.Errorf("expected error message 'config invalid', got %q", err.Error())
	}
}

// testError is a simple error type for testing.
type testError struct {
	msg string
}

func (e *testError) Error() string { return e.msg }
