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

var (
	quota           int
	storageBasePath string
)

var mailboxAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a mailbox",
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
			color.Red(err.Error())
			os.Exit(1)
		}

		iredmail.PrintMailboxes(iredmail.Mailboxes{mailbox})
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAddCmd)

	mailboxCmd.PersistentFlags().IntVarP(&quota, "quota", "", 2048, "Quota (default 2048 MB)")
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
