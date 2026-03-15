package generate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"

	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider"
	"github.com/JulienQNN/comai/internal/theme"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.callLLM())
}

func (m model) callLLM() tea.Cmd {
	return func() tea.Msg {
		ch := make(chan string, 64)
		errCh := make(chan error, 1)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			defer cancel()
			defer close(ch)
			err := m.provider.Stream(ctx, m.params, ch)
			errCh <- err
			close(errCh)
		}()

		return llmStreamStartMsg{ch: ch, errCh: errCh}
	}
}

func waitForToken(ch <-chan string, errCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		token, ok := <-ch
		if !ok {
			return llmDoneMsg{err: <-errCh}
		}
		return llmTokenMsg(token)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case llmStreamStartMsg:
		m.tokenCh = msg.ch
		m.errCh = msg.errCh
		return m, waitForToken(m.tokenCh, m.errCh)
	case llmTokenMsg:
		m.partial += string(msg)
		return m, waitForToken(m.tokenCh, m.errCh)
	case llmDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	elapsed := time.Since(m.start).Truncate(100 * time.Millisecond)

	if m.done {
		return tea.NewView("")
	}

	status := fmt.Sprintf("%s Generating commit message... %s", m.spinner.View(), elapsed.String())
	if m.partial != "" {
		return tea.NewView(status + "\n" + m.partial)
	}

	return tea.NewView(status)
}

// Start launches the streaming TUI and calls the LLM provider.
func Start(p provider.Provider, params prompt.CompletionParams) (Result, error) {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = theme.Default().Spinner
	m := model{
		spinner:  s,
		start:    time.Now(),
		provider: p,
		params:   params,
	}

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return Result{}, err
	}

	gm, ok := result.(model)
	if !ok {
		return Result{}, fmt.Errorf("unexpected result type from TUI program")
	}

	if gm.err != nil {
		return Result{}, gm.err
	}

	return Result{
		CommitMsg: strings.TrimSpace(gm.partial),
		Elapsed:   time.Since(m.start),
	}, nil
}
