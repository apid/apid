package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/apid/util"
	"time"
)


var _ = Describe("timeout", func() {

	retry := 50 * time.Millisecond
	timeout := 200 * time.Millisecond
	approximately := 10 * time.Millisecond

	It("should finish after retries without error", func(done Done) {

		count := 0
		action := func() (ok bool, err error) {
			count++
			ok = count > 2
			return
		}

		start := time.Now()
		ok, err := util.RetryEveryUntil(action, retry, timeout)
		Expect(ok).To(BeTrue())
		Expect(err).ShouldNot(HaveOccurred())

		Expect(count).To(Equal(3))
		finish := time.Now()
		Expect(finish).Should(BeTemporally("~", start.Add(time.Duration(count) * retry), approximately))

		close(done)
	})

	It("should timeout with error", func(done Done) {

		action := func() (ok bool, err error) {
			time.Sleep(500 * time.Millisecond)
			return
		}

		start := time.Now()
		ok, err := util.RetryEveryUntil(action, retry, timeout)
		Expect(ok).To(BeFalse())
		Expect(err).Should(HaveOccurred())

		finish := time.Now()
		Expect(finish).Should(BeTemporally("~", start.Add(timeout), approximately))

		close(done)
	})
})
