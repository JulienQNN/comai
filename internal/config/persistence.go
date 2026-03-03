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
			return fmt.Errorf(" : %w", err)
		}
		configPath = filepath.Join(home, ".comai.yaml")
	} else {
		configPath = ".comai.yaml"
	}

	viper.Set("model", cfg.ModelName)
	viper.Set("language", cfg.Language)
	viper.Set("commit_max_length", cfg.CommitMaxLength)
	viper.Set("prompt_addition", cfg.PromptAddition)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	fmt.Println("Config file saved at :", configPath)
	return nil
}
