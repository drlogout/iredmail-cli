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

// mailboxAliasDeleteCmd represents the delete-alias command
var mailboxAliasDeleteCmd = &cobra.Command{
	Use:   "delete-alias",
	Short: "Delete mailbox alias",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [ALIAS_EMAIL]")
		}

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [ALIAS_EMAIL] format: %s", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailboxAliasEmail := args[0]

		err = server.MailboxAliasDelete(mailboxAliasEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted mailbox alias %s\n", mailboxAliasEmail)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAliasDeleteCmd)

	mailboxAliasDeleteCmd.SetUsageTemplate(usageTemplate("mailbox delete-alias [ALIAS_EMAIL]"))
}
