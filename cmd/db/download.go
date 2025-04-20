package db

import (
	"github.com/hibare/GoGeoIP/internal/maxmind"
	"github.com/spf13/cobra"
)

var dBDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Geo IP DB",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		maxmind.DownloadAllDB()
	},
}
