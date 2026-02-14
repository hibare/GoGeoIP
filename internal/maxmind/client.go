package maxmind

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/oschwald/geoip2-golang"
)

// DBType represents the type of MaxMind database.
type DBType string

const (
	DBTypeCountry DBType = constants.DBTypeCountry
	DBTypeCity    DBType = constants.DBTypeCity
	DBTypeASN     DBType = constants.DBTypeASN
)

// Client handles MaxMind database operations and lookups.
type Client struct {
	config  *config.MaxMindConfig
	dataDir string
	readers map[DBType]*geoip2.Reader
	mu      sync.RWMutex
}

// NewClient creates a new MaxMind client.
func NewClient(cfg *config.MaxMindConfig, dataDir string) *Client {
	return &Client{
		config:  cfg,
		dataDir: dataDir,
		readers: make(map[DBType]*geoip2.Reader),
	}
}

// Close closes all open database readers.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, reader := range c.readers {
		if reader != nil {
			_ = reader.Close()
		}
	}
	c.readers = make(map[DBType]*geoip2.Reader)
}

// Load loads all databases from disk.
func (c *Client) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	types := []DBType{DBTypeCountry, DBTypeCity, DBTypeASN}

	for _, t := range types {
		path := c.getDBPath(t)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			slog.Warn("Database file not found, skipping load", "type", t, "path", path)
			continue
		}

		reader, err := geoip2.Open(path)
		if err != nil {
			return fmt.Errorf("%w: type=%s path=%s err=%w", ErrDBOpenFailed, t, path, err)
		}

		// Close existing reader if any
		if oldReader, ok := c.readers[t]; ok && oldReader != nil {
			_ = oldReader.Close()
		}

		c.readers[t] = reader
		slog.Info("Loaded database", "type", t, "path", path)
	}

	return nil
}

// getDBPath returns the full path for a database file.
func (c *Client) getDBPath(t DBType) string {
	return filepath.Join(c.dataDir, fmt.Sprintf("%s.mmdb", t))
}
