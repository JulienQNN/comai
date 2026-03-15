package wizard

import (
	"context"
	"time"

	"charm.land/huh/v2"

	copilotprovider "github.com/JulienQNN/comai/internal/provider/copilot"
	"github.com/JulienQNN/comai/internal/provider/ollama"
)

func listCopilotModels() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := copilotprovider.New("")
	if err != nil {
		return fallbackCopilotModels
	}

	defer func() {
		if closeErr := c.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	models, err := c.ListModels(ctx)
	if err != nil || len(models) == 0 {
		return fallbackCopilotModels
	}
	return models
}

func Start(isGlobal bool) (Result, error) {
	var result Result

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	copilotModels := fallbackCopilotModels
	copilotCh := make(chan []string, 1)
	go func() {
		select {
		case copilotCh <- listCopilotModels():
		case <-ctx.Done():
		}
	}()

	getModels := func(provider string) []string {
		if provider == "copilot" {
			select {
			case models := <-copilotCh:
				copilotModels = models
			case <-time.After(5 * time.Second):
			}
			return copilotModels
		}
		return ollama.RecommendedModels
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Provider").
				Options(huh.NewOptions("ollama", "copilot")...).
				Value(&result.ProviderName),

			huh.NewSelect[string]().
				Title("Model").
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(getModels(result.ProviderName)...)
				}, &result.ProviderName).
				Value(&result.ModelName),

			huh.NewInput().
				Title("Language").
				Placeholder("en (optional)").
				Value(&result.Language),

			huh.NewInput().
				Title("Commit Max Length").
				Placeholder("50 (optional)").
				Value(&result.MaxLength),

			huh.NewInput().
				Title("Custom instructions").
				Placeholder("e.g. use conventional commits (optional)").
				Value(&result.CustomInstructions),
		),
	)

	if err := form.Run(); err != nil {
		return Result{}, err
	}

	if result.Language == "" {
		result.Language = "en"
	}

	if result.MaxLength == "" {
		result.MaxLength = "50"
	}

	return result, nil
}
