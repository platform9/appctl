package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	appCmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create an app",
		Long:  `Create an app`,
		Run:   appCmdCreateRun,
	}
	appCmdGetAppByName = &cobra.Command{
		Use:   "get",
		Short: "Get an app info",
		Long:  `Get an app info`,
		Run:   appCmdGetAppByNameRun,
	}
)

var AppName string

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appCmdList)
	appCmd.AddCommand(appCmdCreate)
	appCmd.AddCommand(appCmdGetAppByName)
	appCmdCreate.Flags().StringVarP(&CreateApp.Name, "app_name", "a", "", `set app name to create 
(lowercase alphanumeric characters, '_' or '.', must start with alphanumeric characters only)`)
	appCmdCreate.Flags().StringVarP(&CreateApp.Image, "image", "i", "", "set app source image to create")
	appCmdGetAppByName.Flags().StringVarP(&AppName, "app_name", "a", "", "set app name to get info")

}

type App struct {
	Name  string
	Image string
}

var CreateApp App

// To list apps running in given namespace.
func appCmdListRun(cmd *cobra.Command, args []string) {
	// Call function to get user namespace from login info.
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	errapi := appManageAPI.ListAppsInfo(nameSpace)
	if errapi != nil {
		fmt.Printf("Not able to list apps from namespace %v: Error %v\n", nameSpace, err)
	}
}

func appCmdCreateRun(cmd *cobra.Command, args []string) {
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	reader := bufio.NewReader(os.Stdin)
	for true {
		if CreateApp.Name == "" {
			fmt.Printf("App Name: ")
			appName, _ := reader.ReadString('\n')
			CreateApp.Name = strings.TrimSuffix(appName, "\n")
		}
		if CreateApp.Image == "" {
			fmt.Printf("Source Image: ")
			appSourceImage, _ := reader.ReadString('\n')
			CreateApp.Image = strings.TrimSuffix(appSourceImage, "\n")
		}
		if CreateApp.Name == "" && CreateApp.Image != "" {
			fmt.Printf("App Name is found empty, give valid app name\n")
		}
		if CreateApp.Name != "" && CreateApp.Image == "" {
			fmt.Printf("Source Image is found empty, give valid image\n")
		}
		if CreateApp.Name == "" && CreateApp.Image == "" {
			fmt.Printf("Both App Name, Source Image are found empty, give valid information\n")
		}
		if CreateApp.Name != "" && CreateApp.Image != "" {
			break
		}
	}
	errapi := appManageAPI.CreateApp(CreateApp.Name, nameSpace, CreateApp.Image)
	if errapi != nil {
		fmt.Printf("Not able to create app of source image %v in namespace %v: Error %v\n", CreateApp.Image, nameSpace, err)
	}
}

// To get app information by its name
func appCmdGetAppByNameRun(cmd *cobra.Command, args []string) {
	// Call function to get user namespace from login info.
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	errapi := appManageAPI.GetAppByNameInfo(AppName, nameSpace)
	if errapi != nil {
		fmt.Printf("\nNot able to get app information from namespace %v: Error %v\n", nameSpace, err)
	}
}
