package statisticsserver

import (
	"html/template"
	"net/http"
	"strings"
)

// handleRoot is the handler for the root path.
func (ss *StatisticsServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("t").Parse(ss.htmlTemplates["index.gohtml"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := &strings.Builder{}
	ss.Statistics.PromptsLock.Lock()
	defer ss.Statistics.PromptsLock.Unlock()
	// TODO: sort by count
	byUserAgent := map[string]int{}
	for userAgent, requests := range ss.Statistics.GetRequestsGroupedByUserAgent() {
		byUserAgent[userAgent] = len(requests)
	}

	// TODO: sort by count
	byIpAddress := map[string]int{}
	for ipAddress, requests := range ss.Statistics.GetRequestsGroupedByIpAddress() {
		byIpAddress[ipAddress] = len(requests)
	}

	err = tpl.Execute(buffer, struct {
		ConfigurationInfo string
		Prompts           map[string]int
		ByUserAgent       map[string]int
		ByIpAddress       map[string]int
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
