package integrationTest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	cli = ""
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	err := setupDB()
	Expect(err).NotTo(HaveOccurred())

	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	cli = filepath.Join(cwd, "../", "iredmail-cli")

	cmd := exec.Command("go", "build", "-o", cli)

	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
})
