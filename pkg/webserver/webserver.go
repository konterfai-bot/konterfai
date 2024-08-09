package webserver

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
)

// WebServer is the structure for the WebServer.
type WebServer struct {
	Host                  string
	Port                  int
	Hallucinator          *hallucinator.Hallucinator
	Statistics            *statistics.Statistics
	HttpOkProbability     float64
	Uncertainty           float64
	HttpResponseCache     []WebServerCacheItem
	HttpResponseCacheSize int
	HttpResponseCacheLock sync.Mutex
	HttpBaseUrl           url.URL
}

// WebServerCacheItem is the structure for the WebServer cache item.
type WebServerCacheItem struct {
	Url  string
	Code int
}

// NewWebServer creates a new WebServer instance.
func NewWebServer(host string, port int, hallucinator *hallucinator.Hallucinator, statistics *statistics.Statistics, baseUrl url.URL, HttpOkProbability, Uncertainty float64, errorCacheSize int) *WebServer {
	return &WebServer{
		Host:                  host,
		Port:                  port,
		Hallucinator:          hallucinator,
		Statistics:            statistics,
		HttpOkProbability:     HttpOkProbability,
		Uncertainty:           Uncertainty,
		HttpResponseCache:     []WebServerCacheItem{},
		HttpResponseCacheSize: errorCacheSize,
		HttpBaseUrl:           baseUrl,
	}
}

// Serve starts the web server.
func (ws *WebServer) Serve() error {
	server := http.NewServeMux()
	server.HandleFunc("/robots.txt", ws.handleRobotsTxt)

	server.HandleFunc("/", ws.handleRoot)

	err := http.ListenAndServe(ws.Host+":"+strconv.Itoa(ws.Port), server)
	if err != nil {
		return err
	}
	return nil
}
