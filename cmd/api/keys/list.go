package keys

import (
	"fmt"
	"strings"

	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var ListKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List API Keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nAvaialable API Keys")
		fmt.Println("--------------------")
		fmt.Println(strings.Join(config.Current.API.APIKeys, "\n"))
	},
}
