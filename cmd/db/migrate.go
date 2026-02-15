package db

import (
	"log/slog"

	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate DB",
	Long:  `Execute database migrations to apply schema changes, create new tables, or modify existing structures in the Axon database.`,
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

		// if err := dbInstance.RunSQLFromFS(db.ViewsFS, db.ViewsPath); err != nil {
		// 	return err
		// }

		// slog.InfoContext(ctx, "DB views created")
		return nil
	},
}
