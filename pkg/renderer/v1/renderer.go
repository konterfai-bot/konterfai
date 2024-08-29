package renderer

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"go.opentelemetry.io/otel"
)

//go:embed assets
var assets embed.FS

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/renderer/v1")

var RandomizableTemplates = []string{"head.gohtml"}

type HTMLTemplate struct {
	Template       *template.Template
	IsRandomizable bool
}

type Renderer struct {
	htmlTemplates map[string]*HTMLTemplate
}

// RenderData is the structure for the RenderData.
func NewRenderer(ctx context.Context, logger *slog.Logger, _ []string) *Renderer {
	// not sure if we should generate the headline links here or pass them in
	_, span := tracer.Start(ctx, "NewRenderer")
	defer span.End()
	htmlTemplates := make(map[string]*HTMLTemplate)
	templates, err := assets.ReadDir("assets")
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("could not read assets (%v)", err))
	}
	for _, file := range templates {
		if file.IsDir() {
			continue
		}
		f, err := assets.ReadFile("assets/" + file.Name())
		if err != nil {
			logger.ErrorContext(ctx, fmt.Sprintf("could not read file (%v)", err))
			defer os.Exit(1)
			runtime.Goexit()
		}
		var isRandomizable bool
		if slices.Contains(RandomizableTemplates, file.Name()) {
			isRandomizable = true
		}
		htmlTemplates[file.Name()] = &HTMLTemplate{
			Template:       template.Must(template.New(file.Name()).Parse(string(f))),
			IsRandomizable: isRandomizable,
		}
	}

	return &Renderer{
		htmlTemplates: htmlTemplates,
	}
}

// RenderInRandomTemplate renders the data in a random template.
// Deprecated: use RenderRandomSite instead.
func (r *Renderer) RenderInRandomTemplate(ctx context.Context, _ interface{}) (string, error) {
	_, span := tracer.Start(ctx, "RenderInRandomTemplate")
	defer span.End()

	return r.RenderRandomSite(ctx)
}

func (r *Renderer) RenderRandomSite(ctx context.Context) (string, error) {
	_, span := tracer.Start(ctx, "Render")
	defer span.End()
	data := struct {
		LanguageCode string
		Title        string
	}{
		LanguageCode: functions.PickRandomStringFromSlice(ctx, &dictionaries.LanguageCodes),
		Title:        textblocks.RandomNewsPaperName(ctx) + " - " + textblocks.RandomHeadline(ctx),
	}

	return r.render(ctx, data)
}

func (r *Renderer) render(ctx context.Context, data interface{}) (string, error) {
	_, span := tracer.Start(ctx, "Render")
	defer span.End()
	buffer := &strings.Builder{}
	tpl := template.Must(template.ParseFS(assets, "assets/root.gohtml", "assets/head.gohtml"))
	for _, t := range tpl.Templates() {
		if slices.Contains(RandomizableTemplates, t.Name()) {
			rt, err := r.readTemplate(ctx, t.Name())
			if err != nil {
				return "", err
			}
			rt = r.randomizeLines(ctx, rt)
			_, err = t.Parse(rt)
			if err != nil {
				return "", err
			}
		}
	}
	err := tpl.Execute(buffer, data)

	return buffer.String(), err
}

func (r *Renderer) readTemplate(ctx context.Context, name string) (string, error) {
	_, span := tracer.Start(ctx, "ReadTemplate")
	defer span.End()
	f, err := assets.ReadFile("assets/" + name)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func (r *Renderer) randomizeLines(ctx context.Context, text string) string {
	_, span := tracer.Start(ctx, "RandomizeLines")
	defer span.End()
	lines := strings.Split(text, "\n")
	for i := range lines {
		j := rand.Intn(i + 1) //nolint:gosec
		lines[i], lines[j] = lines[j], lines[i]
	}

	return strings.Join(lines, "\n")
}
