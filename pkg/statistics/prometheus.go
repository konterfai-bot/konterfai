package statistics

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestsTotal is the total number of requests.
	RequestTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konterfai_requests_total",
		Help: "The total number of requests processed.",
	})

	// PromptsGeneratedTotal is the total number of prompts generated.
	PromptsGeneratedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konterfai_prompts_generated_total",
		Help: "The total number of prompts generated.",
	})

	// DataFedTotal is the total amount of data fed.
	DataFedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konterfai_data_fed_total_bytes",
		Help: "The total amount of data fed in bytes.",
	})

	// RobotsTxtViolatorsTotal is the total number of violators of robots.txt.
	RobotsTxtViolatorsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "konterfai_robots_txt_violators_total",
		Help: "The total number of violators of robots.txt.",
	})

	// AgentTraffic is the traffic per user agent.
	AgentTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "konterfai_agent_traffic_bytes",
		Help: "The traffic per user agent in bytes.",
	}, []string{"user_agent"})

	// AgentRequests is the requests per user agent.
	AgentRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "konterfai_agent_requests_total",
		Help: "The requests per user agent.",
	}, []string{"user_agent"})

	// IpTraffic is the traffic per IP address.
	IpTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "konterfai_ip_traffic_bytes",
		Help: "The traffic per IP address in bytes.",
	}, []string{"ip_address"})

	// IpRequests is the requests per IP address.
	IpRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "konterfai_ip_requests_total",
		Help: "The requests per IP address.",
	}, []string{"ip_address"})
)

// recordStatistics records the statistics.
func (s *Statistics) recordStatistics(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Statistics.recordStatistics")
	defer span.End()

	isProcessing := sync.Mutex{}
	go func() {
		for {
			select {
			case <-s.Context.Done():
				return
			case <-time.After(5 * time.Second):
				if isProcessing.TryLock() {
					RobotsTxtViolatorsTotal.Set(float64(s.GetTotalRobotsTxtViolators(ctx)))

					for _, agent := range s.GetAgents(ctx) {
						AgentTraffic.WithLabelValues(agent).Set(float64(s.GetTotalDataSizeServedByAgent(ctx, agent)))
						AgentRequests.WithLabelValues(agent).Set(float64(s.GetTotalRequestsByAgent(ctx, agent)))
					}

					for _, ip := range s.GetIpAddresses(ctx) {
						IpTraffic.WithLabelValues(ip).Set(float64(s.GetTotalDataSizeServedByIpAddress(ctx, ip)))
						IpRequests.WithLabelValues(ip).Set(float64(s.GetTotalRequestsByIpAddress(ctx, ip)))
					}

					isProcessing.Unlock()
				}
			}
		}
	}()
}
