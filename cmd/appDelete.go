package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/spf13/cobra"
)

// usage example
var deleteExample = `
  # Delete an app using app-name.
  appctl delete -n <appname>

  # Force delete an app using app-name and force flag.
  appctl delete -n <appname> -f
 `

// appCmdDelete -- To delete an existing app.
var (
	appCmdDelete = &cobra.Command{
		Use:     "delete",
		Short:   "Delete an existing app",
		Example: deleteExample,
		Long:    `Delete an existing app`,
		Run:     appCmdDeleteRun,
	}
)

// command variables
var (
	// App name to delete.
	appNameDelete string
	// To force delete an app.
	force bool
	//Choice to delete app
	deleteConfirmChoice string
)

func init() {
	rootCmd.AddCommand(appCmdDelete)
	appCmdDelete.Flags().StringVarP(&appNameDelete, "app-name", "n", "", "Provide the name of app to be deleted")
	appCmdDelete.Flags().BoolVarP(&force, "force", "f", false, "To force delete an app")
}

// To delete an app by its name.
func appCmdDeleteRun(cmd *cobra.Command, args []string) {
	if appNameDelete == "" {
		fmt.Printf("App name not specified.\n")
		return
	}
	// To ask user if to delete app when force delete is false.
	if !(force) {
		var count = 0
		for count < 3 {
			count++
			// To make sure delete the app
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Are you sure you want to delete app (y/n)? ")
			deleteConfirmChoice, _ = reader.ReadString('\n')
			deleteConfirmChoice = strings.TrimSuffix(deleteConfirmChoice, "\n")
			deleteConfirmChoice = strings.TrimSuffix(deleteConfirmChoice, "\r")

			// If response is other than "y" or "n"
			if deleteConfirmChoice != "y" && deleteConfirmChoice != "n" {
				fmt.Printf("Please enter correct input (y/n).\n")
				continue
			}
			// To delete app if Yes
			if deleteConfirmChoice == "y" {
				// Validate app name.
				if !constants.RegexValidate(appNameDelete, constants.ValidAppNameRegex) {
					fmt.Printf("Invalid app name.\n")
					return
				}
				errapi := appManageAPI.DeleteApp(appNameDelete)
				if errapi != nil {
					fmt.Printf("%v", errapi)
					return
				}
				fmt.Printf("Successfully deleted the app: %v\n", appNameDelete)
				break
			}
			// To stop delete app process if No
			if deleteConfirmChoice == "n" {
				fmt.Printf("You have cancelled the app deletion activity!!\n")
				break
			}
		}
	} else {
		// Validate app name.
		if !constants.RegexValidate(appNameDelete, constants.ValidAppNameRegex) {
			fmt.Printf("Invalid app name.\n")
			return
		}
		// If force delete an app.
		errapi := appManageAPI.DeleteApp(appNameDelete)
		if errapi != nil {
			fmt.Printf("%v", errapi)
			return
		}
		fmt.Printf("Successfully deleted the app: %v\n", appNameDelete)
	}

}
