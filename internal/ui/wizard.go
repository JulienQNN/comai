package ui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("196"))
)

func StartWizard(isGlobal bool) (WizardResult, error) {
	p := tea.NewProgram(InitialModel(isGlobal))
	m, err := p.Run()
	if err != nil {
		return WizardResult{}, err
	}

	model := m.(model)

	return WizardResult{
		ModelName:      model.inputs[0].Value(),
		Language:       model.inputs[1].Value(),
		MaxLength:      model.inputs[2].Value(),
		PromptAddition: model.inputs[3].Value(),
	}, nil
}

func InitialModel(isGlobal bool) model {
	m := model{
		inputs:   make([]textinput.Model, 4),
		isGlobal: isGlobal,
	}

	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 64
		t.SetWidth(64)
		s := t.Styles()
		s.Cursor.Color = lipgloss.Color("202")
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)

		switch i {
		case 0:
			t.Placeholder = "Model name (e.g. gemini-3.1-pro-preview)"
			t.Focus()
		case 1:
			t.Placeholder = "Language to use (e.g. en, fr, es)"
			t.CharLimit = 2
		case 2:
			t.Placeholder = "Maximum length of the commit message by characters (e.g. 200)"
			t.CharLimit = 3
			t.Validate = func(s string) error {
				if s == "" {
					return nil
				}
				if _, err := strconv.Atoi(s); err != nil {
					return fmt.Errorf("max length must be a number")
				}
				return nil
			}
		case 3:
			t.Placeholder = "Any addition to the prompt ? (e.g. 'Write a short commit with sementic commit message style')"
			t.CharLimit = 600
			t.SetWidth(600)
		}
		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				hasErrors := false
				allFilled := true
				for _, in := range m.inputs {
					if in.Err != nil {
						hasErrors = true
						break
					}
					if in.Value() == "" {
						allFilled = false
						break
					}
				}
				if !hasErrors && allFilled {
					return m, tea.Quit
				}
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i, in := range m.inputs {
		inputView := in.View()
		if in.Err != nil {
			inputView += "\n" + errorStyle.Render(fmt.Sprintf("%s", in.Err.Error()))
		}

		b.WriteString(inputView)

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
		if m.cursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	label := "Save Local Config"
	if m.isGlobal {
		label = "Save Global Config"
	}
	focusedButton := fmt.Sprintf("[ %s ]", focusedStyle.Render(label))
	blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render(label))
	button := blurredButton

	if m.focusIndex == len(m.inputs) {
		hasErrors := false
		allFilled := true
		for _, in := range m.inputs {
			if in.Err != nil {
				hasErrors = true
				break
			}
			if in.Value() == "" {
				allFilled = false
				break
			}
		}
		if !hasErrors && allFilled {
			button = focusedButton
		}
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)
	if m.quitting {
		b.WriteRune('\n')
	}

	v := tea.NewView(b.String())
	v.Cursor = c

	return v
}
