package config

import (
	"log"
	"path/filepath"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/env"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/GoGeoIP/internal/constants"
)

// LoggerConfig defines logging configuration parameters.
type LoggerConfig struct {
	Level string
	Mode  string
}

// APIConfig defines API configuration parameters.
type DBConfig struct {
	LicenseKey         string
	AutoUpdateEnabled  bool
	AutoUpdateInterval time.Duration
}

// ServerConfig defines API configuration parameters.
type ServerConfig struct {
	ListenAddr   string
	ListenPort   int
	APIKeys      []string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Config holds the configuration for the application.
type Config struct {
	Logger LoggerConfig
	DB     DBConfig
	Server ServerConfig
}

func (c Config) GetAssetDirPath() (string, error) {
	return filepath.Abs(constants.AssetDir)
}

var Current *Config

func Load() {
	env.Load()

	Current = &Config{
		Logger: LoggerConfig{
			Level: env.MustString("GO_GEOIP_LOG_LEVEL", commonLogger.DefaultLoggerLevel),
			Mode:  env.MustString("GO_GEOIP_LOG_MODE", commonLogger.DefaultLoggerMode),
		},
		DB: DBConfig{
			LicenseKey:         env.MustString("GO_GEOIP_DB_LICENSE_KEY", ""),
			AutoUpdateEnabled:  env.MustBool("GO_GEOIP_DB_AUTOUPDATE_ENABLED", constants.DefaultDBAutoUpdate),
			AutoUpdateInterval: env.MustDuration("GO_GEOIP_DB_AUTOUPDATE_INTERVAL", constants.DefaultDBAutoUpdateInterval),
		},
		Server: ServerConfig{
			ListenAddr:   env.MustString("GO_GEOIP_SERVER_LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort:   env.MustInt("GO_GEOIP_SERVER_LISTEN_PORT", constants.DefaultAPIListenPort),
			APIKeys:      env.MustStringSlice("GO_GEOIP_SERVER_API_KEYS", []string{}),
			ReadTimeout:  env.MustDuration("GO_GEOIP_SERVER_READ_TIMEOUT", commonHttp.DefaultServerTimeout),
			WriteTimeout: env.MustDuration("GO_GEOIP_SERVER_WRITE_TIMEOUT", commonHttp.DefaultServerWriteTimeout),
			IdleTimeout:  env.MustDuration("GO_GEOIP_SERVER_IDLE_TIMEOUT", commonHttp.DefaultServerIdleTimeout),
		},
	}

	if len(Current.DB.LicenseKey) <= 0 {
		log.Fatal("GO_GEOIP_DB_LICENSE_KEY env var is required")
	}

	if !commonLogger.IsValidLogLevel(Current.Logger.Level) {
		log.Fatal("Error invalid logger level")
	}

	if !commonLogger.IsValidLogMode(Current.Logger.Mode) {
		log.Fatal("Error invalid logger mode")
	}

	commonLogger.InitLogger(&Current.Logger.Level, &Current.Logger.Mode)
}
