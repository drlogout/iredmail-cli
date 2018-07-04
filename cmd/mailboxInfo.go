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
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show mailbox info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [MAILBOX_EMAIL] as argument")
		}

		var err error

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [MAILBOX_EMAIL] format: %s", args[0])
		}
		args[0], err = govalidator.NormalizeEmail(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailbox, err := server.Mailbox(args[0])
		if err != nil {
			fatal("%v\n", err)
		}

		printUserInfo(mailbox)
	},
}

func init() {
	mailboxCmd.AddCommand(infoCmd)
}

func printUserInfo(mailbox iredmail.Mailbox) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"MAILBOX", mailbox.Email})
	table.SetAutoFormatHeaders(false)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{})

	table.Append([]string{"Quota", strconv.Itoa(mailbox.Quota)})

	if len(mailbox.MailboxAliases) > 0 {
		table.Append([]string{"Aliases", ""})
		for _, a := range mailbox.MailboxAliases {
			rightColumn := fmt.Sprintf("%s %s", a.Forwarding, arrowRight)
			table.Append([]string{"", rightColumn})
		}
	}

	if len(mailbox.Forwardings) > 0 {
		table.Append([]string{"Forwardings", ""})
		for _, f := range mailbox.Forwardings {
			rightColumn := fmt.Sprintf("%s %s", arrowRight, f.Forwarding)
			table.Append([]string{"", rightColumn})
		}
		keepCopy := "no"
		if mailbox.IsCopyKept() {
			keepCopy = "yes"
		}
		table.Append([]string{"Keep copy in mailbox", keepCopy})
	}
	table.Render()

	// bold := color.New(color.Bold).SprintfFunc()
	// w := new(tabwriter.Writer)
	// w.Init(os.Stdout, 40, 8, 0, ' ', 0)

	// fmt.Fprintf(w, "%v\t%v\n", bold("Mailbox"), mailbox.Email)
	// fmt.Fprintf(w, "%v\t%v KB\n", bold("Quota"), mailbox.Quota)

	// keepCopy := "no"
	// if mailbox.IsCopyKept() {
	// 	keepCopy = "yes"
	// }

	// forwardings := mailbox.Forwardings
	// if len(forwardings) > 0 {
	// 	fmt.Fprintf(w, "%v\n", bold("Forwardings"))
	// 	fmt.Fprintf(w, "%v  %v\t%v\n", bold(""), "Keep copy in mailbox", keepCopy)
	// 	for _, f := range forwardings {
	// 		fmt.Fprintf(w, "%v\t%v -> %v\n", bold(""), f.Address, f.Forwarding)
	// 	}
	// }

	// w.Flush()
}
