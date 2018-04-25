package cmd

import (
	"github.com/spf13/cobra"
)

var forwardingCmd = &cobra.Command{
	Use:   "forwarding",
	Short: "Add, remove, list forwardings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(forwardingCmd)
}
