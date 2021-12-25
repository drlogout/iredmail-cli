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

	"github.com/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// forwardingListCmd represents the list command
var forwardingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forwardings",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n")
		}
		defer server.Close()

		forwardings, err := server.Forwardings()
		if err != nil {
			fatal("%v\n")
		}

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			forwardings = forwardings.FilterBy(filter)
		}

		printForwardings(forwardings)
	},
}

func init() {
	forwardingCmd.AddCommand(forwardingListCmd)

	forwardingListCmd.Flags().StringP("filter", "f", "", "Filter result")

	forwardingListCmd.SetUsageTemplate(usageTemplate("forwarding list", printFlags))
}

func printForwardings(forwardings iredmail.Forwardings) {
	if len(forwardings) == 0 {
		info("No forwardings\n")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Mailbox Email", "Destination Email", "Keep copy in mailbox"})

	var lastAddress string

	for _, f := range forwardings {
		currentAddress, copyLeft := f.Address, "no"

		if f.IsCopyKeptInMailbox {
			copyLeft = "yes"
		}

		if lastAddress == f.Address {
			currentAddress, copyLeft = "", ""
		}

		table.Append([]string{currentAddress, f.Forwarding, copyLeft})

		lastAddress = f.Address
	}
	table.Render()
}
