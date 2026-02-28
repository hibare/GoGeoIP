package cmd

import (
	"os"

	"github.com/hibare/Waypoint/cmd/db"
	"github.com/hibare/Waypoint/cmd/lookup"
	"github.com/hibare/Waypoint/cmd/maxmind"
	"github.com/hibare/Waypoint/cmd/server"
	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/constants"
	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

var rootCmd = &cobra.Command{
	Use:     "waypoint",
	Short:   "IP Geo location Service",
	Long:    "Waypoint is an IP geo location service that provides geographic information for any IP address using MaxMind GeoIP databases.",
	Version: constants.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load(cmd.Context(), ConfigPath)
		return err
	},
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
}
