// Copyright © 2018 Christian Nolte
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/noxpost/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// domainListCmd represents the list command
var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		domains, err := server.Domains()
		if err != nil {
			fatal("%v\n", err)
		}

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			domains = domains.FilterBy(filter)
		}

		printDomains(domains)
	},
}

func init() {
	domainCmd.AddCommand(domainListCmd)

	domainListCmd.Flags().StringP("filter", "f", "", "Filter result")

	domainListCmd.SetUsageTemplate(usageTemplate("domain list", printFlags))
}

func printDomains(domains iredmail.Domains) {
	if len(domains) == 0 {
		info("No domains\n")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Alias", "Catch-all address", "Description"})

	for _, d := range domains {
		firstAlias := ""
		firstCatchall := ""

		if len(d.Aliases) > 0 {
			firstAlias = d.Aliases[0].AliasDomain
		}

		if len(d.Catchalls) > 0 {
			firstCatchall = d.Catchalls[0].Forwarding
		}

		table.Append([]string{d.Domain, firstAlias, firstCatchall, d.Description})

		aliases := []string{}
		for i := range d.Aliases {
			if (i + 1) < len(d.Aliases) {
				aliases = append(aliases, d.Aliases[i+1].AliasDomain)
			}
		}

		catchalls := []string{}
		for i := range d.Catchalls {
			if (i + 1) < len(d.Catchalls) {
				catchalls = append(catchalls, d.Catchalls[i+1].Forwarding)
			}
		}

		var longerSlice []string
		aliasLength := len(aliases)
		catchallLength := len(catchalls)

		if aliasLength >= catchallLength {
			longerSlice = aliases
		} else {
			longerSlice = catchalls
		}

		for i := range longerSlice {
			alias, catchall := "", ""
			if aliasLength > i {
				alias = aliases[i]
			}
			if catchallLength > i {
				catchall = catchalls[i]
			}
			table.Append([]string{"", alias, catchall, ""})
		}
	}
	table.Render()

}
