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
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// domainAliasListCmd represents the 'domain list-alias' command
var domainAliasListCmd = &cobra.Command{
	Use:   "list-alias",
	Short: "List [ALIAS_DOMAIN]s",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}

		aliasDomains, err := server.DomainAliasList()
		if err != nil {
			fatal("%v\n", err)
		}

		// filter := cmd.Flag("filter").Value.String()
		// if filter != "" {
		// 	domains = domains.FilterBy(filter)
		// }

		printAliasDomains(aliasDomains)
	},
}

func init() {
	domainCmd.AddCommand(domainAliasListCmd)

	domainAliasListCmd.Flags().StringP("filter", "f", "", "Filter result")
}

func printAliasDomains(aliasDomains iredmail.AliasDomains) {
	if len(aliasDomains) == 0 {
		info("No [ALIAS_DOMAIN]s\n")
		return
	}

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "Alias domain", "Domain")
	fmt.Fprintf(w, "%v\t%v\n", "------------", "------")
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)
	for _, a := range aliasDomains {
		fmt.Fprintf(w, "%v\t-> %v\n", a.AliasDomain, a.Domain)
	}

	w.Flush()
}
