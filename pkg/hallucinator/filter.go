package hallucinator

import (
	"context"
	"regexp"
)

// isValidResult checks if the result is valid
func (h *Hallucinator) isValidResult(ctx context.Context, txt string) bool {
	ctx, span := tracer.Start(ctx, "Hallucinator.isValidResult")
	defer span.End()
	for _, re := range invalidResultsRegexps {
		r := regexp.MustCompile(re)
		if r.MatchString(txt) {
			return false
		}
	}
	return true
}
