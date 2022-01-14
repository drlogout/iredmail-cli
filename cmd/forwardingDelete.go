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

// forwardingDeleteCmd represents the delete command
var forwardingDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete forwarding",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires [MAILBOX_EMAIL] and [DESTINATION_EMAIL]")
		}

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [MAILBOX_EMAIL] format: %s", args[0])
		}

		if !govalidator.IsEmail(args[1]) {
			return fmt.Errorf("Invalid [DESTINATION_EMAIL] format: %s", args[1])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailboxEmail, destinationEmail := args[0], args[1]

		err = server.ForwardingDelete(mailboxEmail, destinationEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted forwarding %s %s %s\n", mailboxEmail, arrowRight, destinationEmail)
	},
}

func init() {
	forwardingCmd.AddCommand(forwardingDeleteCmd)

	forwardingDeleteCmd.SetUsageTemplate(usageTemplate("forwarding delete [MAILBOX_EMAIL] [DESTINATION_EMAIL]"))
}
