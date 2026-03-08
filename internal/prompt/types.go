package prompt

type CompletionParams struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
}
