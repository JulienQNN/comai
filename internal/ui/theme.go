package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title       lipgloss.Style
	ErrorHeader lipgloss.Style
	Info        lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			MarginTop(1).
			MarginBottom(1).
			Transform(strings.ToUpper),

		ErrorHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("204")).
			Bold(true).
			Padding(0, 1).
			SetString(" ERROR "),

		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true),
	}
}

