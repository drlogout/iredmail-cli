package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mailbox info", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can show mailbox info", func() {
		if skipMailboxInfo && !isCI {
			Skip("can show mailbox info")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias1, mailboxName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add-alias", mailboxAlias2, mailboxName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress2)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "info", mailboxName1, "--pretty=false")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_show_mailbox_info")

		// test fails because date is always different, hence set the same date
		date := "2018.07.05.06.12.36"
		re := regexp.MustCompile(`\d\d\d\d\.\d\d\.\d\d\.\d\d\.\d\d\.\d\d`)
		actual = re.ReplaceAllString(actual, date)
		expected = re.ReplaceAllString(expected, date)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
