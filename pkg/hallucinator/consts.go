package hallucinator

const dreamString = "We are sorry, but the requested article could be not found!"
const backToStartString = "Back to start."
const continueString = "Continue reading..."

// invalidResultsRegexps is a list of regular expressions that are used to filter out invalid results.
var invalidResultsRegexps = []string{
	"/^I apologize, but I'm sorry to inform you that as an AI language model,*/",
	"/^I cannot proceed as this question pertains to sensitive topics .*/",
}
