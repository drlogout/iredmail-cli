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

	It("can list mailboxes", func() {
		if skipMailboxList && !isCI {
			Skip("can list mailboxes")
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
		expected := loadGolden("can_list_mailboxes")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can list mailboxes and filter result", func() {
		if skipMailboxList && !isCI {
			Skip("can list mailboxes and filter result")
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

		cli := exec.Command(cliPath, "mailbox", "list", "-f", "domain.com")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_mailboxes_and_filter_result")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
