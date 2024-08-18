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

// getRandomHttpResonseCode returns a random http response code.
func getRandomHttpResonseCode(ctx context.Context, okProbability float64) int {
	ctx, span := tracer.Start(ctx, "WebServer.getRandomHttpResonseCode")
	defer span.End()

	if rand.Float64() < okProbability {
		return http.StatusOK
	}
	seed := time.Now().UTC().UnixNano()
	return ValidHttpStatusCodes[rand.New(rand.NewSource(seed)).Intn(len(ValidHttpStatusCodes))]
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
		ws.Statistics.AppendRequest(ctx, statistics.Request{
			IpAddress:   r.RemoteAddr,
			Timestamp:   time.Now(),
			UserAgent:   r.Header.Get("User-Agent"),
			IsRobotsTxt: false,
			Size:        len(hallucination),
		})
	}()
	_, err := w.Write([]byte(hallucination))
	if err != nil {
		fmt.Println(fmt.Errorf("error writing hallucination: %w", err))
	}
}

// getErrorFromCache returns the error code from the cache.
func (ws *WebServer) getErrorFromCache(ctx context.Context, requestUrl *url.URL) int {
	ctx, span := tracer.Start(ctx, "WebServer.getErrorFromCache")
	defer span.End()

	ws.HttpResponseCacheLock.Lock()
	defer ws.HttpResponseCacheLock.Unlock()
	for _, item := range ws.HttpResponseCache {
		if item.Url == fmt.Sprintf("%s%s", ws.HttpBaseUrl.String(), requestUrl.Path) {
			return item.Code
		}
	}
	return 0
}

// putErrorToCache puts the error code to the cache.
func (ws *WebServer) putErrorToCache(ctx context.Context, requestUrl *url.URL, errorCode int) {
	ctx, span := tracer.Start(ctx, "WebServer.putErrorToCache")
	defer span.End()

	ws.HttpResponseCacheLock.Lock()
	defer ws.HttpResponseCacheLock.Unlock()
	if len(ws.HttpResponseCache) >= ws.HttpResponseCacheSize {
		ws.HttpResponseCache = append(ws.HttpResponseCache[:0], ws.HttpResponseCache[1:]...)
	}
	ws.HttpResponseCache = append(ws.HttpResponseCache, WebServerCacheItem{Url: fmt.Sprintf("%s%s", ws.HttpBaseUrl.String(), requestUrl.String()), Code: errorCode})
}
