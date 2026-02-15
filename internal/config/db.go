package config

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	commonDB "github.com/hibare/GoCommon/v2/pkg/db"
)

var (
	// ErrDBUsernameEmpty indicates that the database username is empty.
	ErrDBUsernameEmpty = errors.New("database username is empty")

	// ErrDBPasswordEmpty indicates that the database password is empty.
	ErrDBPasswordEmpty = errors.New("database password is empty")

	// ErrDBHostEmpty indicates that the database host is empty.
	ErrDBHostEmpty = errors.New("database host is empty")

	// ErrDBPortInvalid indicates that the database port is invalid.
	ErrDBPortInvalid = errors.New("database port is invalid")

	// ErrDBNameEmpty indicates that the database name is empty.
	ErrDBNameEmpty = errors.New("database name is empty")

	// ErrMissingDBPath indicates that the database path is missing (SQLite).
	ErrMissingDBPath = errors.New("database path is missing (SQLite)")

	// ErrUnsupportedDB indicates that the database type is unsupported.
	ErrUnsupportedDB = errors.New("unsupported database type")
)

const (
	// DefaultDBHost is the default database host.
	DefaultDBHost = "127.0.0.1"

	// DefaultDBPort is the default database port.
	DefaultDBPort = 5432

	// DefaultDBName is the default database name.
	DefaultDBName = "tasks"

	// DefaultDBHost is the default database idle connection count.
	DefaultDBMaxIdleConn = 5

	// DefaultDBMaxOpenConn is the default maximum number of open database connections.
	DefaultDBMaxOpenConn = 10

	// DefaultDBConnMaxLifetime is the default maximum lifetime of a database connection.
	DefaultDBConnMaxLifetime = 10 * time.Minute // 10 minutes
)

// DBType represents the type of the DB.
type DBType string

const (
	// DBTypePostgres represents a PostgreSQL database.
	DBTypePostgres DBType = "postgres"
)

// DBConfig holds database-related configuration.
type DBConfig struct {
	DBType          DBType        `mapstructure:"type"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (d *DBConfig) Validate() error {
	switch d.DBType {
	case DBTypePostgres:
		if d.Username == "" {
			return ErrDBUsernameEmpty
		}
		if d.Password == "" {
			return ErrDBPasswordEmpty
		}
		if d.Host == "" {
			return ErrDBHostEmpty
		}
		if d.Port <= 0 {
			return ErrDBPortInvalid
		}
		if d.Name == "" {
			return ErrDBNameEmpty
		}

	default:
		return ErrUnsupportedDB
	}
	return nil
}

// GetDSN constructs the Data Source Name based on DB type.
func (d *DBConfig) GetDSN() (string, error) {
	switch d.DBType {
	case DBTypePostgres:
		return fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			d.Username, d.Password, net.JoinHostPort(d.Host, strconv.Itoa(d.Port)), d.Name,
		), nil
	default:
		return "", ErrUnsupportedDB
	}
}

// GetDatabase returns the database implementation based on DB type.
func (d *DBConfig) GetDatabase() (commonDB.Database, error) {
	switch d.DBType {
	case DBTypePostgres:
		return &commonDB.PostgresDatabase{}, nil
	default:
		return nil, ErrUnsupportedDB
	}
}
