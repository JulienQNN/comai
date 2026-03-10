package prompt

import (
	"fmt"

	"github.com/JulienQNN/comai/internal/config"
)

func Build(diff string, cfg config.Config) CompletionParams {
	const maxDiffLen = 4000
	if len(diff) > maxDiffLen {
		diff = diff[:maxDiffLen] + "\n...(truncated)"
	}

	system := fmt.Sprintf(
		"Output ONLY a git commit message in lowercase. MaxLength: %v Language: %s.",
		cfg.MaxLength, cfg.Language)
	if cfg.CustomInstructions != "" {
		system += " " + cfg.CustomInstructions
	}

	return CompletionParams{
		SystemPrompt: system,
		UserPrompt:   "Diff:\n" + diff,
		MaxTokens:    128,
	}
}
