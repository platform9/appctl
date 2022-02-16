package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/spf13/cobra"

	"golang.org/x/crypto/ssh/terminal"
)

var deploy_example = `
  # Deploy an app using app-name and container image (public registry path)
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i gcr.io/knative-samples/helloworld-go
  
  # Deploy an app using app-name and container image (private registry path)
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i <private registry image path> -u <container registry username> -P <container registry password>
  
  # Deploy an app using app-name and container image, and pass environment variables.
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i <image> -e key1=value1 -e key2=value2

  # Deploy an app using app-name, container image and pass environment variables and set port where application listens on.
  appctl deploy -n <appname> -i <image> -e key1=value1 -e key2=value2 -p <port>
  Ex: appctl deploy -n hello -i gcr.io/knative-samples/helloworld-go -e TARGET="appctler" -p 7893
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
	Env   []string
	Port  string
	Username string
	Password string
}

// This deployApp is of type App to take app name and app image from user.
var deployApp App

func init() {
	rootCmd.AddCommand(appCmdDeploy)
	appCmdDeploy.Flags().StringVarP(&deployApp.Name, "app-name", "n", "", `Name of the app to be deployed 
(lowercase alphanumeric characters, '-' or '.', must start with alphanumeric characters only)`)
	appCmdDeploy.Flags().StringVarP(&deployApp.Image, "image", "i", "", "Container image of the app (public / private registry path)")
	appCmdDeploy.Flags().StringVarP(&deployApp.Username, "username", "u", "", "Username of private container registry")
	appCmdDeploy.Flags().StringVarP(&deployApp.Password, "password", "P", "", "Password of private container registry")
	appCmdDeploy.Flags().StringArrayVarP(&deployApp.Env, "env", "e", nil, "Environment variable to set, as key=value pair")
	appCmdDeploy.Flags().StringVarP(&deployApp.Port, "port", "p", "", "The port where app server listens, set as '--port <port>'")
}

func appCmdDeployRun(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	if deployApp.Name == "" {
		fmt.Printf("App Name: ")
		appName, _ := reader.ReadString('\n')
		deployApp.Name = strings.TrimSuffix(appName, "\n")
		deployApp.Name = strings.TrimSuffix(deployApp.Name, "\r")
	}

	// Validate app name.
	if !constants.RegexValidate(deployApp.Name, constants.ValidAppNameRegex) {
		fmt.Printf("Invalid app name.\n")
		fmt.Printf("Name of the app to be deployed must contain a lowercase alphanumeric characters, '-' or '.'\nand must start with alphanumeric characters only.\n")
		return
	}

	if deployApp.Image == "" {
		fmt.Printf("Source Image: ")
		appSourceImage, _ := reader.ReadString('\n')
		deployApp.Image = strings.TrimSuffix(appSourceImage, "\n")
		deployApp.Image = strings.TrimSuffix(deployApp.Image, "\r")
	}

	var isPrivateReg bool = true;
	if deployApp.Username == "" && deployApp.Password == "" {
		fmt.Printf("Is the image from a private registry (Y/n)? [n]: ")
		readerChar := bufio.NewReader(os.Stdin)
		char, _, _ := readerChar.ReadRune()
		if char == 'y' || char == 'Y' {
			fmt.Printf("Username: ")
			appSourceUsername, _ := reader.ReadString('\n')
			deployApp.Username = strings.TrimSuffix(appSourceUsername, "\n")
			deployApp.Username = strings.TrimSuffix(deployApp.Username, "\r")

			fmt.Printf("Password: ")
			appPassword, _ := terminal.ReadPassword(0)
			fmt.Printf("\n")
			appSourcePassword := string(appPassword)
			deployApp.Password = strings.TrimSuffix(appSourcePassword, "\n")
			deployApp.Password = strings.TrimSuffix(deployApp.Password, "\r")
		} else {
			isPrivateReg = false;
		}
	}

	//App to be deployed from private registry. Check if required options are provided
	if (isPrivateReg) {
		if deployApp.Username != "" && deployApp.Password != "" {
			//Continue in this case
		} else {
			//incorrect options specified. Either both or none of the Username and Password should be specified.
			fmt.Printf("Incorrect options specified. Either both or none of the Username and Password should be specified.\n")
			os.Exit(0)
		}
	}

	if deployApp.Port == "" {
		fmt.Printf("Port [8080]: ")
		port, _ := reader.ReadString('\n')
		deployApp.Port = strings.TrimSuffix(port, "\n")
		deployApp.Port = strings.TrimSuffix(deployApp.Port, "\r")
	}

	if deployApp.Port != "" {
		// Check if port given is valid i.e numeric only.
		_, err := strconv.Atoi(deployApp.Port)
		if err != nil {
			fmt.Printf("Invalid port. Please enter a valid port\n")
			return
		}
	}

	errapi := appManageAPI.CreateApp(deployApp.Name, deployApp.Image, deployApp.Username,
		deployApp.Password, deployApp.Env, deployApp.Port)
	if errapi != nil {
		fmt.Printf("\nNot able to deploy app: %v.\nError: %v", deployApp.Name, errapi)
	}
}
