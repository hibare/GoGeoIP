package maxmind

import (
	"github.com/spf13/cobra"
)

var MaxmindCmd = &cobra.Command{
	Use:   "maxmind",
	Short: "MaxMind database management",
	Long:  "Manage MaxMind GeoIP databases including downloading and updating database files.",
}

func init() {
	MaxmindCmd.AddCommand(downloadCmd)
}
