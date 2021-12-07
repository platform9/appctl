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

var DeployApp App

func init() {
	rootCmd.AddCommand(appCmdDeploy)
	appCmdDeploy.Flags().StringVarP(&DeployApp.Name, "app-name", "n", "", `set app name to create 
(lowercase alphanumeric characters, '_' or '.', must start with alphanumeric characters only)`)
	appCmdDeploy.Flags().StringVarP(&DeployApp.Image, "image", "i", "", "set app source image to create")
}

func appCmdDeployRun(cmd *cobra.Command, args []string) {
	nameSpace, err := appManageAPI.GetNameSpace()
	if err != nil {
		fmt.Printf("Not able to get namespace. Error %v\n", err)
	}
	reader := bufio.NewReader(os.Stdin)

	if DeployApp.Name == "" {
		fmt.Printf("App Name: ")
		appName, _ := reader.ReadString('\n')
		DeployApp.Name = strings.TrimSuffix(appName, "\n")
	}
	if DeployApp.Image == "" {
		fmt.Printf("Source Image: ")
		appSourceImage, _ := reader.ReadString('\n')
		DeployApp.Image = strings.TrimSuffix(appSourceImage, "\n")
	}
	errapi := appManageAPI.CreateApp(DeployApp.Name, nameSpace, DeployApp.Image)
	if errapi != nil {
		fmt.Printf("\nNot able to create app of source image %v in namespace %v: Error %v\n", DeployApp.Image, nameSpace, errapi)
	}
}
