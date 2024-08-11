package statisticsserver

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"codeberg.org/konterfai/konterfai/pkg/statistics"
)

//go:embed assets
var assets embed.FS

// StatisticsServer is the structure for the StatisticsServer.
type StatisticsServer struct {
	Host       string
	Port       int
	Statistics *statistics.Statistics

	htmlTemplates map[string]string
}

// NewStatisticsServer creates a new StatisticsServer instance.
func NewStatisticsServer(host string, port int, statistics *statistics.Statistics) *StatisticsServer {
	htmlTemplates := map[string]string{}
	templates, err := assets.ReadDir("assets")
	if err != nil {
		fmt.Println(fmt.Errorf("could not read assets directory (%v)", err))
		os.Exit(1)
	}
	for _, file := range templates {
		if file.IsDir() {
			continue
		}
		f, err := assets.ReadFile(fmt.Sprintf("assets/%s", file.Name()))
		if err != nil {
			fmt.Println(fmt.Errorf("could not read asset file (%v)", err))
			os.Exit(1)
		}
		htmlTemplates[file.Name()] = string(f)
	}
	return &StatisticsServer{
		Host:          host,
		Port:          port,
		Statistics:    statistics,
		htmlTemplates: htmlTemplates,
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
