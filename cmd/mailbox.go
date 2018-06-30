package cmd

import (
	"github.com/spf13/cobra"
)

var (
	quota           int
	storageBasePath string
)

var mailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "add, delete, list mailbox | mailbox-alias | mailbox-forwarding",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mailboxCmd)
}
