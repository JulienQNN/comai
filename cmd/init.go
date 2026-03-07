package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/JulienQNN/comai/internal/config"
)

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		isGlobal, _ := cmd.Flags().GetBool("global")
		err := config.RunSetupWizard(isGlobal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur : %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().
		BoolP("global", "g", false, "Generate the comai.yaml file in the global configuration directory.")
}
