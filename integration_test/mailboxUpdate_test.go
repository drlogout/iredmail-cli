package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mailbox update", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can disable keep-copy", func() {
		if skipMailboxUpdate && !isCI {
			Skip("can disable keep-copy")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-k", "no", "--pretty=false")
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_disable_keep-copy")

		// test fails because date is always different, hence set the same date
		date := "2018.07.05.06.12.36"
		re := regexp.MustCompile(`\d\d\d\d\.\d\d\.\d\d\.\d\d\.\d\d\.\d\d`)
		actual = re.ReplaceAllString(actual, date)
		expected = re.ReplaceAllString(expected, date)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0 AND active = 1 );`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can enable keep-copy", func() {
		if skipMailboxUpdate && !isCI {
			Skip("can enable keep-copy")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-k", "no")
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-k", "yes", "--pretty=false")
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_enable_keep-copy")

		// test fails because date is always different, hence set the same date
		date := "2018.07.05.06.12.36"
		re := regexp.MustCompile(`\d\d\d\d\.\d\d\.\d\d\.\d\d\.\d\d\.\d\d`)
		actual = re.ReplaceAllString(actual, date)
		expected = re.ReplaceAllString(expected, date)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0 AND active = 1 );`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't disable keep-copy if no forwarding exists", func() {
		if skipMailboxUpdate && !isCI {
			Skip("can't disable keep-copy if no forwarding exists")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-k", "no", "--pretty=false")
		output, err = cli.CombinedOutput()
		if err == nil {
			Fail("Expect an error")
		}

		actual := string(output)
		expected := fmt.Sprintf("No forwardings exist for mailbox %s\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0 AND active = 1 );`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can enable keep-copy automatically if all forwardings are deleted", func() {
		if skipMailboxUpdate && !isCI {
			Skip("can enable keep-copy automatically if all forwardings are deleted")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "add", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-k", "no")
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "forwarding", "delete", mailboxName1, forwardingAddress1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT address FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND is_alias = 0 AND is_list = 0 AND active = 1 );`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can set quota", func() {
		if skipMailboxUpdate && !isCI {
			Skip("can set quota")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var quota int

		sqlQuery := `SELECT quota FROM mailbox WHERE username = ?;`
		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&quota)
		Expect(err).NotTo(HaveOccurred())
		Expect(quota).To(Equal(2048))

		cli = exec.Command(cliPath, "mailbox", "update", mailboxName1, "-q", strconv.Itoa(customQuota), "--pretty=false")
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_set_quota")

		// test fails because date is always different, hence set the same date
		date := "2018.07.05.06.12.36"
		re := regexp.MustCompile(`\d\d\d\d\.\d\d\.\d\d\.\d\d\.\d\d\.\d\d`)
		actual = re.ReplaceAllString(actual, date)
		expected = re.ReplaceAllString(expected, date)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		sqlQuery = `SELECT quota FROM mailbox WHERE username = ?;`
		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&quota)
		Expect(err).NotTo(HaveOccurred())
		Expect(quota).To(Equal(customQuota))
	})
})
