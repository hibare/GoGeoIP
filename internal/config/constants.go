package config

import (
	"time"

	"github.com/hibare/GoCommon/v2/pkg/logger"
)

// Default configuration values.
const (
	DefaultServerListenAddr          = "0.0.0.0"
	DefaultServerListenPort          = 5000
	DefaultMaxMindAutoUpdate         = true
	DefaultMaxMindAutoUpdateInterval = 24 * time.Hour
	DefaultAssetDirPath              = "./data"
	DefaultLoggerLevel               = logger.LogLevelInfo
	DefaultLoggerMode                = logger.LogModePretty
)
