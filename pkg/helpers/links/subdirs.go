package links

import (
	"math/rand"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/types"
	"github.com/google/uuid"
)

// generateSubDirectories generates a random number of subdirectories.
func generateSubDirectories(subdirectories int) string {
	sd := []string{}
	subcount := rand.Intn(subdirectories) + 1
	for i := 0; i < subcount; i++ {
		sd = append(sd, getSubDirectoryString())
	}
	return strings.Join(sd, "/")
}

// getSubDirectoryString returns a random subdirectory string.
func getSubDirectoryString() string {
	typeRand := rand.Intn(types.PathTypesCount) + 1
	switch types.PathTypes(typeRand) {
	case types.UUIDPath:
		return uuid.NewString()
	case types.NounPath:
		return functions.PickRandomFromSlice(&dictionaries.Nouns)
	case types.TwoNounPath:
		return strings.Join([]string{
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
		}, "")
	case types.ThreeNounPath:
		return strings.Join([]string{
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
		}, "")
	case types.TwoNounDashedPath:
		return strings.Join([]string{
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.ThreeNounDashedPath:
		return strings.Join([]string{
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.VerbPath:
		return functions.PickRandomFromSlice(&dictionaries.Verbs)
	case types.VerbNounPath:
		return strings.Join([]string{
			functions.PickRandomFromSlice(&dictionaries.Verbs),
			functions.PickRandomFromSlice(&dictionaries.Nouns),
		}, "")
	case types.DatePath:
		return functions.PickRandomDate()
	case types.YearPath:
		return functions.PickRandomYear()
	case types.MonthPath:
		return functions.PickRandomFromSlice(&dictionaries.Months)
	case types.DayPath:
		return functions.PickRandomFromSlice(&dictionaries.Weekdays)
	default:
		return "default"
	}
}
