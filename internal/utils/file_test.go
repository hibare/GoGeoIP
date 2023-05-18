package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hibare/GoGeoIP/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func CreateTestFile() ([]byte, string, error) {
	file, err := os.CreateTemp("", "test-file-*.txt")
	if err != nil {
		return []byte{}, "", err
	}
	defer file.Close()

	content := []byte("This is a test file.\nIt contains some sample content.")
	_, err = file.Write(content)
	if err != nil {
		return []byte{}, "", err
	}

	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return []byte{}, "", err
	}

	return content, absPath, err
}

func TestCalculateFileSHA256Pass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

	calculatedSHA256, err := CalculateFileSHA256(absPath)

	assert.NoError(t, err)
	assert.Equal(t, expectedSHA256, calculatedSHA256)
}

func TestCalculateFileSHA256Fail(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

	calculatedSHA256, err := CalculateFileSHA256(absPath)

	assert.NoError(t, err)
	assert.NotEqual(t, expectedSHA256, calculatedSHA256)
}

func TestValidateFileSHA256Pass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

	err = ValidateFileSha256(absPath, expectedSHA256)
	assert.NoError(t, err)
}

func TestValidateFileSHA256Fail(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

	err = ValidateFileSha256(absPath, expectedSHA256)
	assert.Error(t, err)
}

func TestReadFilePass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	lines, err := ReadFile(absPath)
	assert.NoError(t, err)
	assert.Len(t, lines, 2)
}

func TestReadFileFail(t *testing.T) {
	lines, err := ReadFile("some/random/path")
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestDownloadFilePass(t *testing.T) {
	// Create a test file
	_, absPath, err := CreateTestFile()
	assert.NoError(t, err)

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a test file
		http.ServeFile(w, r, absPath)
	}))
	defer server.Close()

	downloadFilePath := filepath.Join(os.TempDir(), "test-file.txt")
	defer os.Remove(downloadFilePath)

	// Download the file using the download function
	err = DownloadFile(server.URL, downloadFilePath)
	assert.NoError(t, err)

	lines, err := ReadFile(downloadFilePath)
	assert.NoError(t, err)
	assert.Len(t, lines, 2)

}

func TestDownloadFileFail(t *testing.T) {

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a test file
		http.ServeFile(w, r, "some/random/path")
	}))
	defer server.Close()

	downloadFilePath := filepath.Join(os.TempDir(), "test-file.txt")
	defer os.Remove(downloadFilePath)

	// Download the file using the download function
	err := DownloadFile(server.URL, downloadFilePath)
	assert.Error(t, err)

	lines, err := ReadFile(downloadFilePath)
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestExtractFileFromTarGzPass(t *testing.T) {
	archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
	targetFilename := "sample.txt"
	extractedFilePath := filepath.Join(os.TempDir(), targetFilename)

	extractedPath, err := ExtractFileFromTarGz(archivePath, targetFilename)
	assert.NoError(t, err)
	assert.Equal(t, extractedFilePath, extractedPath)
}

func TestExtractFileFromTarGzFail(t *testing.T) {
	archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
	targetFilename := "sample1.txt"

	_, err := ExtractFileFromTarGz(archivePath, targetFilename)
	assert.Error(t, err)
}
