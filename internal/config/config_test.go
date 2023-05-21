package config

import (
	"os"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/stretchr/testify/assert"
)

const (
	testDBLicenseKey         = "test-license"
	testDBAutoUpdate         = false
	testDBAutoUpdateInterval = 2 * time.Hour
	testAPIListenAddr        = "127.0.0.1"
	testAPIListenPort        = 10000
	testAPIKeys              = "test-api-key"
	testIsDev                = true
)

func TestEnvLoadedConfig(t *testing.T) {
	// Set test env vars
	os.Setenv("DB_LICENSE_KEY", testDBLicenseKey)
	os.Setenv("DB_AUTOUPDATE", strconv.FormatBool(testDBAutoUpdate))
	os.Setenv("DB_AUTOUPDATE_INTERVAL", testDBAutoUpdateInterval.String())
	os.Setenv("API_LISTEN_ADDR", testAPIListenAddr)
	os.Setenv("API_LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	os.Setenv("API_KEYS", testAPIKeys)
	os.Setenv("IS_DEV", strconv.FormatBool(testIsDev))

	Load()
	defer os.RemoveAll(constants.AssetDir)

	assert.Equal(t, testAPIListenAddr, Current.API.ListenAddr)
	assert.Equal(t, testAPIListenPort, Current.API.ListenPort)
	assert.Equal(t, []string{testAPIKeys}, Current.API.APIKeys)
	assert.Equal(t, testDBLicenseKey, Current.DB.LicenseKey)
	assert.Equal(t, testDBAutoUpdate, Current.DB.AutoUpdate)
	assert.Equal(t, testDBAutoUpdateInterval, Current.DB.AutoUpdateInterval)
	assert.Equal(t, testIsDev, Current.Util.IsDev)

	// Check asset dir
	_, err := os.Stat(Current.Util.AssetDirPath)
	assert.NoError(t, err)
	assert.NotErrorIs(t, err, os.ErrNotExist)

	// Unset all env vars except
	os.Unsetenv("DB_LICENSE_KEY")
	os.Unsetenv("DB_AUTOUPDATE")
	os.Unsetenv("DB_AUTOUPDATE_INTERVAL")
	os.Unsetenv("API_LISTEN_ADDR")
	os.Unsetenv("API_LISTEN_PORT")
	os.Unsetenv("API_KEYS")
	os.Unsetenv("IS_DEV")

}

func TestDefaultConfig(t *testing.T) {
	// Unset all env vars except DB_LICENSE_KEY
	os.Unsetenv("DB_AUTOUPDATE")
	os.Unsetenv("DB_AUTOUPDATE_INTERVAL")
	os.Unsetenv("API_LISTEN_ADDR")
	os.Unsetenv("API_LISTEN_PORT")
	os.Unsetenv("API_KEYS")
	os.Unsetenv("IS_DEV")

	os.Setenv("DB_LICENSE_KEY", testDBLicenseKey)

	Load()
	defer os.RemoveAll(constants.AssetDir)

	assert.Equal(t, constants.DefaultAPIListenAddr, Current.API.ListenAddr)
	assert.Equal(t, constants.DefaultAPIListenPort, Current.API.ListenPort)
	assert.NotEmpty(t, Current.API.APIKeys)
	assert.Len(t, Current.API.APIKeys, 1)
	assert.True(t, true, Current.DB.AutoUpdate)
	assert.Equal(t, constants.DefaultDBAutoUpdateInterval, Current.DB.AutoUpdateInterval)
	assert.Equal(t, false, Current.Util.IsDev)

	os.Unsetenv("DB_LICENSE_KEY")
}

func TestDefaultConfigFail(t *testing.T) {
	// Unset all env vars except DB_LICENSE_KEY
	os.Unsetenv("DB_AUTOUPDATE")
	os.Unsetenv("DB_AUTOUPDATE_INTERVAL")
	os.Unsetenv("API_LISTEN_ADDR")
	os.Unsetenv("API_LISTEN_PORT")
	os.Unsetenv("API_KEYS")
	os.Unsetenv("IS_DEV")

	defer func() { log.StandardLogger().ExitFunc = nil }()
	var fatal bool
	log.StandardLogger().ExitFunc = func(int) { fatal = true }

	Load()
	defer os.RemoveAll(constants.AssetDir)

	assert.Equal(t, true, fatal)
}
