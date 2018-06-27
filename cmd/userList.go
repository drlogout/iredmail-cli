package cmd

import (
	"bytes"
	"fmt"
	"os"
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

		users, err := server.Users()
		if err != nil {
			fatal("%v", err)
		}

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			users = users.FilterBy(filter)
		}

		printUserList(users)
	},
}

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().StringP("filter", "f", "", "Filter result")
}

func printUserList(users iredmail.Users) {
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "User", "Quota (KB)")
	fmt.Fprintf(w, "%v\t%v\n", "----", "----------")
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)
	for _, u := range users {
		fmt.Fprintf(w, "%v\t%v\n", u.Email, u.Quota)
	}

	w.Flush()
}
