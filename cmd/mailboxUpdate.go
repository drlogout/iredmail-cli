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
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

var (
	keepCopyInMailbox = "yes"
)

// mailboxUpdateCmd represents the 'mailbox update' command
var mailboxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update keep-copy and quota",
	Long: `Update keep-copy and quota.

-k, --keep-copy:
If mailboxes with forwardings should not keep a copy of the forwarded email use "--keep-copy no".
This is only possible if at least one forwarding for [MAILBOX_EMAIL] exists.
By default copies are kept in the mailbox.

-p, --password:
Set password. Min length: ` + strconv.Itoa(passwordMinLength) + ` characters.

-q, --quota:
Set quota "--quota 4096" (in MB).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [MAILBOX_EMAIL]")
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

		updated := false
		mailboxEmail := args[0]

		if cmd.Flag("quota").Changed {
			err = server.MailboxSetQuota(mailboxEmail, quota)
			if err != nil {
				fatal("%v\n", err)
			}
			info("Updating quota...\n")
			updated = true
		}

		if cmd.Flag("keep-copy").Changed {
			err := server.MailboxSetKeepCopy(mailboxEmail, keepCopyInMailbox == "yes")
			if err != nil {
				fatal("%v\n", err)
			}
			info("Updating keep-copy...\n")
			updated = true
		}

		if cmd.Flag("password").Changed {
			pw := cmd.Flag("password").Value.String()
			if len(pw) < passwordMinLength {
				fatal("password length to short (min length " + strconv.Itoa(passwordMinLength) + ")")
			}
			err := server.MailboxSetPassword(mailboxEmail, pw)
			if err != nil {
				fatal("ERROR %v\n", err)
			}
			info("Updating password...\n")
			updated = true
		}

		if updated {
			success("Successfully updated mailbox %s\n", mailboxEmail)
		} else {
			info("No changes, nothing updated\n")
		}
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxUpdateCmd)

	mailboxUpdateCmd.Flags().IntVarP(&quota, "quota", "q", 2048, "Sets quota (in MB)")
	mailboxUpdateCmd.Flags().StringVarP(&keepCopyInMailbox, "keep-copy", "k", "yes", "Sets keep-copy of forwardings")
	mailboxUpdateCmd.Flags().StringP("password", "p", "", "Set password")

	mailboxUpdateCmd.SetUsageTemplate(usageTemplate("mailbox update [MAILBOX_EMAIL]", printFlags))
}
