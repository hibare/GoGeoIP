package db

import (
	"github.com/spf13/cobra"
)

var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "IP DB related commands",
	Long:  "",
}

func init() {
	DBCmd.AddCommand(dBDowloadCmd)
}
