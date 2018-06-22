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
	userName          = "post@domain.com"
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
		cli := exec.Command(cliPath, "user", "add", userName, userPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: 2048 KB)\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + userName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete an user", func() {
		cli := exec.Command(cliPath, "user", "add", userName, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "delete", "--force", userName)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted user %v\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + userName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't add an existing user", func() {
		cli := exec.Command(cliPath, "user", "add", userName, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add", userName, userPW)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("User %v already exists\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can add an user with custom quota", func() {
		cli := exec.Command(cliPath, "user", "add", "--quota", strconv.Itoa(customQuota), userName, userPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: "+strconv.Itoa(customQuota)+" KB)\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + userName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		var quota int

		query = `SELECT quota FROM mailbox WHERE username = '` + userName + `'`
		err = db.QueryRow(query).Scan(&quota)
		Expect(err).NotTo(HaveOccurred())
		Expect(quota).To(Equal(customQuota))
	})

	It("can add an user with custom storage path", func() {
		cli := exec.Command(cliPath, "user", "add", "--storage-path", customStoragePath, userName, userPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user %v (quota: 2048 KB)\n", userName)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		query := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = '` + userName + `');`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		query = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + userName + `' AND forwarding = '` + userName + `' 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		var storageBaseDirectory, storageNode string

		query = `SELECT storagebasedirectory, storagenode FROM mailbox WHERE username = '` + userName + `'`
		err = db.QueryRow(query).Scan(&storageBaseDirectory, &storageNode)
		Expect(err).NotTo(HaveOccurred())
		Expect(storageBaseDirectory).To(Equal(filepath.Dir(customStoragePath)))
		Expect(storageNode).To(Equal(filepath.Base(customStoragePath)))
	})
})
