package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/hibare/GoCommon/v2/pkg/env"
	"github.com/hibare/GoGeoIP/internal/constants"
)

type UtilConfig struct {
	AssetDirPath string
	IsDev        bool
}

type DBConfig struct {
	LicenseKey         string
	AutoUpdate         bool
	AutoUpdateInterval time.Duration
}

type APIConfig struct {
	ListenAddr string
	ListenPort int
	APIKeys    []string
}

type Config struct {
	DB   DBConfig
	API  APIConfig
	Util UtilConfig
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
		DB: DBConfig{
			LicenseKey:         env.MustString("DB_LICENSE_KEY", ""),
			AutoUpdate:         env.MustBool("DB_AUTOUPDATE", constants.DefaultDBAutoUpdate),
			AutoUpdateInterval: env.MustDuration("DB_AUTOUPDATE_INTERVAL", constants.DefaultDBAutoUpdateInterval),
		},
		API: APIConfig{
			ListenAddr: env.MustString("API_LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort: env.MustInt("API_LISTEN_PORT", constants.DefaultAPIListenPort),
			APIKeys:    env.MustStringSlice("API_KEYS", token),
		},
		Util: UtilConfig{
			AssetDirPath: assetDir,
			IsDev:        env.MustBool("IS_DEV", false),
		},
	}

	if len(Current.DB.LicenseKey) <= 0 {
		log.Fatal("DB_LICENSE_KEY is required")
	}

	// Create asset dir
	if err := os.MkdirAll(constants.AssetDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create asset dir: %s", err)
	}

	log.Info("Loaded config")
}
