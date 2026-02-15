package db

import (
	"context"
	"embed"

	commonDB "github.com/hibare/GoCommon/v2/pkg/db"
	"github.com/hibare/GoGeoIP/internal/config"
)

//go:embed migrations/table/*.sql
var TablesFS embed.FS
var TablesPath = "migrations/table"

func New(ctx context.Context, cfg *config.Config) (*commonDB.DB, error) {
	var (
		dsn      string
		database commonDB.Database
		err      error
	)

	if dsn, err = cfg.DB.GetDSN(); err != nil {
		return nil, err
	}

	if database, err = cfg.DB.GetDatabase(); err != nil {
		return nil, err
	}

	return commonDB.NewClient(ctx, commonDB.DatabaseConfig{
		DSN:            dsn,
		MigrationsFS:   TablesFS,
		MigrationsPath: TablesPath,
		DBType:         database,
	})
}
