/*
Copyright © 2025 MetaDandy benitezarroyojoseph@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "goce",
	Short: "GoCelerator (alias goce): CLI for bootstrapping and managing Go projects",
	Long: `goce is the command-line interface for generating and managing Go applications with a modular architecture.
	Available commands:
	init      Initialize a new Go project
				Usage: goce init <name> [--fiber] [--goversion <ver>] [--no-air] [--no-docker]
	serve     Start the development server (with optional hot-reload)
	docker    Build Docker image and start services with Docker Compose
	generate  Create new modules or core packages with stubs
	migrate   Manage database migrations (Gorm/SQL)
	seed      Run data seeders
	env       Generate a .env file and guide for API keys
	test      Run tests with coverage report
	doctor    Verify prerequisites (Go, PostgreSQL, environment variables)
	version   Show goce and Go versions installed

	Examples:
	# Scaffold a standard net/http project
	goce init myapp

	# Scaffold a Fiber project
	goce init myapp --fiber

	# Skip Air and Docker setup
	goce init myapp --no-air --no-docker
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.GoCelerator.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".GoCelerator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".GoCelerator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
