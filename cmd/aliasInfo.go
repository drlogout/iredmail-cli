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
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// aliasInfoCmd represents the 'alias info' command
var aliasInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show alias info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires alias email as argument")
		}

		var err error

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid alias email format: \"%v\"", args[0])
		}
		args[0], err = govalidator.NormalizeEmail(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}

		alias, err := server.Alias(args[0])
		if err != nil {
			fatal("%v\n", err)
		}

		printAlias(alias)
	},
}

func init() {
	aliasCmd.AddCommand(aliasInfoCmd)

	aliasInfoCmd.SetUsageTemplate(usageTemplate("alias info [ALIAS_EMAIL]"))
}

func printAlias(alias iredmail.Alias) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Alias", "Forwardings"})
	table.SetAutoMergeCells(true)

	firstForwarding := ""
	if len(alias.Forwardings) > 0 {
		firstForwarding = alias.Forwardings[0].Forwarding
	}
	table.Append([]string{alias.Address, firstForwarding})
	for i := range alias.Forwardings {
		if (i + 1) < len(alias.Forwardings) {
			table.Append([]string{alias.Address, alias.Forwardings[i+1].Forwarding})
		}
	}
	table.Render()
}
