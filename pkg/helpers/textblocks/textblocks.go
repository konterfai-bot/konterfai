package textblocks

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"fmt"
)

// RandomHeadline returns a random headline.
func RandomHeadline() string {
	return fmt.Sprintf("%s: %s %s %s",
		functions.PickRandomFromSlice(&dictionaries.HeadlineStarters),
		functions.PickRandomFromSlice(&dictionaries.Nouns),
		functions.PickRandomFromSlice(&dictionaries.Verbs),
		functions.PickRandomFromSlice(&dictionaries.Nouns),
	)
}

// RandomNewsPaperName returns a random newspaper name.
func RandomNewsPaperName() string {
	return fmt.Sprintf("%s %s",
		functions.PickRandomFromSlice(&dictionaries.Cities),
		functions.PickRandomFromSlice(&dictionaries.NewsPaperNames),
	)
}

// RandomTopic returns a random topic.
func RandomTopic() string {
	return fmt.Sprintf("%s %s %s",
		functions.PickRandomFromSlice(&dictionaries.Nouns),
		functions.PickRandomFromSlice(&dictionaries.Verbs),
		functions.PickRandomFromSlice(&dictionaries.Nouns),
	)
}
