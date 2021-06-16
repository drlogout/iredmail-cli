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
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/KostaGorod/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// domainDeleteCmd represents the 'domain delete' command
var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires [DOMAIN] as argument")
		}

		if !govalidator.IsDNSName(args[0]) {
			return fmt.Errorf("Invalid [DOMAIN] format: %s", args[0])
		}
		args[0] = strings.ToLower(args[0])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		domain := args[0]

		if !forceDelete {
			fmt.Printf("Do you really want to delete the domain %s (with all its alias domains and catch-all forwardings)? ", domain)
			delete := askForConfirmation()

			if !delete {
				fatal("cancelled\n")
			}
		}

		err = server.DomainDelete(domain)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted domain %s\n", domain)
	},
}

func init() {
	domainCmd.AddCommand(domainDeleteCmd)
	domainDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "force deletion")

	domainDeleteCmd.SetUsageTemplate(usageTemplate("domain delete [DOMAIN]"))
}
