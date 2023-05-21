package keys

import (
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage API Keys",
	Long:  ``,
}

func init() {
	KeysCmd.AddCommand(ListKeysCmd)
}
