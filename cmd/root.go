package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var verbosity bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "pf9-appctl",
	Long: `CLI tool for Platform9 app run management.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Fatalf(err.Error())
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	// To tell Cobra not to provide the default completion command.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&verbosity, "verbose", false, "print verbose logs to console")
}
