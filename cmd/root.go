package cmd

import (
	"log/slog"
	"os"

	"github.com/hibare/GoGeoIP/cmd/db"
	"github.com/hibare/GoGeoIP/cmd/lookup"
	"github.com/hibare/GoGeoIP/cmd/maxmind"
	"github.com/hibare/GoGeoIP/cmd/server"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

var rootCmd = &cobra.Command{
	Use:     "go_geo_ip",
	Short:   "API to fetch Geo information for an IP",
	Long:    "",
	Version: constants.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "Path to config file")

	// Add subcommands
	rootCmd.AddCommand(db.DBCmd)
	rootCmd.AddCommand(maxmind.MaxmindCmd)
	rootCmd.AddCommand(lookup.LookupCmd)
	rootCmd.AddCommand(server.ServeCmd)

	// Load config with context
	ctx := rootCmd.Context()
	if _, err := config.Load(ctx, ConfigPath); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
}
