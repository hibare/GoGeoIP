package maxmind

import (
	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/maxmind"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download MaxMind database files",
	Long:  "Download and verify MaxMind GeoIP2/GeoLite2 databases. Requires a MaxMind license key to be configured.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Core.DataDir)
		return mmClient.DownloadAllDB()
	},
}
