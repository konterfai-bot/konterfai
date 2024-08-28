package renderer

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
)

//go:embed assets
var assets embed.FS

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/renderer/v0")

// Renderer is the structure for the Renderer.
type Renderer struct {
	htmlTemplates     []string
	htmlTemplatesLock sync.Mutex
	headlineLinks     []string
}

// MetaData is the structure for the meta-tags.
type MetaData struct {
	Description string
	Keywords    string
	Charset     string
}

// RenderData is the structure for the RenderData.
type RenderData struct {
	NewsAnchor    string
	Headline      string
	Content       template.HTML
	FollowUpLink  template.HTML
	HeadlineLinks []string
	RandomTopics  []RandomTopic
	Year          string
	CurrentYear   string
	MetaData      MetaData
	LanguageCode  string
}

// RandomTopic is the structure for the RandomTopic.
type RandomTopic struct {
	Topic string
	Link  string
}

func NewRenderer(ctx context.Context, logger *slog.Logger, headLineLinks []string) *Renderer {
	_, span := tracer.Start(ctx, "NewRenderer")
	defer span.End()

	htmlTemplates := []string{}
	templates, err := assets.ReadDir("assets")
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("could not read assets directory (%v)", err))
		defer os.Exit(1)
		runtime.Goexit()
	}
	for _, file := range templates {
		if file.IsDir() {
			continue
		}
		f, err := assets.ReadFile("assets/" + file.Name())
		if err != nil {
			logger.ErrorContext(ctx, fmt.Sprintf("could not read asset file (%v)", err))
			defer os.Exit(1)
			runtime.Goexit()
		}
		htmlTemplates = append(htmlTemplates, string(f))
	}

	return &Renderer{htmlTemplates: htmlTemplates, headlineLinks: headLineLinks}
}

// RenderInRandomTemplate renders the given text in a random template using go templates.
func (r *Renderer) RenderInRandomTemplate(ctx context.Context, rd RenderData) (string, error) {
	ctx, span := tracer.Start(ctx, "Renderer.RenderInRandomTemplate")
	defer span.End()

	tplContent, err := r.getRandomTemplate(ctx)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("t").Parse(tplContent)
	if err != nil {
		return "", err
	}
	if rd.HeadlineLinks == nil || len(rd.HeadlineLinks) < 10 {
		if r.headlineLinks == nil || len(r.headlineLinks) < 10 {
			return "", errors.New("headlineLinks is nil or has less than 10 elements, is empty or unset")
		}
		rd.HeadlineLinks = r.headlineLinks
	}

	year, _, _ := time.Now().Date()
	rd.CurrentYear = strconv.Itoa(year)
	buffer := &strings.Builder{}
	err = tpl.Execute(buffer, rd)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// SetTemplates sets the templates, at the moment only used for testing.
func (r *Renderer) SetTemplates(templates []string) {
	_, span := tracer.Start(context.Background(), "Renderer.SetTemplates")
	defer span.End()
	r.htmlTemplatesLock.Lock()
	defer r.htmlTemplatesLock.Unlock()
	r.htmlTemplates = templates
}

// getRandomTemplate returns a random template.
func (r *Renderer) getRandomTemplate(ctx context.Context) (string, error) {
	_, span := tracer.Start(ctx, "Renderer.getRandomTemplate")
	defer span.End()

	r.htmlTemplatesLock.Lock()
	defer r.htmlTemplatesLock.Unlock()

	if len(r.htmlTemplates) < 1 {
		return "", errors.New("no templates found")
	}

	return r.htmlTemplates[rand.Intn(len(r.htmlTemplates))], nil //nolint:gosec
}
