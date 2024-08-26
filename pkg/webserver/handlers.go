package webserver

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/helpers/robots"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"go.opentelemetry.io/otel/attribute"
)

// handleRobotsTxt handles the /robots.txt request.
func (ws *WebServer) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "WebServer.handleRobotsTxt")
	defer span.End()
	r = r.WithContext(ctx)

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.user-agent", r.UserAgent()),
		attribute.String("http.remote-addr", r.RemoteAddr),
	)

	responseData := robots.RobotsTxt(r)
	go func() {
		ws.Statistics.AppendRequest(ctx, statistics.Request{
			IPAddress:   r.RemoteAddr,
			Timestamp:   time.Now(),
			UserAgent:   r.Header.Get("User-Agent"),
			IsRobotsTxt: true,
			Size:        len(responseData),
		})
	}()
	_, err := w.Write(responseData)
	if err != nil {
		ws.Logger.ErrorContext(ctx, fmt.Sprintf("error writing robots.txt (%v)", err.Error()))
	}
}

// handleRoot handles the root request.
func (ws *WebServer) handleRoot(w http.ResponseWriter, r *http.Request) { //nolint:cyclop
	ctx, span := tracer.Start(r.Context(), "WebServer.handleRoot")
	defer span.End()
	span.SetAttributes(attribute.String("http.method", r.Method), attribute.String("http.url", r.URL.String()),
		attribute.String("http.user-agent", r.UserAgent()), attribute.String("http.remote-addr", r.RemoteAddr))
	r = r.WithContext(ctx)

	httpCode := ws.getErrorFromCache(ctx, r.URL)
	if httpCode < 1 {
		if r.URL.Path == "/" || r.URL.Path == ws.HTTPBaseURL.Path || r.URL.Path == "" {
			httpCode = http.StatusOK
		} else {
			// We generate a random response code.
			httpCode = getRandomHTTPResonseCode(ctx,
				functions.RecalculateProbabilityWithUncertainity(ctx, ws.HTTPOkProbability, ws.Uncertainty, 0))
			if r.URL.Path != "/" &&
				r.URL.Path != ws.HTTPBaseURL.Path &&
				r.URL.Path != "" &&
				// we do not want to store 200 OK responses
				httpCode != http.StatusOK &&
				// these are non-persistent errors we do not want to save
				httpCode != http.StatusTooManyRequests &&
				httpCode != http.StatusServiceUnavailable &&
				httpCode != http.StatusTooEarly {
				ws.putErrorToCache(ctx, r.URL, httpCode)
			}
		}
	}
	if httpCode != http.StatusOK {
		if httpCode == http.StatusMovedPermanently {
			http.Redirect(w, r, links.RandomSimpleLink(ctx, ws.HTTPBaseURL), httpCode)

			return
		}
		// We write the response code to the response.
		w.WriteHeader(httpCode)
		_, err := w.Write([]byte(http.StatusText(httpCode)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}
	ws.handleHallucination(w, r)
}
