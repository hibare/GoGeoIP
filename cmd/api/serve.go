package api

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/hibare/GoGeoIP/internal/api"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API Server",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		app := &api.App{}
		app.Initialize()
		app.Serve()

		log.Info("shutting down")
		os.Exit(0)
	},
}
