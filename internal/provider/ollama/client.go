package ollama

import (
	"context"
	"fmt"
	"strings"

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

func (c *Client) Complete(ctx context.Context, params prompt.CompletionParams) (string, error) {
	think := api.ThinkValue{Value: false}
	stream := false

	req := &api.ChatRequest{
		Model:  c.model,
		Think:  &think,
		Stream: &stream,
		Options: map[string]any{
			"num_ctx":     2048,
			"num_predict": params.MaxTokens,
			"temperature": 0.2,
			"seed":        42,
		},
		Messages: []api.Message{
			{Role: "system", Content: params.SystemPrompt},
			{Role: "user", Content: params.UserPrompt},
		},
	}

	var result strings.Builder
	err := c.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		result.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(result.String())
	if content == "" {
		return "", fmt.Errorf(
			"ollama returned empty content (model: %s) — check that thinking is disabled and num_predict is sufficient",
			c.model,
		)
	}
	return content, nil
}
