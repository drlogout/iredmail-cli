package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mailbox list", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can list mailbox", func() {
		if skipMailboxrList && !isCI {
			Skip("can list mailbox")
		}

		mailboxes := []string{
			mailboxName1,
			mailboxName2,
			mailboxName3,
			mailboxName4,
			mailboxName5,
		}

		for _, mailbox := range mailboxes {
			cli := exec.Command(cliPath, "mailbox", "add", mailbox, mailboxPW)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "mailbox", "list")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_users")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can list users and filter result", func() {
		if skipUserList && !isCI {
			Skip("can list users and filter result")
		}

		users := []string{
			mailboxName1,
			mailboxName2,
			mailboxName3,
			mailboxName4,
			mailboxName5,
		}

		for _, user := range users {
			cli := exec.Command(cliPath, "mailbox", "add", user, mailboxPW)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "mailbox", "list", "-f", "domain.com")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_users_and_filter_result")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
