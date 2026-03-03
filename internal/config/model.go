package config

import (
	"fmt"
	"os"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/ui"
)

func RunSetupWizard(isGlobal bool) error {
	result, err := ui.StartWizard(isGlobal)
	if err != nil {
		return err
	}
	fmt.Printf("Is Global config : %v\n", isGlobal)

	commitMaxLength, err := strconv.Atoi(result.MaxLength)
	if err != nil {
		return fmt.Errorf("invalid max length : %w", err)
	}

	cfg := Config{
		ModelName:       result.ModelName,
		Language:        result.Language,
		CommitMaxLength: commitMaxLength,
		PromptAddition:  result.PromptAddition,
	}

	err = SaveConfig(cfg, isGlobal)
	fmt.Println("Setup complete!")
	return nil
}

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

	content := fmt.Sprintf(
		"%s %s\n%s %s\n%s %v\n%s %s",
		keyStyle.Render("Modèle:"), valStyle.Render(cfg.ModelName),
		keyStyle.Render("Langue:"), valStyle.Render(cfg.Language),
		keyStyle.Render("Max Len:"), valStyle.Render(fmt.Sprintf("%v", cfg.CommitMaxLength)),
		keyStyle.Render("Prompt:"), valStyle.Render(cfg.PromptAddition),
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

	fmt.Println(borderStyle.Render(content))
}
