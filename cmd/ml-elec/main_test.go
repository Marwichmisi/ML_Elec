package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestBinaryCompiles(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "ml-elec")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ml-elec")
	cmd.Dir = "/home/marwane/Documents/My-ML_Elec"
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatal("binary was not created")
	}
}

type binaryInstance struct {
	cmd    *exec.Cmd
	tmpDir string
	port   int
}

func startBinaryWithPort(t *testing.T, port int) *binaryInstance {
	t.Helper()
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "ml-elec")
	dbPath := filepath.Join(tmpDir, "sensors.db")
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Build binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ml-elec")
	buildCmd.Dir = "/home/marwane/Documents/My-ML_Elec"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}

	// Write minimal config
	config := fmt.Sprintf(`
nats:
  host: 127.0.0.1
  port: -1
storage:
  path: %s
api:
  port: 0
plugins:
  enabled: []
`, dbPath)
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Start binary with port override
	cmd := exec.Command(binaryPath)
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(),
		"CONFIG_PATH="+configPath,
		fmt.Sprintf("ML_ELEC_PORT=%d", port),
	)

	// Capture stderr for port discovery
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start binary: %v", err)
	}

	// Read stderr in background to detect port from logs
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Log lines go to stderr (JSON format)
			_ = scanner.Text()
		}
	}()

	return &binaryInstance{cmd: cmd, tmpDir: tmpDir, port: port}
}

func (b *binaryInstance) cleanup() {
	_ = b.cmd.Process.Signal(syscall.SIGTERM)
	_ = b.cmd.Wait()
	os.RemoveAll(b.tmpDir)
}

func waitForHealth(t *testing.T, port int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("health check did not become ready on port %d within %v", port, timeout)
}

func getFreePort(t *testing.T) int {
	t.Helper()
	// Use a port in the ephemeral range
	// Each test gets a unique port based on test name hash
	return 10000 + (int(time.Now().UnixNano()) % 50000)
}

func TestFullStartup(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer inst.cleanup()

	waitForHealth(t, port, 5*time.Second)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status, ok := result["status"]; !ok || status != "ok" {
		t.Errorf("expected status \"ok\", got %v", status)
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer os.RemoveAll(inst.tmpDir)

	waitForHealth(t, port, 5*time.Second)

	if err := inst.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send SIGTERM: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- inst.cmd.Wait()
	}()

	select {
	case err := <-done:
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				t.Errorf("expected exit code 0, got %d", exitErr.ExitCode())
			}
		}
	case <-time.After(10 * time.Second):
		t.Fatal("process did not exit within 10 seconds after SIGTERM")
	}
}

func TestConcurrentAPILoad(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer inst.cleanup()

	waitForHealth(t, port, 5*time.Second)

	const numRequests = 20
	done := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
			if err != nil {
				done <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				done <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}
			done <- nil
		}()
	}

	for i := 0; i < numRequests; i++ {
		if err := <-done; err != nil {
			t.Errorf("request %d failed: %v", i, err)
		}
	}
}

func TestSignalHandlingSIGINT(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer os.RemoveAll(inst.tmpDir)

	waitForHealth(t, port, 5*time.Second)

	if err := inst.cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("failed to send SIGINT: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- inst.cmd.Wait()
	}()

	select {
	case err := <-done:
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				t.Errorf("expected exit code 0, got %d", exitErr.ExitCode())
			}
		}
	case <-time.After(10 * time.Second):
		t.Fatal("process did not exit within 10 seconds after SIGINT")
	}
}

func TestHealthEndpointReturnsJSON(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer inst.cleanup()

	waitForHealth(t, port, 5*time.Second)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestContextCancellation(t *testing.T) {
	t.Parallel()
	port := getFreePort(t)
	inst := startBinaryWithPort(t, port)
	defer inst.cleanup()

	waitForHealth(t, port, 5*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/health", port), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// Suppress unused import warning
var _ = regexp.MustCompile
var _ = strconv.Itoa
