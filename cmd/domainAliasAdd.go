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

// domainAliasAddCmd represents the add command
var domainAliasAddCmd = &cobra.Command{
	Use:   "add-alias",
	Short: "Add an alias domain",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires an alias domain and a destination domain")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		aliasDomain := args[0]
		domain := args[1]

		err = server.DomainAliasAdd(aliasDomain, domain)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added alias domain %v -> %v\n", aliasDomain, domain)
	},
}

func init() {
	domainCmd.AddCommand(domainAliasAddCmd)

	domainAliasAddCmd.Flags().StringP("description", "d", "", "domain description (default: none)")
	domainAliasAddCmd.Flags().StringP("settings", "s", "", "domain settings (default: default_user_quota:2048)")
	domainAliasAddCmd.SetUsageTemplate(usageTemplate("domain add [domain]"))

}
