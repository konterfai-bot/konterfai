package statisticsserver

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

//go:embed assets
var assets embed.FS

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/statisticsserver")

// StatisticsServer is the structure for the StatisticsServer.
type StatisticsServer struct {
	Host       string
	Port       int
	Statistics *statistics.Statistics

	htmlTemplates map[string]string
}

// NewStatisticsServer creates a new StatisticsServer instance.
func NewStatisticsServer(ctx context.Context, host string, port int, st *statistics.Statistics) *StatisticsServer {
	_, span := tracer.Start(ctx, "NewStatisticsServer")
	defer span.End()

	htmlTemplates := map[string]string{}
	templates, err := assets.ReadDir("assets")
	if err != nil {
		fmt.Printf("could not read assets directory (%v)\n", err)
		defer os.Exit(1)
		runtime.Goexit()
	}
	for _, file := range templates {
		if file.IsDir() {
			continue
		}
		f, err := assets.ReadFile("assets/" + file.Name())
		if err != nil {
			fmt.Printf("could not read asset file (%v)\n", err)
			defer os.Exit(1)
			runtime.Goexit()
		}
		htmlTemplates[file.Name()] = string(f)
	}

	return &StatisticsServer{
		Host:          host,
		Port:          port,
		Statistics:    st,
		htmlTemplates: htmlTemplates,
	}
}

// Serve starts the statistics server.
func (ss *StatisticsServer) Serve(ctx context.Context) error {
	_, span := tracer.Start(ctx, "StatisticsServer.Serve")
	defer span.End()

	serverMux := http.NewServeMux()
	serverMux.Handle("/metrics", promhttp.Handler())
	serverMux.HandleFunc("/", ss.handleRoot)
	server := &http.Server{
		Addr:              ss.Host + ":" + strconv.Itoa(ss.Port),
		Handler:           serverMux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
