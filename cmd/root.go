package cmd

import (
	"fmt"
	"os"

	"github.com/platform9/appctl/pkg/color"
	"github.com/platform9/appctl/pkg/constants"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

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
	requiredSecrets := map[string]string{
		"APPURL":     constants.APPURL,
		"DOMAIN":     constants.DOMAIN,
		"CLIENTID":   constants.CLIENTID,
		"GRANT_TYPE": constants.GrantType,
	}
	missing := ""
	for secret, val := range requiredSecrets {
		if val == "" {
			missing = fmt.Sprintf("%s%s\n", missing, secret)
		}
	}
	if missing != "" {
		fmt.Println(color.Red("appctl secrets not set.\n",
			"Please ensure the following values are set:\n",
			missing,
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
