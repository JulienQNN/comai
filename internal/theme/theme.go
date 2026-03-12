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
	CommitTitle  lipgloss.Style
	CommitBorder lipgloss.Style
	Italic       lipgloss.Style
	Muted        lipgloss.Style

	ConfigTitle  lipgloss.Style
	ConfigBorder lipgloss.Style
	ConfigKey    lipgloss.Style
	ConfigValue  lipgloss.Style

	Spinner lipgloss.Style
}

func Default() Theme {
	return Theme{
		CommitTitle: lipgloss.NewStyle().Bold(true).Foreground(Primary),
		CommitBorder: lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(Primary).
			PaddingLeft(1),
		Italic: lipgloss.NewStyle().Italic(true).Foreground(Muted),
		Muted:  lipgloss.NewStyle().Foreground(Muted),

		ConfigTitle: lipgloss.NewStyle().Bold(true).Foreground(Primary).MarginTop(1),
		ConfigBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1),
		ConfigKey:   lipgloss.NewStyle().Foreground(Muted),
		ConfigValue: lipgloss.NewStyle().Foreground(White),

		Spinner: lipgloss.NewStyle().Foreground(Primary),
	}
}
