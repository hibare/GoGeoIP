package db

import (
	"log/slog"

	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  "Execute database migrations to apply schema changes, create new tables, or modify existing structures in the Waypoint database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		dbInstance, err := db.New(ctx, config.Current)
		if err != nil {
			return err
		}
		defer func() { _ = dbInstance.Close() }()

		if err := dbInstance.Migrate(); err != nil {
			return err
		}
		slog.InfoContext(ctx, "DB tables migrated successfully")

		return nil
	},
}
