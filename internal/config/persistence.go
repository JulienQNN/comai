package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func SaveConfig(cfg Config, isGlobal bool) error {
	var configPath string
	if isGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(home, ".comai.yaml")
	} else {
		configPath = ".comai.yaml"
	}

	v := viper.New()
	v.Set("provider", cfg.ProviderName)
	v.Set("model", cfg.ModelName)
	v.Set("language", cfg.Language)
	v.Set("max_length", cfg.MaxLength)
	v.Set("custom_instructions", cfg.CustomInstructions)

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	fmt.Println("Config file saved at :", configPath)
	return nil
}
