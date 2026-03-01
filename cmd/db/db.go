package db

import "github.com/spf13/cobra"

var DBCmd = &cobra.Command{
	Use:          "db",
	Short:        "Database management commands",
	Long:         "Manage database operations including migrations and schema management for Waypoint.",
	SilenceUsage: true,
}

func init() {
	DBCmd.AddCommand(migrateCmd)
}
