package statisticsserver

import (
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strings"
)

type Data struct {
	Count               int
	Size                string
	IsRobotsTxtViolator string
}

// handleRoot is the handler for the root path.
func (ss *StatisticsServer) handleRoot(w http.ResponseWriter, _ *http.Request) {
	tpl, err := template.New("t").Parse(ss.htmlTemplates["index.gohtml"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := &strings.Builder{}
	ss.Statistics.PromptsLock.Lock()
	defer ss.Statistics.PromptsLock.Unlock()
	// TODO: sort by count
	byUserAgent := analyseStatistics(ss.Statistics.GetRequestsGroupedByUserAgent())

	// TODO: sort by count
	byIpAddress := analyseStatistics(ss.Statistics.GetRequestsGroupedByIpAddress())

	err = tpl.Execute(buffer, struct {
		ConfigurationInfo string
		Prompts           map[string]int
		ByUserAgent       map[string]Data
		ByIpAddress       map[string]Data
	}{
		ConfigurationInfo: ss.Statistics.ConfigurationInfo,
		Prompts:           ss.Statistics.Prompts,
		ByUserAgent:       byUserAgent,
		ByIpAddress:       byIpAddress,
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

// analyseStatistics is a helper function to analyze the statistics.
func analyseStatistics(requestData map[string][]statistics.Request) map[string]Data {
	data := map[string]Data{}
	for ipAddress, requests := range requestData {
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
		data[ipAddress] = Data{
			Count:               len(requests),
			Size:                convertByteSizeToSIUnits(size),
			IsRobotsTxtViolator: isRobotsTxtViolator,
		}
	}
	return data
}

// convertByteSizeToSIUnits converts the byte size to SI units.
func convertByteSizeToSIUnits(bytes int) string {
	byteToFloat64 := float64(bytes)
	for _, siUnits := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(byteToFloat64) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", byteToFloat64, siUnits)
		}
		byteToFloat64 /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", byteToFloat64)
}
