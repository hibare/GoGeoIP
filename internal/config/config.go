package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/spf13/viper"
)

// ServerConfig holds API server-related configuration.
type ServerConfig struct {
	ListenAddr   string   `mapstructure:"listen_addr"`
	ListenPort   int      `mapstructure:"listen_port"`
	APIKeys      []string `mapstructure:"api_keys"`
	AssetDirPath string   `mapstructure:"asset_dir_path"`
	IsDev        bool     `mapstructure:"is_dev"`
}

// GetAddr returns the API server's listen address in "host:port" format.
func (s *ServerConfig) GetAddr() string {
	return net.JoinHostPort(s.ListenAddr, strconv.Itoa(s.ListenPort))
}

// PostProcess performs post-processing on the server configuration.
func (s *ServerConfig) PostProcess() {
	// Resolve asset dir path to absolute path
	absPath, err := filepath.Abs(s.AssetDirPath)
	if err == nil {
		s.AssetDirPath = absPath
	}
}

// Validate checks if the server configuration is valid.
func (s *ServerConfig) Validate() error {
	if s.ListenPort <= 0 || s.ListenPort > 65535 {
		return ErrAPIListenPortInvalid
	}
	if len(s.APIKeys) == 0 {
		return ErrAPIKeysEmpty
	}
	if s.AssetDirPath == "" {
		return ErrAssetDirEmpty
	}
	return nil
}

// LoggerConfig holds logger-related configuration.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Mode  string `mapstructure:"mode"`
}

// Validate checks if the logger configuration is valid.
func (l *LoggerConfig) Validate() error {
	if !commonLogger.IsValidLogLevel(l.Level) {
		return ErrInvalidLogLevel
	}
	if !commonLogger.IsValidLogMode(l.Mode) {
		return ErrInvalidLogMode
	}
	return nil
}

// MaxMindConfig holds MaxMind GeoIP database-related configuration.
type MaxMindConfig struct {
	LicenseKey         string        `mapstructure:"license_key"`
	AutoUpdate         bool          `mapstructure:"auto_update"`
	AutoUpdateInterval time.Duration `mapstructure:"auto_update_interval"`
}

// Validate checks if the MaxMind configuration is valid.
func (m *MaxMindConfig) Validate() error {
	if m.LicenseKey == "" {
		return ErrMaxMindLicenseKeyEmpty
	}
	if m.AutoUpdate && m.AutoUpdateInterval <= 0 {
		return ErrMaxMindAutoUpdateIntervalInvalid
	}
	return nil
}

// Config holds the entire application configuration.
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MaxMind MaxMindConfig `mapstructure:"maxmind"`
	Logger  LoggerConfig  `mapstructure:"logger"`
}

// Validate validates the entire configuration.
func (c *Config) Validate() error {
	var vFuncs = []func() error{
		c.MaxMind.Validate,
		c.Server.Validate,
		c.Logger.Validate,
	}

	for _, vf := range vFuncs {
		if err := vf(); err != nil {
			return err
		}
	}

	return nil
}

// PostProcess performs post-processing on the configuration.
func (c *Config) PostProcess() {
	c.Server.PostProcess()
}

func (c *Config) getViper(ctx context.Context, configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Config search paths.
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/gogeoip/")
	}

	// Environment variable binding.
	v.SetEnvPrefix("GOGEOIP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Environment variable bindings.
	envBindings := map[string]string{
		"server.listen_addr":           "SERVER_LISTEN_ADDR",
		"server.listen_port":           "SERVER_LISTEN_PORT",
		"server.api_keys":              "API_KEYS",
		"server.asset_dir_path":        "ASSET_DIR_PATH",
		"server.is_dev":                "IS_DEV",
		"logger.level":                 "LOG_LEVEL",
		"logger.mode":                  "LOG_MODE",
		"maxmind.license_key":          "MAXMIND_LICENSE_KEY",
		"maxmind.auto_update":          "MAXMIND_AUTOUPDATE",
		"maxmind.auto_update_interval": "MAXMIND_AUTOUPDATE_INTERVAL",
	}

	for key, envVar := range envBindings {
		if err := v.BindEnv(key, envVar); err != nil {
			slog.WarnContext(ctx, "Failed to bind environment variable",
				slog.String("config", key),
				slog.String("env", envVar),
				slog.String("error", err.Error()))
		}
	}

	// Set default values.
	v.SetDefault("server.listen_addr", DefaultServerListenAddr)
	v.SetDefault("server.listen_port", DefaultServerListenPort)
	v.SetDefault("server.api_keys", []string{uuid.New().String()})
	v.SetDefault("server.asset_dir_path", constants.AssetDir)
	v.SetDefault("server.is_dev", false)
	v.SetDefault("logger.level", commonLogger.LogLevelInfo)
	v.SetDefault("logger.mode", commonLogger.LogModePretty)
	v.SetDefault("maxmind.license_key", "")
	v.SetDefault("maxmind.auto_update", DefaultMaxMindAutoUpdate)
	v.SetDefault("maxmind.auto_update_interval", DefaultMaxMindAutoUpdateInterval)

	return v
}

// Current holds the current application configuration.
var Current *Config

// Load loads the configuration from environment variables and/or config file.
func Load(ctx context.Context, configPath string) (*Config, error) {
	cfg := &Config{}
	v := cfg.getViper(ctx, configPath)

	// Try read config.
	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundErr) {
			slog.WarnContext(ctx, "No config file found, relying on env vars/defaults")
		} else {
			return nil, err
		}
	} else {
		slog.DebugContext(ctx, "Using config file", slog.String("file", v.ConfigFileUsed()))
	}

	// Unmarshal into config.
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Handle comma-separated API keys from environment variable
	if apiKeysStr := v.GetString("server.api_keys"); apiKeysStr != "" {
		// If it's a comma-separated string, split it
		if strings.Contains(apiKeysStr, ",") {
			keys := strings.Split(apiKeysStr, ",")
			for i, key := range keys {
				keys[i] = strings.TrimSpace(key)
			}
			cfg.Server.APIKeys = keys
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Post-process config.
	cfg.PostProcess()

	// Initialize logger.
	commonLogger.InitLogger(&cfg.Logger.Level, &cfg.Logger.Mode)

	// Create asset dir
	if err := os.MkdirAll(cfg.Server.AssetDirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAssetDirCreation, err)
	}

	Current = cfg

	slog.InfoContext(ctx, "Loaded config")
	return cfg, nil
}

// LoadDefault loads the configuration with default settings (no context required).
func LoadDefault() error {
	ctx := context.Background()
	_, err := Load(ctx, "")
	return err
}

// GenerateConfigFile generates a config file with default values.
func GenerateConfigFile(ctx context.Context, configPath string) (string, error) {
	cfg := &Config{}
	v := cfg.getViper(ctx, configPath)

	// Unmarshal viper's defaults into the config struct
	if err := v.Unmarshal(cfg); err != nil {
		return "", fmt.Errorf("%w: %w", ErrConfigUnmarshal, err)
	}

	return v.ConfigFileUsed(), v.WriteConfig()
}
