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
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var (
	forceDelete = false
)

// mailboxDeleteCmd represents the delete command
var mailboxDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a mailbox",
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

		mailboxEmail := args[0]

		if !forceDelete {
			fmt.Printf("Do you really want to delete the mailbox %s (with all its alias mailboxes and forwardings)? ", mailboxEmail)
			delete := askForConfirmation()

			if !delete {
				fatal("cancelled\n")
			}
		}

		err = server.MailboxDelete(mailboxEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted mailbox %s\n", mailboxEmail)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxDeleteCmd)

	mailboxDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "force deletion")

	mailboxDeleteCmd.SetUsageTemplate(usageTemplate("mailbox delete [MAILBOX_EMAIL]", printFlags))
}
