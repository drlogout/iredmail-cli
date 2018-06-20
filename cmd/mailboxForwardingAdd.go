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
	"log"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

// mailboxForwardingAddCmd represents the add command
var mailboxForwardingAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add mailbox forwarding",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires mailbox email and destination email")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid mailbox email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		err = emailx.Validate(args[1])
		if err != nil {
			return fmt.Errorf("Invalid destination email format: \"%v\"", args[1])
		}

		args[1] = emailx.Normalize(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		err = server.MailboxForwardingAdd(args[0], args[1])
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxForwardingAddCmd)
}
