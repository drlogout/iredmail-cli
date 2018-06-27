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

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// domainDeleteCmd represents the add command
var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a domain name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		domain := args[0]

		err = server.DomainDelete(domain)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted domain %v\n", domain)
	},
}

func init() {
	domainCmd.AddCommand(domainDeleteCmd)

	domainDeleteCmd.Flags().StringP("description", "d", "", "domain description (default: none)")
	domainDeleteCmd.Flags().StringP("settings", "s", "", "domain settings (default: default_user_quota:2048)")
	domainDeleteCmd.SetUsageTemplate(usageTemplate("domain add [domain]"))

}
