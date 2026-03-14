package ollama

import "github.com/ollama/ollama/api"

// Default inference parameters.
const (
	defaultContextSize  = 2048
	defaultTemperature  = 0.2
	defaultSeed         = 42
)

type Client struct {
	model  string
	client *api.Client
}
