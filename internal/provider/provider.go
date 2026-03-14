package provider

import (
	"context"

	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider/copilot"
	"github.com/JulienQNN/comai/internal/provider/ollama"
)

type Provider interface {
	Stream(ctx context.Context, params prompt.CompletionParams, ch chan<- string) error
	Close() error
}

func New(providerName, model string) (Provider, error) {
	switch providerName {
	case "copilot":
		return copilot.New(model)
	case "ollama":
		return ollama.New(model)
	default:
		return ollama.New(model)
	}
}
