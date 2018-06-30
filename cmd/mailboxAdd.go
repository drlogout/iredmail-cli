package cmd

import (
	"errors"
	"fmt"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

var mailboxAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a mailbox (e.g. post@domain.com)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires mailbox (email) and password as arguments")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid mailbox format: \"%v\"", args[0])
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
			fatal("%v\n", err)
		}
		defer server.Close()

		mailbox, err := server.MailboxAdd(args[0], args[1], quota, storageBasePath)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added mailbox %v (quota: %v KB)\n", mailbox.Email, mailbox.Quota)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAddCmd)

	mailboxAddCmd.Flags().IntVarP(&quota, "quota", "q", 2048, "Quota (default 2048 MB)")
	mailboxAddCmd.Flags().StringVarP(&storageBasePath, "storage-path", "s", "/var/vmail/vmail1", "Storage base path")

	mailboxAddCmd.SetUsageTemplate(usageTemplate("mailbox add [mailbox_email] [plain_password]"))
}
