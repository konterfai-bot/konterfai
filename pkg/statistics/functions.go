package statistics

import (
	"strings"
	"time"
)

// AppendRequest appends a request to the statistics.
func (s *Statistics) AppendRequest(r Request) {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	r.IpAddress = strings.Split(r.IpAddress, ":")[0]
	s.Requests = append(s.Requests, r)
}

// Clear clears the statistics.
func (s *Statistics) Clear() {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	s.Requests = []Request{}
}

// GetRequests returns the requests.
func (s *Statistics) GetRequests() []Request {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	return s.Requests
}

// GetRequestsByIpAddress returns the requests by IP address.
func (s *Statistics) GetRequestsByIpAddress(ipAddress string) []Request {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		if r.IpAddress == ipAddress {
			requests = append(requests, r)
		}
	}
	return requests
}

// GetRequestsByTimeRange returns the requests by time range.
func (s *Statistics) GetRequestsByTimeRange(start, end time.Time) []Request {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		if r.Timestamp.After(start) && r.Timestamp.Before(end) {
			requests = append(requests, r)
		}
	}
	return requests
}

// GetRequestsByUserAgent returns the requests by user agent.
func (s *Statistics) GetRequestsByUserAgent(userAgent string) []Request {
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

// GetRequestsGroupedByIpAddress returns the requests grouped by IP address.
func (s *Statistics) GetRequestsGroupedByIpAddress() map[string][]Request {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.IpAddress] = append(grouped[r.IpAddress], r)
	}
	return grouped
}

// GetRequestsGroupedByUserAgent returns the requests grouped by user agent.
func (s *Statistics) GetRequestsGroupedByUserAgent() map[string][]Request {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.UserAgent] = append(grouped[r.UserAgent], r)
	}
	return grouped
}

// GetTotalDataSizeServed returns the data size served.
func (s *Statistics) GetTotalDataSizeServed() int {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		size += r.Size
	}
	return size
}

// GetTotalDataSizeServedByTimeRange returns the data size served by time range.
func (s *Statistics) GetTotalDataSizeServedByTimeRange(start, end time.Time) int {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	var size int
	for _, r := range s.Requests {
		if r.Timestamp.After(start) && r.Timestamp.Before(end) {
			size += r.Size
		}
	}
	return size
}

// GetTotalRequests returns the total requests.
func (s *Statistics) GetTotalRequests() int {
	s.StatisticsLock.Lock()
	defer s.StatisticsLock.Unlock()
	return len(s.Requests)
}

// UpdatePrompts updates the prompts.
func (s *Statistics) UpdatePrompts(prompts map[string]int) {
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
