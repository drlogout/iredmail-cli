package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/spf13/cobra"
)

// mailboxAddCmd represents the add command
var mailboxAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a mailbox",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires [MAILBOX_EMAIL] and [PLAIN_PASSWORD] as arguments")
		}

		var err error

		if !govalidator.IsEmail(args[0]) {
			return fmt.Errorf("Invalid [MAILBOX_EMAIL] format: %s", args[0])
		}
		args[0], err = govalidator.NormalizeEmail(args[0])
		if err != nil {
			return err
		}

		if len(args[1]) < passwordMinLength {
			return errors.New("[PLAIN_PASSWORD] length to short (min length " + strconv.Itoa(passwordMinLength) + ")")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		mailboxEmail, password := args[0], args[1]
		err = server.MailboxAdd(mailboxEmail, password, quota, storageBasePath)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added mailbox %s\n", mailboxEmail)
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAddCmd)

	mailboxAddCmd.Flags().IntVarP(&quota, "quota", "q", 2048, "Quota (default 2048 MB)")
	mailboxAddCmd.Flags().StringVarP(&storageBasePath, "storage-path", "s", "/var/vmail/vmail1", "Storage base path")

	mailboxAddCmd.SetUsageTemplate(usageTemplate("mailbox add [MAILBOX_EMAIL] [PLAIN_PASSWORD]", printFlags))
}
