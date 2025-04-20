package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/GoGeoIP/cmd/api"
	"github.com/hibare/GoGeoIP/cmd/db"
	"github.com/hibare/GoGeoIP/cmd/geoip"
	"github.com/hibare/GoGeoIP/cmd/keys"
	"github.com/hibare/GoGeoIP/internal/config"
)

var Version = "0.0.0"

var rootCmd = &cobra.Command{
	Use:     "go_geo_ip",
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

func preChecks() {
	ctx := context.Background()
	assetDir, err := config.Current.GetAssetDirPath()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get asset dir path", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(assetDir, os.ModePerm); err != nil {
		slog.ErrorContext(ctx, "Failed to create asset dir", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(db.DBCmd)
	rootCmd.AddCommand(geoip.GeoIPCmd)
	rootCmd.AddCommand(api.ServeCmd)
	rootCmd.AddCommand(keys.KeysCmd)

	config.Load()

	initFuncs := []func(){
		commonLogger.InitDefaultLogger,
		preChecks,
	}

	// if !config.Current.Util.IsDev {
	// 	initFuncs = append(initFuncs, maxmind.DownloadAllDB)
	// }

	cobra.OnInitialize(initFuncs...)
}
