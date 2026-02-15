package testhelpers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hibare/GoGeoIP/internal/constants"
)

const (
	TestDataDir = "../testhelper/test_data"
)

var TestDBFiles = []string{
	filepath.Join(TestDataDir, "GeoLite2-City.tar.gz"),
	filepath.Join(TestDataDir, "GeoLite2-Country.tar.gz"),
	filepath.Join(TestDataDir, "GeoLite2-ASN.tar.gz"),
}

var TestDBSHA256Files = []string{
	filepath.Join(TestDataDir, "GeoLite2-City.tar.gz.sha256"),
	filepath.Join(TestDataDir, "GeoLite2-Country.tar.gz.sha256"),
	filepath.Join(TestDataDir, "GeoLite2-ASN.tar.gz.sha256"),
}

func LoadTestDB() error {
	err := os.MkdirAll(constants.AssetDir, os.ModePerm)
	if err != nil {
		return err
	}
	dbFiles := []string{
		filepath.Join(TestDataDir, "GeoLite2-City.mmdb"),
		filepath.Join(TestDataDir, "GeoLite2-Country.mmdb"),
		filepath.Join(TestDataDir, "GeoLite2-ASN.mmdb"),
	}

	// loop through all dbfiles and copy them to AssetDir
	for _, dbFile := range dbFiles {
		sourceFile := filepath.Join(TestDataDir, dbFile)
		src, err := os.Open(sourceFile)
		if err != nil {
			return err
		}

		// Create the destination file
		destFile := filepath.Join(constants.AssetDir, filepath.Base(sourceFile))
		dest, err := os.Create(destFile)
		if err != nil {
			_ = src.Close()
			return err
		}

		// Copy the contents from the source file to the destination file
		_, err = io.Copy(dest, src)
		if err != nil {
			_ = dest.Close()
			_ = src.Close()
			return err
		}
		_ = dest.Close()
		_ = src.Close()
	}

	return nil
}
