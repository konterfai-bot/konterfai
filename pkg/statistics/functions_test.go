package statistics_test

import (
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Functions", func() {
	var ctx context.Context
	var s *statistics.Statistics
	var r statistics.Request

	BeforeEach(func() {
		ctx = context.Background()
		s = &statistics.Statistics{}
		r = statistics.Request{
			Context:   ctx,
			IpAddress: "127.0.0.1",
			UserAgent: "Mozilla/5.0",
			Size:      1024,
		}
	})

	Context("AppendRequest", func() {
		It("should append a request to the statistics", func() {
			s.AppendRequest(ctx, r)
			cmp := s.GetRequests(ctx)[0]
			Expect(cmp.UserAgent).To(Equal(r.UserAgent))
			Expect(cmp.IpAddress).To(Equal(r.IpAddress))
			Expect(cmp.Size).To(Equal(r.Size))
			Expect(cmp.Context).To(Equal(r.Context))
		})

		It("should contain two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			Expect(len(s.GetRequests(ctx))).To(Equal(2))
		})
	})

	Context("GetAgents", func() {
		It("should return one agent", func() {
			s.AppendRequest(ctx, r)
			agents := s.GetAgents(ctx)
			Expect(agents).To(ContainElement(r.UserAgent))
		})

		It("should return one agent when two requests with the same agent are appended", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			agents := s.GetAgents(ctx)
			Expect(len(agents)).To(Equal(1))
		})

		It("should return two agents", func() {
			s.AppendRequest(ctx, r)
			r.UserAgent = "Mozilla/4.0"
			s.AppendRequest(ctx, r)
			agents := s.GetAgents(ctx)
			Expect(len(agents)).To(Equal(2))
		})
	})

	Context("GetIpAddresses", func() {
		It("should return one IP address", func() {
			s.AppendRequest(ctx, r)
			ips := s.GetIpAddresses(ctx)
			Expect(ips).To(ContainElement(r.IpAddress))
		})

		It("should return one IP address when two requests with the same IP address are appended", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			ips := s.GetIpAddresses(ctx)
			Expect(len(ips)).To(Equal(1))
		})

		It("should return two IP addresses", func() {
			s.AppendRequest(ctx, r)
			r.IpAddress = "172.0.0.1"
			s.AppendRequest(ctx, r)
			ips := s.GetIpAddresses(ctx)
			Expect(len(ips)).To(Equal(2))
		})
	})

	Context("GetRequests", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequests(ctx)
			Expect(len(requests)).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequests(ctx)
			Expect(len(requests)).To(Equal(2))
		})
	})

	Context("GetRequestsByIpAddress", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByIpAddress(ctx, r.IpAddress)
			Expect(len(requests)).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByIpAddress(ctx, r.IpAddress)
			Expect(len(requests)).To(Equal(2))
		})

		It("should return zero requests", func() {
			requests := s.GetRequestsByIpAddress(ctx, r.IpAddress)
			Expect(len(requests)).To(Equal(0))
		})

		It("should return two requests with the same IP address", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByIpAddress(ctx, r.IpAddress)
			Expect(len(requests)).To(Equal(2))
		})
	})

	Context("GetRequestsByTimeRange", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByTimeRange(ctx, r.Timestamp.Add(-10*time.Second), r.Timestamp.Add(10*time.Second))
			Expect(len(requests)).To(Equal(1))
			requests = s.GetRequestsByTimeRange(ctx, r.Timestamp, r.Timestamp)
			Expect(len(requests)).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByTimeRange(ctx, r.Timestamp, r.Timestamp)
			Expect(len(requests)).To(Equal(2))
		})

		It("should return zero requests", func() {
			requests := s.GetRequestsByTimeRange(ctx, r.Timestamp, r.Timestamp)
			Expect(len(requests)).To(Equal(0))
		})

		It("should return two requests with the same time range", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByTimeRange(ctx, r.Timestamp, r.Timestamp)
			Expect(len(requests)).To(Equal(2))
		})
	})

	Context("GetRequestsByUserAgent", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByUserAgent(ctx, r.UserAgent)
			Expect(len(requests)).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByUserAgent(ctx, r.UserAgent)
			Expect(len(requests)).To(Equal(2))
		})

		It("should return zero requests", func() {
			requests := s.GetRequestsByUserAgent(ctx, r.UserAgent)
			Expect(len(requests)).To(Equal(0))
		})

		It("should return two requests with the same user agent", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsByUserAgent(ctx, r.UserAgent)
			Expect(len(requests)).To(Equal(2))
		})
	})

	Context("GetRequestsGroupedByIpAddress", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByIpAddress(ctx)
			Expect(len(requests[r.IpAddress])).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByIpAddress(ctx)
			Expect(len(requests[r.IpAddress])).To(Equal(2))
		})

		It("should return zero requests", func() {
			requests := s.GetRequestsGroupedByIpAddress(ctx)
			Expect(len(requests)).To(Equal(0))
		})

		It("should return two requests with the same IP address", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByIpAddress(ctx)
			Expect(len(requests[r.IpAddress])).To(Equal(2))
		})
	})

	Context("GetRequestsGroupedByUserAgent", func() {
		It("should return one request", func() {
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByUserAgent(ctx)
			Expect(len(requests[r.UserAgent])).To(Equal(1))
		})

		It("should return two requests", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByUserAgent(ctx)
			Expect(len(requests[r.UserAgent])).To(Equal(2))
		})

		It("should return zero requests", func() {
			requests := s.GetRequestsGroupedByUserAgent(ctx)
			Expect(len(requests)).To(Equal(0))
		})

		It("should return two requests with the same user agent", func() {
			s.AppendRequest(ctx, r)
			s.AppendRequest(ctx, r)
			requests := s.GetRequestsGroupedByUserAgent(ctx)
			Expect(len(requests[r.UserAgent])).To(Equal(2))
		})
	})

	Context("GetTotalDataSizeServed", func() {
		It("should return the total data served size", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServed(ctx)).To(Equal(r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServed(ctx)).To(Equal(2 * r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServed(ctx)).To(Equal(3 * r.Size))
		})
	})

	Context("GetTotalDataSizeServedByAgent", func() {
		It("should return the total data size served by agent", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByAgent(ctx, r.UserAgent)).To(Equal(r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByAgent(ctx, r.UserAgent)).To(Equal(2 * r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByAgent(ctx, r.UserAgent)).To(Equal(3 * r.Size))
			r.UserAgent = "Mozilla/4.0"
			s.AppendRequest(ctx, r)
			r.UserAgent = "Mozilla/5.0"
			Expect(s.GetTotalDataSizeServedByAgent(ctx, r.UserAgent)).To(Equal(3 * r.Size))
			Expect(len(s.GetRequests(ctx))).To(Equal(4))
		})
	})

	Context("GetTotalDataSizeServedByIpAddress", func() {
		It("should return the total data size served by IP address", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByIpAddress(ctx, r.IpAddress)).To(Equal(r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByIpAddress(ctx, r.IpAddress)).To(Equal(2 * r.Size))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalDataSizeServedByIpAddress(ctx, r.IpAddress)).To(Equal(3 * r.Size))
			r.IpAddress = "172.0.0.1"
			s.AppendRequest(ctx, r)
			r.IpAddress = "127.0.0.1"
			Expect(s.GetTotalDataSizeServedByIpAddress(ctx, r.IpAddress)).To(Equal(3 * r.Size))
		})

		Context("GetTotalDataSizeServedByTimeRange", func() {
			It("should return the total data size served by time range", func() {
				s.AppendRequest(ctx, r)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp.Add(-100*time.Second), r.Timestamp.Add(100*time.Second))).To(Equal(r.Size))
				r.Timestamp = r.Timestamp.Add(10 * time.Second)
				s.AppendRequest(ctx, r)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp.Add(-100*time.Second), r.Timestamp.Add(100*time.Second))).To(Equal(2 * r.Size))
				r.Timestamp = r.Timestamp.Add(10 * time.Second)
				s.AppendRequest(ctx, r)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp.Add(-100*time.Second), r.Timestamp.Add(100*time.Second))).To(Equal(3 * r.Size))
				r.Timestamp = r.Timestamp.Add(1000 * time.Second)
				s.AppendRequest(ctx, r)
				r.Timestamp = r.Timestamp.Add(-1000 * time.Second)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp.Add(-100*time.Second), r.Timestamp.Add(100*time.Second))).To(Equal(3 * r.Size))
			})

			It("should return the total data size served by time range when start and end are equal", func() {
				s.AppendRequest(ctx, r)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp, r.Timestamp)).To(Equal(r.Size))
				r.Timestamp = r.Timestamp.Add(10 * time.Second)
				s.AppendRequest(ctx, r)
				Expect(s.GetTotalDataSizeServedByTimeRange(ctx, r.Timestamp, r.Timestamp)).To(Equal(r.Size))
			})
		})
	})

	Context("GetTotalRequestsByAgent", func() {
		It("should return the total requests by agent", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByAgent(ctx, r.UserAgent)).To(Equal(1))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByAgent(ctx, r.UserAgent)).To(Equal(2))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByAgent(ctx, r.UserAgent)).To(Equal(3))
			r.UserAgent = "Mozilla/4.0"
			s.AppendRequest(ctx, r)
			r.UserAgent = "Mozilla/5.0"
			Expect(s.GetTotalRequestsByAgent(ctx, r.UserAgent)).To(Equal(3))
		})
	})

	Context("GetTotalRequestsByIpAddress", func() {
		It("should return the total requests by IP address", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByIpAddress(ctx, r.IpAddress)).To(Equal(1))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByIpAddress(ctx, r.IpAddress)).To(Equal(2))
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByIpAddress(ctx, r.IpAddress)).To(Equal(3))
			r.IpAddress = "172.0.0.1"
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequestsByIpAddress(ctx, r.IpAddress)).To(Equal(1))
			r.IpAddress = "127.0.0.1"
			Expect(s.GetTotalRequestsByIpAddress(ctx, r.IpAddress)).To(Equal(3))
		})
	})

	Context("GetTotalRequests", func() {
		It("should return the total requests", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequests(ctx)).To(Equal(1))
			r.UserAgent = "Mozilla/4.0"
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequests(ctx)).To(Equal(2))
			r.UserAgent = "Mozilla/99.0"
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRequests(ctx)).To(Equal(3))
		})
	})

	Context("GetTotalRobotsTxtViolators", func() {
		It("should return the total robots.txt violators", func() {
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRobotsTxtViolators(ctx)).To(Equal(0))
			r.IsRobotsTxt = true
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRobotsTxtViolators(ctx)).To(Equal(1))
			r.UserAgent = "Brave AI which respects robots"
			r.IsRobotsTxt = false
			s.AppendRequest(ctx, r)
			Expect(s.GetTotalRobotsTxtViolators(ctx)).To(Equal(1))
		})
	})

	Context("UpdatePrompts", func() {
		It("should update the prompts", func() {
			prompts := map[string]int{
				"test":   1,
				"foobar": 2,
				"barfoo": 3,
			}
			s.UpdatePrompts(ctx, prompts)
			Expect(len(s.Prompts)).To(Equal(3))
			Expect(s.Prompts["test"]).To(Equal(1))
			Expect(s.Prompts["foobar"]).To(Equal(2))
			Expect(s.Prompts["barfoo"]).To(Equal(3))

			prompts["baz"] = 4
			s.UpdatePrompts(ctx, prompts)
			Expect(len(s.Prompts)).To(Equal(4))
			Expect(s.Prompts["test"]).To(Equal(1))
			Expect(s.Prompts["foobar"]).To(Equal(2))
			Expect(s.Prompts["barfoo"]).To(Equal(3))
			Expect(s.Prompts["baz"]).To(Equal(4))
		})
	})
})
