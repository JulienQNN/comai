package theme

import (
	"charm.land/huh/v2"
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

func FormhuhTheme() huh.ThemeFunc {
	t := huh.ThemeBase(false)
	lightDark := lipgloss.LightDark(true)

	var (
		normalFg = lightDark(lipgloss.Color("252"), lipgloss.Color("235"))
		fuchsia  = lipgloss.Color("#F780E2")
		green    = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
		red      = lightDark(lipgloss.Color("#FF4672"), lipgloss.Color("#ED567A"))
	)
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
		t.Focused.Card = t.Focused.Base
		t.Focused.Title = t.Focused.Title.Foreground(Secondary).Bold(false)
		t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(Secondary).Bold(false).MarginBottom(1)
		t.Focused.Directory = t.Focused.Directory.Foreground(Secondary)
		t.Focused.Description = t.Focused.Description.Foreground(
			lightDark(lipgloss.Color(""), lipgloss.Color("243")),
		)
		t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
		t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
		t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(fuchsia)
		t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(fuchsia)
		t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(fuchsia)
		t.Focused.Option = t.Focused.Option.Foreground(normalFg)
		t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(fuchsia)
		t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
		t.Focused.SelectedPrefix = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#02CF92"), lipgloss.Color("#02A877"))).
			SetString("✓ ")
		t.Focused.UnselectedPrefix = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color(""), lipgloss.Color("243"))).
			SetString("• ")
		t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)

		t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
		t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(
			lightDark(lipgloss.Color("248"), lipgloss.Color("238")),
		)

		t.Blurred.NextIndicator = lipgloss.NewStyle()
		t.Blurred.PrevIndicator = lipgloss.NewStyle()

		t.Group.Title = t.Focused.Title
		t.Group.Description = t.Focused.Description
		t.Focused.Base = t.Focused.Base.MarginTop(1).MarginBottom(0)
		t.Blurred.Base = t.Blurred.Base.MarginTop(1).MarginBottom(0)
		t.FieldSeparator = lipgloss.NewStyle().SetString(" ")
		return t
	})
}
