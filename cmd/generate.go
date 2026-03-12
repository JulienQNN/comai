package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"charm.land/huh/v2"
	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/generate"
	"github.com/JulienQNN/comai/internal/git"
	"github.com/JulienQNN/comai/internal/prompt"
	"github.com/JulienQNN/comai/internal/provider"
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

		if dateFlag != "" {
			if _, err := git.ParseDate(dateFlag); err != nil {
				log.Error("parsing date", "err", err)
				return
			}
		}

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

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Error("loading config", "err", err)
			return
		}

		config.PrintConfig(cfg, verbose)
		diff, err := git.GetStagedDiff()
		if err != nil {
			log.Error("getting git diff", "err", err)
			return
		}

		fmt.Println(diff.Stats)

		for _, f := range diff.Files {
			fmt.Printf("[%s] %s\n", f.Status, f.Path)
		}

		p, err := provider.New(cfg.ModelName)
		if err != nil {
			log.Error("creating provider", "err", err)
			return
		}

		result, err := generate.Start(p, prompt.Build(diff.RawDiff, cfg))
		if err != nil {
			log.Error("generating commit message", "err", err)
			return
		}

		commitMsg := strings.ToLower(result.CommitMsg)

		author, err := git.GetAuthorInfo()
		if err != nil {
			log.Error("getting author info", "err", err)
			return
		}

		t := theme.Default()
		titleCommit := t.CommitTitle.Render("Commit")
		titleSep := t.Muted.Render(
			" - Generated in " + result.Elapsed.Truncate(10*time.Millisecond).String(),
		)
		dateDisplay := "now"
		if dateInteractive {
			dateDisplay = "to be defined"
		} else if dateFlag != "" {
			parsed, err := git.ParseDate(dateFlag)
			if err != nil {
				log.Error("parsing date", "err", err)
				return
			}
			dateDisplay = parsed.Format("2006-01-02 15:04:05")
		}

		fmt.Println(titleCommit + titleSep)
		fmt.Println(t.CommitBorder.Render(commitMsg))
		fmt.Println(t.Italic.PaddingBottom(1).Render(fmt.Sprintf(" %s <%s>", author.Name, author.Email)))
		fmt.Println(t.Italic.Render(fmt.Sprintf(" %s", dateDisplay)))

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
			log.Error("running confirmation form", "err", err)
			return
		}
		if !confirmed {
			log.Info("Commit cancelled by user")
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
			log.Error("committing changes", "err", err)
			return
		}
		log.Info("Changes committed successfully")
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
