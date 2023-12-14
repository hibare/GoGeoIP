package maxmind

import (
	"sync"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

var (
	dbCountryReader *geoip2.Reader
	dbCityReader    *geoip2.Reader
	dbAsnReader     *geoip2.Reader
	dbLock					sync.Mutex
)

func LoadAllDB() error {
	var err error

	dbCountryReader, err = loadDB(constants.DBTypeCountry)
	if err != nil {
		return err
	}

	dbCityReader, err = loadDB(constants.DBTypeCity)
	if err != nil {
		return err
	}

	dbAsnReader, err = loadDB(constants.DBTypeASN)
	if err != nil {
		return err
	}

	return nil
}

func loadDB(dbType string) (reader *geoip2.Reader, err error) {
	oldDB := GetDB(dbType)
	if oldDB != nil {
		oldDB.Close()
	}

	dbReader, err := geoip2.Open(GetDBFilePath(dbType))
	if err != nil {
		log.Fatalf("Error opening DB file %s, %s", dbType, err)
		return nil, err
	}

	return dbReader, nil
}

func GetDB(dbType string) *geoip2.Reader {
	dbLock.Lock()
	defer dbLock.Unlock()

	switch dbType {
	case constants.DBTypeCountry:
		return dbCountryReader
	case constants.DBTypeCity:
		return dbCityReader
	case constants.DBTypeASN:
		return dbAsnReader
	default:
		return nil
	}
}
