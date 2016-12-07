package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	"io/ioutil"
	"os"
)

var tmpDir string

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())

	var err error
	config := apid.Config()
	tmpDir, err = ioutil.TempDir("", "apid_test")
	Expect(err).NotTo(HaveOccurred())
	config.Set("local_storage_path", tmpDir)
})

var _ = AfterSuite(func() {
	os.RemoveAll(tmpDir)
})

func TestEvents(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}
