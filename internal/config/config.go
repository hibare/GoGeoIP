package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

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

// DBConfig defines database configuration parameters.
type UtilConfig struct {
	AssetDirPath string
	IsDev        bool
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
	Util   UtilConfig
}

func (c Config) GetAssetDirPath() (string, error) {
	return filepath.Abs(constants.AssetDir)
}

var Current *Config

func Load() {
	env.Load()

	token := []string{
		uuid.New().String(),
	}

	assetDir, err := filepath.Abs(constants.AssetDir)

	if err != nil {
		log.Fatalf("Unable to load config, %s", err.Error())
	}

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
			ListenAddr:   env.MustString("GO_GEOIP_API_LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort:   env.MustInt("GO_GEOIP_API_LISTEN_PORT", constants.DefaultAPIListenPort),
			APIKeys:      env.MustStringSlice("GO_GEOIP_API_KEYS", token),
			ReadTimeout:  env.MustDuration("GO_GEOIP_API_READ_TIMEOUT", commonHttp.DefaultServerTimeout),
			WriteTimeout: env.MustDuration("GO_GEOIP_API_WRITE_TIMEOUT", commonHttp.DefaultServerWriteTimeout),
			IdleTimeout:  env.MustDuration("GO_GEOIP_API_IDLE_TIMEOUT", commonHttp.DefaultServerIdleTimeout),
		},
		Util: UtilConfig{
			AssetDirPath: assetDir,
			IsDev:        env.MustBool("IS_DEV", false),
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

	// Create asset dir
	if err := os.MkdirAll(constants.AssetDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create asset dir: %s", err)
	}

	log.Info("Loaded config")
}
