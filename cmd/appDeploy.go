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

var deploy_example = `
  # Deploy an app using app-name and container image (public registory path) 
  appctl deploy -n <appname> -i gcr.io/knative-samples/helloworld-go
  
  # Deploy an app using app-name and container image, and pass environment variables.
  appctl deploy -n <appname> -i <image> -e key1=value1 -e key2=value2
 `

// appCmdDeploy - To deploy an app.
var (
	appCmdDeploy = &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy an app",
		Example: deploy_example,
		Long:    `Deploy an app`,
		Run:     appCmdDeployRun,
	}
)

type App struct {
	Name  string
	Image string
	Env   map[string]string
}

// This deployApp is of type App to take app name and app image from user.
var deployApp App

func init() {
	rootCmd.AddCommand(appCmdDeploy)
	appCmdDeploy.Flags().StringVarP(&deployApp.Name, "app-name", "n", "", `Name of the app to be deployed 
(lowercase alphanumeric characters, '-' or '.', must start with alphanumeric characters only)`)
	appCmdDeploy.Flags().StringVarP(&deployApp.Image, "image", "i", "", "Container image of the app (public registry path)")
	appCmdDeploy.Flags().StringToStringVarP(&deployApp.Env, "env", "e", nil, "Environment variable to set, as key=value pair")
}

func appCmdDeployRun(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	if deployApp.Name == "" {
		fmt.Printf("App Name: ")
		appName, _ := reader.ReadString('\n')
		deployApp.Name = strings.TrimSuffix(appName, "\n")
		deployApp.Name = strings.TrimSuffix(deployApp.Name, "\t")
	}

	// Validate app name.
	if !constants.RegexValidate(deployApp.Name) {
		fmt.Printf("Invalid App name.\n")
		fmt.Printf("Name of the app to be deployed must contain a lowercase alphanumeric characters, '-' or '.'\nand must start with alphanumeric characters only.\n")
		return
	}

	if deployApp.Image == "" {
		fmt.Printf("Source Image: ")
		appSourceImage, _ := reader.ReadString('\n')
		deployApp.Image = strings.TrimSuffix(appSourceImage, "\n")
		deployApp.Image = strings.TrimSuffix(deployApp.Image, "\t")
	}

	errapi := appManageAPI.CreateApp(deployApp.Name, deployApp.Image, deployApp.Env)
	if errapi != nil {
		fmt.Printf("\nNot able to deploy app: %v. Error: %v\n", deployApp.Name, errapi)
	}
}
