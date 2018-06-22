package cmd

import (
	"github.com/spf13/cobra"
)

var mailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "add, delete, list mailbox, mailbox-aliase & mailbox-forwardings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mailboxCmd)
}
