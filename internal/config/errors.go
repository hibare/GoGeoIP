package config

import "errors"

var (
	// ErrAPIListenPortInvalid indicates that the API listen port is invalid.
	ErrAPIListenPortInvalid = errors.New("API listen port must be between 1 and 65535")

	// ErrAPIKeysEmpty indicates that no API keys were provided.
	ErrAPIKeysEmpty = errors.New("at least one API key is required")

	// ErrDBLicenseKeyEmpty indicates that the database license key is empty.
	ErrAssetDirEmpty = errors.New("asset directory cannot be empty")
)

var (
	// ErrMaxMindLicenseKeyEmpty indicates that the MaxMind license key is empty.
	ErrMaxMindLicenseKeyEmpty = errors.New("MaxMind license key is required")

	// ErrMaxMindAutoUpdateIntervalInvalid indicates that the MaxMind auto-update interval is invalid.
	ErrMaxMindAutoUpdateIntervalInvalid = errors.New("MaxMind auto-update interval must be positive")
)

var (
	// ErrInvalidLogLevel indicates that the log level is invalid.
	ErrInvalidLogLevel = errors.New("invalid log level")

	// ErrInvalidLogMode indicates that the log mode is invalid.
	ErrInvalidLogMode = errors.New("invalid log mode")

	// ErrAssetDirCreation indicates that the asset directory could not be created.
	ErrAssetDirCreation = errors.New("failed to create asset directory")

	// ErrConfigUnmarshal indicates that the configuration could not be unmarshaled.
	ErrConfigUnmarshal = errors.New("failed to unmarshal configuration")
)
