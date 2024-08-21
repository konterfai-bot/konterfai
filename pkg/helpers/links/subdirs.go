package links

import (
	"context"
	"math/rand"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/types"
	"github.com/google/uuid"
)

// generateSubDirectories generates a random number of subdirectories.
func generateSubDirectories(ctx context.Context, subdirectories int) string {
	ctx, span := tracer.Start(ctx, "generateSubDirectories")
	defer span.End()

	if subdirectories < 1 {
		return ""
	}

	sd := []string{}
	subcount := rand.Intn(subdirectories) + 1 //nolint:gosec
	for range subcount {
		sd = append(sd, getSubDirectoryString(ctx))
	}

	return strings.Join(sd, "/")
}

// getSubDirectoryString returns a random subdirectory string.
func getSubDirectoryString(ctx context.Context) string { //nolint:cyclop
	ctx, span := tracer.Start(ctx, "getSubDirectoryString")
	defer span.End()

	typeRand := rand.Intn(types.PathTypesCount) + 1 //nolint:gosec
	switch types.PathTypes(typeRand) {
	case types.UUIDPath:
		return uuid.NewString()
	case types.NounPath:
		return functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns)
	case types.TwoNounPath:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		}, "")
	case types.ThreeNounPath:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		}, "")
	case types.TwoNounDashedPath:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		}, "-")
	case types.ThreeNounDashedPath:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		}, "-")
	case types.VerbPath:
		return functions.PickRandomStringFromSlice(ctx, &dictionaries.Verbs)
	case types.VerbNounPath:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Verbs),
			functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns),
		}, "")
	case types.DatePath:
		return functions.PickRandomDate(ctx)
	case types.YearPath:
		return functions.PickRandomYear(ctx)
	case types.MonthPath:
		return functions.PickRandomStringFromSlice(ctx, &dictionaries.Months)
	case types.DayPath:
		return functions.PickRandomStringFromSlice(ctx, &dictionaries.Weekdays)
	default:
		return fallbackDefaultWord
	}
}
