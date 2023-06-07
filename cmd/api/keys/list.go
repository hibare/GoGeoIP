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
		fmt.Println("\nAvaialable API Keys")
		fmt.Println("--------------------")

		for _, apikey := range config.Current.API.APIKeys {
			qrCode, err := qrcode.New(apikey, qrcode.Medium)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Println(qrCode.ToString(false))
		}
	},
}
