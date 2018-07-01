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

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

var (
	keepCopyInMailbox = "yes"
)

// mailboxUpdateCmd represents the 'mailbox update' command
var mailboxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update quota and \"keep copy in mailbox\"",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires mailbox email")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid mailbox email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		updated := false
		keepCopy := keepCopyInMailbox == "yes"
		mailboxEmail := args[0]
		mailbox, err := server.Mailbox(mailboxEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		if cmd.Flag("quota").Changed && quota != mailbox.Quota {
			info("Udating quota...\n")
			mailbox.Quota = quota
			err = server.MailboxUpdate(mailbox)
			if err != nil {
				fatal("%v\n", err)
			}
			updated = true
		}

		if cmd.Flag("keep-copy").Changed && mailbox.IsCopyKept() != keepCopy {
			info("Udating keep-copy...\n")
			err := server.MailboxKeepCopy(mailbox, keepCopy)
			if err != nil {
				fatal("%v\n", err)
			}
			updated = true
		}

		if updated {
			success("Successfully updated mailbox\n")
			mailbox, err = server.Mailbox(mailboxEmail)
			if err != nil {
				fatal("%v\n", err)
			}
			fmt.Println()
			printUserInfo(mailbox)
		} else {
			info("No changes, nothing updated\n")
		}
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxUpdateCmd)

	mailboxUpdateCmd.Flags().IntVarP(&quota, "quota", "q", 0, "Quota")
	mailboxUpdateCmd.Flags().StringVarP(&keepCopyInMailbox, "keep-copy", "k", "yes", "Keep copy in mailbox")
}
