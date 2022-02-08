# appctl
**Run apps, Not Clusters.** Deploy your app on kubernetes in seconds,with no clusters required. Check for more information at https://platform9.com/appctl/

* Read the docs: [getting started with appctl](https://platform9.com/docs/appctl/getting-started)


### Purpose
* Appctl is a CLI that can be installed on Windows, MacOS and Linux, which connects to a Platform9 Managed Kubernetes Cluster running in AWS and enables users to deploy containerized applications in seconds.

### A better way to run apps on K8s
*  **Appctl** exposes the high value app orchestration capabilities available from Kubernetes and k-native, while hiding infrastructure complexity. 

* As a result, it is much faster to run apps while also running them more cost effectively in the cloud

### How appctl works
![flow-diagram](images/graphic_how-appctl-works.png)

### Pre-requisites
The CLI currently supports
* Linux (64 bit)
* Windows (64 bit)
* MacOS (64 bit)


### Installation and Usage
- Downloading the CLI can be done from [appctl website](https://platform9.com/appctl/) and from the command line. 

To install from the command line of host machine, run the following commands to download the appctl CLI and give executable permission to use it.

**For Linux**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/linux/appctl

chmod +x appctl
```

**For Mac**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/macos/appctl

chmod +x appctl
```

**For Windows**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/windows/appctl
```
After successfull download give the executable permission to appctl.

Once the CLI is successfully downloaded, run the Login command to authenticate to Platform9 and deploy applications.


### Appctl all commands
Below are all supported commands for Appctl.

```sh
#./appctl --help

CLI to deploy & manage apps in Platform9 environment.
Login first using $appctl login to use available commands.

Usage:
  appctl [command]

Available Commands:
  delete      Delete an existing app
  deploy      Deploy an app
  describe    Provide detailed app information in json format
  help        Help about any command
  list        Show all the running apps
  login       Login using Google account/Github account to use appctl
  version     Current version of appctl CLI being used

Flags:
  -h, --help      help for appctl

Use "appctl [command] --help" for more information about a command.
```

### Login 
To appctl you must first login by running ```./appctl login```

```sh
./appctl login --help
Login using Google account/Github account to use appctl

Usage:
  appctl login [flags]

Examples:

  # Login using Google account/Github account to use appctl.
  appctl login
 

Flags:
  -h, --help   help for login

```

**Example Login**
```sh
./appctl login
Starting login process.
Device verification is required to continue login.
Your Device Confirmation code is: TX2KW-BNPW6%
- Waiting for login to complete in browser... 
✔ Successfully Logged in!!
```

**Interactive login** The login command requires internet access and a web browser.

Appctl login is a two step process:

1. Device Verification: First you must verify where appctl is being run.
2. Login: Login using one of the supported federated identities (Google and Github).

When you rung ```appctl login``` a browser window will automatically open and prompt for the device confirmation code.

Confirm the device code displayed in the browser is identical to the code displayed by appctl, if it is correct click "Confirm" and the browser will redirected to _appctl log in _page.

**Appctl device confirmation**

![appctl_device_confirmation](images/appctl_device_confirmation.png)

Next, login using your **Google or Github account**

![login_using_google_github](images/Login_using_google_github_account.png)

Once you have successfully logged in, you will now be able to use appctl to deploy applications.

### Version

  **This command is used to get the current version of the CLI**
```sh
#appctl version

appctl version: v1.1

```

### Deploy

To deploy an app, run ```./appctl deploy```

The deploy command will deploy the specified container image using the provided name into Platform9 and automatically provision a public fully qualified domain with a unique port to access the application.

```sh
./appctl deploy --help
Deploy an app

Usage:
  appctl deploy [flags]

Examples:

  # Deploy an app using app-name and container image (public registry path)
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i gcr.io/knative-samples/helloworld-go
  
  # Deploy an app using app-name and container image, and pass environment variables.
  # Assumes the container has a server that will listen on port 8080
  appctl deploy -n <appname> -i <image> -e key1=value1 -e key2=value2

  # Deploy an app using app-name, container image and pass environment variables and set port where application listens on.
  appctl deploy -n <appname> -i <image> -e key1=value1 -e key2=value2 -p <port>
  Ex: appctl deploy -n hello -i gcr.io/knative-samples/helloworld-go -e TARGET="appctler" -p 7893
  

Flags:
  -n, --app-name string   Name of the app to be deployed 
                          (lowercase alphanumeric characters, '-' or '.', must start with alphanumeric characters only)
  -e, --env stringArray   Environment variable to set, as key=value pair
  -h, --help              help for deploy
  -i, --image string      Container image of the app (public registry path)
  -p, --port string       The port where app server listens, set as '--port <port>'
```

- **Example Deploy**
```sh
./appctl deploy --app-name <name> --image <docker-image path>

Example:
/appctl deploy --app-name hello --image gcr.io/knative-samples/helloworld-go
```

- **Specifying Ports**
If your application server listens on a specific port, then you can specify that while deploying the app using --port flag.

```sh 
./appctl deploy --app-name <name> --image <docker-image path> --port <port-value>

Example:
./appctl deploy --app-name hello --image gcr.io/knative-samples/helloworld-go --port 7893
```

- **Using Environment Variables**
```sh
./appctl deploy --app-name <name> --image <docker-image path> --env key1=value1

Example:
./appctl deploy --app-name hello --image gcr.io/knative-samples/helloworld-go --env TARGET=appctler
```

Appctl supports multiple --env variables

```sh
./appctl deploy --app-name <name> --image <docker-image path> --env key1=value1 --env key2=value2
```