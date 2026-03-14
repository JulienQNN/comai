package wizard

type Result struct {
	ProviderName       string
	ModelName          string
	Language           string
	MaxLength          string
	CustomInstructions string
}

var modelsByProvider = map[string][]string{
	"ollama": {
		"gemma3:1b",
		"gemma3:4b",
		"gemma3:12b",
		"gemma3:27b",
		"qwen3.5:35b",
		"qwen3.5:27b",
		"qwen3.5:9b",
		"qwen3.5:2b",
		"qwen3.5:0.8b",
		"qwen3.5:4b",
		"qwen3:14b",
		"qwen3:8b",
		"qwen3:4b",
	},
	"copilot": {
		"gpt-4.1",
		"gpt-5-mini",
		"gpt-5.1",
		"gpt-5.1-codex",
		"gpt-5.1-codex-mini",
		"gpt-5.1-codex-max",
		"gpt-5.2",
		"gpt-5.2-codex",
		"gpt-5.3-codex",
		"gpt-5.4",
		"claude-haiku-4.5",
		"claude-opus-4.5",
		"claude-opus-4.6",
		"claude-opus-4.6-fast-preview",
		"claude-sonnet-4",
		"claude-sonnet-4.5",
		"claude-sonnet-4.6",
		"gemini-2.5-pro",
		"gemini-3-flash",
		"gemini-3-pro",
		"gemini-3.1-pro",
		"grok-code-fast-1",
		"raptor-mini",
		"goldeneye",
	},
}
