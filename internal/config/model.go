package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

func PrintConfig(cfg Config, verbose bool) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginTop(1)
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n", keyStyle.Render("Provider:"), valStyle.Render(cfg.ProviderName))
	fmt.Fprintf(&b, "%s %s\n", keyStyle.Render("Model:"), valStyle.Render(cfg.ModelName))
	fmt.Fprintf(
		&b,
		"%s %s\n",
		keyStyle.Render("Commit Max Length:"),
		valStyle.Render(cfg.MaxLength),
	)
	fmt.Fprintf(&b, "%s %s\n", keyStyle.Render("Language:"), valStyle.Render(cfg.Language))
	fmt.Fprintf(
		&b,
		"%s %s",
		keyStyle.Render("Instructions:"),
		valStyle.Render(cfg.CustomInstructions),
	)

	fmt.Println(titleStyle.Render(" ComAI Generate"))
	if verbose {
		usedFile := viper.ConfigFileUsed()
		if usedFile != "" {
			fmt.Fprintln(os.Stderr, subtitleStyle.Render(" Using config file:", usedFile))
		} else {
			fmt.Fprintln(
				os.Stderr,
				subtitleStyle.Render(" No config file found, using defaults or env vars"),
			)
		}
	}

	fmt.Println(borderStyle.Render(b.String()))
}
