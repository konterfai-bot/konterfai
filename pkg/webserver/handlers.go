package webserver

import (
	"fmt"
	"net/http"

	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/helpers/robots"
)

// handleRobotsTxt handles the /robots.txt request.
func (ws *WebServer) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(robots.RobotsTxt(r))
	if err != nil {
		fmt.Println(fmt.Errorf("error writing robots.txt: %w", err))
	}
}

// handleRoot handles the root request.
func (ws *WebServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	httpCode := ws.getErrorFromCache(r.URL)
	if httpCode < 1 {
		if r.URL.Path == "/" || r.URL.Path == ws.HttpBaseUrl.Path || r.URL.Path == "" {
			httpCode = http.StatusOK
		} else {
			// We generate a random response code.
			httpCode = getRandomHttpResonseCode(functions.RecalculateProbabilityWithUncertainity(ws.HttpOkProbability, ws.Uncertainty))
			if r.URL.Path != "/" &&
				r.URL.Path != ws.HttpBaseUrl.Path &&
				r.URL.Path != "" &&
				// we do not want to store 200 OK responses
				httpCode != http.StatusOK &&
				// these are non persistant errors we do not want to save
				httpCode != http.StatusTooManyRequests &&
				httpCode != http.StatusServiceUnavailable &&
				httpCode != http.StatusTooEarly {
				ws.putErrorToCache(r.URL, httpCode)
			}
		}
	}
	if httpCode != http.StatusOK {
		if httpCode == http.StatusMovedPermanently {
			http.Redirect(w, r, links.RandomSimpleLink(ws.HttpBaseUrl), httpCode)
			return
		}
		// We write the response code to the response.
		w.WriteHeader(httpCode)
		// We write the response code to the response.
		_, err := w.Write([]byte(http.StatusText(httpCode)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	ws.handleHallucination(w, r)
}
