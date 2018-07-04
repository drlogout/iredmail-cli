package integrationTest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain add/delete-alias", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	// It("can add an domain", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can add an domain")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "add", domain1)
	// 	output, err := cli.CombinedOutput()
	// 	if err != nil {
	// 		Fail(string(output))
	// 	}

	// 	actual := string(output)
	// 	expected := fmt.Sprintf("Successfully added domain %s\n", domain1)

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}

	// 	db, err := sql.Open("mysql", dbConnectionString)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	defer db.Close()

	// 	var exists bool

	// 	query := `SELECT exists
	// 	(SELECT domain FROM domain
	// 	WHERE domain = '` + domain1 + `');`

	// 	err = db.QueryRow(query).Scan(&exists)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(exists).To(Equal(true))
	// })

	// It("can't add an existing domain", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can't add an existing domain")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "add", domain1)
	// 	err := cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "domain", "add", domain1)
	// 	output, err := cli.CombinedOutput()
	// 	Expect(err).To(HaveOccurred())

	// 	actual := string(output)
	// 	expected := fmt.Sprintf("Domain %s already exists\n", domain1)

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}
	// })

	// It("can create a domain automatically while adding a mailbox", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can create a domain automatically while adding a mailbox")
	// 	}

	// 	cli := exec.Command(cliPath, "mailbox", "add", mailboxName3, mailboxPW)
	// 	err := cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	db, err := sql.Open("mysql", dbConnectionString)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	defer db.Close()

	// 	var exists bool

	// 	query := `SELECT exists
	// 	(SELECT domain FROM domain
	// 	WHERE domain = '` + domain1 + `');`

	// 	err = db.QueryRow(query).Scan(&exists)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(exists).To(Equal(true))
	// })

	// It("can delete an domain", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can delete an domain")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "add", domain1)
	// 	err := cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "domain", "delete", domain1)
	// 	output, err := cli.CombinedOutput()
	// 	if err != nil {
	// 		Fail(string(output))
	// 	}

	// 	actual := string(output)
	// 	expected := fmt.Sprintf("Successfully deleted domain %s\n", domain1)

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}

	// 	db, err := sql.Open("mysql", dbConnectionString)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	defer db.Close()

	// 	var exists bool

	// 	query := `SELECT exists
	// 	(SELECT domain FROM domain
	// 	WHERE domain = '` + domain1 + `');`

	// 	err = db.QueryRow(query).Scan(&exists)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(exists).To(Equal(false))
	// })

	// It("can't delete an domain with mailboxes", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can delete an domain with mailboxes")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "add", domain1)
	// 	err := cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "mailbox", "add", mailboxName3, mailboxPW)
	// 	err = cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "domain", "delete", domain1)
	// 	output, err := cli.CombinedOutput()
	// 	if err == nil {
	// 		Fail("It should exit because mailbox exists")
	// 	}

	// 	actual := string(output)
	// 	expected := fmt.Sprintf("The domain %s still has mailboxes you need to delete them before\n", domain1)

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}

	// 	db, err := sql.Open("mysql", dbConnectionString)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	defer db.Close()

	// 	var exists bool

	// 	query := `SELECT exists
	// 	(SELECT domain FROM domain
	// 	WHERE domain = '` + domain1 + `');`

	// 	err = db.QueryRow(query).Scan(&exists)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(exists).To(Equal(true))
	// })

	// It("can't delete an non existing domain", func() {
	// 	if skipDomainAddDelete && !isCI {
	// 		Skip("can't delete an non existing domain")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "delete", domain1)
	// 	output, err := cli.CombinedOutput()
	// 	Expect(err).To(HaveOccurred())

	// 	actual := string(output)
	// 	expected := fmt.Sprintf("Domain %s doesn't exist\n", domain1)

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}
	// })
})
