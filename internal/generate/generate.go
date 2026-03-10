package generate

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.callLLM())
}

func (m model) callLLM() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		msg, err := m.provider.Complete(ctx, m.params)
		return llmResponseMsg{msg: msg, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case llmResponseMsg:
		m.done = true
		m.result = msg.msg
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	elapsed := time.Since(m.start).Truncate(100 * time.Millisecond)
	timerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	if m.done {
		if m.err != nil {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			return fmt.Sprintf("%s %s\n", errStyle.Render("✗ Error:"), m.err.Error())
		}
		return fmt.Sprintf("Generated in %s\n", timerStyle.Render(elapsed.String()))
	}

	return fmt.Sprintf("%s Generating commit message... %s",
		m.spinner.View(), timerStyle.Render(elapsed.String()))
}

// Start launches the spinner TUI and calls the LLM.
func Start(p provider.Provider, params prompt.CompletionParams) (Result, error) {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

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
		return Result{}, fmt.Errorf("unexpected internal state")
	}
	return Result{
		CommitMsg: gm.result,
		Elapsed:   time.Since(m.start),
		Err:       gm.err,
	}, nil
}
