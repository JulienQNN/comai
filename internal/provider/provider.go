package provider

import (
	"context"

	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider/ollama"
)

type Provider interface {
	Complete(ctx context.Context, params prompt.CompletionParams) (string, error)
}

func New(model string) (Provider, error) {
	return ollama.New(model)
}
