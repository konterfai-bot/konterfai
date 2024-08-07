package textblocks

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"fmt"
	"strings"
)

// RandomHeadline returns a random headline.
func RandomHeadline() string {
	return fmt.Sprintf("%s: %s %s %s",
		functions.PickRandomStringFromSlice(&dictionaries.HeadlineStarters),
		functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		functions.PickRandomStringFromSlice(&dictionaries.Verbs),
		functions.PickRandomStringFromSlice(&dictionaries.Nouns),
	)
}

// RandomKeywords returns n random keywords.
func RandomKeywords(n int) string {
	keygroup := functions.PickRandomSliceFromSlice(&dictionaries.MetaKeywordsGroups)
	keywords := []string{}
	for i := 0; i < n; i++ {
		keywords = append(keywords, functions.PickRandomStringFromSlice(&keygroup))
	}
	return strings.Join(keywords, ", ")
}

// RandomNewsPaperName returns a random newspaper name.
func RandomNewsPaperName() string {
	return fmt.Sprintf("%s %s",
		functions.PickRandomStringFromSlice(&dictionaries.Cities),
		functions.PickRandomStringFromSlice(&dictionaries.NewsPaperNames),
	)
}

// RandomTopic returns a random topic.
func RandomTopic() string {
	return fmt.Sprintf("%s %s %s",
		functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		functions.PickRandomStringFromSlice(&dictionaries.Verbs),
		functions.PickRandomStringFromSlice(&dictionaries.Nouns),
	)
}
