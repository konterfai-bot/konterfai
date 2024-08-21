package statistics

import (
	"context"
	"strings"
	"time"
)

// AppendRequest appends a request to the statistics.
func (s *Statistics) AppendRequest(ctx context.Context, r Request) {
	_, span := tracer.Start(ctx, "Statistics.AppendRequest")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	r.IPAddress = strings.Split(r.IPAddress, ":")[0]
	s.Requests = append(s.Requests, r)

	// Update Prometheus metrics
	RequestTotal.Inc()
	DataFedTotal.Add(float64(r.Size))
}

// GetAgents returns the agents.
func (s *Statistics) GetAgents(ctx context.Context) []string {
	_, span := tracer.Start(ctx, "Statistics.GetAgents")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	agents := map[string]struct{}{}
	for _, r := range s.Requests {
		agents[r.UserAgent] = struct{}{}
	}
	agentsList := make([]string, 0, len(agents))
	for agent := range agents {
		agentsList = append(agentsList, agent)
	}

	return agentsList
}

// GetIPAddresses returns the IP addresses.
func (s *Statistics) GetIPAddresses(ctx context.Context) []string {
	_, span := tracer.Start(ctx, "Statistics.GetIPAddresses")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	ips := map[string]struct{}{}
	for _, r := range s.Requests {
		ips[r.IPAddress] = struct{}{}
	}
	ipsList := make([]string, 0, len(ips))
	for ip := range ips {
		ipsList = append(ipsList, ip)
	}

	return ipsList
}

// GetRequests returns the requests.
func (s *Statistics) GetRequests(ctx context.Context) []Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequests")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()

	return s.Requests
}

// GetRequestsByIPAddress returns the requests by IP address.
func (s *Statistics) GetRequestsByIPAddress(ctx context.Context, ipAddress string) []Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequestsByIPAddress")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		if r.IPAddress == ipAddress {
			requests = append(requests, r)
		}
	}

	return requests
}

// GetRequestsByTimeRange returns the requests by time range.
func (s *Statistics) GetRequestsByTimeRange(ctx context.Context, start, end time.Time) []Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequestsByTimeRange")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		if r.Timestamp.After(start) && r.Timestamp.Before(end) {
			requests = append(requests, r)
		}
		if r.Timestamp.Equal(start) || r.Timestamp.Equal(end) {
			requests = append(requests, r)
		}
	}

	return requests
}

// GetRequestsByUserAgent returns the requests by user agent.
func (s *Statistics) GetRequestsByUserAgent(ctx context.Context, userAgent string) []Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequestsByUserAgent")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		if r.UserAgent == userAgent {
			requests = append(requests, r)
		}
	}

	return requests
}

// GetRequestsGroupedByIPAddress returns the requests grouped by IP address.
func (s *Statistics) GetRequestsGroupedByIPAddress(ctx context.Context) map[string][]Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequestsGroupedByIPAddress")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.IPAddress] = append(grouped[r.IPAddress], r)
	}

	return grouped
}

// GetRequestsGroupedByUserAgent returns the requests grouped by user agent.
func (s *Statistics) GetRequestsGroupedByUserAgent(ctx context.Context) map[string][]Request {
	_, span := tracer.Start(ctx, "Statistics.GetRequestsGroupedByUserAgent")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.UserAgent] = append(grouped[r.UserAgent], r)
	}

	return grouped
}

// GetTotalDataSizeServed returns the data size served.
func (s *Statistics) GetTotalDataSizeServed(ctx context.Context) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalDataSizeServed")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		size += r.Size
	}

	return size
}

// GetTotalDataSizeServedByAgent returns the data size served by agent.
func (s *Statistics) GetTotalDataSizeServedByAgent(ctx context.Context, agent string) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalDataSizeServedByAgent")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		if r.UserAgent == agent {
			size += r.Size
		}
	}

	return size
}

// GetTotalDataSizeServedByIPAddress returns the data size served by IP address.
func (s *Statistics) GetTotalDataSizeServedByIPAddress(ctx context.Context, ipAddress string) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalDataSizeServedByIPAddress")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		if r.IPAddress == ipAddress {
			size += r.Size
		}
	}

	return size
}

// GetTotalRequestsByAgent returns the total requests by agent.
func (s *Statistics) GetTotalRequestsByAgent(ctx context.Context, agent string) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalRequestsByAgent")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var count int
	for _, r := range s.Requests {
		if r.UserAgent == agent {
			count++
		}
	}

	return count
}

// GetTotalRequestsByIPAddress returns the total requests by IP address.
func (s *Statistics) GetTotalRequestsByIPAddress(ctx context.Context, ipAddress string) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalRequestsByIPAddress")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var count int
	for _, r := range s.Requests {
		if r.IPAddress == ipAddress {
			count++
		}
	}

	return count
}

// GetTotalDataSizeServedByTimeRange returns the data size served by time range.
func (s *Statistics) GetTotalDataSizeServedByTimeRange(ctx context.Context, start, end time.Time) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalDataSizeServedByTimeRange")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		if r.Timestamp.After(start) && r.Timestamp.Before(end) {
			size += r.Size
		}
		if r.Timestamp.Equal(start) || r.Timestamp.Equal(end) {
			size += r.Size
		}
	}

	return size
}

// GetTotalRequests returns the total requests.
func (s *Statistics) GetTotalRequests(ctx context.Context) int {
	_, span := tracer.Start(ctx, "Statistics.GetTotalRequests")
	defer span.End()

	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()

	return len(s.Requests)
}

// GetTotalRobotsTxtViolators returns the total robots.txt violators.
func (s *Statistics) GetTotalRobotsTxtViolators(ctx context.Context) int {
	ctx, span := tracer.Start(ctx, "Statistics.GetTotalRobotsTxtViolators")
	defer span.End()

	requests := s.GetRequestsGroupedByUserAgent(ctx)
	violators := map[string]struct{}{}
	for identifier, requests := range requests {
		robotsTxtCounter := 0
		for _, request := range requests {
			if request.IsRobotsTxt {
				robotsTxtCounter++
			}
		}
		if robotsTxtCounter > 0 && robotsTxtCounter < len(requests) {
			violators[identifier] = struct{}{}
		}
	}

	return len(violators)
}

// UpdatePrompts updates the prompts.
func (s *Statistics) UpdatePrompts(ctx context.Context, prompts map[string]int) {
	_, span := tracer.Start(ctx, "Statistics.UpdatePrompts")
	defer span.End()

	s.PromptsLock.Lock()
	defer s.PromptsLock.Unlock()
	delta := 0
	for prompt := range prompts {
		if _, contains := s.Prompts[prompt]; !contains {
			delta++
		}
	}
	s.PromptsCount += delta
	s.Prompts = prompts
}
