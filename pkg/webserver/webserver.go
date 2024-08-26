package webserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"go.opentelemetry.io/otel"
)

// WebServer is the structure for the WebServer.
type WebServer struct {
	Host                  string
	Port                  int
	Hallucinator          *hallucinator.Hallucinator
	Statistics            *statistics.Statistics
	HTTPOkProbability     float64
	Uncertainty           float64
	HTTPResponseCache     []ErrorCacheItem
	HTTPResponseCacheSize int
	HTTPResponseCacheLock sync.Mutex
	HTTPBaseURL           url.URL
	ServeMux              *http.ServeMux
	Logger                *slog.Logger
}

// ErrorCacheItem is the structure for the WebServer cache item.
type ErrorCacheItem struct {
	URL  string
	Code int
}

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/webserver")

// NewWebServer creates a new WebServer instance.
func NewWebServer(ctx context.Context, logger *slog.Logger, host string, port int,
	hallucinator *hallucinator.Hallucinator, statistics *statistics.Statistics, baseURL url.URL, httpOkProbability,
	uncertainty float64, errorCacheSize int,
) *WebServer {
	_, span := tracer.Start(ctx, "NewWebServer")
	defer span.End()

	return &WebServer{
		Host:                  host,
		Port:                  port,
		Hallucinator:          hallucinator,
		Statistics:            statistics,
		HTTPOkProbability:     httpOkProbability,
		Uncertainty:           uncertainty,
		HTTPResponseCache:     []ErrorCacheItem{},
		HTTPResponseCacheSize: errorCacheSize,
		HTTPBaseURL:           baseURL,
		Logger:                logger,
	}
}

// Serve starts the web server.
func (ws *WebServer) Serve(ctx context.Context) error {
	_, span := tracer.Start(ctx, "WebServer.Serve")
	defer span.End()

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/robots.txt", ws.handleRobotsTxt)
	serverMux.HandleFunc("/", ws.handleRoot)
	server := &http.Server{
		Addr:              ws.Host + ":" + strconv.Itoa(ws.Port),
		Handler:           serverMux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
