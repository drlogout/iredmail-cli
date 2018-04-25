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
	"log"
	"sort"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// forwardingListCmd represents the list command
var forwardingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forwardings",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		forwardings, err := server.ForwardingList()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(forwardings)

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			forwardings = forwardings.FilterBy(filter)
		}

		iredmail.PrintForwardings(forwardings)
	},
}

func init() {
	forwardingCmd.AddCommand(forwardingListCmd)

	forwardingListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
