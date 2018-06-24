package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	userName1         = "post@web.de"
	userName2         = "info@domain.com"
	userName3         = "webmaster@example.com"
	userName4         = "abuse@domain.com"
	userName5         = "support@wurst.de"
	userPW            = "alskdlqkdjalskd"
	forwardingAddress = "info@example.com"
	customQuota       = 4096
	customStoragePath = "/var/mail/custom"
)

var _ = Describe("user", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an user", func() {
		if skipUser && !isCI {
			Skip("can add an user")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: 2048 KB)\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + userName1 + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete an user", func() {
		if skipUser && !isCI {
			Skip("can delete an user")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "delete", "--force", userName1)
		output, err = cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted user %v\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + userName1 + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't add an existing user", func() {
		if skipUser && !isCI {
			Skip("can't add an existing user")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		cli = exec.Command(cliPath, "user", "add", userName1, userPW)
		output, err = cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("User %v already exists\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can add an user with custom quota", func() {
		if skipUser && !isCI {
			Skip("can add an user with custom quota")
		}

		cli := exec.Command(cliPath, "user", "add", "--quota", strconv.Itoa(customQuota), userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: "+strconv.Itoa(customQuota)+" KB)\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + userName1 + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		var quota int

		query = `SELECT quota FROM mailbox WHERE username = '` + userName1 + `'`
		err = db.QueryRow(query).Scan(&quota)
		Expect(err).NotTo(HaveOccurred())
		Expect(quota).To(Equal(customQuota))
	})

	It("can add an user with custom storage path", func() {
		if skipUser && !isCI {
			Skip("can add an user with custom storage path")
		}

		cli := exec.Command(cliPath, "user", "add", "--storage-path", customStoragePath, userName1, userPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: 2048 KB)\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName1 + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName1 + `' AND forwarding = '` + userName1 + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		var storageBaseDirectory, storageNode string

		query = `SELECT storagebasedirectory, storagenode FROM mailbox WHERE username = '` + userName1 + `'`
		err = db.QueryRow(query).Scan(&storageBaseDirectory, &storageNode)
		Expect(err).NotTo(HaveOccurred())
		Expect(storageBaseDirectory).To(Equal(filepath.Dir(customStoragePath)))
		Expect(storageNode).To(Equal(filepath.Base(customStoragePath)))
	})
})
