package generate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/git"
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
func Start(t theme.Theme, cfg config.Config, date string, dateInteractive bool) error {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = theme.Default().Spinner

	p, err := provider.New(cfg.ProviderName, cfg.ModelName)
	defer func() {
		if closeErr := p.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("initializing provider: %w", err)
	}

	params := prompt.Build(diff.RawDiff, cfg)
	m := model{
		spinner:  s,
		start:    time.Now(),
		provider: p,
		params:   params,
	}

	for _, f := range diff.Files {
		fmt.Println(t.RenderFileChange(f.Status, f.Path))
	}

	output := lipgloss.JoinVertical(lipgloss.Left, diff.Stats)
	fmt.Println(output)

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return err
	}

	gm, ok := result.(model)
	if !ok {
		return fmt.Errorf("unexpected result type from TUI program")
	}
	if gm.err != nil {
		return gm.err
	}

	author, err := git.GetAuthorInfo()
	if err != nil {
		return fmt.Errorf("getting author info: %w", err)
	}

	partial := strings.TrimSpace(gm.partial)
	commitMsg := strings.ToLower(partial)

	elasped := time.Since(m.start)

	formattedDate, _ := git.FormatDate(date)
	renderPreview(t, cfg, commitMsg, author, formattedDate, elasped)

	confirmed := true
	form := huh.NewForm(
		huh.NewGroup(huh.NewConfirm().
			Affirmative("Commit").
			Negative("Cancel").
			Value(&confirmed),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("getting confirmation: %w", err)
	}

	if !confirmed {
		fmt.Println(t.Muted.Render("Commit cancelled."))
		return nil
	}

	if dateInteractive && date == "" {
		if err := huh.NewForm(huh.NewGroup(huh.NewInput().Title("Commit date").
			Placeholder("e.g. yesterday, 2024-01-01, last friday").
			Value(&date),
		),
		).Run(); err != nil {
			return fmt.Errorf("getting commit date: %w", err)
		}
		if err != nil {
			return err
		}
	}

	if err := git.Commit(commitMsg, git.CommitOptions{Date: formattedDate}); err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}

	return nil
}
