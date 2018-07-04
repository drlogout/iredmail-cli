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
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// domainAliasAddCmd represents the 'domain add-alias' command
var domainAliasAddCmd = &cobra.Command{
	Use:   "add-alias",
	Short: "Add an alias domain",
	Long:  "Emails sent to user@[ALIAS_DOMAIN] will be delivered to user@[DOMAIN]",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires an [ALIAS_DOMAIN] and a [DOMAIN]")
		}

		if !govalidator.IsDNSName(args[0]) {
			return fmt.Errorf("Invalid [ALIAS_DOMAIN] name format: \"%v\"", args[0])
		}
		args[0] = strings.ToLower(args[0])

		if !govalidator.IsDNSName(args[1]) {
			return fmt.Errorf("Invalid [DOMAIN] name format: \"%v\"", args[1])
		}
		args[1] = strings.ToLower(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		aliasDomain := args[0]
		targetDomain := args[1]

		err = server.DomainAliasAdd(aliasDomain, targetDomain)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added alias domain %v %v %v\n", aliasDomain, arrowRight, targetDomain)
	},
}

func init() {
	domainCmd.AddCommand(domainAliasAddCmd)

	domainAliasAddCmd.SetUsageTemplate(usageTemplate("domain add-alias [ALIAS_DOMAIN] [DOMAIN]"))
}
