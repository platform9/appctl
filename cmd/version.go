package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/constants"
	"github.com/spf13/cobra"
)

// usage example
var versionExample = `
  # Check the current version of appctl CLI being used.
  appctl version
 `

// versionCmd represents "Version of appctl being used.".
var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Current version of appctl CLI being used",
		Example: versionExample,
		Long:    `Current version of appctl CLI being used`,
		Run: func(cmd *cobra.Command, args []string) {
			//Prints the current version of appctl being used.
			fmt.Println(constants.CLIVersion)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
