package config

import (
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name: "all defaults",
			envVars: map[string]string{
				"GO_GEOIP_DB_LICENSE_KEY": "test-secret",
			},
			expected: &Config{
				Logger: LoggerConfig{
					Level: commonLogger.DefaultLoggerLevel,
					Mode:  commonLogger.DefaultLoggerMode,
				},
				DB: DBConfig{
					LicenseKey:         "test-secret",
					AutoUpdateEnabled:  true,
					AutoUpdateInterval: 24 * time.Hour,
				},
				Server: ServerConfig{
					ListenAddr:   "0.0.0.0",
					ListenPort:   5000,
					APIKeys:      []string{},
					ReadTimeout:  60 * time.Second,
					WriteTimeout: 15 * time.Second,
					IdleTimeout:  60 * time.Second,
				},
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"GO_GEOIP_LOG_LEVEL":              "debug",
				"GO_GEOIP_LOG_MODE":               "json",
				"GO_GEOIP_DB_LICENSE_KEY":         "test-secret",
				"GO_GEOIP_DB_AUTOUPDATE_ENABLED":  "false",
				"GO_GEOIP_DB_AUTOUPDATE_INTERVAL": "1h",
				"GO_GEOIP_SERVER_LISTEN_ADDR":     "127.0.0.1",
				"GO_GEOIP_SERVER_LISTEN_PORT":     "8080",
				"GO_GEOIP_SERVER_API_KEYS":        "test-key-1,test-key-2",
				"GO_GEOIP_SERVER_READ_TIMEOUT":    "30s",
				"GO_GEOIP_SERVER_WRITE_TIMEOUT":   "10s",
				"GO_GEOIP_SERVER_IDLE_TIMEOUT":    "120s",
			},
			expected: &Config{
				Logger: LoggerConfig{
					Level: "debug",
					Mode:  "json",
				},
				DB: DBConfig{
					LicenseKey:         "test-secret",
					AutoUpdateEnabled:  false,
					AutoUpdateInterval: 1 * time.Hour,
				},
				Server: ServerConfig{
					ListenAddr:   "127.0.0.1",
					ListenPort:   8080,
					APIKeys:      []string{"test-key-1", "test-key-2"},
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Load config
			Load()

			// Verify config
			assert.Equal(t, tt.expected, Current)
		})
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	// Set invalid log level
	t.Setenv("GO_GEOIP_LOG_LEVEL", "invalid-level")

	// Test that Load exits with invalid log level
	if os.Getenv("TEST_EXIT") == "1" {
		Load()
		return
	}
	const testName = "TestLoad_InvalidLogLevel"
	// #nosec G204
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLoad_InvalidLogMode(t *testing.T) {
	// Set invalid log mode
	t.Setenv("GO_GEOIP_LOG_MODE", "invalid-mode")

	// Test that Load exits with invalid log mode
	if os.Getenv("TEST_EXIT") == "1" {
		Load()
		return
	}
	const testName = "TestLoad_InvalidLogMode"
	// #nosec G204
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLoad_MissingLicenseKey(t *testing.T) {
	// Test that Load exits when license key is not set
	if os.Getenv("TEST_EXIT") == "1" {
		Load()
		return
	}
	_ = os.Unsetenv("GO_GEOIP_DB_LICENSE_KEY")

	const testName = "TestLoad_MissingLicenseKey"
	// #nosec G204
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
