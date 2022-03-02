package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(configCmd)

}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure GREST",
	Long:  `Configure aspects of GREST functioning such as server auto reload, declaration mode, etc`,
}
