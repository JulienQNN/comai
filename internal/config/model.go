package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/theme"
)

func PrintConfig(cfg Config, verbose bool) {
	t := theme.Default()

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"%s %s\n",
		t.ConfigKey.Render("Provider:"),
		t.ConfigValue.Render(cfg.ProviderName),
	)
	fmt.Fprintf(&b, "%s %s\n", t.ConfigKey.Render("Model:"), t.ConfigValue.Render(cfg.ModelName))
	fmt.Fprintf(
		&b,
		"%s %s\n",
		t.ConfigKey.Render("Commit Max Length:"),
		t.ConfigValue.Render(cfg.MaxLength),
	)
	fmt.Fprintf(&b, "%s %s\n", t.ConfigKey.Render("Language:"), t.ConfigValue.Render(cfg.Language))
	fmt.Fprintf(
		&b,
		"%s %s",
		t.ConfigKey.Render("Instructions:"),
		t.ConfigValue.Render(cfg.CustomInstructions),
	)

	fmt.Println(t.ConfigTitle.Render(" ComAI Generate"))
	if verbose {
		subtitleStyle := t.Muted
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

	fmt.Println(t.ConfigBorder.Render(b.String()))
}
