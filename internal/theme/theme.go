package theme

import (
	"charm.land/lipgloss/v2"
)

var (
	Primary   = lipgloss.Color("205") // pink/purple — brand color
	Secondary = lipgloss.Color("62")  // purple — accents
	Muted     = lipgloss.Color("240") // gray — secondary text
	White     = lipgloss.Color("15")  // white — values
)

type Theme struct {
	Title        lipgloss.Style
	CommitBorder lipgloss.Style
	Italic       lipgloss.Style
	Muted        lipgloss.Style
	MutedItalic  lipgloss.Style

	ConfigBorder lipgloss.Style
	ConfigKey    lipgloss.Style
	ConfigValue  lipgloss.Style

	Spinner lipgloss.Style
}

func Default() Theme {
	baseLeftBlock := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		PaddingLeft(1)

	baseText := lipgloss.NewStyle().Foreground(Muted)

	return Theme{
		// commits
		CommitBorder: baseLeftBlock.Copy().BorderForeground(Primary),
		// config
		ConfigBorder: baseLeftBlock.Copy().BorderForeground(Secondary),
		ConfigKey:    baseText.Copy(),
		ConfigValue:  lipgloss.NewStyle().Foreground(White),
		// shared styles
		Title:       lipgloss.NewStyle().Bold(true).Foreground(Primary).MarginTop(1),
		Italic:      baseText.Copy().Italic(true),
		Muted:       baseText.Copy(),
		MutedItalic: baseText.Copy().Italic(true),
		Spinner:     lipgloss.NewStyle().Foreground(Primary),
	}
}
