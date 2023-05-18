package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hibare/GoGeoIP/internal/constants"
)

func CalculateFileSHA256(p string) (string, error) {
	f, err := os.Open(p)

	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func ValidateFileSha256(p string, sha256Str string) error {
	calculatedSha256, err := CalculateFileSHA256(p)

	if err != nil {
		return err
	}

	if calculatedSha256 != sha256Str {
		return constants.ErrChecksumMismatch
	}
	return nil
}

func DownloadFile(url string, destination string) error {
	response, err := http.Get(url)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", response.Status)
	}

	out, err := os.Create(destination)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, response.Body)

	defer out.Close()

	if err != nil {
		return err
	}

	return nil
}

func ReadFile(p string) ([]string, error) {
	f, err := os.Open(p)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	var lines []string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func ExtractFileFromTarGz(archivePath, targetFilename string) (string, error) {
	var targetFilePath string

	file, err := os.Open(archivePath)
	if err != nil {
		return targetFilePath, err
	}
	defer file.Close()

	// Create a gzip reader for the file
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return targetFilePath, err
	}
	defer gzipReader.Close()

	// Create a tar reader for the gzip reader
	tarReader := tar.NewReader(gzipReader)

	// Find the target file in the archive

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return targetFilePath, os.ErrNotExist
		}
		if err != nil {
			return targetFilePath, err
		}

		if strings.HasSuffix(header.Name, targetFilename) {

			targetFilePath = filepath.Join(os.TempDir(), targetFilename)

			// Create the target file and copy the content of the file from the archive
			targetFile, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return targetFilePath, err
			}

			if err != nil {
				return targetFilePath, err
			}

			if _, err := io.Copy(targetFile, tarReader); err != nil {
				targetFile.Close()
				os.Remove(targetFilename)
				return targetFilePath, err
			}
			targetFile.Close()
			break
		}
	}
	return targetFilePath, nil
}
