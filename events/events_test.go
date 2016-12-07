package events_test

import (
	"github.com/30x/apid"
	"github.com/30x/apid/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
)

var _ = Describe("Events Service", func() {

	It("should ignore event with no listeners", func() {
		em := events.CreateService()
		defer em.Close()
		em.Emit("no listeners", &test_event{"test"})
	})

	It("should publish an event to a listener", func(done Done) {
		em := events.CreateService()

		h := test_handler{
			"handler",
			func(event apid.Event) {
				defer GinkgoRecover()

				em.Close()
				close(done)
			},
		}

		em.Listen("selector", &h)
		em.Emit("selector", &test_event{"test"})
	})

	It("should publish an event to a listener func", func(done Done) {
		em := events.CreateService()

		h := func(event apid.Event) {
			defer GinkgoRecover()

			em.Close()
			close(done)
		}

		em.ListenFunc("selector", h)
		em.Emit("selector", &test_event{"test"})
	})

	It("should publish multiple events to a listener", func(done Done) {
		em := events.CreateService()

		count := int32(0)
		h := test_handler{
			"handler",
			func(event apid.Event) {
				defer GinkgoRecover()

				c := atomic.AddInt32(&count, 1)
				if c > 1 {
					em.Close()
					close(done)
				}
			},
		}

		em.Listen("selector", &h)
		em.Emit("selector", &test_event{"test1"})
		em.Emit("selector", &test_event{"test2"})
	})

	It("EmitWithCallback should call the callback when done with delivery", func(done Done) {
		em := events.CreateService()

		delivered := func(event apid.Event) {
			defer GinkgoRecover()
			close(done)
		}

		em.EmitWithCallback("selector", &test_event{"test1"}, delivered)
	})

	It("should publish only one event to a listenOnce", func(done Done) {
		em := events.CreateService()

		count := 0
		h := func(event apid.Event) {
			defer GinkgoRecover()
			count++
		}

		delivered := func(event apid.Event) {
			defer GinkgoRecover()
			Expect(count).To(Equal(1))
			em.Close()
			close(done)
		}

		em.ListenOnceFunc("selector", h)
		em.Emit("selector", &test_event{"test1"})
		em.EmitWithCallback("selector", &test_event{"test2"}, delivered)
	})

	It("should publish an event to multiple listeners", func(done Done) {
		em := events.CreateService()
		defer em.Close()

		hitH1 := false
		hitH2 := false
		h1 := test_handler{
			"handler 1",
			func(event apid.Event) {
				defer GinkgoRecover()

				hitH1 = true
				if hitH1 && hitH2 {
					em.Close()
					close(done)
				}
			},
		}
		h2 := test_handler{
			"handler 2",
			func(event apid.Event) {
				defer GinkgoRecover()

				hitH2 = true
				if hitH1 && hitH2 {
					em.Close()
					close(done)
				}
			},
		}

		em.Listen("selector", &h1)
		em.Listen("selector", &h2)
		em.Emit("selector", &test_event{"test"})
	})

	It("should publish an event delivered event", func(done Done) {
		em := events.CreateService()
		testEvent := &test_event{"test"}
		var testSelector apid.EventSelector = "selector"

		dummy := func(event apid.Event) {}
		em.ListenFunc(testSelector, dummy)

		h := test_handler{
			"event delivered handler",
			func(event apid.Event) {
				defer GinkgoRecover()

				e, ok := event.(apid.EventDeliveryEvent)

				Expect(ok).To(BeTrue())
				Expect(e.Event).To(Equal(testEvent))
				Expect(e.Selector).To(Equal(testSelector))

				em.Close()
				close(done)
			},
		}

		em.Listen(apid.EventDeliveredSelector, &h)
		em.Emit(testSelector, testEvent)
	})

	It("should be able to remove a listener", func(done Done) {
		em := events.CreateService()

		event1 := &test_event{"test1"}
		event2 := &test_event{"test2"}
		event3 := &test_event{"test3"}

		dummy := func(event apid.Event) {}
		em.ListenFunc("selector", dummy)

		h := test_handler{
			"handler",
			func(event apid.Event) {
				defer GinkgoRecover()

				Expect(event).NotTo(Equal(event2))
				if event == event3 {
					em.Close()
					close(done)
				}
			},
		}
		em.Listen("selector", &h)

		// need to drive test like this because of async delivery
		td := test_handler{
			"test driver",
			func(event apid.Event) {
				defer GinkgoRecover()

				e := event.(apid.EventDeliveryEvent)
				if e.Event == event1 {
					em.StopListening("selector", &h)
					em.Emit("selector", event2)
				} else if e.Event == event2 {
					em.Listen("selector", &h)
					em.Emit("selector", event3)
				}
			},
		}
		em.Listen(apid.EventDeliveredSelector, &td)

		em.Emit("selector", event1)
	})

	It("should deliver events according selector", func(done Done) {
		em := events.CreateService()

		e1 := &test_event{"test1"}
		e2 := &test_event{"test2"}

		count := int32(0)

		h1 := test_handler{
			"handler1",
			func(event apid.Event) {
				defer GinkgoRecover()

				c := atomic.AddInt32(&count, 1)
				Expect(event).Should(Equal(e1))
				if c == 2 {
					em.Close()
					close(done)
				}
			},
		}

		h2 := test_handler{
			"handler2",
			func(event apid.Event) {
				defer GinkgoRecover()

				c := atomic.AddInt32(&count, 1)
				Expect(event).Should(Equal(e2))
				if c == 2 {
					em.Close()
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
