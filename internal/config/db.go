package config

import (
	"errors"
	"strings"
	"time"

	commonDB "github.com/hibare/GoCommon/v2/pkg/db"
)

var (
	// ErrDBDSNEmpty indicates that the database DSN is empty.
	ErrDBDSNEmpty = errors.New("database DSN is empty")

	// ErrUnsupportedDB indicates that the database type is unsupported.
	ErrUnsupportedDB = errors.New("unsupported database type")
)

const (
	// DefaultDBMaxIdleConn is the default database idle connection count.
	DefaultDBMaxIdleConn = 5

	// DefaultDBMaxOpenConn is the default maximum number of open database connections.
	DefaultDBMaxOpenConn = 10

	// DefaultDBConnMaxLifetime is the default maximum lifetime of a database connection.
	DefaultDBConnMaxLifetime = 10 * time.Minute
)

// DBConfig holds database-related configuration.
type DBConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// Validate checks that the DSN is set.
func (d *DBConfig) Validate() error {
	if d.DSN == "" {
		return ErrDBDSNEmpty
	}
	return nil
}

// GetDSN returns the configured DSN.
func (d *DBConfig) GetDSN() (string, error) {
	if d.DSN == "" {
		return "", ErrDBDSNEmpty
	}
	return d.DSN, nil
}

// GetDatabase returns the database implementation inferred from the DSN scheme.
func (d *DBConfig) GetDatabase() (commonDB.Database, error) {
	switch {
	case strings.HasPrefix(d.DSN, "postgres://"), strings.HasPrefix(d.DSN, "postgresql://"):
		return &commonDB.PostgresDatabase{}, nil
	default:
		return nil, ErrUnsupportedDB
	}
}
