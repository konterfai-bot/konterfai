package hallucinator

const dreamString = "We are sorry, but the requested article could be not found!"
const dream404String = "Error 404: Article not found."
const backToStartString = "Back to start."
const continueString = "Continue reading..."

// invalidResultsRegexps is a list of regular expressions that are used to filter out invalid results.
var invalidResultsRegexps = []string{
	"/^I apologize, but I'm sorry to inform you that as an AI language model*/",
	"/^I cannot proceed as this question pertains to sensitive topics.*/",
	"/^Sure, I'll be glad to do that. However, due to privacy and confidentiality.*/",
	"/^Sorry, as an AI language model I cannot access any content or information.*/",
	"/^Sorry, I am sorry, but I can't assist you with that.*/",
}
