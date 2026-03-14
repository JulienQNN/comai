package ollama

import (
	"context"

	"github.com/ollama/ollama/api"

	"github.com/JulienQNN/comai/internal/prompt"
)

func New(model string) (*Client, error) {
	c, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &Client{model: model, client: c}, nil
}

func (c *Client) Close() error { return nil }

func (c *Client) Stream(
	ctx context.Context,
	params prompt.CompletionParams,
	ch chan<- string,
) error {
	req := &api.ChatRequest{
		Model:  c.model,
		Think:  &api.ThinkValue{Value: false},
		Stream: new(true),
		Options: map[string]any{
			"num_ctx":     defaultContextSize,
			"num_predict": params.MaxTokens,
			"temperature": defaultTemperature,
			"seed":        defaultSeed,
		},
		Messages: []api.Message{
			{Role: "system", Content: params.SystemPrompt},
			{Role: "user", Content: params.UserPrompt},
		},
	}

	return c.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if resp.Message.Content != "" {
			ch <- resp.Message.Content
		}
		return nil
	})
}
