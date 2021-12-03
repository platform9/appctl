package cmd

import (
	"fmt"

	"github.com/platform9/pf9-appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// appCmd represents the app commands can be run.
var (
	appCmdList = &cobra.Command{
		Use:   "list",
		Short: "List apps running",
		Long:  `List apps running`,
		Run:   appCmdListRun,
	}
)

func init() {
	rootCmd.AddCommand(appCmdList)
}

// To list apps running in given namespace.
func appCmdListRun(cmd *cobra.Command, args []string) {
	// Call function to get user namespace from login info.
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	errapi := appManageAPI.ListAppsInfo(nameSpace)
	if errapi != nil {
		fmt.Printf("%v\n", errapi)
	}
}
