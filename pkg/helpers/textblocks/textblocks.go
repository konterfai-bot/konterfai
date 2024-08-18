package textblocks

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/helpers/textblocks")

// RandomHeadline returns a random headline.
func RandomHeadline(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "textblocks.RandomHeadline")
	defer span.End()

	return fmt.Sprintf("%s: %s %s %s",
		functions.PickRandomStringFromSlice(ctx, &dictionaries.HeadlineStarters),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Verbs),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
	)
}

// RandomKeywords returns n random keywords.
func RandomKeywords(ctx context.Context, n int) string {
	ctx, span := tracer.Start(ctx, "textblocks.RandomKeywords")
	defer span.End()

	keygroup := functions.PickRandomSliceFromSlice(ctx, &dictionaries.MetaKeywordsGroups)
	keywords := []string{}
	for i := 0; i < n; i++ {
		keywords = append(keywords, functions.PickRandomStringFromSlice(ctx, &keygroup))
	}
	return strings.Join(keywords, ", ")
}

// RandomNewsPaperName returns a random newspaper name.
func RandomNewsPaperName(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "textblocks.RandomNewsPaperName")
	defer span.End()
	return fmt.Sprintf("%s %s",
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Cities),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.NewsPaperNames),
	)
}

// RandomTopic returns a random topic.
func RandomTopic(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "textblocks.RandomTopic")
	defer span.End()
	return fmt.Sprintf("%s %s %s",
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Verbs),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
	)
}
