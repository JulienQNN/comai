package generate

import (
	"time"

	"charm.land/bubbles/v2/spinner"

	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider"
)

type Result struct {
	CommitMsg string
	Elapsed   time.Duration
	Err       error
}

type model struct {
	spinner  spinner.Model
	start    time.Time
	done     bool
	result   string
	err      error
	provider provider.Provider
	params   prompt.CompletionParams
}

type llmResponseMsg struct {
	msg string
	err error
}
