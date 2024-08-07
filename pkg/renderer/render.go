package renderer

import (
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed assets
var assets embed.FS

// Renderer is the structure for the Renderer.
type Renderer struct {
	htmlTemplates []string
	headlineLinks []string
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

func NewRenderer(headLineLinks []string) *Renderer {
	htmlTemplates := []string{}
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
		htmlTemplates = append(htmlTemplates, string(f))
	}
	return &Renderer{
		htmlTemplates: htmlTemplates,
		headlineLinks: headLineLinks,
	}
}

// RenderInRandomTemplate renders the given text in a random template using go templates.
func (r *Renderer) RenderInRandomTemplate(rd RenderData) (string, error) {
	tpl, err := template.New("t").Parse(r.getRandomTemplate())
	if err != nil {
		return "", err
	}
	if rd.HeadlineLinks == nil || len(rd.HeadlineLinks) < 10 {
		rd.HeadlineLinks = r.headlineLinks
	}
	year, _, _ := time.Now().Date()
	rd.CurrentYear = fmt.Sprintf("%d", year)
	buffer := &strings.Builder{}
	err = tpl.Execute(buffer, rd)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// getRandomTemplate returns a random template.
func (r *Renderer) getRandomTemplate() string {
	return r.htmlTemplates[rand.Intn(len(r.htmlTemplates))]
}
