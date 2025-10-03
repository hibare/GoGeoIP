package maxmind

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hibare/GoCommon/v2/pkg/crypto/hash"
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

func downloadDB(ctx context.Context, dbType string) error {
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
	if err := file.DownloadFile(ctx, dbUrl, downloadArchivePath); err != nil {
		return err
	}
	log.Infof("Downloaded DB file, path=%s", downloadArchivePath)

	log.Infof("Downloading sha256 file, path=%s", sha256FilePath)
	if err := file.DownloadFile(ctx, dbSha256Url, sha256FilePath); err != nil {
		return err
	}
	log.Infof("Downloaded sha256 file, path=%s", sha256FilePath)

	sha256, extractDBFilename, err := parseSHA256File(sha256FilePath)
	if err != nil {
		return err
	}

	hasher := hash.NewSHA256Hasher()
	if b, err := hasher.VerifyFile(downloadArchivePath, sha256); err != nil {
		return err
	} else if !b {
		return fmt.Errorf("checksum mismatch for archive %s", downloadArchivePath)
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

func dbFileExists(dbType string) bool {
	dbFilePath := GetDBFilePath(dbType)

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func DownloadCountryDB(ctx context.Context) error {
	return downloadDB(ctx, constants.DBTypeCountry)
}

func DownloadCityDB(ctx context.Context) error {
	return downloadDB(ctx, constants.DBTypeCity)
}

func DownloadASNDB(ctx context.Context) error {
	return downloadDB(ctx, constants.DBTypeASN)
}

func checkAllDBFilesExist() bool {
	return dbFileExists(constants.DBTypeCountry) && dbFileExists(constants.DBTypeCity) && dbFileExists(constants.DBTypeASN)
}

func DownloadAllDB() {
	log.Infof("Downloading all DB files")

	ctx := context.Background()

	errors := []error{
		DownloadCountryDB(ctx),
		DownloadCityDB(ctx),
		DownloadASNDB(ctx),
	}
	for _, err := range errors {
		if err != nil {
			exists := checkAllDBFilesExist()
			if !exists {
				log.Fatalf("Error downloading DB: %v", err)
			}
			log.Errorf("Error downloading DB: %v", err)
			log.Warn("Continuing with existing DB files")
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
