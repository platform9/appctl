package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/spf13/cobra"
)

var describe_example = `
  # Get detailed information about an app deployed through app-name in json format.
  appctl describe -n <appname>
 `

// appCmdDescribe -- To describe an app running.
var (
	appCmdDescribe = &cobra.Command{
		Use:     "describe",
		Short:   "Provide detailed app information in json format",
		Example: describe_example,
		Long:    `Provide detailed app information in json format`,
		Run:     appCmdDescribeRun,
	}
)

var AppName string

func init() {
	rootCmd.AddCommand(appCmdDescribe)
	appCmdDescribe.Flags().StringVarP(&AppName, "app-name", "n", "", "Name of app to be described")
}

// To get app information by its name
func appCmdDescribeRun(cmd *cobra.Command, args []string) {
	// Check if App name provided.
	if AppName == "" {
		fmt.Printf("App Name not specified.\n")
		return
	}

	// Validate app name.
	if !constants.RegexValidate(AppName) {
		fmt.Printf("Invalid App name.\n")
		return
	}

	errapi := appManageAPI.GetAppByNameInfo(AppName)
	if errapi != nil {
		fmt.Printf("%v\n", errapi)
	}
}
