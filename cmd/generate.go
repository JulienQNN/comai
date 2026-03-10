package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/generate"
	"github.com/JulienQNN/comai/internal/git"
	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider"
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
				fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
				return
			}
			viper.SetConfigFile(filepath.Join(home, ".comai.yaml"))
			if err := viper.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading global config: %v\n", err)
				return
			}
		}

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			return
		}

		config.PrintConfig(cfg, verbose)

		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Println()
		fmt.Println(diff.Stats)
		fmt.Println()
		for _, f := range diff.Files {
			fmt.Printf("[%s] %s\n", f.Status, f.Path)
		}

		p, err := provider.New(cfg.ModelName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating provider: %v\n", err)
			return
		}

		result, err := generate.Start(p, prompt.Build(diff.RawDiff, cfg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", result.Err)
			return
		}
		if result.CommitMsg == "" {
			fmt.Fprintf(os.Stderr, "Error: empty response from LLM\n")
			return
		}

		commitMsg := strings.ToLower(result.CommitMsg)

		fmt.Println()
		fmt.Println(commitMsg)

		confirmed := false
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Commit with this message ?").
					Affirmative("Commit").
					Negative("Cancel").
					Value(&confirmed),
			),
		)
		if err := form.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		if !confirmed {
			fmt.Println("Cancelled")
			return
		}

		if err := git.Commit(commitMsg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Println("Committed !")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().
		BoolP("verbose", "v", false, "Show verbose output including config file path")
	generateCmd.Flags().BoolP("global", "g", false, "Use global configuration")
}
