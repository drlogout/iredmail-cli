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
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// mailboxInfoCmd represents the info command
var mailboxInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show mailbox info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [MAILBOX_EMAIL] as argument")
		}

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [MAILBOX_EMAIL] format: %s", args[0])
		}

		return nil
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

		printMailboxInfo(mailbox, prettyPrint)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxInfoCmd)
}

func printMailboxInfo(mailbox iredmail.Mailbox, prettyPrint bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"MAILBOX", mailbox.Email})
	table.SetAutoFormatHeaders(false)

	if prettyPrint {
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
		table.SetColumnColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{})
	}

	table.Append([]string{"Display Name", mailbox.Name})
	table.Append([]string{"Quota", fmt.Sprintf("%v MB", strconv.Itoa(mailbox.Quota))})

	if len(mailbox.MailboxAliases) > 0 {
		name := strings.Split(mailbox.MailboxAliases[0].Address, "@")[0]
		table.Append([]string{"Mailbox aliases", name})
		for i := range mailbox.MailboxAliases {
			if (i + 1) < len(mailbox.MailboxAliases) {
				name = strings.Split(mailbox.MailboxAliases[i+1].Address, "@")[0]
				table.Append([]string{"", name})
			}
		}
	}

	if len(mailbox.Forwardings) > 0 {
		table.Append([]string{"Forwardings", mailbox.Forwardings[0].Forwarding})
		for i := range mailbox.Forwardings {
			if (i + 1) < len(mailbox.Forwardings) {
				table.Append([]string{"", mailbox.Forwardings[i+1].Forwarding})
			}
		}
		keepCopy := "no"
		if mailbox.Forwardings[0].IsCopyKeptInMailbox {
			keepCopy = "yes"
		}
		table.Append([]string{"Keep copy in mailbox", keepCopy})
	}

	table.Append([]string{"Maildir", mailbox.MailDir})

	table.Render()
}
