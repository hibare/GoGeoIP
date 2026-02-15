package db

import (
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/maxmind"
	"github.com/spf13/cobra"
)

var dBDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Geo IP DB",
	Long:  "Download and verify MaxMind GeoIP databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Server.AssetDirPath)
		return mmClient.DownloadAllDB()
	},
}
