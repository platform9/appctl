package cmd

import (
	"fmt"
	"os"

	"github.com/platform9/appctl/pkg/color"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/platform9/appctl/pkg/segment"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var verbosity bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "appctl",
	Long: `CLI to deploy & manage apps in Platform9 environment.
Login first using "appctl login" to use available commands.`,
	PersistentPreRun: ensureAppSecrets,
}

func ensureAppSecrets(cmd *cobra.Command, args []string) {
	if cmd.Name() == "help" || cmd.Name() == "version" {
		return
	}
	if constants.APPURL == "" || constants.DOMAIN == "" ||
	constants.CLIENTID == "" || constants.GrantType == "" || segment.APPCTL_SEGMENT_WRITE_KEY == "" {
		fmt.Println(color.Red("appctl secrets not set.\n",
		"Please ensure all of APPURL, DOMAIN, CLIENTID, GRANT_TYPE and APPCTL_SEGMENT_WRITE_KEY are set during the build.",
		))
		os.Exit(1)
	}

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
	//rootCmd.PersistentFlags().BoolVar(&verbosity, "verbose", false, "print verbose logs to console")
}
