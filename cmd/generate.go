package cmd

import (
	"os"
	"path/filepath"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/generate"
	"github.com/JulienQNN/comai/internal/theme"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a commit message using AI models.",
	Long:  `The generate command analyzes the git diff in your git repository and generates a commit message using AI models.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		isGlobal, _ := cmd.Flags().GetBool("global")
		dateFlag, _ := cmd.Flags().GetString("date")
		dateInteractive, _ := cmd.Flags().GetBool("date-interactive")
		t := theme.Default()
		var cfg config.Config

		if isGlobal {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Error("getting home directory", "err", err)
				return
			}
			viper.SetConfigFile(filepath.Join(home, ".comai.yaml"))
			if err := viper.ReadInConfig(); err != nil {
				log.Error("getting home directory", "err", err)
				return
			}
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			log.Error("loading config", "err", err)
			return
		}

		config.PrintConfig(cfg, verbose)
		err := generate.Start(t, cfg, dateFlag, dateInteractive)
		if err != nil {
			log.Error("generating commit message", "err", err)
			os.Exit(1)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().
		BoolP("verbose", "v", false, "Show verbose output including config file path")
	generateCmd.Flags().BoolP("global", "g", false, "Use global configuration")
	generateCmd.Flags().
		StringP("date", "d", "", "Override the commit date (e.g. yesterday, \"2024-01-01\", \"last friday\")")
	generateCmd.Flags().
		BoolP("date-interactive", "D", false, "Prompt for commit date interactively after confirming")
}
