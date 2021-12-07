package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// appCmdDescribe -- To describe an app running.
var (
	appCmdDescribe = &cobra.Command{
		Use:   "describe",
		Short: "Provide detailed app information in json format",
		Long:  `Provide detailed app information in json format`,
		Run:   appCmdDescribeRun,
	}
)

var AppName string

func init() {
	rootCmd.AddCommand(appCmdDescribe)
	appCmdDescribe.Flags().StringVarP(&AppName, "app-name", "n", "", "Name of app to be described")
}

// To get app information by its name
func appCmdDescribeRun(cmd *cobra.Command, args []string) {
	// Call function to get user namespace from login info.
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	errapi := appManageAPI.GetAppByNameInfo(AppName, nameSpace)
	if errapi != nil {
		fmt.Printf("%v\n", errapi)
	}
}
