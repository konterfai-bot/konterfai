package robots

import (
	"fmt"
	"math/rand"
	"net/http"
	"slices"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/helpers/robots")

// RobotsTxt generates the robots.txt content.
func RobotsTxt(r *http.Request) []byte {
	ctx, span := tracer.Start(r.Context(), "RobotsTxt")
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.user-agent", r.UserAgent()),
		attribute.String("http.remote-addr", r.RemoteAddr),
	)
	defer span.End()

	r = r.WithContext(ctx)

	// This list has been inspired by https://hellocoding.de/blog/seo/ki-ausschliessen-von-webseite
	// and https://www.cyberciti.biz/web-developer/block-openai-bard-bing-ai-crawler-bots-using-robots-txt-file/
	// We print that out to tell the ai crawlers to not index this site.
	// If they do not comply, they had it coming.
	const disallow = "Disallow: /\n\n"
	robotsTxt := [][]byte{
		[]byte("# AI Bots are not welcome here\n"),
		[]byte("User-Agent: GPTBot\n" +
			disallow),
		[]byte("# OpenAI ChatGPT bot\n"),
		[]byte("User-Agent: ChatGPT-User\n" +
			disallow),
		[]byte("# Google Bots\n"),
		[]byte("User-Agent: Google-Extended\n" +
			disallow),
		[]byte("# Apple Bots\n"),
		[]byte("User-Agent: Applebot-Extended\n" +
			disallow),
		[]byte("# Microsoft Bots\n"),
		[]byte("User-Agent: CCBot\n" +
			disallow),

		[]byte("User-agent: PerplexityBot\n" +
			disallow),
		[]byte("User-agent: anthropic-ai\n" +
			disallow),
		[]byte("User-agent: Claude-Web\n" +
			disallow),
		[]byte("User-agent: ClaudeBot\n" +
			disallow),
		[]byte("User-agent: Amazonbot\n" +
			disallow),
		[]byte("User-agent: Omgilibot\n" +
			disallow),
		[]byte("User-agent: Omgili\n" +
			disallow),
		[]byte("User-Agent: FacebookBot\n" +
			disallow),
		[]byte("User-agent: anthropic-ai\n" +
			disallow),
		[]byte("User-agent: Bytespider\n" +
			disallow),
		[]byte("User-agent: YouBot\n" +
			disallow),
		[]byte("User-agent: ImagesiftBot\n" +
			disallow),
		// We add the user agent of the request to the robots.txt.
		// Since you landed here there is probably a reason for that.
		[]byte(fmt.Sprintf("User-Agent: %s\n", r.UserAgent()) +
			disallow),
	}
	// We shuffle the robots.txt to make it harder for the ai to learn the pattern.
	rand.Shuffle(len(robotsTxt), func(i, j int) {
		robotsTxt[i], robotsTxt[j] = robotsTxt[j], robotsTxt[i]
	})
	return slices.Concat(robotsTxt...)
}
