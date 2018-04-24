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
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// mailboxListCmd represents the list command
var mailboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		mailboxes, err := server.MailboxList()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(mailboxes)

		domainFilter := cmd.Flag("filter").Value.String()
		if domainFilter != "" {
			mailboxes = mailboxes.FilterBy(domainFilter)
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 16, 8, 0, '\t', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "Email (user name)", "Quota", "Name", "Domain")
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "-----------------", "-----", "----", "------")
		for _, m := range mailboxes {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", m.Email, m.Quota, m.Name, m.Domain)
		}
		w.Flush()
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mailboxListCmd.PersistentFlags().String("foo", "", "A help for foo")

	mailboxListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
