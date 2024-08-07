package links

import (
	"fmt"
	"math/rand"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/types"
)

// generateVariables generates a random number of variables.
func generateVariables(variablesCount int, linkHasVariablesProbability float64) string {
	variables := []string{}
	variablesValue := []string{}
	varcount := rand.Intn(variablesCount) + 1
	for i := 0; i < varcount; i++ {
		variables = append(variables, getVariableNameString())
		variablesValue = append(variablesValue, getVariableValueString())
	}

	variablesString := ""
	if rand.Float64() < linkHasVariablesProbability {
		for i := 0; i < len(variables); i++ {
			if variablesString != "" {
				variablesString = fmt.Sprintf("%s&%s=%s", variablesString, variables[i], variablesValue[i])
			} else {
				variablesString = fmt.Sprintf("%s=%s", variables[i], variablesValue[i])
			}
		}
	}
	return variablesString
}

// getVariableNameString returns a random variable name string.
func getVariableNameString() string {
	typeRand := rand.Intn(types.VariableNamesCount) + 1
	switch types.VariableNames(typeRand) {
	case types.SingleCharacterVariable:
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		return string(charset[rand.Intn(len(charset))])
	case types.VerbVariable:
		return functions.PickRandomStringFromSlice(&dictionaries.Verbs)
	case types.NounVariable:
		return functions.PickRandomStringFromSlice(&dictionaries.Nouns)
	case types.TwoNounVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.ThreeNounVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.TwoNounDashedVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.ThreeNounDashedVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.VerbNounCombinationVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Verbs),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.NounVerbCombinationVariable:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			strings.ToUpper(functions.PickRandomStringFromSlice(&dictionaries.Verbs)),
		}, "")
	case types.MonthVariable:
		return functions.PickRandomStringFromSlice(&dictionaries.Months)
	case types.DayVariable:
		return functions.PickRandomStringFromSlice(&dictionaries.Weekdays)
	default:
		return "default"
	}
}

// getVariableValueString returns a random variable value string.
func getVariableValueString() string {
	typeRand := rand.Intn(types.VariableValuesCount) + 1
	switch types.VariableValues(typeRand) {
	case types.VerbValue:
		return functions.PickRandomStringFromSlice(&dictionaries.Verbs)
	case types.NounValue:
		return functions.PickRandomStringFromSlice(&dictionaries.Nouns)
	case types.TwoNounValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.ThreeNounValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.TwoNounDashedValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.ThreeNounDashedValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "-")
	case types.VerbNounCombinationValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Verbs),
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
		}, "")
	case types.NounVerbCombinationValue:
		return strings.Join([]string{
			functions.PickRandomStringFromSlice(&dictionaries.Nouns),
			strings.ToUpper(functions.PickRandomStringFromSlice(&dictionaries.Verbs)),
		}, "")
	case types.DateValue:
		return functions.PickRandomDate()
	case types.YearValue:
		return functions.PickRandomYear()
	case types.MonthValue:
		return functions.PickRandomStringFromSlice(&dictionaries.Months)
	case types.DayValue:
		return functions.PickRandomStringFromSlice(&dictionaries.Weekdays)
	case types.Base64Value:
		return functions.RandomBase64String()
	default:
		return "default"
	}
}
