package maxmind

import (
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/maxmind"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Geo IP DB",
	Long:  "Download and verify MaxMind GeoIP databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Core.DataDir)
		return mmClient.DownloadAllDB()
	},
}
