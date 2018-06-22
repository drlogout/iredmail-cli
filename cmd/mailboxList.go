package cmd

import (
	"log"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var mailboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List mailboxes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}
		defer server.Close()

		mailboxes, err := server.MailboxList()
		if err != nil {
			log.Fatal(err)
		}

		domainFilter := cmd.Flag("filter").Value.String()
		if domainFilter != "" {
			mailboxes = mailboxes.FilterBy(domainFilter)
		}

		mailboxes.Print()
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxListCmd)

	mailboxListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
