package maxmind

import (
	"github.com/spf13/cobra"
)

var MaxmindCmd = &cobra.Command{
	Use:   "maxmind",
	Short: "Manage MaxMind GeoIP databases",
	Long:  "",
}

func init() {
	MaxmindCmd.AddCommand(downloadCmd)
}
