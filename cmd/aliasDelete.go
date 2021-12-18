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

// aliasDeleteCmd represents the 'alias delete' command
var aliasDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an alias",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [ALIAS_EMAIL] as argument")
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

		aliasEmail := args[0]

		if !forceDelete {
			fmt.Printf("Do you really want to delete the alias %s (with all its forwardings)? ", aliasEmail)
			delete := askForConfirmation()

			if !delete {
				fatal("cancelled\n")
			}
		}

		err = server.AliasDelete(aliasEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted alias %s\n", aliasEmail)
	},
}

func init() {
	aliasCmd.AddCommand(aliasDeleteCmd)

	aliasDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "force deletion")

	aliasDeleteCmd.SetUsageTemplate(usageTemplate("alias delete [ALIAS_EMAIL]", printFlags))
}
