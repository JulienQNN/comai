/*
Copyright © 2026 Julien QUENNEHEN <julienqhn@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JulienQNN/comai/internal/config"
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

		// 3. Suite de ta logique (Git diff, Appel IA...)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().
		BoolP("verbose", "v", false, "Show verbose output including config file path")
	generateCmd.Flags().BoolP("global", "g", false, "Use global configuration")
}
