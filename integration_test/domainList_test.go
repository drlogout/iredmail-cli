package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain list", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can list domains", func() {
		if skipDomainList && !isCI {
			Skip("can list domains")
		}

		cli := exec.Command(cliPath, "domain", "add", domain1)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add", domain2, "-s", "default_user_quota:4096")
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add", domain3)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain1, domain3)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "add-alias", aliasDomain2, domain3)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "domain", "list")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_domains")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	// It("can list domains and filter result", func() {
	// 	if skipDomainList && !isCI {
	// 		Skip("can list domains and filter result")
	// 	}

	// 	cli := exec.Command(cliPath, "domain", "add", domain1)
	// 	err := cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "domain", "add", domain2, "-s", "default_user_quota:4096")
	// 	err = cli.Run()
	// 	Expect(err).NotTo(HaveOccurred())

	// 	cli = exec.Command(cliPath, "domain", "add", domain3)
	// 	err = cli.Run()

	// 	cli = exec.Command(cliPath, "domain", "list", "-f", "4096")
	// 	output, err := cli.CombinedOutput()
	// 	if err != nil {
	// 		Fail(string(output))
	// 	}

	// 	actual := string(output)
	// 	expected := loadGolden("can_list_domains_and_filter_result")

	// 	if !reflect.DeepEqual(actual, expected) {
	// 		Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
	// 	}
	// })
})
