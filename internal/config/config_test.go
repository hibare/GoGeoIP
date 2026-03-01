package config

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hibare/Waypoint/internal/constants"
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
	testSecretKey                 = "test-secret-key"
	testDataDir                   = "./data"
	testDBDSN                     = "postgres://user:testpwd@localhost:5432/waypoint?sslmode=disable" //nolint:gosec // test DSN with dummy credentials
)

func TestEnvLoadedConfig(t *testing.T) {
	ctx := context.Background()

	// Set test env vars
	t.Setenv("WAYPOINT_MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE", strconv.FormatBool(testMaxMindAutoUpdate))
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE_INTERVAL", testMaxMindAutoUpdateInterval.String())
	t.Setenv("WAYPOINT_SERVER_LISTEN_ADDR", testAPIListenAddr)
	t.Setenv("WAYPOINT_SERVER_LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	t.Setenv("WAYPOINT_SERVER_API_KEYS", testAPIKeys)
	t.Setenv("WAYPOINT_CORE_SECRET_KEY", testSecretKey)
	t.Setenv("WAYPOINT_CORE_DATA_DIR", testDataDir)
	t.Setenv("WAYPOINT_DB_DSN", testDBDSN)

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
	assert.Equal(t, testSecretKey, Current.Core.SecretKey)
	assert.True(t, strings.HasSuffix(Current.Core.DataDir, "data"))

	// Check data dir
	_, err = os.Stat(Current.Core.DataDir)
	require.NoError(t, err)
	require.NotErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultConfig(t *testing.T) {
	ctx := context.Background()

	// Set minimal env vars
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE", "")
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE_INTERVAL", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_ADDR", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_PORT", "")
	t.Setenv("WAYPOINT_SERVER_API_KEYS", "")
	t.Setenv("WAYPOINT_CORE_SECRET_KEY", "test-secret-key")
	t.Setenv("WAYPOINT_CORE_DATA_DIR", "./data")
	t.Setenv("WAYPOINT_MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)
	t.Setenv("WAYPOINT_DB_DSN", testDBDSN)

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
	assert.NotEmpty(t, Current.Core.SecretKey)
}

func TestConfigValidationFail(t *testing.T) {
	ctx := context.Background()

	// Unset all env vars
	t.Setenv("WAYPOINT_MAXMIND_LICENSE_KEY", "")
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE", "")
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE_INTERVAL", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_ADDR", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_PORT", "")
	t.Setenv("WAYPOINT_SERVER_API_KEYS", "")
	t.Setenv("WAYPOINT_CORE_SECRET_KEY", "")
	t.Setenv("WAYPOINT_CORE_DATA_DIR", "")
	t.Setenv("WAYPOINT_DB_DSN", testDBDSN)

	// Load without CORE_SECRET_KEY should fail
	_, err := Load(ctx, "")
	require.Error(t, err)
	assert.Equal(t, ErrSecretKeyEmpty, err)
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
				ListenAddr: "0.0.0.0",
				ListenPort: 5000,
				APIKeys:    []string{"test-key"},
			},
			expectErr: nil,
		},
		{
			name: "invalid port - too low",
			config: ServerConfig{
				ListenAddr: "0.0.0.0",
				ListenPort: 0,
				APIKeys:    []string{"test-key"},
			},
			expectErr: ErrAPIListenPortInvalid,
		},
		{
			name: "invalid port - too high",
			config: ServerConfig{
				ListenAddr: "0.0.0.0",
				ListenPort: 70000,
				APIKeys:    []string{"test-key"},
			},
			expectErr: ErrAPIListenPortInvalid,
		},
		{
			name: "empty api keys",
			config: ServerConfig{
				ListenAddr: "0.0.0.0",
				ListenPort: 5000,
				APIKeys:    []string{},
			},
			expectErr: ErrAPIKeysEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestCoreConfigValidation(t *testing.T) {
	testCases := []struct {
		name      string
		config    CoreConfig
		expectErr error
	}{
		{
			name: "valid config",
			config: CoreConfig{
				SecretKey: "test-secret",
				DataDir:   "./data",
			},
			expectErr: nil,
		},
		{
			name: "empty secret key",
			config: CoreConfig{
				SecretKey: "",
				DataDir:   "./data",
			},
			expectErr: ErrSecretKeyEmpty,
		},
		{
			name: "empty data dir",
			config: CoreConfig{
				SecretKey: "test-secret",
				DataDir:   "",
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

	// Set minimal env vars
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE", "")
	t.Setenv("WAYPOINT_MAXMIND_AUTO_UPDATE_INTERVAL", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_ADDR", "")
	t.Setenv("WAYPOINT_SERVER_LISTEN_PORT", "")
	t.Setenv("WAYPOINT_CORE_SECRET_KEY", "test-secret-key")
	t.Setenv("WAYPOINT_CORE_DATA_DIR", "./data")
	t.Setenv("WAYPOINT_MAXMIND_LICENSE_KEY", testMaxMindLicenseKey)
	t.Setenv("WAYPOINT_DB_DSN", testDBDSN)
	t.Setenv("WAYPOINT_SERVER_API_KEYS", "key1,key2,key3")

	_, err := Load(ctx, "")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(constants.AssetDir)
	}()

	assert.Equal(t, []string{"key1", "key2", "key3"}, Current.Server.APIKeys)
}
