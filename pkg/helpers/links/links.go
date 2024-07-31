package links

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"fmt"
	"net/url"
	"strings"
)

// RandomLink generates a random link based on the given base URL, the number of random subdirectories and the number of random variables.
func RandomLink(baseUrl url.URL, subdirectories, variablesCount int, linkHasVariablesProbability float64) string {
	subDirectoryPath := generateSubDirectories(subdirectories)

	variables := generateVariables(variablesCount, linkHasVariablesProbability)

	if variables != "" {
		return fmt.Sprintf("%s://%s/%s?%s", baseUrl.Scheme, baseUrl.Host, subDirectoryPath, variables)
	}
	return fmt.Sprintf("%s://%s/%s", baseUrl.Scheme, baseUrl.Host, subDirectoryPath)
}

// RandomSimpleLink generates a random link based on the given base URL.
func RandomSimpleLink(baseUrl url.URL) string {
	name := strings.ToLower(strings.Join([]string{
		functions.PickRandomFromSlice(&dictionaries.Nouns),
		functions.PickRandomFromSlice(&dictionaries.Nouns),
	}, "-"))
	return fmt.Sprintf("%s://%s/%s", baseUrl.Scheme, baseUrl.Host, name)
}
