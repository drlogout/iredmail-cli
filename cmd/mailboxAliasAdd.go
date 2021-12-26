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

	"github.com/asaskevich/govalidator"
	"github.com/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// mailboxAliasAddCmd represents the add-alias command
var mailboxAliasAddCmd = &cobra.Command{
	Use:   "add-alias",
	Short: "Add a mailbox alias",
	Long: `Add a mailbox alias
	
A mailbox [MAILBOX_EMAIL] can have additional email addresses;
 [ALIAS]@[DOMAIN|ALIAS_DOMAIN], all emails sent to these addresses 
 will be delivered to the same mailbox.
 Only this type of alias gets SENDING permission. Even with domain alias, names must be declared.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires [ALIAS_EMAIL] and [MAILBOX_EMAIL]")
		}

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [ALIAS_EMAIL] format: %s, to support alias domains iredmail-cli now needs full alias email.", args[0])
		}

		if !govalidator.IsEmail(args[1]) {
			return fmt.Errorf("Invalid [MAILBOX_EMAIL] format: %s", args[1])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		err = server.MailboxAliasAdd(args[0], args[1])
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added mailbox alias %s %s %s\n", args[0], arrowRight, args[1])
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAliasAddCmd)

	mailboxAliasAddCmd.SetUsageTemplate(usageTemplate("mailbox add-alias [ALIAS] [MAILBOX_EMAIL]"))
}
