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

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

var (
	quota           int
	storageBasePath string
)

// mailboxAddCmd represents the add command
var mailboxAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a mailbox",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires email and password as arguments")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		if len(args[1]) < 10 {
			return errors.New("Password length to short (min length 10)")
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			log.Fatal(err)
		}

		mailbox, err := server.MailboxAdd(args[0], args[1], quota, storageBasePath)
		if err != nil {
			log.Fatal(err)
		}

		iredmail.PrintMailboxes(iredmail.Mailboxes{mailbox})
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAddCmd)

	mailboxCmd.PersistentFlags().IntVarP(&quota, "quota", "q", 2048, "Quota (default 2048 MB)")
	mailboxCmd.PersistentFlags().StringVarP(&storageBasePath, "storage-path", "s", "/var/vmail/vmail1", "Storage base path (default /var/vmail/vmail1)")

	mailboxAddCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
	iredmail-cli mailbox add user@example.com plain_password{{end}}{{if .HasAvailableSubCommands}}
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
