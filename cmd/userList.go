package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v", err)
		}
		defer server.Close()

		users, err := server.UserList()
		if err != nil {
			fatal("%v", err)
		}

		domainFilter := cmd.Flag("filter").Value.String()
		if domainFilter != "" {
			users = users.FilterBy(domainFilter)
		}

		printUserListHeading()
		fmt.Printf(users.String())
	},
}

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().StringP("filter", "f", "", "Filter result")
}

func printUserListHeading() {
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "User", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "----", "----------")
	w.Flush()
	info(buf.String())
}
