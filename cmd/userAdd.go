package cmd

import (
	"errors"
	"fmt"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

var (
	quota           int
	storageBasePath string
)

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a user (e.g. post@domain.com)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires user email and password as arguments")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid user email format: \"%v\"", args[0])
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

		user, err := server.UserAdd(args[0], args[1], quota, storageBasePath)
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added user %v (quota: %v KB)\n", user.Email, user.Quota)
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)

	userAddCmd.Flags().IntVarP(&quota, "quota", "", 2048, "Quota (default 2048 MB)")
	userAddCmd.Flags().StringVarP(&storageBasePath, "storage-path", "s", "/var/vmail/vmail1", "Storage base path")

	userAddCmd.SetUsageTemplate(usageTemplate("user add [user_email] [plain_password]"))
}
