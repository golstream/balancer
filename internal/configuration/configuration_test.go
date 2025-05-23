package configuration

import (
	"os"
	"reflect"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("failed to set env %s: %v", key, err)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatalf("failed to unset env %s: %v", key, err)
	}
}

func TestInit_Success(t *testing.T) {
	setEnv(t, "PORT", "8080")
	setEnv(t, "METHOD", "round_robin")
	setEnv(t, "SERVERS", "127.0.0.1:8000,127.0.0.1:8001")
	setEnv(t, "WEIGHTS", "1,2")
	setEnv(t, "WITH_LOG", "true")
	setEnv(t, "HEALTHCHECK_INTERVAL", "30")
	setEnv(t, "HEALTHCHECK_TIMEOUT", "10")

	defer func() {
		unsetEnv(t, "PORT")
		unsetEnv(t, "METHOD")
		unsetEnv(t, "SERVERS")
		unsetEnv(t, "WEIGHTS")
		unsetEnv(t, "WITH_LOG")
		unsetEnv(t, "HEALTHCHECK_INTERVAL")
		unsetEnv(t, "HEALTHCHECK_TIMEOUT")
	}()

	cfg, err := Init()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}

	if cfg.Method != "round_robin" {
		t.Errorf("expected method round_robin, got %s", cfg.Method)
	}

	expectedServers := servers{"127.0.0.1:8000", "127.0.0.1:8001"}
	if !reflect.DeepEqual(cfg.Servers, expectedServers) {
		t.Errorf("expected servers %v, got %v", expectedServers, cfg.Servers)
	}

	expectedWeights := weights{"1", "2"}
	if !reflect.DeepEqual(cfg.Weights, expectedWeights) {
		t.Errorf("expected weights %v, got %v", expectedWeights, cfg.Weights)
	}

	if !cfg.WithLog {
		t.Errorf("expected WithLog true, got false")
	}

	if cfg.HealthCheckInterval != 30 {
		t.Errorf("expected HealthCheckInterval 30, got %d", cfg.HealthCheckInterval)
	}

	if cfg.HealthCheckTimeout != 10 {
		t.Errorf("expected HealthCheckTimeout 10, got %d", cfg.HealthCheckTimeout)
	}
}

func TestInit_MissingRequiredField(t *testing.T) {
	unsetEnv(t, "PORT") // required
	setEnv(t, "METHOD", "least_connections")
	setEnv(t, "SERVERS", "127.0.0.1:8000")
	setEnv(t, "WEIGHTS", "1")

	defer func() {
		unsetEnv(t, "METHOD")
		unsetEnv(t, "SERVERS")
		unsetEnv(t, "WEIGHTS")
	}()

	_, err := Init()
	if err == nil {
		t.Fatal("expected error due to missing required PORT")
	}
}
