package events_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
)

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())
})


func TestEvents(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Events Suite")
}
