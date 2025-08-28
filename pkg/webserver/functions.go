package webserver

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"go.opentelemetry.io/otel/attribute"
)

// getRandomHTTPResonseCode returns a random http response code.
func getRandomHTTPResonseCode(ctx context.Context, okProbability float64) int {
	_, span := tracer.Start(ctx, "WebServer.getRandomHTTPResonseCode")
	defer span.End()

	if rand.Float64() < okProbability { //nolint:gosec
		return http.StatusOK
	}
	seed := time.Now().UTC().UnixNano()

	return ValidHTTPStatusCodes[rand.New(rand.NewSource(seed)).Intn(len(ValidHTTPStatusCodes))] //nolint:gosec
}

// handleHallucination handles the hallucination request.
func (ws *WebServer) handleHallucination(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "WebServer.handleHallucination")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.user-agent", r.UserAgent()),
		attribute.String("http.remote-addr", r.RemoteAddr),
	)

	r = r.WithContext(ctx)

	hallucination := ws.Hallucinator.PopRandomHallucination(ctx)
	go func() {
		ipAddr := r.RemoteAddr
		// Check for X-Forwarded-For header if behind a proxy
		if r.Header.Get("X-Forwarded-For") != "" {
			ipAddr = r.Header.Get("X-Forwarded-For")
		}
		ws.Statistics.AppendRequest(ctx, statistics.Request{
			IPAddress:   ipAddr,
			Timestamp:   time.Now(),
			UserAgent:   r.Header.Get("User-Agent"),
			IsRobotsTxt: false,
			Size:        len(hallucination),
		})
	}()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	actualSleep := time.Duration(rnd.Intn(int(ws.MaxPageDeliveryDelay.Nanoseconds() + 1)))
	time.Sleep(actualSleep)
	_, err := w.Write([]byte(hallucination))
	if err != nil {
		ws.Logger.ErrorContext(ctx, fmt.Sprintf("error writing hallucination (%v)", err.Error()))
	}
}

// getErrorFromCache returns the error code from the cache.
func (ws *WebServer) getErrorFromCache(ctx context.Context, requestURL *url.URL) int {
	_, span := tracer.Start(ctx, "WebServer.getErrorFromCache")
	defer span.End()

	ws.HTTPResponseCacheLock.Lock()
	defer ws.HTTPResponseCacheLock.Unlock()
	for _, item := range ws.HTTPResponseCache {
		if item.URL == fmt.Sprintf("%s%s", ws.HTTPBaseURL.String(), requestURL.Path) {
			return item.Code
		}
	}

	return 0
}

// putErrorToCache puts the error code to the cache.
func (ws *WebServer) putErrorToCache(ctx context.Context, requestURL *url.URL, errorCode int) {
	_, span := tracer.Start(ctx, "WebServer.putErrorToCache")
	defer span.End()

	ws.HTTPResponseCacheLock.Lock()
	defer ws.HTTPResponseCacheLock.Unlock()
	if len(ws.HTTPResponseCache) >= ws.HTTPResponseCacheSize {
		ws.HTTPResponseCache = append(ws.HTTPResponseCache[:0], ws.HTTPResponseCache[1:]...)
	}
	ws.HTTPResponseCache = append(ws.HTTPResponseCache, ErrorCacheItem{URL: fmt.Sprintf("%s%s",
		ws.HTTPBaseURL.String(), requestURL.String()), Code: errorCode})
}
