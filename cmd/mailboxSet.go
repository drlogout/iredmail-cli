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
	keepCopyInMailbox = true
)

// mailboxSetCmd represents the set command
var mailboxSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set quota and \"keep copy in mailbox\"",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires mailbox email as argument")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid mailbox email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailboxEmail := args[0]
		mailbox, err := server.Mailbox(mailboxEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		if cmd.Flag("quota").Changed {
			if quota != mailbox.Quota {
				fmt.Println(quota, mailbox.Quota)
				mailbox.Quota = quota
			}
		}

		if cmd.Flag("keep-copy").Changed {

		}

		err = server.MailboxSet(mailbox)
		if err != nil {
			fatal("%v\n", err)
		}

	},
}

func init() {
	mailboxCmd.AddCommand(mailboxSetCmd)

	mailboxSetCmd.Flags().IntVarP(&quota, "quota", "q", 0, "Quota")
	mailboxSetCmd.Flags().BoolVarP(&keepCopyInMailbox, "keep-copy", "-k", true, "Keep copy in mailbox")
}
