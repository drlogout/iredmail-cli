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

// mailboxAliasDeleteCmd represents the add-alias command
var mailboxAliasDeleteCmd = &cobra.Command{
	Use:   "delete-alias",
	Short: "Delete mailbox alias (e.g. abuse@domain.com)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires alias (email)")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid alias email format: \"%v\"", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		err = server.MailboxAliasDelete(args[0])
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted mailbox alias %v\n", args[0])
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAliasDeleteCmd)
	mailboxAliasDeleteCmd.SetUsageTemplate(usageTemplate("mailbox delete-alias [alias_email]"))
}
