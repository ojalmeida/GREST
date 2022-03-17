package cmd

import (
	"bufio"
	"fmt"
	"github.com/ojalmeida/GREST/src/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func init() {

	configCmd.AddCommand(declarationModeCmd)

}

var declarationModeCmd = &cobra.Command{
	Use:   "declaration-mode",
	Short: "Toggles endpoints declaration mode",
	Long:  "Toggles the way from endpoints will be configured: via API or file declaration",
	Run: func(cmd *cobra.Command, args []string) {

		config.LoadConfigs()
		reader := bufio.NewScanner(os.Stdin)

		if config.Conf.Configuration.DeclarationMode == "api" {

			fmt.Println("Current value: API")

			for {

				fmt.Printf("Would you like to change it to file mode ? [yes/No]")
				reader.Scan()
				choose := strings.ToLower(reader.Text())

				if choose == "yes" || choose == "y" {

					err := config.ChangeConfig("DeclarationMode", "file")

					if err != nil {
						log.Printf("error: %s", err.Error())
					}
					break

				} else if choose == "no" || choose == "n" || choose == "" {

					break

				} else {

					fmt.Println("Invalid input")

				}

			}

		} else if config.Conf.Configuration.DeclarationMode == "file" {

			fmt.Println("Current value: file")

			for {

				fmt.Println("Would you like to change it to API ? [yes/No]")
				reader.Scan()
				choose := strings.ToLower(reader.Text())

				if choose == "yes" || choose == "y" {

					err := config.ChangeConfig("DeclarationMode", "api")

					if err != nil {
						log.Printf("error: %s", err.Error())
					}
					break

				} else if choose == "no" || choose == "n" || choose == "" {

					break

				} else {

					fmt.Println("Invalid input")

				}

			}

		}

	},
}
