package config

import (
	"errors"
	"time"
)

var (
	// ErrMaxMindLicenseKeyEmpty indicates that the MaxMind license key is empty.
	ErrMaxMindLicenseKeyEmpty = errors.New("MaxMind license key is required")

	// ErrMaxMindAutoUpdateIntervalInvalid indicates that the MaxMind auto-update interval is invalid.
	ErrMaxMindAutoUpdateIntervalInvalid = errors.New("MaxMind auto-update interval must be positive")
)

const (
	// DefaultMaxMindAutoUpdate is the default value for automatically updating the MaxMind GeoIP database.
	DefaultMaxMindAutoUpdate = true

	// DefaultMaxMindAutoUpdateInterval is the default interval for automatically updating the MaxMind GeoIP database.
	DefaultMaxMindAutoUpdateInterval = 24 * time.Hour
)

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
