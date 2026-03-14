package copilot

import (
	"context"
	"fmt"

	copilot "github.com/github/copilot-sdk/go"

	"github.com/JulienQNN/comai/internal/prompt"
)

func New(model string) (*Client, error) {
	if model == "" {
		model = defaultModel
	}

	c := copilot.NewClient(&copilot.ClientOptions{
		CLIPath:         "gh",
		CLIArgs:         []string{"copilot"},
		UseLoggedInUser: new(true),
	})
	if err := c.Start(context.Background()); err != nil {
		return nil, fmt.Errorf(
			"copilot client start (is 'gh' installed and authenticated?): %w",
			err,
		)
	}

	return &Client{model: model, client: c}, nil
}

func (c *Client) Stream(
	ctx context.Context,
	params prompt.CompletionParams,
	ch chan<- string,
) error {
	session, err := c.client.CreateSession(ctx, &copilot.SessionConfig{
		Model:               c.model,
		ClientName:          "comai",
		Streaming:           true,
		OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
		SystemMessage: &copilot.SystemMessageConfig{
			Mode:    "replace",
			Content: params.SystemPrompt,
		},
	})
	if err != nil {
		return fmt.Errorf("copilot create session: %w", err)
	}

	defer func() {
		dErr := session.Disconnect()
		if dErr != nil && err == nil {
			err = fmt.Errorf("copilot disconnect session: %w", dErr)
		}
	}()

	done := make(chan error, 1)

	session.On(func(event copilot.SessionEvent) {
		switch event.Type {
		case copilot.AssistantMessageDelta:
			if event.Data.DeltaContent != nil {
				ch <- *event.Data.DeltaContent
			}
		case copilot.SessionIdle:
			done <- nil
		}
	})

	if _, err := session.Send(ctx, copilot.MessageOptions{Prompt: params.UserPrompt}); err != nil {
		return fmt.Errorf("copilot send message: %w", err)
	}

	return <-done
}

func (c *Client) Close() error {
	return c.client.Stop()
}
