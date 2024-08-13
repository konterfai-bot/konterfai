package statistics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
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
)

// recordStatistics records the statistics.
func (s *Statistics) recordStatistics() {
	go func() {
		for {
			select {
			case <-s.Context.Done():
				return
			case <-time.After(5 * time.Second):
				RobotsTxtViolatorsTotal.Set(float64(s.GetTotalRobotsTxtViolators()))
			}
		}
	}()
}
