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
	"github.com/kostagorod/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// aliasAddForwardingCmd represents the 'alias add-forwarding' command
var aliasAddForwardingCmd = &cobra.Command{
	Use:   "add-forwarding",
	Short: "Add forwarding to an alias",
	Long: `Add forwarding to an alias.

Emails sent to [ALIAS_EMAIL] will be delivered to [DESTINATION_EMAIL].
An alias can have multiple forwardings.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires [ALIAS_EMAIL] and [DESTINATION_EMAIL]")
		}

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [ALIAS_EMAIL] format: %s", args[0])
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

		aliasEmail := args[0]
		forwardingEmail := args[1]

		err = server.AliasForwardingAdd(aliasEmail, forwardingEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added alias forwarding %s %s %s\n", aliasEmail, arrowRight, forwardingEmail)
	},
}

func init() {
	aliasCmd.AddCommand(aliasAddForwardingCmd)

	aliasAddForwardingCmd.SetUsageTemplate(usageTemplate("alias add-forwarding [ALIAS_EMAIL] [DESTINATION_EMAIL]"))
}
