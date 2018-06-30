package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// aliasListCmd represents the list command
var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}

		aliases, err := server.Aliases()
		if err != nil {
			fatal("%v\n", err)
		}

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			aliases = aliases.FilterBy(filter)
		}

		printAliases(aliases)
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)

	aliasListCmd.Flags().StringP("filter", "f", "", "Filter result")
}

func printAliases(aliases iredmail.Aliases) {
	if len(aliases) == 0 {
		info("No aliases\n")
		return
	}

	bold := color.New(color.Bold).SprintfFunc()
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)

	// fmt.Fprintf(w, "%v\t%v\n", bold("Mailbox"), mailbox.Email)
	// fmt.Fprintf(w, "%v\t%v KB\n", bold("Quota"), mailbox.Quota)

	fmt.Fprintf(w, "%v\t%v\n", bold("Aliases"), bold("Forwardings"))
	for _, a := range aliases {

		fmt.Fprintf(w, "%v\n", a.Address)
		for _, f := range a.Forwardings {
			fmt.Fprintf(w, "\t➞ %v\n", f.Forwarding)
		}
	}

	w.Flush()
}
