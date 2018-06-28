package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

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
			fatal("%v", err)
		}
		defer server.Close()

		mailboxes, err := server.Mailboxes()
		if err != nil {
			fatal("%v", err)
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
}

func printUserList(mailboxes iredmail.Mailboxes) {
	if len(mailboxes) == 0 {
		info("No mailboxes\n")
		return
	}

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Mailbox", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "----", "----------")
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)
	for _, u := range mailboxes {
		fmt.Fprintf(w, "%v\t%v\n", u.Email, u.Quota)
	}

	w.Flush()
}
