package testhelpers

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gorm.io/gorm"
)

const (
	seedsDir     = "seeds"
	sqlExtension = ".sql"
)

// PopulateSeeds populates the test database with seed data from SQL files.
// This data is used by tests and provides a consistent baseline for testing.
func PopulateSeeds(tx *gorm.DB) error {
	// Find the seeds directory relative to this file's location
	// This works regardless of the current working directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return os.ErrNotExist
	}

	// Get the directory containing this file
	dir := filepath.Dir(filename)
	seedsDir := filepath.Join(dir, seedsDir)

	// Read all seed files from the filesystem
	entries, err := os.ReadDir(seedsDir)
	if err != nil {
		return err
	}

	// Sort files by name to ensure consistent execution order
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), sqlExtension) {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// Execute each seed file
	for _, filename := range files {
		sqlFilePath := filepath.Join(seedsDir, filename)

		sqlContent, err := os.ReadFile(sqlFilePath)
		if err != nil {
			return err
		}

		// Execute the SQL
		if err := tx.Exec(string(sqlContent)).Error; err != nil {
			return err
		}
	}

	return nil
}
