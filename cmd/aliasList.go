package cmd

import (
	"log"
	"sort"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// aliasListCmd represents the list command
var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		aliases, err := server.AliasList()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(aliases)

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			aliases = aliases.FilterBy(filter)
		}

		iredmail.PrintAliases(aliases)
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)

	aliasListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
