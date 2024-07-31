package hallucinator

import "regexp"

// isValidResult checks if the result is valid
func (h *Hallucinator) isValidResult(txt string) bool {
	for _, re := range invalidResultsRegexps {
		r := regexp.MustCompile(re)
		if r.MatchString(txt) {
			return false
		}
	}
	return true
}
