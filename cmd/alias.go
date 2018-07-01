package cmd

import (
	"github.com/spf13/cobra"
)

// aliasCmd represents the 'alias' command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Add, delete, list aliases and their forwardings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
