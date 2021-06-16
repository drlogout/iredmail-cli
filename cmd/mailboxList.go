package cmd

import (
	"os"
	"strconv"

	"github.com/KostaGorod/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// mailboxListCmd represents the 'mailbox list' command
var mailboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List mailboxes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailboxes, err := server.Mailboxes()
		if err != nil {
			fatal("%v\n", err)
		}

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			mailboxes = mailboxes.FilterBy(filter)
		}

		printUserList(mailboxes)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxListCmd)

	mailboxListCmd.Flags().StringP("filter", "f", "", "Filter result")

	mailboxListCmd.SetUsageTemplate(usageTemplate("mailbox list", printFlags))
}

func printUserList(mailboxes iredmail.Mailboxes) {
	if len(mailboxes) == 0 {
		info("No mailboxes\n")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Display Name", "Mailbox", "Quota (MB)"})

	for _, m := range mailboxes {
		table.Append([]string{m.Name, m.Email, strconv.Itoa(m.Quota)})
	}
	table.Render()
}
