package wizard

type Result struct {
	ProviderName       string
	ModelName          string
	Language           string
	CustomInstructions string
}

var modelsByProvider = map[string][]string{
	"ollama": {"gpt-oss:20b", "qwen3.5:35b", "qwen3.5:27b", "qwen3.5:9b", "qwen3.5:2b", "qwen3.5:0.8b", "qwen3.5:4b", "qwen3:14b", "qwen3:8b", "qwen3:4b"},
}
