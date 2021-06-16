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

// domainAddCmd represents the 'domain add' command
var domainAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a domain",
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
		description := cmd.Flag("description").Value.String()
		settings := cmd.Flag("settings").Value.String()

		if settings == "" {
			settings = iredmail.DomainDefaultSettings
		}

		err = server.DomainAdd(iredmail.Domain{
			Domain:      domain,
			Description: description,
			Settings:    settings,
		})
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added domain %s\n", domain)
	},
}

func init() {
	domainCmd.AddCommand(domainAddCmd)

	domainAddCmd.Flags().StringP("description", "d", "", "domain description (default: none)")
	domainAddCmd.Flags().StringP("settings", "s", "", "domain settings (default: default_user_quota:2048)")

	domainAddCmd.SetUsageTemplate(usageTemplate("domain add [DOMAIN]", printFlags))
}
