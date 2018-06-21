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
	forceDelete = false
)

// mailboxDeleteCmd represents the delete command
var mailboxDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a mailbox",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires email as single argument")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}

		email := args[0]

		if !forceDelete {
			fmt.Printf("Do you really want to delete the mailbox %v? ", email)
			delete := askForConfirmation()

			if !delete {
				fatal("cancelled\n")
			}
		}

		err = server.MailboxDelete(email)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted mailbox %v\n", email)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mailboxDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mailboxDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "force deletion")
}
