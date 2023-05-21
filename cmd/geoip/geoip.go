package geoip

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/hibare/GoGeoIP/internal/maxmind"
)

var GeoIPCmd = &cobra.Command{
	Use:   "geoip",
	Short: "Lookup Geo information for an IP",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ip, err := cmd.Flags().GetString("ip")
		if err != nil {
			log.Fatalf("Error fetching flags %s", err)
		}

		if ip == "" {
			cmd.Help()
			return
		}

		record, err := maxmind.IP2Geo(ip)
		if err != nil {
			log.Fatalf("Error fetching record: %s", err)
		}

		b, err := json.MarshalIndent(record, "", "    ")
		if err != nil {
			log.Fatalf("Error parsing record: %s", err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	GeoIPCmd.PersistentFlags().String("ip", "", "IP to lookup")
}
