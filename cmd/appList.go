package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// usage example
var listExample = `
  # Get all the apps deployed.
  appctl list
 `

// appCmdList -- To list all apps running.
var (
	appCmdList = &cobra.Command{
		Use:     "list",
		Short:   "Show all the running apps",
		Example: listExample,
		Long:    `Show all the running apps`,
		Run:     appCmdListRun,
	}
)

func init() {
	rootCmd.AddCommand(appCmdList)
}

// To list apps running in given namespace.
func appCmdListRun(cmd *cobra.Command, args []string) {
	errapi := appManageAPI.ListAppsInfo()
	if errapi != nil {
		fmt.Printf("%v", errapi)
	}
}
