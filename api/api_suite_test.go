package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

import (
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	"net/http/httptest"
	"os"
)

var (
	testDir    string
	testServer *httptest.Server
)

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())

	apid.Config().Set("api_expvar_path", "/exp/vars")

	// get the router - this will have the /exp/vars route registered
	router := apid.API().Router()

	// create our test server
	testServer = httptest.NewServer(router)

})

var _ = AfterSuite(func() {
	apid.Events().Close()
	if testServer != nil {
		testServer.Close()
	}
	os.RemoveAll(testDir)
})

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
