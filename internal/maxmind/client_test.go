package maxmind_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/maxmind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_DownloadAndLoad(t *testing.T) {
	// Skip if no license key provided (integration test)
	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	if licenseKey == "" {
		t.Skip("MAXMIND_LICENSE_KEY not set, skipping integration test")
	}

	// Create temp dir for data
	tempDir := t.TempDir()

	cfg := &config.MaxMindConfig{
		LicenseKey:         licenseKey,
		AutoUpdate:         true,
		AutoUpdateInterval: 24 * time.Hour,
	}

	client := maxmind.NewClient(cfg, tempDir)

	// Test DownloadAllDB
	err := client.DownloadAllDB()
	require.NoError(t, err)

	// Verify files exist
	assert.FileExists(t, filepath.Join(tempDir, "GeoLite2-Country.mmdb"))
	assert.FileExists(t, filepath.Join(tempDir, "GeoLite2-City.mmdb"))
	assert.FileExists(t, filepath.Join(tempDir, "GeoLite2-ASN.mmdb"))

	// Test Lookups
	// 8.8.8.8 should be Google in US
	ip := "8.8.8.8"

	// Test Country
	country, err := client.IP2Country(ip)
	require.NoError(t, err)
	assert.Equal(t, "United States", country.Country)
	assert.Equal(t, "US", country.ISOCountryCode)

	// Test City
	city, err := client.IP2City(ip)
	require.NoError(t, err)
	assert.Equal(t, "United States", city.Country)
	// City might be blank or approximate, but Country should be consistent

	// Test ASN
	asn, err := client.IP2ASN(ip)
	require.NoError(t, err)
	assert.Equal(t, uint(15169), asn.ASN)
	assert.Equal(t, "GOOGLE", asn.Organization)

	// Test IP2Geo (Combined)
	geo, err := client.IP2Geo(ip)
	require.NoError(t, err)
	assert.Equal(t, "United States", geo.Country)
	assert.Equal(t, uint(15169), geo.ASN)
}

func TestClient_InvalidIP(t *testing.T) {
	client := maxmind.NewClient(&config.MaxMindConfig{}, "")

	_, err := client.IP2Country("invalid")
	require.Error(t, err)

	_, err = client.IP2City("invalid")
	require.Error(t, err)

	_, err = client.IP2ASN("invalid")
	require.Error(t, err)

	_, err = client.IP2Geo("invalid")
	require.Error(t, err)
}
