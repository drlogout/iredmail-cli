package cmd

import (
	"os"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// aliasListCmd represents the 'alias list' command
var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	Args:  cobra.NoArgs,
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

	aliasListCmd.SetUsageTemplate(usageTemplate("alias list", printFlags))
}

func printAliases(aliases iredmail.Aliases) {
	if len(aliases) == 0 {
		info("No aliases\n")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Alias", "Forwardings"})

	for _, a := range aliases {
		firstForwarding := ""
		if len(a.Forwardings) > 0 {
			firstForwarding = a.Forwardings[0].Forwarding
		}
		table.Append([]string{a.Address, firstForwarding})
		for i := range a.Forwardings {
			if (i + 1) < len(a.Forwardings) {
				table.Append([]string{"", a.Forwardings[i+1].Forwarding})
			}
		}
	}
	table.Render()
}
