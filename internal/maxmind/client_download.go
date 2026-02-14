package maxmind

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/crypto/hash"
	"github.com/hibare/GoCommon/v2/pkg/file"
	"github.com/hibare/GoGeoIP/internal/constants"
)

const minSHA256FileParts = 2

// RunDBDownloadJob starts the background job for downloading databases.
func (c *Client) RunDBDownloadJob() {
	ticker := time.NewTicker(c.config.AutoUpdateInterval)
	slog.Info("Scheduling DB update job")

	for {
		<-ticker.C
		if err := c.DownloadAllDB(); err != nil {
			slog.Error("Background DB update failed", "error", err)
		}
	}
}

// DownloadAllDB downloads all configured databases.
func (c *Client) DownloadAllDB() error {
	slog.Info("Downloading all DB files")

	ctx := context.Background()
	types := []DBType{DBTypeCountry, DBTypeCity, DBTypeASN}

	var hasError bool
	for _, t := range types {
		if err := c.downloadDB(ctx, t); err != nil {
			slog.Error("Error downloading DB", "type", t, "error", err)
			hasError = true
		}
	}

	if hasError {
		if !c.checkAllDBFilesExist() {
			return ErrDBDownloadFailed
		}
		slog.Warn("Continuing with existing DB files despite download errors")
	}

	if err := c.Load(); err != nil {
		return fmt.Errorf("failed to reload databases: %w", err)
	}

	return nil
}

func (c *Client) downloadDB(ctx context.Context, dbType DBType) error {
	if len(c.config.LicenseKey) == 0 {
		return ErrLicenseKeyRequired
	}

	if err := os.MkdirAll(c.dataDir, os.ModePerm); err != nil {
		return err
	}

	tmpDir := os.TempDir()
	archivePath := filepath.Join(tmpDir, fmt.Sprintf("%s.%s", dbType, constants.DBArchiveDownloadSuffix))
	sha256Path := filepath.Join(tmpDir, fmt.Sprintf("%s.%s", dbType, constants.DBSHA256FileDownloadSuffix))

	dbURL := fmt.Sprintf(constants.MaxMindDownloadURL, dbType, c.config.LicenseKey, constants.DBArchiveDownloadSuffix)
	sha256URL := fmt.Sprintf(constants.MaxMindDownloadURL, dbType, c.config.LicenseKey, constants.DBSHA256FileDownloadSuffix)

	finalDBPath := c.getDBPath(dbType)

	// Clean up temp files on exit
	defer func() {
		_ = os.Remove(archivePath)
		_ = os.Remove(sha256Path)
	}()

	slog.Info("Downloading DB file", "path", archivePath)
	if err := file.DownloadFile(ctx, dbURL, archivePath); err != nil {
		return err
	}

	slog.Info("Downloading sha256 file", "path", sha256Path)
	if err := file.DownloadFile(ctx, sha256URL, sha256Path); err != nil {
		return err
	}

	sha256Sum, extractName, err := c.parseSHA256File(sha256Path)
	if err != nil {
		return err
	}

	hasher := hash.NewSHA256Hasher()
	valid, verifyErr := hasher.VerifyFile(archivePath, sha256Sum)
	if verifyErr != nil {
		return verifyErr
	} else if !valid {
		return fmt.Errorf("%w for archive %s", ErrChecksumMismatch, archivePath)
	}

	slog.Info("Checksum validated", "path", archivePath)

	slog.Info("Extracting file", "file", extractName, "archive", archivePath)
	extractedPath, err := file.ExtractFileFromTarGz(archivePath, extractName)
	if err != nil {
		return err
	}

	slog.Info("Loading new DB file", "path", finalDBPath)
	_ = os.Remove(finalDBPath) // Remove old DB file

	if err := os.Rename(extractedPath, finalDBPath); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseSHA256File(path string) (string, string, error) {
	lines, err := file.ReadFileLines(path)
	if err != nil {
		return "", "", err
	}
	if len(lines) == 0 {
		return "", "", ErrEmptySHA256File
	}

	parts := strings.Split(lines[0], "  ")
	if len(parts) < minSHA256FileParts {
		return "", "", ErrInvalidSHA256File
	}

	sha256Sum := parts[0]
	fileName := parts[1]

	baseName := strings.Split(strings.Split(fileName, ".")[0], "_")[0]
	dbFilename := fmt.Sprintf("%s.%s", baseName, constants.DBSuffix)

	return sha256Sum, dbFilename, nil
}

func (c *Client) checkAllDBFilesExist() bool {
	types := []DBType{DBTypeCountry, DBTypeCity, DBTypeASN}
	for _, t := range types {
		if _, err := os.Stat(c.getDBPath(t)); os.IsNotExist(err) {
			return false
		}
	}
	return true
}
