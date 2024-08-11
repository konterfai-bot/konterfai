package statistics

import (
	"sync"
	"time"
)

// Statistics is the structure for the Statistics.
type Statistics struct {
	Requests          []Request
	Mutex             sync.Mutex
	ConfigurationInfo string
	// TODO: maybe use map map[string]int here to count the remaining requests
	Prompts      []string
	PromptsMutex sync.Mutex
}

// Request is the structure for the Request.
type Request struct {
	UserAgent   string    `yaml:"userAgent"`
	IpAddress   string    `yaml:"ipAddress"`
	Timestamp   time.Time `yaml:"timestamp"`
	IsRobotsTxt bool      `yaml:"isRobotsTxt"`
}

// NewStatistics creates a new Statistics instance.
func NewStatistics(configurationInfo string) *Statistics {
	return &Statistics{
		Requests:          []Request{},
		ConfigurationInfo: configurationInfo,
	}
}
