package statistics

import "time"

// AppendRequest appends a request to the statistics.
func (s *Statistics) AppendRequest(r Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Requests = append(s.Requests, r)
}

// Clear clears the statistics.
func (s *Statistics) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Requests = []Request{}
}

// GetRequests returns the requests.
func (s *Statistics) GetRequests() []Request {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Requests
}

// GetRequestsByIpAddress returns the requests by IP address.
func (s *Statistics) GetRequestsByIpAddress(ipAddress string) []Request {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var requests []Request
	for _, r := range s.Requests {
		t, err := time.Parse(time.RFC3339, r.Timestamp)
		if err != nil {
			continue
		}
		if t.After(start) && t.Before(end) {
			requests = append(requests, r)
		}
	}
	return requests
}

// GetRequestsByUserAgent returns the requests by user agent.
func (s *Statistics) GetRequestsByUserAgent(userAgent string) []Request {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.IpAddress] = append(grouped[r.IpAddress], r)
	}
	return grouped
}

// GetRequestsGroupedByUserAgent returns the requests grouped by user agent.
func (s *Statistics) GetRequestsGroupedByUserAgent() map[string][]Request {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	grouped := make(map[string][]Request)
	for _, r := range s.Requests {
		grouped[r.UserAgent] = append(grouped[r.UserAgent], r)
	}
	return grouped
}

// ClearPersistent clears the persistent statistics.
func (s *Statistics) ClearPersistent() error {
	// TODO: Implement
	return nil
}

// Persist persists the statistics.
func (s *Statistics) Persist() error {
	// TODO: Implement
	return nil
}

// Load loads the statistics.
func (s *Statistics) Load() error {
	// TODO: Implement
	return nil
}
