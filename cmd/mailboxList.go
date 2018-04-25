package cmd

import (
	"log"
	"sort"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var mailboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		mailboxes, err := server.MailboxList()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(mailboxes)

		domainFilter := cmd.Flag("filter").Value.String()
		if domainFilter != "" {
			mailboxes = mailboxes.FilterBy(domainFilter)
		}

		iredmail.PrintMailboxes(mailboxes)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxListCmd)

	mailboxListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
