package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/wizard"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize comai configuration interactively",
	Run: func(cmd *cobra.Command, args []string) {
		isGlobal, _ := cmd.Flags().GetBool("global")

		result, err := wizard.Start(isGlobal)
		if err != nil {
			log.Fatal("Error during configuration wizard", "err", err)
		}

		cfg := config.Config{
			ProviderName:       result.ProviderName,
			ModelName:          result.ModelName,
			Language:           result.Language,
			MaxLength:          result.MaxLength,
			CustomInstructions: result.CustomInstructions,
		}

		if err := config.SaveConfig(cfg, isGlobal); err != nil {
			log.Fatal("Error saving configuration", "err", err)
		}

		//TODO print success message with config path
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().
		BoolP("global", "g", false, "Generate the comai.yaml file in the global configuration directory.")
}
