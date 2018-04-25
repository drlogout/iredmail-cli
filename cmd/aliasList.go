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
	"log"
	"sort"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// aliasListCmd represents the list command
var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		aliases, err := server.AliasList()
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(aliases)

		filter := cmd.Flag("filter").Value.String()
		if filter != "" {
			aliases = aliases.FilterBy(filter)
		}

		iredmail.PrintAliases(aliases)
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)

	aliasListCmd.Flags().StringP("filter", "f", "", "Filter result")
}
