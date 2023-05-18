package testhelper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hibare/GoGeoIP/internal/constants"
)

const (
	TestDataDir = "../testhelper/test_data"
)

var TestDBFiles = map[string]string{
	constants.DBTypeCity:    filepath.Join(TestDataDir, "GeoLite2-City.tar.gz"),
	constants.DBTypeCountry: filepath.Join(TestDataDir, "GeoLite2-Country.tar.gz"),
	constants.DBTypeASN:     filepath.Join(TestDataDir, "GeoLite2-ASN.tar.gz"),
}

var TestDBSHA256Files = map[string]string{
	constants.DBTypeCity:    filepath.Join(TestDataDir, "GeoLite2-City.tar.gz.sha256"),
	constants.DBTypeCountry: filepath.Join(TestDataDir, "GeoLite2-Country.tar.gz.sha256"),
	constants.DBTypeASN:     filepath.Join(TestDataDir, "GeoLite2-ASN.tar.gz.sha256"),
}

func LoadTestDB() error {
	err := os.MkdirAll(constants.AssetDir, os.ModePerm)
	if err != nil {
		return err
	}
	dbFiles := []string{
		fmt.Sprintf("%s.%s", constants.DBTypeCity, constants.DBSuffix),
		fmt.Sprintf("%s.%s", constants.DBTypeCountry, constants.DBSuffix),
		fmt.Sprintf("%s.%s", constants.DBTypeASN, constants.DBSuffix),
	}

	// loop through all dbfiles and copy them to AssetDir
	for _, dbFile := range dbFiles {
		sourceFile := filepath.Join(TestDataDir, dbFile)
		src, err := os.Open(sourceFile)
		if err != nil {
			return err
		}
		defer src.Close()

		// Create the destination file
		destFile := filepath.Join(constants.AssetDir, filepath.Base(sourceFile))
		dest, err := os.Create(destFile)
		if err != nil {
			return err
		}
		defer dest.Close()

		// Copy the contents from the source file to the destination file
		_, err = io.Copy(dest, src)
		if err != nil {
			return err
		}
	}
	return nil
}
