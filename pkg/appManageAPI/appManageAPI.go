package appManageAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// Config structure for configfile.
type CONFIG struct {
	Domain    string
	IDToken   string
	ExpiresAt time.Time
	Scope     []string
}

// To list apps.
func ListAppsInfo(
	nameSpace string, // namespace to list apps.
) error {
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.\n")
	}

	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to list apps. Please login using command `appctl login`.\n")
	}

	if config.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To list and store output.
	var list ListAppInfo
	var Output []string
	Output = append(Output, constants.TABLEFORMAT)

	// Fetch the running apps in given namespace.
	list_apps, err := appAPIs.ListApps(nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to list apps with error: %v\n", err)
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
	nameSpace := "default"
	return nameSpace, nil
}

// To create an app.
func CreateApp(
	name string, // App name to create.
	nameSpace string, // namespace to list apps.
	image string, // Source Image to create app.
) error {
	if name == "" || image == "" {
		return fmt.Errorf("Either or both of App Name and Image not specified.\n")
	}
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.\n")
	}
	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to create app. Please login using command `appctl login`.\n")
	}

	if config.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To check if app with same name already exists.
	appExists, err := appAPIs.GetAppByName(name, nameSpace)
	if err == nil && appExists != nil {
		return fmt.Errorf("App with same name already exists!! Please use different name.")
	}

	fmt.Printf("Deploying app..\n")
	// Fetch the running apps in given namespace.

	errcreate := appAPIs.CreateApp(name, nameSpace, image)
	if errcreate != nil {
		return fmt.Errorf("Failed to create app with error: %v", errcreate)
	}

	time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
	// Polling to fetch URL if app is deployed.
	var count = 0
	for count <= constants.APPDEPLOYINTERVAL {
		count++
		// Fetch the detailedapp information for given name from given namespace.
		get_app, err := appAPIs.GetAppByName(name, nameSpace)
		if err != nil {
			time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
			continue
		}

		// URL/ Endpoint where the app service is available.
		url := (get_app["status"]).(map[string]interface{})["url"]
		if url != nil {
			//Since creation of App takes some time.
			fmt.Printf("App " + color.Yellow(name) + " is deployed and can be accessed at URL: " + color.Yellow(url) + "\n")
			return nil
		}
	}
	fmt.Printf("App deploy taking time. Check latest status by running command `appctl list`.\n")

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

	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to get app information. Please login using command `appctl login`.\n")
	}

	if config.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// Fetch the detailedapp information for given name from given namespace.
	get_app, err := appAPIs.GetAppByName(name, nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to get app information with error: %v\nCheck 'appctl list' for more information on apps running.", err)
	}
	jsonFormatted, err := json.MarshalIndent(get_app, "", "  ")
	fmt.Printf("%v\n", string(jsonFormatted))
	return nil
}

// To login using Device authentication and access appctl.
func LoginApp() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("red")

	fmt.Printf(color.Blue("Starting login process.") + "\n")
	// Get the device code.
	deviceCode, err := appAPIs.GetDeviceCode()
	if err != nil {
		fmt.Printf("Unable to generate device code.")
	}

	fmt.Printf("Device verification is required to continue login.\n")
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
		// Request for token polling.
		Token, err = appAPIs.RequestToken(deviceCode.DeviceCode)
		if err != nil {
			fmt.Printf("Falied to fetch token Error:%s", err)
		}

		// If token is fetched, then next write to config.
		if Token.IdToken != "" {
			s.Stop()
			break
		}

		// If authorization is still pending in browser.
		if Token.Error == "authorization_pending" {
			// This is time interval we can poll for token as per auth0 docs.
			time.Sleep(constants.TOKENPOLLINTERVAL * time.Second)
			continue
		}

		// If device code is expired.
		if Token.Error == "expired_token" {
			s.Stop()
			return fmt.Errorf("\nThe device code was expired as the app was not authorized in time!\n" +
				"Login again using `appctl login`!!\n")
		}

		// If access is Denied.
		if Token.Error == "access_denied" {
			s.Stop()
			return fmt.Errorf("\n" + color.Red("Cannot login. Please try again.") + "\n")
		}
	}
	// To create and write to config file.
	var config = CONFIG{Domain: constants.DOMAIN,
		IDToken:   Token.IdToken,
		ExpiresAt: time.Now().Add(time.Duration(Token.ExpiresIn) * time.Second),
		Scope:     constants.RequiredScopes,
	}

	errConfig := CreateConfig(config)
	if errConfig != nil {
		return fmt.Errorf("Cannot login. Please try again.\n")
	}
	// Send info to fast-path api.
	appAPIs.Login(config.IDToken)

	fmt.Printf("\n" + color.Green("✔ ") + "Successfully Logged in!!\n")

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

	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to delete app. Please login using command `appctl login`.\n")
	}

	if config.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To check if app exists.
	_, errApp := appAPIs.GetAppByName(name, nameSpace)
	if errApp != nil {
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.", errApp)
	}

	// Fetch the detailedapp information for given name from given namespace.
	errDel := appAPIs.DeleteAppByName(name, nameSpace)
	if errDel != nil {
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.", errDel)
	}

	return nil
}

func CreateConfig(config CONFIG) error {
	//Create the pf9 config directory to store configfile
	err := CreateDirectoryIfNotExist()
	if err != nil {
		return fmt.Errorf("Failed to create config directory!!")
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	// Write data in to file
	if err := ioutil.WriteFile(constants.CONFIGFILEPATH, data, 0600); err != nil {
		return err
	}
	return nil
}

func CreateDirectoryIfNotExist() error {
	// Create a pf9 directory
	var err error
	if _, err := os.Stat(constants.CONFIGDIR); os.IsNotExist(err) {
		if errdir := os.MkdirAll(constants.CONFIGDIR, 0700); errdir != nil {
			return errdir
		}
	}
	return err
}

func LoadConfig() (*CONFIG, error) {
	config, _ := ioutil.ReadFile(constants.CONFIGFILEPATH)

	readConfig := CONFIG{}

	err := json.Unmarshal([]byte(config), &readConfig)
	if err != nil {
		return &CONFIG{}, fmt.Errorf("Failed to parse config with error: %s", err)
	}
	return &readConfig, nil
}
