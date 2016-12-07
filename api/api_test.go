package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var _ = Describe("API Service", func() {

	It("should return vars from /exp/vars with request counter", func() {

		uri, err := url.Parse(testServer.URL)
		Expect(err).NotTo(HaveOccurred())
		uri.Path = "/exp/vars"

		resp, err := http.Get(uri.String())
		defer resp.Body.Close()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(resp.Body)
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		Expect(err).ShouldNot(HaveOccurred())

		requests := m["requests"].(map[string]interface{})
		Expect(requests["/exp/vars"]).Should(Equal(float64(1)))
	})
})
