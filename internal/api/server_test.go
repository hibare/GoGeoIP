package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hibare/GoCommon/pkg/errors"
	commonMiddleware "github.com/hibare/GoCommon/pkg/http/middleware"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

var (
	app App
)

func TestMain(m *testing.M) {
	os.Setenv("DB_LICENSE_KEY", "test-license")

	config.Load()
	app.Initialize()
	code := m.Run()
	os.Exit(code)
}

func TestGeoIP400(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		expectedBody errors.Error
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/api/v1/ip/8.8.8",
			expectedBody: errors.Error{
				Code:    http.StatusBadRequest,
				Message: constants.ErrInvalidIP.Error(),
			},
		}, {
			Name: "URL with trailing slash",
			URL:  "/api/v1/ip/8.8.8/",
			expectedBody: errors.Error{
				Code:    http.StatusBadRequest,
				Message: constants.ErrInvalidIP.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)

			r.Header.Add(commonMiddleware.AuthHeaderName, config.Current.API.APIKeys[0])

			w := httptest.NewRecorder()

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			responseBody := errors.Error{}
			err = json.NewDecoder(w.Body).Decode(&responseBody)
			assert.NoError(t, err)

			if assert.NotNil(t, responseBody) {
				assert.NotEmpty(t, responseBody)
				assert.Equal(t, responseBody, tc.expectedBody)
			}
		})
	}
}

func TestGeoIP401(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		expectedBody errors.Error
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/api/v1/ip/8.8.8.8",
			expectedBody: errors.Error{
				Code:    http.StatusUnauthorized,
				Message: errors.ErrUnauthorized.Error(),
			},
		}, {
			Name: "URL with trailing slash",
			URL:  "/api/v1/ip/8.8.8.8/",
			expectedBody: errors.Error{
				Code:    http.StatusUnauthorized,
				Message: errors.ErrUnauthorized.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			responseBody := errors.Error{}
			err = json.NewDecoder(w.Body).Decode(&responseBody)
			assert.NoError(t, err)

			if assert.NotNil(t, responseBody) {
				assert.NotEmpty(t, responseBody)
				assert.Equal(t, responseBody, tc.expectedBody)
			}
		})
	}
}

func TestGeoIP500(t *testing.T) {
	os.RemoveAll(constants.AssetDir)

	testCases := []struct {
		Name string
		URL  string
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/api/v1/ip/8.8.8.8",
		}, {
			Name: "URL with trailing slash",
			URL:  "/api/v1/ip/8.8.8.8/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)

			r.Header.Add(commonMiddleware.AuthHeaderName, config.Current.API.APIKeys[0])

			w := httptest.NewRecorder()

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

func TestGeoIP200(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	testCases := []struct {
		Name string
		URL  string
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/api/v1/ip/8.8.8.8",
		}, {
			Name: "URL with trailing slash",
			URL:  "/api/v1/ip/8.8.8.8/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)

			r.Header.Add(commonMiddleware.AuthHeaderName, config.Current.API.APIKeys[0])

			w := httptest.NewRecorder()

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestMyIP200(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	testCases := []struct {
		Name string
		URL  string
	}{
		{
			Name: "URL without trailing slash",
			URL:  "/api/v1/ip",
		}, {
			Name: "URL with trailing slash",
			URL:  "/api/v1/ip/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)

			r.Header.Add(commonMiddleware.AuthHeaderName, config.Current.API.APIKeys[0])

			w := httptest.NewRecorder()

			app.Router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
