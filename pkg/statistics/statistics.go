package statistics

import (
	"context"
	"go.opentelemetry.io/otel"
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
	PromptsCount      int
	Context           context.Context
}

// Request is the structure for the Request.
type Request struct {
	UserAgent   string    `yaml:"userAgent"`
	IpAddress   string    `yaml:"ipAddress"`
	Timestamp   time.Time `yaml:"timestamp"`
	IsRobotsTxt bool      `yaml:"isRobotsTxt"`
	Size        int       `yaml:"size"`
}

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/statistics")

// NewStatistics creates a new Statistics instance.
func NewStatistics(ctx context.Context, configurationInfo string) *Statistics {
	st := &Statistics{
		Requests:          []Request{},
		ConfigurationInfo: configurationInfo,
		Context:           ctx,
	}
	st.recordStatistics()
	return st
}
