package keys

import (
	"fmt"

	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var ListKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List API Keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nAvailable API Keys")
		fmt.Println("--------------------")

		for _, apiKey := range config.Current.API.APIKeys {
			qrCode, err := qrcode.New(apiKey, qrcode.Medium)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Println(qrCode.ToString(false))
		}
	},
}
