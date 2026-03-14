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
	partial  string // accumulates streaming tokens
	err      error
	tokenCh  <-chan string
	errCh    <-chan error
	provider provider.Provider
	params   prompt.CompletionParams
}

type llmStreamStartMsg struct {
	ch    <-chan string
	errCh <-chan error
}

type llmTokenMsg string

type llmDoneMsg struct{ err error }
