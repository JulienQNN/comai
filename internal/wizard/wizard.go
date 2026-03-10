package wizard

import (
	"github.com/charmbracelet/huh"
)

func Start(isGlobal bool) (Result, error) {
	var result Result

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Provider").
				Options(huh.NewOptions("ollama")...).
				Value(&result.ProviderName),

			huh.NewSelect[string]().
				Title("Model").
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(modelsByProvider[result.ProviderName]...)
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

	if result.CustomInstructions == "" {
		result.CustomInstructions = "<type>(<optional scope>): <description>"
	}

	return result, nil
}
