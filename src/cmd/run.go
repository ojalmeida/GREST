package cmd

import (
	"fmt"
	"github.com/ojalmeida/GREST/src/server"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func init() {

	rootCmd.AddCommand(runCmd)

}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs GREST",
	Long:  `Start serving requests`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("PID:", os.Getpid())

		stopChannel := make(chan os.Signal, 1)
		signal.Notify(stopChannel, os.Interrupt)

		go func() {

			<-stopChannel
			os.Exit(1)

		}()

		for {

			server.Start() // blocks execution

		}

	},
}
