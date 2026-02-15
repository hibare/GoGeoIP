package db

import "github.com/spf13/cobra"

var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage DB",
	Long:  `Manage database operations including migrations and schema management for the Axon application.`,
}

func init() {
	DBCmd.AddCommand(migrateCmd)
}
