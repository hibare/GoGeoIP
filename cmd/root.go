package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hibare/GoGeoIP/cmd/api"
	"github.com/hibare/GoGeoIP/cmd/api/keys"
	"github.com/hibare/GoGeoIP/cmd/db"
	"github.com/hibare/GoGeoIP/cmd/geoip"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/maxmind"
)

var Version = "0.0.0"

var rootCmd = &cobra.Command{
	Use:     "GoGeoIP",
	Short:   "API to fetch Geo information for an IP",
	Long:    "",
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ScheduleBackgroundJobs() {
	// Schedule regular DB update job
	go maxmind.RunDBDownloadJob()
}

func init() {
	rootCmd.AddCommand(db.DBCmd)
	rootCmd.AddCommand(geoip.GeoIPCmd)
	rootCmd.AddCommand(api.ServeCmd)
	rootCmd.AddCommand(keys.KeysCmd)

	config.Load()

	initFuncs := []func(){
		ScheduleBackgroundJobs,
	}

	if !config.Current.Util.IsDev {
		initFuncs = append(initFuncs, maxmind.DownloadAllDB)
	}

	cobra.OnInitialize(initFuncs...)
}
