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

// usage example
var deployExample = `
  # Deploy an app using app-name and container image (public registry path)
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i gcr.io/knative-samples/helloworld-go
  
  # Deploy an app using app-name and container image (private registry path)
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i <private registry image path> -u <container registry username> -P <container registry password>

  	  # Sample command to deploy an app from a docker private registry path
  	  appctl deploy -n <appname> -i docker.io/<username>/<image>:<tag> -u <Docker username> -P <Docker password>

  	  # Sample command to deploy an app from an AWS ECR private registry path
  	  appctl deploy -n <appname> -i <aws_account_id>.dkr.ecr.<region>.amazonaws.com/<image>:<tag> -u AWS -P <Password obtained from AWS CLI>

  	  # Sample command to deploy an app from a GCR private registry path
  	  appctl deploy -n <appname> -i gcr.io/<GCP_projectID>/<image> -u oauth2accesstoken -P <Token obtained from gcloud CLI>


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
		Example: deployExample,
		Long:    `Deploy an app`,
		Run:     appCmdDeployRun,
	}
)

type App struct {
	name     string
	image    string
	env      []string
	port     string
	userName string
	password string
}

// command variables
// This deployApp is of type App to take app name and app image from user.
var deployApp App

func init() {
	rootCmd.AddCommand(appCmdDeploy)
	appCmdDeploy.Flags().StringVarP(&deployApp.name, "app-name", "n", "", `Name of the app to be deployed 
(lowercase alphanumeric characters, '-' or '.', must start with alphanumeric characters only)`)
	appCmdDeploy.Flags().StringVarP(&deployApp.image, "image", "i", "", "Container image of the app (public / private registry path)")
	appCmdDeploy.Flags().StringVarP(&deployApp.userName, "username", "u", "", "Username of private container registry")
	appCmdDeploy.Flags().StringVarP(&deployApp.password, "password", "P", "", "Password of private container registry")
	appCmdDeploy.Flags().StringArrayVarP(&deployApp.env, "env", "e", nil, "Environment variable to set, as key=value pair")
	appCmdDeploy.Flags().StringVarP(&deployApp.port, "port", "p", "", "The port where app server listens, set as '--port <port>'")
}

func appCmdDeployRun(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	if deployApp.name == "" {
		fmt.Printf("App Name: ")
		appName, _ := reader.ReadString('\n')
		deployApp.name = strings.TrimSuffix(appName, "\n")
		deployApp.name = strings.TrimSuffix(deployApp.name, "\r")
	}

	// Validate app name.
	if !constants.RegexValidate(deployApp.name, constants.ValidAppNameRegex) {
		fmt.Printf("Invalid app name.\n")
		fmt.Printf("Name of the app to be deployed must contain a lowercase alphanumeric characters, '-' or '.'\nand must start with alphanumeric characters only.\n")
		return
	}

	if deployApp.image == "" {
		fmt.Printf("Source Image: ")
		appSourceImage, _ := reader.ReadString('\n')
		deployApp.image = strings.TrimSuffix(appSourceImage, "\n")
		deployApp.image = strings.TrimSuffix(deployApp.image, "\r")
	}

	var isPrivateReg bool = true
	if deployApp.userName == "" && deployApp.password == "" {
		fmt.Printf("Is the image from a private registry (Y/n)? [n]: ")
		readerChar := bufio.NewReader(os.Stdin)
		char, _, _ := readerChar.ReadRune()
		if char == 'y' || char == 'Y' {
			fmt.Printf("Username: ")
			appSourceUsername, _ := reader.ReadString('\n')
			deployApp.userName = strings.TrimSuffix(appSourceUsername, "\n")
			deployApp.userName = strings.TrimSuffix(deployApp.userName, "\r")

			fmt.Printf("Password: ")
			appPassword, _ := terminal.ReadPassword(0)
			fmt.Printf("\n")
			appSourcePassword := string(appPassword)
			deployApp.password = strings.TrimSuffix(appSourcePassword, "\n")
			deployApp.password = strings.TrimSuffix(deployApp.password, "\r")
		} else {
			isPrivateReg = false
		}
	}

	//App to be deployed from private registry. Check if required options are provided
	if isPrivateReg {
		if deployApp.userName != "" && deployApp.password != "" {
			//Continue in this case
		} else {
			//incorrect options specified. Either both or none of the Username and Password should be specified.
			fmt.Printf("Incorrect options specified. Either both or none of the Username and Password should be specified.\n")
			os.Exit(0)
		}
	}

	if deployApp.port == "" {
		fmt.Printf("Port [8080]: ")
		port, _ := reader.ReadString('\n')
		deployApp.port = strings.TrimSuffix(port, "\n")
		deployApp.port = strings.TrimSuffix(deployApp.port, "\r")
	}

	if deployApp.port != "" {
		// Check if port given is valid i.e numeric only.
		_, err := strconv.Atoi(deployApp.port)
		if err != nil {
			fmt.Printf("Invalid port. Please enter a valid port\n")
			return
		}
	}

	errapi := appManageAPI.CreateApp(deployApp.name, deployApp.image, deployApp.userName,
		deployApp.password, deployApp.env, deployApp.port)
	if errapi != nil {
		fmt.Printf("\nNot able to deploy app: %v.\nError: %v", deployApp.name, errapi)
	}
}
