package events_test

import (
	"github.com/30x/apid"
	"github.com/30x/apid/events"
	"github.com/30x/apid/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
)

var _ = Describe("Events Service", func() {

	BeforeSuite(func() {
		apid.Initialize(factory.DefaultServicesFactory())
	})

	It("should ignore event with no listeners", func() {
		em := events.CreateService()
		defer em.Close()
		em.Emit("no listeners", &test_event{"test"})
	})

	It("should publish an event to a listener", func(done Done) {
		em := events.CreateService()
		defer em.Close()

		h := test_handler{
			"handler",
			func(event apid.Event) {
				close(done)
			},
		}

		em.Listen("selector", &h)
		em.Emit("selector", &test_event{"test"})
	})

	It("should publish an event multiple listeners", func(done Done) {
		em := events.CreateService()
		defer em.Close()

		hitH1 := false
		hitH2 := false
		h1 := test_handler{
			"handler 1",
			func(event apid.Event) {
				hitH1 = true
				if hitH1 && hitH2 {
					close(done)
				}
			},
		}
		h2 := test_handler{
			"handler 2",
			func(event apid.Event) {
				hitH2 = true
				if hitH1 && hitH2 {
					close(done)
				}
			},
		}

		em.Listen("selector", &h1)
		em.Listen("selector", &h2)
		em.Emit("selector", &test_event{"test"})
	})

	It("should deliver events according selector", func(done Done) {
		em := events.CreateService()
		defer em.Close()

		e1 := &test_event{"test1"}
		e2 := &test_event{"test2"}

		count := int32(0)

		h1 := test_handler{
			"handler1",
			func(event apid.Event) {
				c := atomic.AddInt32(&count, 1)
				Expect(event).To(BeIdenticalTo(e1))
				if c == 2 {
					close(done)
				}
			},
		}

		h2 := test_handler{
			"handler2",
			func(event apid.Event) {
				c := atomic.AddInt32(&count, 1)
				Expect(event).To(BeIdenticalTo(e2))
				if c == 2 {
					close(done)
				}
			},
		}

		em.Listen("selector1", &h1)
		em.Listen("selector2", &h2)

		em.Emit("selector2", e2)
		em.Emit("selector1", e1)
	})
})

type test_handler struct {
	description string
	f           func(event apid.Event)
}

func (t *test_handler) String() string {
	return t.description
}

func (t *test_handler) Handle(event apid.Event) {
	t.f(event)
}

type test_event struct {
	description string
}
