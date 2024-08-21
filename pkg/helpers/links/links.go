package links

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/helpers/links")

// RandomLink generates a random link based on given base URL, the number of random subdirectories and random variables.
func RandomLink(ctx context.Context, baseURL url.URL, subdirectories,
	variablesCount int, linkHasVariablesProbability float64,
) string {
	ctx, span := tracer.Start(ctx, "RandomLink")
	defer span.End()

	subDirectoryPath := generateSubDirectories(ctx, subdirectories)

	variables := generateVariables(ctx, variablesCount, linkHasVariablesProbability)

	if variables != "" {
		return fmt.Sprintf("%s://%s/%s?%s", baseURL.Scheme, baseURL.Host, subDirectoryPath, variables)
	}

	return fmt.Sprintf("%s://%s/%s", baseURL.Scheme, baseURL.Host, subDirectoryPath)
}

// RandomSimpleLink generates a random link based on the given base URL.
func RandomSimpleLink(ctx context.Context, baseURL url.URL) string {
	ctx, span := tracer.Start(ctx, "RandomSimpleLink")
	defer span.End()

	name := strings.ToLower(strings.Join([]string{
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
	}, "-"))

	return fmt.Sprintf("%s://%s/%s", baseURL.Scheme, baseURL.Host, name)
}
