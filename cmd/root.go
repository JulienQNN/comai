package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "comai",
	Short: "Comai is a CLI tool for generating commits using AI models.",
	Long:  `Comai is a CLI tool for generating commits using AI models.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.comai.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home)

		viper.SetConfigName(".comai")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur Viper : %v\n", err)
	}
}
