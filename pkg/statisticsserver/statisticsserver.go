package statisticsserver

import (
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"net/http"
	"strconv"
)

// StatisticsServer is the structure for the StatisticsServer.
type StatisticsServer struct {
	Host       string
	Port       int
	Statistics *statistics.Statistics
}

// NewStatisticsServer creates a new StatisticsServer instance.
func NewStatisticsServer(host string, port int, statistics *statistics.Statistics) *StatisticsServer {
	return &StatisticsServer{
		Host:       host,
		Port:       port,
		Statistics: statistics,
	}
}

// Serve starts the statistics server.
func (ss *StatisticsServer) Serve() error {
	server := http.NewServeMux()
	server.HandleFunc("/", ss.handleRoot)
	err := http.ListenAndServe(ss.Host+":"+strconv.Itoa(ss.Port), server)
	if err != nil {
		return err
	}
	return nil
}
