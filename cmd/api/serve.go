package api

import (
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API Server",
	Long:  "",
	RunE: func(_ *cobra.Command, _ []string) error {
		srv, err := NewServer()
		if err != nil {
			return err
		}
		return srv.Start()
	},
}
