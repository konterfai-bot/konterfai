package statistics

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
)

// Statistics is the structure for the Statistics.
type Statistics struct {
	Requests          []Request
	StatisticsLock    sync.Mutex
	ConfigurationInfo string
	Prompts           map[string]int
	PromptsLock       sync.Mutex
	PromptsCount      int
}

// Request is the structure for the Request.
type Request struct {
	UserAgent   string    `yaml:"userAgent"`
	IPAddress   string    `yaml:"ipAddress"`
	Timestamp   time.Time `yaml:"timestamp"`
	IsRobotsTxt bool      `yaml:"isRobotsTxt"`
	Size        int       `yaml:"size"`
}

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/statistics")

// NewStatistics creates a new Statistics instance.
func NewStatistics(ctx context.Context, configurationInfo string) *Statistics {
	ctx, span := tracer.Start(ctx, "Statistics.NewStatistics")
	defer span.End()

	st := &Statistics{
		Requests:          []Request{},
		ConfigurationInfo: configurationInfo,
	}
	st.recordStatistics(ctx)

	return st
}
