package robots_test

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/robots"
	"net/http"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRobots(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functions Suite")
}

var _ = Describe("Robots", func() {
	var r *http.Request
	BeforeEach(func() {
		r = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
		}
	})
	Context("RobotsTxt", func() {
		It("should return a robots.txt file", func() {
			Expect(robots.RobotsTxt(r)).NotTo(BeEmpty())
		})

		It("should not return an empty robots.txt file", func() {
			Expect(robots.RobotsTxt(r)).NotTo(Equal([]byte("")))
		})

		It("should not return the same robots.txt file", func() {
			Expect(robots.RobotsTxt(r)).NotTo(Equal(robots.RobotsTxt(r)))
		})
	})
})
