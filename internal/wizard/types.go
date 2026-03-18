package wizard

type Result struct {
	ProviderName       string
	ModelName          string
	Language           string
	CommitMaxLength    string
	CustomInstructions string
}

var fallbackCopilotModels = []string{
	"claude-sonnet-4.5",
	"claude-3.5-haiku",
	"gpt-4o",
	"gpt-4o-mini",
}
