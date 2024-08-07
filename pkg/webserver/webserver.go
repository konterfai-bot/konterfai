package webserver

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
)

// WebServer is the structure for the WebServer.
type WebServer struct {
	Host                  string
	Port                  int
	Hallucinator          *hallucinator.Hallucinator
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
func NewWebServer(host string, port int, hallucinator *hallucinator.Hallucinator, baseUrl url.URL, HttpOkProbability, Uncertainty float64, errorCacheSize int) *WebServer {
	return &WebServer{
		Host:                  host,
		Port:                  port,
		Hallucinator:          hallucinator,
		HttpOkProbability:     HttpOkProbability,
		Uncertainty:           Uncertainty,
		HttpResponseCache:     []WebServerCacheItem{},
		HttpResponseCacheSize: errorCacheSize,
		HttpBaseUrl:           baseUrl,
	}
}

// Serve starts the web server.
func (ws *WebServer) Serve() error {
	http.HandleFunc("/robots.txt", ws.handleRobotsTxt)

	http.HandleFunc("/", ws.handleRoot)

	err := http.ListenAndServe(ws.Host+":"+strconv.Itoa(ws.Port), nil)
	if err != nil {
		return err
	}
	return nil
}
