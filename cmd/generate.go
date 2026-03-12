package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
		dateFlag, _ := cmd.Flags().GetString("date")
		dateInteractive, _ := cmd.Flags().GetBool("date-interactive")

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

		author, err := git.GetAuthorInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		titleCommit := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Render("Commit")
		titleSep := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
			" - Generated in " + result.Elapsed.Truncate(10*time.Millisecond).String(),
		)
		commitStyle := lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("205")).
			PaddingLeft(1)
		italicStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("240"))

		dateDisplay := "now"
		if dateInteractive {
			dateDisplay = "to be defined"
		} else if dateFlag != "" {
			dateDisplay = dateFlag
		}

		fmt.Println()
		fmt.Println(titleCommit + titleSep)
		fmt.Println()
		fmt.Println(commitStyle.Render(commitMsg))
		fmt.Println()
		fmt.Println(italicStyle.Render(fmt.Sprintf(" %s <%s>", author.Name, author.Email)))
		fmt.Println(italicStyle.Render(fmt.Sprintf(" %s", dateDisplay)))

		confirmed := false
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Commit with this message ?").
					Affirmative("Commit").
					Negative("Cancel").
					Value(&confirmed),
			),
		).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		if !confirmed {
			fmt.Println("Cancelled")
			return
		}

		if dateInteractive {
			if err := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Commit date").
						Placeholder("e.g. yesterday, 2024-01-01, last friday").
						Value(&dateFlag),
				),
			).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
		}

		if err := git.Commit(commitMsg, git.CommitOptions{Date: dateFlag}); err != nil {
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
	generateCmd.Flags().
		StringP("date", "d", "", "Override the commit date (e.g. yesterday, \"2024-01-01\", \"last friday\")")
	generateCmd.Flags().
		BoolP("date-interactive", "D", false, "Prompt for commit date interactively after confirming")
}
