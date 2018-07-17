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

// domainCatchalldelete represents the 'domain delete-catchall' command
var domainCatchalldelete = &cobra.Command{
	Use:   "delete-catchall",
	Short: "Delete a per-domain catch-all forwarding",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires [DOMAIN] and [DESTINATION_EMAIL] as argument")
		}

		if !govalidator.IsDNSName(args[0]) {
			return fmt.Errorf("Invalid [DOMAIN] format: %s", args[0])
		}
		args[0] = strings.ToLower(args[0])

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
		defer server.Close()

		domain := args[0]
		catchallEmail := args[1]

		err = server.DomainCatchallDelete(domain, catchallEmail)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully deleted catch-all forwarding %s %s %s\n", domain, arrowRight, catchallEmail)
	},
}

func init() {
	domainCmd.AddCommand(domainCatchalldelete)

	domainCatchalldelete.SetUsageTemplate(usageTemplate("domain delete-catchall [DOMAIN] [DESTINATION_EMAIL]"))
}
