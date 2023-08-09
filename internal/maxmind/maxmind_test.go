package maxmind

import (
	"os"
	"testing"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestIP2Geo(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	t.Run("IP2Country", func(t *testing.T) {

		testCases := []struct {
			Name        string
			IP          string
			Expected    IPCountry
			ExpectedErr error
		}{
			{
				Name: "IPv4 Success",
				IP:   "81.2.69.160",
				Expected: IPCountry{
					IP:                  "81.2.69.160",
					Country:             "United Kingdom",
					Continent:           "Europe",
					ISOCountryCode:      "GB",
					ISOContinentCode:    "EU",
					IsAnonymousProxy:    false,
					IsSatelliteProvider: false,
				},
				ExpectedErr: nil,
			},
			{
				Name:        "IPv4 Error",
				IP:          "81.2.69.",
				Expected:    IPCountry{},
				ExpectedErr: constants.ErrInvalidIP,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				ipCountry, err := IP2Country(tc.IP)

				if tc.ExpectedErr != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tc.ExpectedErr)
					assert.Empty(t, ipCountry)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.Expected, ipCountry)
				}
			})
		}
	})

	t.Run("IP2City", func(t *testing.T) {
		testCases := []struct {
			Name        string
			IP          string
			Expected    IPCity
			ExpectedErr error
		}{
			{
				Name: "IPv4 Success",
				IP:   "81.2.69.160",
				Expected: IPCity{
					IP:   "81.2.69.160",
					City: "London",
					IPCountry: IPCountry{
						Country:             "United Kingdom",
						Continent:           "Europe",
						ISOCountryCode:      "GB",
						ISOContinentCode:    "EU",
						IsAnonymousProxy:    false,
						IsSatelliteProvider: false,
					},
					Timezone:  "Europe/London",
					Latitude:  51.5142,
					Longitude: -0.0931,
				},
				ExpectedErr: nil,
			},
			{
				Name:        "IPv4 Success",
				IP:          "81.2.69.",
				Expected:    IPCity{},
				ExpectedErr: constants.ErrInvalidIP,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				IpCity, err := IP2City(tc.IP)

				if tc.ExpectedErr != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tc.ExpectedErr)
					assert.Empty(t, IpCity)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.Expected, IpCity)
				}
			})
		}
	})

	t.Run("IP2ASN", func(t *testing.T) {
		testCases := []struct {
			Name        string
			IP          string
			Expected    IPASN
			ExpectedErr error
		}{
			{
				Name: "IPv4 Success",
				IP:   "149.101.100.0",
				Expected: IPASN{
					IP:           "149.101.100.0",
					ASN:          6167,
					Organization: "CELLCO-PART",
				},
				ExpectedErr: nil,
			},
			{
				Name:        "IPv4 Error",
				IP:          "149.101.100.",
				Expected:    IPASN{},
				ExpectedErr: constants.ErrInvalidIP,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				ipAsn, err := IP2ASN(tc.IP)

				if tc.ExpectedErr != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tc.ExpectedErr)
					assert.Empty(t, ipAsn)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.Expected, ipAsn)
				}
			})
		}
	})

	t.Run("IP2Geo", func(t *testing.T) {
		testCases := []struct {
			Name        string
			IP          string
			Expected    GeoIP
			ExpectedErr error
		}{
			{
				Name: "IPv4 Success",
				IP:   "149.101.100.0",
				Expected: GeoIP{
					IP: "149.101.100.0",
					IPCity: IPCity{
						IP:   "149.101.100.0",
						City: "",
						IPCountry: IPCountry{
							Country:             "United States",
							Continent:           "North America",
							ISOCountryCode:      "US",
							ISOContinentCode:    "NA",
							IsAnonymousProxy:    false,
							IsSatelliteProvider: false,
						},
						Timezone:  "America/Chicago",
						Latitude:  37.751,
						Longitude: -97.822,
					},
					IPASN: IPASN{
						IP:           "149.101.100.0",
						ASN:          6167,
						Organization: "CELLCO-PART",
					},
				},
				ExpectedErr: nil,
			},
			{
				Name:        "IPv4 Error",
				IP:          "149.101.100.",
				Expected:    GeoIP{IP: "149.101.100."},
				ExpectedErr: constants.ErrInvalidIP,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				ipGeo, err := IP2Geo(tc.IP)

				if tc.ExpectedErr != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tc.ExpectedErr)
					assert.Equal(t, tc.Expected, ipGeo)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.Expected, ipGeo)
				}
			})
		}
	})
}
