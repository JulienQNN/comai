package config

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/theme"
)

func PrintConfig(cfg Config, verbose bool) {
	t := theme.Default()
	alignedKey := t.ConfigKey.Align(lipgloss.Right)
	valStyle := t.ConfigValue.MarginLeft(1).Width(55)
	var lines []string
	addLine := func(label string, value any) {
		strValue := fmt.Sprint(value)
		line := lipgloss.JoinHorizontal(lipgloss.Top, alignedKey.Render(label),
			valStyle.Render(strValue),
		)
		lines = append(lines, line)
	}
	addLine("Provider:", cfg.ProviderName)
	addLine("Model:", cfg.ModelName)
	addLine("Language:", cfg.Language)
	addLine("Max Commit Length:", cfg.CommitMaxLength)

	if cfg.CustomInstructions != "" {
		addLine("Custom Instructions:", cfg.CustomInstructions)
	}

	configContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	fmt.Println(t.Title.Render("ComAI Generate"))
	if verbose {
		subtitleStyle := t.MutedItalic
		usedFile := viper.ConfigFileUsed()
		if usedFile != "" {
			fmt.Fprintln(os.Stderr, subtitleStyle.Render("Config file:", usedFile))
			fmt.Println(t.ConfigBorder.Render(configContent))
		} else {
			fmt.Fprintln(
				os.Stderr,
				subtitleStyle.Render(" No config file found, using defaults or env vars"),
			)
		}
	}
}
