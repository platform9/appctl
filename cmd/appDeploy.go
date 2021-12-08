package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// appCmdDeploy - To deploy an app.
var (
	appCmdDeploy = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an app",
		Long:  `Deploy an app`,
		Run:   appCmdDeployRun,
	}
)

type App struct {
	Name  string
	Image string
}

// This deployApp is of type App to take app name and app image from user.
var deployApp App

func init() {
	rootCmd.AddCommand(appCmdDeploy)
	appCmdDeploy.Flags().StringVarP(&deployApp.Name, "app-name", "n", "", `Name of the app to be deployed 
(lowercase alphanumeric characters, '-' or '.', must start with alphanumeric characters only)`)
	appCmdDeploy.Flags().StringVarP(&deployApp.Image, "image", "i", "", "Container image of the app (public registry path)")
}

func appCmdDeployRun(cmd *cobra.Command, args []string) {
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	reader := bufio.NewReader(os.Stdin)

	if deployApp.Name == "" {
		fmt.Printf("App Name: ")
		appName, _ := reader.ReadString('\n')
		deployApp.Name = strings.TrimSuffix(appName, "\n")
		deployApp.Name = strings.TrimSuffix(deployApp.Name, "\t")
	}
	if deployApp.Image == "" {
		fmt.Printf("Source Image: ")
		appSourceImage, _ := reader.ReadString('\n')
		deployApp.Image = strings.TrimSuffix(appSourceImage, "\n")
		deployApp.Image = strings.TrimSuffix(deployApp.Image, "\t")
	}
	errapi := appManageAPI.CreateApp(deployApp.Name, nameSpace, deployApp.Image)
	if errapi != nil {
		fmt.Printf("\nNot able to deploy app: %v. Error: %v\n", deployApp.Name, errapi)
	}
}
