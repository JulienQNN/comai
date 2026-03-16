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

	BadgeMod lipgloss.Style
	BadgeAdd lipgloss.Style
	BadgeDel lipgloss.Style
	BadgeUnk lipgloss.Style

	FileDir  lipgloss.Style
	FileName lipgloss.Style

	Spinner lipgloss.Style
}

func Default() Theme {
	baseText := lipgloss.NewStyle().Foreground(Muted)
	baseLeftBlock := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		PaddingLeft(1)

	return Theme{
		// commits
		CommitBorder: baseLeftBlock.BorderForeground(Primary),
		// config
		ConfigBorder: baseLeftBlock.BorderForeground(Secondary),
		ConfigKey:    baseText,
		ConfigValue:  lipgloss.NewStyle().Foreground(White),
		// diff
		BadgeMod: lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")),
		BadgeAdd: lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")),
		BadgeDel: lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")),
		BadgeUnk: lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")),

		FileDir:  baseText,
		FileName: lipgloss.NewStyle().Foreground(White),
		// shared styles
		Title:       lipgloss.NewStyle().Bold(true).Foreground(Primary).MarginTop(1),
		Italic:      baseText.Italic(true),
		Muted:       baseText,
		MutedItalic: baseText.Italic(true),
		Spinner:     lipgloss.NewStyle().Foreground(Primary),
	}
}
