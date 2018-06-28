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

// forwardingListCmd represents the add command
var forwardingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forwardings",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n")
		}
		defer server.Close()

		forwardings, err := server.ForwardingList()
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
	forwardingListCmd.SetUsageTemplate(usageTemplate("forwarding list"))
}

type fwMap map[string][]string

func (fwMap fwMap) Contains(f iredmail.Forwarding) bool {
	fwEntry, ok := fwMap[f.Address]
	if ok {
		for _, forfwarding := range fwEntry {
			if forfwarding == f.Forwarding {
				return true
			}
		}

		return false
	}
	return false
}

func printForwardings(forwardings iredmail.Forwardings) {
	if len(forwardings) == 0 {
		info("No forwardings\n")
		return
	}

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "%v\t      %v\n", "User Email", "Destination Email")
	fmt.Fprintf(w, "%v\t      %v\n", "----------", "-----------------")
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)

	fw := fwMap{}

	for _, f := range forwardings {
		if _, ok := fw[f.Address]; !ok {
			fw[f.Address] = []string{
				f.Forwarding,
			}
		} else {
			if !fw.Contains(f) {
				fw[f.Address] = append(fw[f.Address], f.Forwarding)
			}
		}
		fmt.Fprintf(w, "%v\t->    %v\n", f.Address, f.Forwarding)
	}

	w.Flush()
}
