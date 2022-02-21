package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/spf13/cobra"
)

// usage example
var describeExample = `
  # Get detailed information about an app deployed through app-name in json format.
  appctl describe -n <appname>
 `

// appCmdDescribe -- To describe an app running.
var (
	appCmdDescribe = &cobra.Command{
		Use:     "describe",
		Short:   "Provide detailed app information in json format",
		Example: describeExample,
		Long:    `Provide detailed app information in json format`,
		Run:     appCmdDescribeRun,
	}
)

// command variables
var appNameDescribe string

func init() {
	rootCmd.AddCommand(appCmdDescribe)
	appCmdDescribe.Flags().StringVarP(&appNameDescribe, "app-name", "n", "", "Name of app to be described")
}

// To get app information by its name
func appCmdDescribeRun(cmd *cobra.Command, args []string) {
	// Check if App name provided.
	if appNameDescribe == "" {
		fmt.Printf("App name not specified.\n")
		return
	}

	// Validate app name.
	if !constants.RegexValidate(appNameDescribe, constants.ValidAppNameRegex) {
		fmt.Printf("Invalid app name.\n")
		return
	}

	errapi := appManageAPI.GetAppByNameInfo(appNameDescribe)
	if errapi != nil {
		fmt.Printf("%v", errapi)
	}
}
