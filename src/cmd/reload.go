package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(reloadCmd)

}

// to work, it will be necessary implement interprocess communication through gRPC
var reloadCmd = &cobra.Command{

	Use:   "reload",
	Short: "Reloads server",
	Long:  "Stops server, reads endpoint definitions, apply them, and them starts the server again",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
