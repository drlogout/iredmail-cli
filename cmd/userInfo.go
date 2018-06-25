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
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show user info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires user email as sole argument")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		user, err := server.User(args[0])
		if err != nil {
			fatal("%v\n", err)
		}

		printUserInfo(user)
	},
}

func init() {
	userCmd.AddCommand(infoCmd)
}

func printUserInfo(user iredmail.User) {
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "User:\t%v\n", user.Email)
	w.Flush()
	info(buf.String())

	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 8, 0, ' ', 0)
	fmt.Fprintf(w, "Quota:\t%v\n", user.Quota)
	// fmt.Fprintf(w, "Maildir:\t%v\n", user.MailDir)

	forwardings := iredmail.Forwardings{}
	for _, f := range user.Forwardings {
		if f.Forwarding != user.Email {
			forwardings = append(forwardings, f)
		}
	}
	if len(forwardings) > 0 {
		fmt.Fprintf(w, "%v\n", "Forwardings:")
		for _, f := range forwardings {
			fmt.Fprintf(w, "\t%v -> %v\n", f.Address, f.Forwarding)
		}
	}

	w.Flush()
}
