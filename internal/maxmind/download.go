package maxmind

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hibare/GoCommon/v2/pkg/file"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
)

var getDBUrl = func(dbType string) string {
	return fmt.Sprintf(constants.MaxMindDownloadURL, dbType, config.Current.DB.LicenseKey, constants.DBArchiveDownloadSuffix)
}

var getDBSha256Url = func(dbType string) string {
	return fmt.Sprintf(constants.MaxMindDownloadURL, dbType, config.Current.DB.LicenseKey, constants.DBSHA256FileDownloadSuffix)
}

func getDBArchiveDownloadPath(dbType string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", dbType, constants.DBArchiveDownloadSuffix))
}

func getDBSHA256DownloadPath(dbType string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", dbType, constants.DBSHA256FileDownloadSuffix))
}

func GetDBFilePath(dbType string) string {
	return filepath.Join(constants.AssetDir, fmt.Sprintf("%s.%s", dbType, constants.DBSuffix))
}

func parseSHA256File(path string) (string, string, error) {
	content, err := file.ReadFileLines(path)
	if err != nil {
		return "", "", err
	}

	firstLine := content[0]
	sha256DirNamePair := strings.Split(firstLine, "  ")
	sha256 := sha256DirNamePair[0]
	dbFilename := fmt.Sprintf("%s.%s", strings.Split(strings.Split(sha256DirNamePair[1], ".")[0], "_")[0], constants.DBSuffix)
	return sha256, dbFilename, nil
}

func downloadDB(dbType string) error {
	downloadArchivePath := getDBArchiveDownloadPath(dbType)
	sha256FilePath := getDBSHA256DownloadPath(dbType)
	dbUrl := getDBUrl(dbType)
	dbSha256Url := getDBSha256Url(dbType)
	dbFilePath := GetDBFilePath(dbType)

	if len(config.Current.DB.LicenseKey) <= 0 {
		log.Fatal("DB_LICENSE_KEY is required")
	}

	if err := os.MkdirAll(constants.AssetDir, os.ModePerm); err != nil {
		return err
	}

	log.Infof("Downloading DB file, path=%s", downloadArchivePath)
	if err := file.DownloadFile(dbUrl, downloadArchivePath); err != nil {
		return err
	}
	log.Infof("Downloaded DB file, path=%s", downloadArchivePath)

	log.Infof("Downloading sha256 file, path=%s", sha256FilePath)
	if err := file.DownloadFile(dbSha256Url, sha256FilePath); err != nil {
		return err
	}
	log.Infof("Downloaded sha256 file, path=%s", sha256FilePath)

	sha256, extractDBFilename, err := parseSHA256File(sha256FilePath)
	if err != nil {
		return err
	}

	if err = file.ValidateFileSha256(downloadArchivePath, sha256); err != nil {
		return err
	}
	log.Infof("Checksum validated for archive %s", downloadArchivePath)

	log.Infof("Extracting file %s from archive %s", extractDBFilename, downloadArchivePath)
	extractDBFilepath, err := file.ExtractFileFromTarGz(downloadArchivePath, extractDBFilename)
	if err != nil {
		return err
	}
	log.Infof("Extracted file %s from archive %s at %s", extractDBFilename, downloadArchivePath, extractDBFilepath)

	log.Infof("Loading new DB file %s", dbFilePath)
	os.Remove(dbFilePath) // Remove old DB file

	if err := os.Rename(extractDBFilepath, dbFilePath); err != nil {
		return err
	}
	log.Infof("New DB file loaded %s", dbFilePath)

	os.Remove(downloadArchivePath)
	os.Remove(sha256FilePath)

	return nil
}

func DownloadCountryDB() error {
	return downloadDB(constants.DBTypeCountry)
}

func DownloadCityDB() error {
	return downloadDB(constants.DBTypeCity)
}

func DownloadASNDB() error {
	return downloadDB(constants.DBTypeASN)
}

func DownloadAllDB() {
	log.Infof("Downloading all DB files")

	errors := []error{
		DownloadCountryDB(),
		DownloadCityDB(),
		DownloadASNDB(),
	}
	for _, err := range errors {
		if err != nil {
			log.Fatalf("Error downloading DB: %v", err)
		}
	}

	LoadAllDB()
}

func RunDBDownloadJob() {
	ticker := time.NewTicker(config.Current.DB.AutoUpdateInterval)
	log.Infof("Scheduling DB update job")

	for {
		<-ticker.C
		DownloadAllDB()
	}
}
