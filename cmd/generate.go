package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/git"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a commit message using AI models.",
	Long:  `The generate command analyzes the git diff in your git repository and generates a commit message using AI models.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		isGlobal, _ := cmd.Flags().GetBool("global")

		if isGlobal {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Error getting home directory: %v\n", err)
				return
			}
			viper.SetConfigFile(filepath.Join(home, ".comai.yaml"))
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("Error loading global config: %v\n", err)
				return
			}
		}

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Printf("Erreur lors du chargement de la config : %v\n", err)
			return
		}

		config.PrintConfig(cfg, verbose)

		// Get staged diff
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println()
		fmt.Println(diff.Stats)
		fmt.Println()
		for _, f := range diff.Files {
			fmt.Printf("[%s] %s\n", f.Status, f.Path)
		}

		// TODO: Send diff.RawDiff to AI for commit message generation
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().
		BoolP("verbose", "v", false, "Show verbose output including config file path")
	generateCmd.Flags().BoolP("global", "g", false, "Use global configuration")
}
