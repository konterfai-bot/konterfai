package statisticsserver

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strings"
)

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
	type Data struct {
		Count int
		Size  string
	}
	byUserAgent := map[string]Data{}
	for userAgent, requests := range ss.Statistics.GetRequestsGroupedByUserAgent() {
		size := 0
		for _, request := range requests {
			size += request.Size
		}
		byUserAgent[userAgent] = Data{
			Count: len(requests),
			Size:  convertByteSizeToSIUnits(size),
		}
	}

	// TODO: sort by count
	byIpAddress := map[string]Data{}
	for ipAddress, requests := range ss.Statistics.GetRequestsGroupedByIpAddress() {
		size := 0
		for _, request := range requests {
			size += request.Size
		}
		byIpAddress[ipAddress] = Data{
			Count: len(requests),
			Size:  convertByteSizeToSIUnits(size),
		}
	}

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
