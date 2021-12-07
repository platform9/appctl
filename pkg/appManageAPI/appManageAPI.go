package appManageAPI

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/platform9/appctl/pkg/appAPIs"
	"github.com/platform9/appctl/pkg/browser"
	"github.com/platform9/appctl/pkg/color"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/ryanuber/columnize"
)

type ListAppInfo struct {
	Name         string
	URL          string
	Image        string
	NameSpace    string
	CreationTime string
}

// To list apps.
func ListAppsInfo(
	nameSpace string, // namespace to list apps.
) error {
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.\n")
	}

	// To list and store output.
	var list ListAppInfo
	var Output []string
	Output = append(Output, constants.TABLEFORMAT)

	// Fetch the running apps in given namespace.
	list_apps, err := appAPIs.ListApps(nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to list apps with Error: %v\n", err)
	}

	// Fetch the App name, namespace deployed in, creationTimestamp
	for _, items := range list_apps["items"].([]interface{}) {
		for key, appInfo := range items.(map[string]interface{}) {
			if key == "metadata" {
				list.Name = fmt.Sprintf("%v", appInfo.(map[string]interface{})["name"])
				list.NameSpace = fmt.Sprintf("%v", appInfo.(map[string]interface{})["namespace"])
				list.CreationTime = fmt.Sprintf("%v", appInfo.(map[string]interface{})["creationTimestamp"])
			}
			// Fetch the Image name.
			if key == "spec" {
				template := (appInfo.(map[string]interface{}))["template"].(map[string]interface{})
				detailedSpec := template["spec"].(map[string]interface{})
				containers := detailedSpec["containers"].([]interface{})[0]
				list.Image = fmt.Sprintf("%v", containers.(map[string]interface{})["image"])
			}

			// Fetch the URL Endpoint
			if key == "status" {
				list.URL = fmt.Sprintf("%v", appInfo.(map[string]interface{})["url"])
			}
		}
		appinfo := fmt.Sprintf("%v | %v | %v | %v | %v", list.Name, list.URL, list.Image, list.NameSpace, list.CreationTime)
		Output = append(Output, appinfo)
	}

	tabularAppInfo := columnize.SimpleFormat(Output)
	fmt.Println(tabularAppInfo)
	return nil
}

func GetNameSpace() (string, error) {
	/*Fetch the user details like login id, email, username
	from config file and get namespace.
	-- Keeping it default, since login logic is not implemented yet.
	*/
	nameSpace := "default"
	return nameSpace, nil
}

// To create an app.
func CreateApp(
	name string, // App name to create.
	nameSpace string, // namespace to list apps.
	image string, // Source Image to create app.
) error {
	if name == "" && image == "" {
		return fmt.Errorf("App Name and Image not specified.\n")
	}
	if name == "" {
		return fmt.Errorf("App Name not specified.\n")
	}
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.\n")
	}
	if image == "" {
		return fmt.Errorf("Image not specified.\n")
	}
	fmt.Printf("Started Creating App\n")
	fmt.Printf("\nName is %v\nNamespace is %v\nImage is %v\n", name, nameSpace, image)
	// Fetch the running apps in given namespace.

	err := appAPIs.CreateApp(name, nameSpace, image)
	if err != nil {
		return fmt.Errorf("\nFailed to create app with error: %v\n", err)
	}

	//Since creation of App takes some time.
	time.Sleep(5 * time.Second)

	// Fetch the detailedapp information for given name from given namespace.
	get_app, err := appAPIs.GetAppByName(name, nameSpace)
	if err != nil {
		return fmt.Errorf("\nFailed to get app information with error: %v\n", err)
	}
	// URL/ Endpoint where the app service is available.
	url := (get_app["status"]).(map[string]interface{})["url"]
	if url != nil {
		fmt.Printf("App created with Name: %v, and is available at URL: %v\n", name, url)
	} else {
		fmt.Printf("App created with Name: %v. Run 'appctl list' to get more information on app\n", name)
	}
	return nil
}

// To get a detailed information of particular app by name.
func GetAppByNameInfo(
	name string, // app name
	nameSpace string, // namespace to list apps.
) error {
	if name == "" {
		return fmt.Errorf("App Name not specified.")
	}
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.")
	}
	// Fetch the detailedapp information for given name from given namespace.
	get_app, err := appAPIs.GetAppByName(name, nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to get app information with error: %v\nCheck 'appctl list' for more information on apps running.", err)
	}
	jsonformated, err := json.MarshalIndent(get_app, "", "  ")
	fmt.Printf("%v\n", string(jsonformated))
	return nil
}

// To login using Device authentication and access appctl.
func LoginApp() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("red")

	fmt.Printf(color.Blue("Starting Appctl Login..") + "\n")
	// Get the device code.
	deviceCode, err := appAPIs.GetDeviceCode()
	if err != nil {
		fmt.Printf("Unable to generate device code.")
	}

	fmt.Printf("Please verify the device to continue for further process of login..\n")
	fmt.Printf("Your Device Confirmation code is: " + color.Yellow(deviceCode.UserCode) + "\n")

	// To open browser, for device verification and SSO.
	err = browser.OpenBrowser(deviceCode.VerificationUrlComplete)
	if err != nil {
		fmt.Printf("\nCouldn't open the URL, kindly do it manually: " + color.Yellow(deviceCode.VerificationUrlComplete) + "\n")
	}

	var Token *appAPIs.TokenInfo

	// Wait for device verification in browser and if its success request the token.
	s.Start()
	s.Suffix = " Waiting for login to complete in browser..."

	for true {
		// Request for token.
		Token, err = appAPIs.RequestToken(deviceCode.DeviceCode)
		if err != nil {
			fmt.Printf("Falied to fetch token Error:%s", err)
		}

		// If authorization is still pending in browser.
		if Token.Error == "authorization_pending" {
			// This is time interval we can poll for token as per auth0 docs.
			time.Sleep(5 * time.Second)
			continue
		}

		// If device code is expired.
		if Token.Error == "expired_token" {
			s.Stop()
			fmt.Printf("\nYou have not authorized the device quickly, and device code expired.\n")
			fmt.Printf("Login again using `appctl login`!!\n")
			break
		}

		// If access is Denied.
		if Token.Error == "access_denied" {
			s.Stop()
			fmt.Printf("\n" + color.Red("Access Denied!!") + "\n")
			break
		}

		// If token is fetched, then write to config.
		if Token.AccessToken != "" {
			s.Stop()
			/* Things Yet to be implemented:
			-- To create a config file to store access token and its expire.*/
			fmt.Printf("\n" + color.Green("✔ ") + "Successfully Logged in!!\n")
			break
		}
	}

	//FetchUserinfo() function needed to be implemented to get user details.

	return nil
}

// To get a detailed information of particular app by name.
func DeleteApp(
	name string, // app name
	nameSpace string, // namespace to list apps.
) error {
	if name == "" {
		return fmt.Errorf("App Name not specified.")
	}
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.")
	}
	// Fetch the detailedapp information for given name from given namespace.
	err := appAPIs.DeleteAppByName(name, nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.", err)
	}

	return nil
}
