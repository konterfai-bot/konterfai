package statisticsserver

import (
	"context"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"sort"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"go.opentelemetry.io/otel/attribute"
)

// analyseStatistics is a helper function to analyze the statistics.
func analyseStatistics(ctx context.Context, rd map[string][]statistics.Request) RequestDataSlice {
	ctx, span := tracer.Start(ctx, "StatisticsServer.analyseStatistics")
	defer span.End()

	data := RequestDataSlice{}
	for identifier, requests := range rd {
		size := 0
		isRobotsTxtViolator := "no"
		robotsTxtCounter := 0
		for _, request := range requests {
			size += request.Size
			if request.IsRobotsTxt {
				robotsTxtCounter++
			}
		}
		if robotsTxtCounter == 0 {
			isRobotsTxtViolator = "ignored"
		}
		if robotsTxtCounter > 0 && robotsTxtCounter < len(requests) {
			isRobotsTxtViolator = "yes"
		}
		data = append(data, &RequestData{
			Identifier:          identifier,
			Count:               len(requests),
			Size:                convertByteSizeToSIUnits(ctx, size),
			IsRobotsTxtViolator: isRobotsTxtViolator,
		})
	}
	sort.Sort(data)

	return data
}

// convertByteSizeToSIUnits converts the byte size to SI units.
func convertByteSizeToSIUnits(ctx context.Context, bytes int) string {
	_, span := tracer.Start(ctx, "StatisticsServer.convertByteSizeToSIUnits")
	defer span.End()

	byteToFloat64 := float64(bytes)
	for _, siUnits := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(byteToFloat64) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", byteToFloat64, siUnits)
		}
		byteToFloat64 /= 1024.0
	}

	return fmt.Sprintf("%.1fYiB", byteToFloat64)
}

// handleRoot is the handler for the root path.
func (ss *StatisticsServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "StatisticsServer.handleRoot")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.user-agent", r.UserAgent()),
		attribute.String("http.remote-addr", r.RemoteAddr),
	)
	_ = r.WithContext(ctx)

	tpl, err := template.New("t").Parse(ss.htmlTemplates["index.gohtml"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	buffer := &strings.Builder{}

	byUserAgent := analyseStatistics(ctx, ss.Statistics.GetRequestsGroupedByUserAgent(ctx))

	byIPAddress := analyseStatistics(ctx, ss.Statistics.GetRequestsGroupedByIPAddress(ctx))

	totalDataSize := convertByteSizeToSIUnits(ctx, ss.Statistics.GetTotalDataSizeServed(ctx))

	totalRequests := len(ss.Statistics.Requests)

	ss.Statistics.PromptsLock.Lock()
	defer ss.Statistics.PromptsLock.Unlock()

	err = tpl.Execute(buffer, struct {
		ConfigurationInfo string
		Prompts           map[string]int
		ByUserAgent       RequestDataSlice
		ByIPAddress       RequestDataSlice
		TotalDataSize     string
		TotalRequests     int
		TotalPrompts      int
	}{
		ConfigurationInfo: ss.Statistics.ConfigurationInfo,
		Prompts:           ss.Statistics.Prompts,
		ByUserAgent:       byUserAgent,
		ByIPAddress:       byIPAddress,
		TotalDataSize:     totalDataSize,
		TotalRequests:     totalRequests,
		TotalPrompts:      ss.Statistics.PromptsCount,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	_, err = w.Write([]byte(buffer.String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
