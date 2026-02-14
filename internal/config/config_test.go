package config

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMaxMindLicenseKey         = "test-license"
	testMaxMindAutoUpdate         = false
	testMaxMindAutoUpdateInterval = 2 * time.Hour
	testAPIListenAddr             = "127.0.0.1"
	testAPIListenPort             = 10000
	testAPIKeys                   = "test-api-key"
	testIsDev                     = true
)

func TestEnvLoadedConfig(t *testing.T) {
	ctx := context.Background()

	// Set test env vars
	t.Setenv("MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)
	t.Setenv("MAXMIND_AUTOUPDATE", strconv.FormatBool(testMaxMindAutoUpdate))
	t.Setenv("MAXMIND_AUTOUPDATE_INTERVAL", testMaxMindAutoUpdateInterval.String())
	t.Setenv("SERVER_LISTEN_ADDR", testAPIListenAddr)
	t.Setenv("SERVER_LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	t.Setenv("API_KEYS", testAPIKeys)
	t.Setenv("IS_DEV", strconv.FormatBool(testIsDev))

	_, err := Load(ctx, "")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(constants.AssetDir)
	}()

	assert.Equal(t, testAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, testAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, []string{testAPIKeys}, Current.Server.APIKeys)
	assert.Equal(t, testMaxMindLicenseKey, Current.MaxMind.LicenseKey)
	assert.Equal(t, testMaxMindAutoUpdate, Current.MaxMind.AutoUpdate)
	assert.Equal(t, testMaxMindAutoUpdateInterval, Current.MaxMind.AutoUpdateInterval)
	assert.Equal(t, testIsDev, Current.Server.IsDev)

	// Check asset dir
	_, err = os.Stat(Current.Server.AssetDirPath)
	require.NoError(t, err)
	require.NotErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultConfig(t *testing.T) {
	ctx := context.Background()

	// Unset all env vars except MAXMIND_LICENSE_KEY
	// Using t.Setenv to empty string to simulate unset for viper (if it treats empty as unset)
	// OR we can just check errors on Unsetenv.
	// Since we are fixing lint errors, let's just use t.Setenv to set clean values.
	t.Setenv("MAXMIND_AUTOUPDATE", "")
	t.Setenv("MAXMIND_AUTOUPDATE_INTERVAL", "")
	t.Setenv("SERVER_LISTEN_ADDR", "")
	t.Setenv("SERVER_LISTEN_PORT", "")
	t.Setenv("API_KEYS", "")
	t.Setenv("IS_DEV", "")

	t.Setenv("MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)

	_, err := Load(ctx, "")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(constants.AssetDir)
	}()

	assert.Equal(t, DefaultServerListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, DefaultServerListenPort, Current.Server.ListenPort)
	assert.NotEmpty(t, Current.Server.APIKeys)
	assert.Len(t, Current.Server.APIKeys, 1)
	assert.True(t, Current.MaxMind.AutoUpdate)
	assert.Equal(t, DefaultMaxMindAutoUpdateInterval, Current.MaxMind.AutoUpdateInterval)
	assert.False(t, Current.Server.IsDev)
}

func TestConfigValidationFail(t *testing.T) {
	ctx := context.Background()

	// Unset all env vars
	t.Setenv("MAXMIND_LICENSE_KEY", "")
	t.Setenv("MAXMIND_AUTOUPDATE", "")
	t.Setenv("MAXMIND_AUTOUPDATE_INTERVAL", "")
	t.Setenv("SERVER_LISTEN_ADDR", "")
	t.Setenv("SERVER_LISTEN_PORT", "")
	t.Setenv("API_KEYS", "")
	t.Setenv("IS_DEV", "")

	// Load without MAXMIND_LICENSE_KEY should fail
	_, err := Load(ctx, "")
	require.Error(t, err)
	assert.Equal(t, ErrMaxMindLicenseKeyEmpty, err)
}

func TestServerConfigValidation(t *testing.T) {
	testCases := []struct {
		name      string
		config    ServerConfig
		expectErr error
	}{
		{
			name: "valid config",
			config: ServerConfig{
				ListenAddr:   "0.0.0.0",
				ListenPort:   5000,
				APIKeys:      []string{"test-key"},
				AssetDirPath: "./data",
			},
			expectErr: nil,
		},
		{
			name: "invalid port - too low",
			config: ServerConfig{
				ListenAddr:   "0.0.0.0",
				ListenPort:   0,
				APIKeys:      []string{"test-key"},
				AssetDirPath: "./data",
			},
			expectErr: ErrAPIListenPortInvalid,
		},
		{
			name: "invalid port - too high",
			config: ServerConfig{
				ListenAddr:   "0.0.0.0",
				ListenPort:   70000,
				APIKeys:      []string{"test-key"},
				AssetDirPath: "./data",
			},
			expectErr: ErrAPIListenPortInvalid,
		},
		{
			name: "empty api keys",
			config: ServerConfig{
				ListenAddr:   "0.0.0.0",
				ListenPort:   5000,
				APIKeys:      []string{},
				AssetDirPath: "./data",
			},
			expectErr: ErrAPIKeysEmpty,
		},
		{
			name: "empty asset dir",
			config: ServerConfig{
				ListenAddr:   "0.0.0.0",
				ListenPort:   5000,
				APIKeys:      []string{"test-key"},
				AssetDirPath: "",
			},
			expectErr: ErrAssetDirEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestMaxMindConfigValidation(t *testing.T) {
	testCases := []struct {
		name      string
		config    MaxMindConfig
		expectErr error
	}{
		{
			name: "valid config",
			config: MaxMindConfig{
				LicenseKey:         "test-key",
				AutoUpdate:         true,
				AutoUpdateInterval: 24 * time.Hour,
			},
			expectErr: nil,
		},
		{
			name: "empty license key",
			config: MaxMindConfig{
				LicenseKey:         "",
				AutoUpdate:         true,
				AutoUpdateInterval: 24 * time.Hour,
			},
			expectErr: ErrMaxMindLicenseKeyEmpty,
		},
		{
			name: "invalid auto-update interval",
			config: MaxMindConfig{
				LicenseKey:         "test-key",
				AutoUpdate:         true,
				AutoUpdateInterval: 0,
			},
			expectErr: ErrMaxMindAutoUpdateIntervalInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestServerConfigGetAddr(t *testing.T) {
	testCases := []struct {
		name     string
		config   ServerConfig
		expected string
	}{
		{
			name: "IPv4 address",
			config: ServerConfig{
				ListenAddr: "0.0.0.0",
				ListenPort: 5000,
			},
			expected: "0.0.0.0:5000",
		},
		{
			name: "IPv6 address",
			config: ServerConfig{
				ListenAddr: "::1",
				ListenPort: 8080,
			},
			expected: "[::1]:8080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr := tc.config.GetAddr()
			assert.Equal(t, tc.expected, addr)
		})
	}
}

func TestMultipleAPIKeys(t *testing.T) {
	ctx := context.Background()

	// Unset all env vars
	t.Setenv("MAXMIND_AUTOUPDATE", "")
	t.Setenv("MAXMIND_AUTOUPDATE_INTERVAL", "")
	t.Setenv("SERVER_LISTEN_ADDR", "")
	t.Setenv("SERVER_LISTEN_PORT", "")
	t.Setenv("IS_DEV", "")

	t.Setenv("MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)
	t.Setenv("API_KEYS", "key1,key2,key3")

	_, err := Load(ctx, "")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(constants.AssetDir)
	}()

	assert.Equal(t, []string{"key1", "key2", "key3"}, Current.Server.APIKeys)
}
