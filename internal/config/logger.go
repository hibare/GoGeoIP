package config

import (
	"errors"

	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
)

var (
	// ErrInvalidLogLevel indicates that the log level is invalid.
	ErrInvalidLogLevel = errors.New("invalid log level")

	// ErrInvalidLogMode indicates that the log mode is invalid.
	ErrInvalidLogMode = errors.New("invalid log mode")

	// ErrDataDirCreation indicates that the data directory could not be created.
	ErrDataDirCreation = errors.New("failed to create data directory")

	// ErrConfigUnmarshal indicates that the configuration could not be unmarshaled.
	ErrConfigUnmarshal = errors.New("failed to unmarshal configuration")
)

const (
	// DefaultLoggerLevel is the default log level for the application.
	DefaultLoggerLevel = commonLogger.LogLevelInfo

	// DefaultLoggerMode is the default log mode for the application.
	DefaultLoggerMode = commonLogger.LogModePretty
)

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
