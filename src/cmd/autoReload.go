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

	configCmd.AddCommand(autoReloadCmd)

}

var autoReloadCmd = &cobra.Command{
	Use:   "auto-reload",
	Short: "Toggle server auto reload",
	Long:  "Toggle server auto reload when an API endpoint change is detected",
	Run: func(cmd *cobra.Command, args []string) {

		reader := bufio.NewScanner(os.Stdin)

		config.LoadConfigs()

		if config.Conf.Configuration.AutoReload {

			fmt.Println("Current value: on")

			for {

				fmt.Printf("Would you like to disable it ? [yes/No] ")
				reader.Scan()
				choose := strings.ToLower(reader.Text())

				if choose == "yes" || choose == "y" {

					err := config.ChangeConfig("AutoReload", "false")

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

		} else {

			fmt.Println("Current value: off")
			fmt.Println("Would you like to enable it ? [yes/No] ")

			for {

				reader.Scan()
				choose := strings.ToLower(reader.Text())

				if choose == "yes" || choose == "y" {

					err := config.ChangeConfig("AutoReload", "true")

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
