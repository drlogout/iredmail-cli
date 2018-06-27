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

// domainListCmd represents the list command
var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}

		domains, err := server.DomainList()
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
}

func printDomains(domains iredmail.Domains) {
	if len(domains) == 0 {
		info("No domains\n")
		return
	}

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\n", "Domain", "Settings", "Description")
	fmt.Fprintf(w, "%v\t%v\t%v\n", "------", "--------", "-----------")
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)
	for _, d := range domains {
		fmt.Fprintf(w, "%v\t%v\t%v\n", d.Domain, d.Settings, d.Description)
	}

	w.Flush()
}
