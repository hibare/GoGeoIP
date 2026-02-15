package config

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrAPIKeysEmpty indicates that no API keys were provided.
	ErrAPIKeysEmpty = errors.New("at least one API key is required")

	//  ErrInvalidPort indicates that the server port is invalid.
	ErrInvalidPort = errors.New("invalid server port. Port must be between 1 and 65535")

	// ErrInvalidBaseURL indicates that the server base URL is invalid.
	ErrInvalidBaseURL = errors.New("invalid server base URL. Base URL must start with http:// or https://")

	// ErrJWTSecretEmpty indicates that the JWT secret is empty.
	ErrJWTSecretEmpty = errors.New("jwt secret is required")

	// ErrAPIListenPortInvalid is an alias for ErrInvalidPort for backwards compatibility.
	ErrAPIListenPortInvalid = ErrInvalidPort
)

const (
	// DefaultServerListenAddr is the default server listen address.
	DefaultServerListenAddr = "0.0.0.0"

	// DefaultServerListenPort is the default server listen port.
	DefaultServerListenPort = 5000

	// DefaultServerReadTimeout is the default server read timeout.
	DefaultServerReadTimeout = 15 * time.Second

	// DefaultServerWriteTimeout is the default server write timeout.
	DefaultServerWriteTimeout = 15 * time.Second

	// DefaultServerIdleTimeout is the default server idle timeout.
	DefaultServerIdleTimeout = 60 * time.Second

	// DefaultServerWaitTimeout is the default server wait time.
	DefaultServerWaitTimeout = 15 * time.Second

	// DefaultServerRequestTimeout is the default request timeout.
	DefaultServerRequestTimeout = 60 * time.Second
)

// ServerConfig holds API server-related configuration.
type ServerConfig struct {
	ListenAddr     string        `mapstructure:"listen_addr"`
	ListenPort     int           `mapstructure:"listen_port"`
	BaseURL        string        `mapstructure:"base_url"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	WaitTimeout    time.Duration `mapstructure:"wait_timeout"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	CertFile       string        `mapstructure:"cert_file"`
	KeyFile        string        `mapstructure:"key_file"`
	APIKeys        []string      `mapstructure:"api_keys"`
}

// GetAddr returns the API server's listen address in "host:port" format.
func (s *ServerConfig) GetAddr() string {
	return net.JoinHostPort(s.ListenAddr, strconv.Itoa(s.ListenPort))
}

// PostProcess performs post-processing on the server configuration.
func (s *ServerConfig) PostProcess() {
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
}

// Validate checks if the server configuration is valid.
func (s *ServerConfig) Validate() error {
	if s.ListenPort <= 0 || s.ListenPort > 65535 {
		return ErrInvalidPort
	}

	if len(s.APIKeys) == 0 {
		return ErrAPIKeysEmpty
	}

	if s.BaseURL != "" {
		// Basic URL validation - should start with http:// or https://
		if !strings.HasPrefix(s.BaseURL, "http://") && !strings.HasPrefix(s.BaseURL, "https://") {
			return ErrInvalidBaseURL
		}
	}

	return nil
}
