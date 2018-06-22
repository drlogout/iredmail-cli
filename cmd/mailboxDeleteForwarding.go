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
	"log"
	"os"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/fatih/color"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

// mailboxDeleteForwardingCmd represents the delete-forwarding command
var mailboxDeleteForwardingCmd = &cobra.Command{
	Use:   "delete-forwarding",
	Short: "Delete a mailbox forwarding",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires mailbox and destination address")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid mailbox address format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		err = emailx.Validate(args[1])
		if err != nil {
			return fmt.Errorf("Invalid destination address format: \"%v\"", args[1])
		}

		args[1] = emailx.Normalize(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}
		defer server.Close()

		mailboxAddress, destinationAddress := args[0], args[1]

		err = server.MailboxDeleteForwarding(mailboxAddress, destinationAddress)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		success("Successfully deleted mailbox-forwarding %v -> %v\n", mailboxAddress, destinationAddress)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxDeleteForwardingCmd)

	mailboxDeleteForwardingCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  iredmail-cli mailbox delete-forwarding [mailbox] [destination address]{{end}}{{if .HasAvailableSubCommands}}
	{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
	
Aliases:
	{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
	{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}
