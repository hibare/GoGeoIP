package maxmind

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("DB_LICENSE_KEY", "test-license")
	config.Load()
	m.Run()
}

func TestGetDBDownloadURL(t *testing.T) {
	testCases := []struct {
		Name        string
		DBType      string
		ExpectedURL string
	}{
		{
			Name:        constants.DBTypeCity,
			DBType:      constants.DBTypeCity,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeCity, config.Current.DB.LicenseKey, constants.DBArchiveDownloadSuffix),
		},
		{
			Name:        constants.DBTypeCountry,
			DBType:      constants.DBTypeCountry,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeCountry, config.Current.DB.LicenseKey, constants.DBArchiveDownloadSuffix),
		},
		{
			Name:        constants.DBTypeASN,
			DBType:      constants.DBTypeASN,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeASN, config.Current.DB.LicenseKey, constants.DBArchiveDownloadSuffix),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := getDBUrl(tc.DBType)
			assert.Equal(t, tc.ExpectedURL, url)
		})
	}
}

func TestDBSha256DownloadURL(t *testing.T) {
	testCases := []struct {
		Name        string
		DBType      string
		ExpectedURL string
	}{
		{
			Name:        constants.DBTypeCity,
			DBType:      constants.DBTypeCity,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeCity, config.Current.DB.LicenseKey, constants.DBSHA256FileDownloadSuffix),
		},
		{
			Name:        constants.DBTypeCountry,
			DBType:      constants.DBTypeCountry,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeCountry, config.Current.DB.LicenseKey, constants.DBSHA256FileDownloadSuffix),
		},
		{
			Name:        constants.DBTypeASN,
			DBType:      constants.DBTypeASN,
			ExpectedURL: fmt.Sprintf(constants.MaxMindDownloadURL, constants.DBTypeASN, config.Current.DB.LicenseKey, constants.DBSHA256FileDownloadSuffix),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := getDBSha256Url(tc.DBType)
			assert.Equal(t, tc.ExpectedURL, url)
		})
	}
}

func TestDBArchiveDownloadPath(t *testing.T) {
	testCases := []struct {
		Name        string
		DBType      string
		ExpectedURL string
	}{
		{
			Name:        constants.DBTypeCity,
			DBType:      constants.DBTypeCity,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeCity, constants.DBArchiveDownloadSuffix)),
		},
		{
			Name:        constants.DBTypeCountry,
			DBType:      constants.DBTypeCountry,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeCountry, constants.DBArchiveDownloadSuffix)),
		},
		{
			Name:        constants.DBTypeASN,
			DBType:      constants.DBTypeASN,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeASN, constants.DBArchiveDownloadSuffix)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := getDBArchiveDownloadPath(tc.DBType)
			assert.Equal(t, tc.ExpectedURL, url)
		})
	}
}

func TestDBSha256DownloadPath(t *testing.T) {
	testCases := []struct {
		Name        string
		DBType      string
		ExpectedURL string
	}{
		{
			Name:        constants.DBTypeCity,
			DBType:      constants.DBTypeCity,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeCity, constants.DBSHA256FileDownloadSuffix)),
		},
		{
			Name:        constants.DBTypeCountry,
			DBType:      constants.DBTypeCountry,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeCountry, constants.DBSHA256FileDownloadSuffix)),
		},
		{
			Name:        constants.DBTypeASN,
			DBType:      constants.DBTypeASN,
			ExpectedURL: filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", constants.DBTypeASN, constants.DBSHA256FileDownloadSuffix)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := getDBSHA256DownloadPath(tc.DBType)
			assert.Equal(t, tc.ExpectedURL, url)
		})
	}
}

func TestParseSHA256File(t *testing.T) {
	testCases := []struct {
		Name        string
		Filepath    string
		ExpectError bool
	}{
		{
			Name:        "Fail",
			Filepath:    "",
			ExpectError: true,
		},
		{
			Name:        "Success",
			Filepath:    filepath.Join(testhelper.TestDataDir, "GeoLite2-ASN.tar.gz.sha256"),
			ExpectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sha256, dbFilename, err := parseSHA256File(tc.Filepath)

			if tc.ExpectError {
				assert.Error(t, err)
				assert.Empty(t, sha256)
				assert.Empty(t, dbFilename)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, sha256)
				assert.NotEmpty(t, dbFilename)
			}
		})
	}
}

func TestDownloadDB(t *testing.T) {
	testCases := []struct {
		Name   string
		DBType string
	}{
		{
			Name:   "City",
			DBType: constants.DBTypeCity,
		},
		{
			Name:   "Country",
			DBType: constants.DBTypeCountry,
		},
		{
			Name:   "ASN",
			DBType: constants.DBTypeASN,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/%s.%s", constants.DBTypeCity, constants.DBArchiveDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBFiles[constants.DBTypeCity])
		case fmt.Sprintf("/%s.%s", constants.DBTypeCountry, constants.DBArchiveDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBFiles[constants.DBTypeCountry])
		case fmt.Sprintf("/%s.%s", constants.DBTypeASN, constants.DBArchiveDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBFiles[constants.DBTypeASN])
		case fmt.Sprintf("/%s.%s", constants.DBTypeCity, constants.DBSHA256FileDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBSHA256Files[constants.DBTypeCity])
		case fmt.Sprintf("/%s.%s", constants.DBTypeCountry, constants.DBSHA256FileDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBSHA256Files[constants.DBTypeCountry])
		case fmt.Sprintf("/%s.%s", constants.DBTypeASN, constants.DBSHA256FileDownloadSuffix):
			http.ServeFile(w, r, testhelper.TestDBSHA256Files[constants.DBTypeASN])
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 Not Found")
		}
	}))

	defer server.Close()

	getDBUrl = func(dbType string) string {
		return fmt.Sprintf("%s/%s.%s", server.URL, dbType, constants.DBArchiveDownloadSuffix)
	}

	getDBSha256Url = func(dbType string) string {
		return fmt.Sprintf("%s/%s.%s", server.URL, dbType, constants.DBSHA256FileDownloadSuffix)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := downloadDB(t.Context(), tc.DBType)
			assert.NoError(t, err)
		})
	}

	t.Run("Download City", func(t *testing.T) {
		err := DownloadCityDB(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Download Country", func(t *testing.T) {
		err := DownloadCountryDB(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Download ASN", func(t *testing.T) {
		err := DownloadASNDB(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Download All", func(t *testing.T) {
		DownloadAllDB()
	})

	os.RemoveAll(constants.AssetDir)
}
