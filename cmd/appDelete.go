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

// appCmdDelete -- To delete an existing app.
var (
	appCmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing app",
		Long:  `Delete an existing app`,
		Run:   appCmdDeleteRun,
	}
)

var (
	// App name to delete.
	AppNameDelete string
	// To force delete an app.
	force bool
	//Choice to delete app
	deleteApp string
)

func init() {
	rootCmd.AddCommand(appCmdDelete)
	appCmdDelete.Flags().StringVarP(&AppNameDelete, "app-name", "n", "", "Provide the name of app to be deleted")
	appCmdDelete.Flags().BoolVarP(&force, "force", "f", false, "To force delete an app")
}

// To delete an app by its name.
func appCmdDeleteRun(cmd *cobra.Command, args []string) {
	if AppNameDelete == "" {
		fmt.Printf("App Name not specified.\n")
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
			deleteApp, _ = reader.ReadString('\n')
			deleteApp = strings.TrimSuffix(deleteApp, "\n")

			// If response is other than "y" or "n"
			if deleteApp != "y" && deleteApp != "n" {
				fmt.Printf("Please enter correct input (y/n).\n")
				continue
			}
			// To delete app if Yes
			if deleteApp == "y" {
				// Validate app name.
				if !constants.RegexValidate(AppNameDelete) {
					fmt.Printf("Invalid App name.\n")
					return
				}
				errapi := appManageAPI.DeleteApp(AppNameDelete)
				if errapi != nil {
					fmt.Printf("%v\n", errapi)
					return
				}
				fmt.Printf("Successfully deleted the app: %v\n", AppNameDelete)
				break
			}
			// To stop delete app process if No
			if deleteApp == "n" {
				fmt.Printf("You have cancelled the app deletion activity!!\n")
				break
			}
		}
	} else {
		// Validate app name.
		if !constants.RegexValidate(AppNameDelete) {
			fmt.Printf("Invalid App name.\n")
			return
		}
		// If force delete an app.
		errapi := appManageAPI.DeleteApp(AppNameDelete)
		if errapi != nil {
			fmt.Printf("%v\n", errapi)
			return
		}
		fmt.Printf("Successfully deleted the app: %v\n", AppNameDelete)
	}

}
