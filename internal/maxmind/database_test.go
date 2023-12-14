package maxmind

import (
	"os"
	"testing"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestLoadAllDB(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	err = LoadAllDB()
	assert.NoError(t, err)

	assert.NotNil(t, dbCountryReader)
	assert.NotNil(t, dbCityReader)
	assert.NotNil(t, dbAsnReader)
}

func TestLoadDB(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	countryDB, err := loadDB(constants.DBTypeCountry)
	assert.NoError(t, err)
	assert.NotNil(t, countryDB)

	cityDB, err := loadDB(constants.DBTypeCity)
	assert.NoError(t, err)
	assert.NotNil(t, cityDB)

	asnDB, err := loadDB(constants.DBTypeASN)
	assert.NoError(t, err)
	assert.NotNil(t, asnDB)
}

func TestGetDB(t *testing.T) {
	err := testhelper.LoadTestDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(constants.AssetDir)
	})

	dbCountryReader = nil
	dbCityReader = nil
	dbAsnReader = nil

	countryDB := GetDB(constants.DBTypeCountry)
	assert.NotNil(t, countryDB)

	cityDB := GetDB(constants.DBTypeCity)
	assert.NotNil(t, cityDB)

	asnDB := GetDB(constants.DBTypeASN)
	assert.NotNil(t, asnDB)
}