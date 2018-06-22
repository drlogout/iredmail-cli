package cmd

import (
	"log"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List useers",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}
		defer server.Close()

		useers, err := server.UserList()
		if err != nil {
			log.Fatal(err)
		}

		domainFilter := cmd.Flag("filter").Value.String()
		if domainFilter != "" {
			useers = useers.FilterBy(domainFilter)
		}

		useers.Print()
	},
}

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
