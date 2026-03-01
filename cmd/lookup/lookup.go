package lookup

import (
	"encoding/json"
	"fmt"

	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/maxmind"
	"github.com/spf13/cobra"
)

var LookupCmd = &cobra.Command{
	Use:   "lookup <ip>",
	Short: "Lookup IP geolocation information",
	Long:  "Lookup geographic information for a given IP address using MaxMind GeoIP databases. Supports both IPv4 and IPv6 addresses.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := args[0]

		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Core.DataDir)
		if err := mmClient.Load(); err != nil {
			return fmt.Errorf("failed to load MaxMind databases: %w", err)
		}

		record, err := mmClient.IP2Geo(ip)
		if err != nil {
			return fmt.Errorf("error fetching record: %w", err)
		}

		b, err := json.MarshalIndent(record, "", "    ")
		if err != nil {
			return fmt.Errorf("error parsing record: %w", err)
		}

		cmd.Println(string(b))
		return nil
	},
	SilenceUsage: true,
}
