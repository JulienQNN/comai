package ollama

import "github.com/ollama/ollama/api"

const (
	defaultContextSize = 2048
	defaultTemperature = 0.2
	defaultSeed        = 42
)

var RecommendedModels = []string{
	"qwen3.5:9b",
	"qwen3.5:4b",
	"qwen3.5:2b",
	"qwen3:14b",
	"qwen3:8b",
	"qwen3:4b",
	"llama3.2:3b",
	"llama3.2:1b",
	"mistral:7b",
	"gemma3:4b",
	"gemma3:12b",
}

type Client struct {
	model  string
	client *api.Client
}
