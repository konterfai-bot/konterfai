package hallucinator

const (
	BackToStartString = "Back to start."
	ContinueString    = "Continue reading..."
	Dream404String    = "Error 404: Article not found."
	DreamString       = "We are sorry, but the requested article could be not found!"
)

// invalidResultsRegexps is a list of regular expressions that are used to filter out invalid results.
var invalidResultsRegexps = []string{
	"/^I apologize, but I'm sorry to inform you that as an AI language model*/",
	"/^I cannot proceed as this question pertains to sensitive topics.*/",
	"/^Sorry, I am sorry, but I can't assist you with that.*/",
	"/^Sorry, I cannot fulfill that request as it is out of scope for me as an AI assistant.*/",
	"/^Sorry, as an AI language model I cannot access any content or information.*/",
	"/^Sorry, as an AI, I'm currently not able to create a written text.*/",
	"/^Sorry, but I can't assist with that.*/",
	"/^Sure, I'll be glad to do that. However, due to privacy and confidentiality.*/",
}
