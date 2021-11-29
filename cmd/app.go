package cmd

import (
	"fmt"

	"github.com/platform9/pf9-appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// appCmd represents the app commands can be run.
var (
	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Create or get or revise app",
		Long:  `Create or get or revise app, app.run controller`,
	}

	appCmdList = &cobra.Command{
		Use:   "list",
		Short: "List apps running",
		Long:  `List apps running`,
		Run:   appCmdListRun,
	}
)

// Flags declaration
var (
	nameSpace string
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appCmdList)
	appCmdList.Flags().StringVarP(&nameSpace, "namespace", "n", "", "set namespace")
}

// To list apps running in given namespace.
func appCmdListRun(cmd *cobra.Command, args []string) {
	err := appManageAPI.ListAppsInfo(nameSpace)
	if err != nil {
		fmt.Printf("Not able to list apps from namespace %v: Error %v", nameSpace, err)
	}
}
