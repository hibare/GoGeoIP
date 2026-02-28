package testhelpers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	commonDB "github.com/hibare/GoCommon/v2/pkg/db"
	"github.com/hibare/Waypoint/internal/db"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	sharedDB     *commonDB.DB
	sharedDBOnce sync.Once
	errSharedDB  error
)

const (
	postgresImage                 = "postgres:18-alpine"
	postgresDB                    = "test"
	postgresUser                  = "user"
	postgresPass                  = "password"
	postgresStartupTimeoutSeconds = 5
	postgresStartupLogOccurrences = 2
)

// SetupSharedTestDB creates a shared test database that persists across all tests.
// The database is created once and reused for all tests in the package.
// This is more efficient than creating a new database for each test.
func SetupSharedTestDB(t *testing.T) *commonDB.DB {
	t.Helper()

	sharedDBOnce.Do(func() {
		ctx := context.Background()

		// Start PostgreSQL container
		pg, err := postgres.Run(
			ctx,
			postgresImage,
			postgres.WithDatabase(postgresDB),
			postgres.WithUsername(postgresUser),
			postgres.WithPassword(postgresPass),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(postgresStartupLogOccurrences).WithStartupTimeout(postgresStartupTimeoutSeconds*time.Second)),
		)
		if err != nil {
			errSharedDB = fmt.Errorf("failed to start postgres container: %w", err)
			return
		}

		// Get connection string
		connStr, err := pg.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			errSharedDB = fmt.Errorf("failed to get connection string: %w", err)
			return
		}

		// Create database client
		client, err := commonDB.NewClient(ctx, commonDB.DatabaseConfig{
			DSN:            connStr,
			MigrationsFS:   db.TablesFS,
			MigrationsPath: db.TablesPath,
			DBType:         &commonDB.PostgresDatabase{},
		})
		if err != nil {
			errSharedDB = fmt.Errorf("failed to create database client: %w", err)
			return
		}

		// Apply migrations
		if err := client.Migrate(); err != nil {
			errSharedDB = fmt.Errorf("failed to run migrations: %w", err)
			return
		}

		// Populate seeds
		if err := PopulateSeeds(client.DB); err != nil {
			errSharedDB = fmt.Errorf("failed to populate seeds: %w", err)
			return
		}

		sharedDB = client
	})

	require.NoError(t, errSharedDB)
	require.NotNil(t, sharedDB)

	return sharedDB
}
