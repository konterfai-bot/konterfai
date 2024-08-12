package statistics

import (
	"sync"
	"time"
)

// Statistics is the structure for the Statistics.
type Statistics struct {
	Requests          []Request
	StatisticsLock    sync.Mutex
	ConfigurationInfo string
	Prompts           map[string]int
	PromptsLock       sync.Mutex
}

// Request is the structure for the Request.
type Request struct {
	UserAgent   string    `yaml:"userAgent"`
	IpAddress   string    `yaml:"ipAddress"`
	Timestamp   time.Time `yaml:"timestamp"`
	IsRobotsTxt bool      `yaml:"isRobotsTxt"`
	Size        int       `yaml:"size"`
}

// NewStatistics creates a new Statistics instance.
func NewStatistics(configurationInfo string) *Statistics {
	return &Statistics{
		Requests:          []Request{},
		ConfigurationInfo: configurationInfo,
	}
}
