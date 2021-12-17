package appManageAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/golang-jwt/jwt"
	"github.com/platform9/appctl/pkg/appAPIs"
	"github.com/platform9/appctl/pkg/browser"
	"github.com/platform9/appctl/pkg/color"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/platform9/appctl/pkg/segment"
	"github.com/ryanuber/columnize"
)

type ListAppInfo struct {
	Name         string
	URL          string
	Image        string
	ReadyStatus  string
	CreationTime string
}

// Config structure for configfile.
type CONFIG struct {
	Domain    string
	IDToken   string
	ExpiresAt time.Time
	Scope     []string
}

type Event struct {
	EventName string
	Status    string
	Data      []ListAppInfo
}

// To list apps.
func ListAppsInfo() error {
	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to list apps. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)

	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To list and store output.
	var list ListAppInfo
	var Output []string
	var event Event
	Output = append(Output, constants.TABLEFORMAT)

	// Fetch the running apps.
	list_apps, err := appAPIs.ListApps(config.IDToken)
	if err != nil {
		return fmt.Errorf("Failed to list apps with error: %v\n", err)
	}

	// Fetch the ListAppInfo for apps deployed.
	for _, items := range list_apps["items"].([]interface{}) {
		for key, appInfo := range items.(map[string]interface{}) {
			if key == "metadata" {
				list.Name = fmt.Sprintf("%v", appInfo.(map[string]interface{})["name"])
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
				conditions := appInfo.(map[string]interface{})["conditions"]
				readyStatus := conditions.([]interface{})[1].(map[string]interface{})["status"]
				list.ReadyStatus = fmt.Sprintf("%v", readyStatus)
			}
		}
		event.Data = append(event.Data, list)
		appinfo := fmt.Sprintf("%v | %v | %v | %v | %v", list.Name, list.URL, list.Image, list.ReadyStatus, list.CreationTime)
		Output = append(Output, appinfo)
	}

	//Event is successfull.
	event.EventName = "List-Apps"
	event.Status = "Success"
	errEvent := Send(event, nil)
	if errEvent != nil {
		//Should add as log message
		//fmt.Printf("%v", errEvent)
	}

	tabularAppInfo := columnize.SimpleFormat(Output)
	fmt.Println(tabularAppInfo)
	return nil
}

// To create an app.
func CreateApp(
	name string, // App name to create.
	image string, // Source Image to create app.
	env []string, // Environment varialbes of app.
	port string, // Port where application listens on.
) error {
	if name == "" || image == "" {
		return fmt.Errorf("Either or both of App Name and Image not specified.\n")
	}
	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to create app. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)

	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To check if app with same name already exists.
	appExists, err := appAPIs.GetAppByName(name, config.IDToken)
	if err == nil && appExists != nil {
		return fmt.Errorf("App with same name already exists!! Please use different name.")
	}

	fmt.Printf("Deploying app..\n")
	errCreate := appAPIs.CreateApp(name, image, env, port, config.IDToken)
	if errCreate != nil {
		if errCreate.Error() == constants.MaxAppDeployLimitError {
			return errCreate
		}
		return fmt.Errorf("Failed to create app with error: %v", errCreate)
	}

	time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
	// Polling to fetch URL if app is deployed.
	var count = 0
	var status = false
	var invalidImage string
	for count <= constants.APPDEPLOYINTERVAL {
		count++
		// Fetch the detailedapp information for given appname.
		get_app, err := appAPIs.GetAppByName(name, config.IDToken)
		if err != nil {
			time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
			continue
		}
		// It takes time to get all routes, configuration, ready state up and running.
		status, invalidImage = checkStatusReady(get_app)
		if invalidImage != "" {
			return fmt.Errorf("%v %v\n", invalidImage, image)
		}
		if !status {
			// Wait until stauts of app deployed is ready and true.
			time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
			continue
		}

		// URL Endpoint where the app service is available.
		url := (get_app["status"]).(map[string]interface{})["url"]
		if url != nil && status {
			fmt.Printf("App " + color.Yellow(name) + " is deployed and can be accessed at URL: " + color.Yellow(url) + "\n")

			// Send Segment Event
			var event Event
			event.EventName = "Deploy-App"
			event.Status = "Success"
			errEvent := Send(event, get_app)
			if errEvent != nil {
				// Should add as log message
				//fmt.Printf("%v", errEvent)
			}
			return nil
		} else {
			fmt.Printf("App deploy taking time. Check latest status by running command `appctl list`.\n")
			return nil
		}
	}
	if !status {
		fmt.Printf("App deploy taking time. Check latest status by running command `appctl list`.\n")
	}
	return nil
}

// Check if all three status are true and ready.
func checkStatusReady(get_app map[string]interface{}) (bool, string) {
	conditions := get_app["status"].(map[string]interface{})["conditions"]
	configurationStatus := fmt.Sprintf("%s", conditions.([]interface{})[0].(map[string]interface{})["status"])
	readyStatus := fmt.Sprintf("%s", conditions.([]interface{})[1].(map[string]interface{})["status"])
	routeStatus := fmt.Sprintf("%s", conditions.([]interface{})[2].(map[string]interface{})["status"])
	if configurationStatus == "True" && readyStatus == "True" && routeStatus == "True" {
		return true, ""
	}
	configurationMessage := fmt.Sprintf("%s", conditions.([]interface{})[0].(map[string]interface{})["message"])
	if strings.Contains(configurationMessage, constants.InvalidImage) {
		return false, constants.InvalidImage
	}
	return false, ""
}

// To get a detailed information of particular app by name.
func GetAppByNameInfo(
	name string, // app name
) error {
	if name == "" {
		return fmt.Errorf("App Name not specified.")
	}
	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to get app information. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)

	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// Fetch the detailedapp information for given appname.
	get_app, err := appAPIs.GetAppByName(name, config.IDToken)
	if err != nil {
		return fmt.Errorf("Failed to get app information with error: %v\nCheck 'appctl list' for more information on apps running.", err)
	}

	// Send Segment Event
	var event Event
	event.EventName = "Describe-App"
	event.Status = "Success"
	errEvent := Send(event, get_app)
	if errEvent != nil {
		// Should add as log message
		//fmt.Printf("%v", errEvent)
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
	errLogin := appAPIs.Login(config.IDToken)
	if errLogin != nil {
		err := RemoveConfig()
		if err != nil {
			//Should add in log message.
			//fmt.Printf("Failed to remove config")
		}
		return fmt.Errorf("\nCannot login!! Backend server is down. Please try later.\n")
	}

	// Send Segment Event
	var event Event
	event.EventName = "Login"
	event.Status = "Success"
	errEvent := Send(event, nil)
	if errEvent != nil {
		// Should add as log message
		//fmt.Printf("%v", errEvent)
	}

	fmt.Printf("\n" + color.Green("✔ ") + "Successfully Logged in!!\n")

	return nil
}

// To get a detailed information of particular app by name.
func DeleteApp(
	name string, // app name
) error {
	if name == "" {
		return fmt.Errorf("App Name not specified.")
	}
	// Load config, and check if id_token expired
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to delete app. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)

	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To check if app exists.

	get_app, errApp := appAPIs.GetAppByName(name, config.IDToken)
	if errApp != nil {
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.", errApp)
	}
	// Fetch app info prior to deletion.
	var event Event
	appInfo, _ := FetchAppInfo(get_app)
	event.Data = append(event.Data, *appInfo)

	// Fetch the detailedapp information for given appname.
	errDel := appAPIs.DeleteAppByName(name, config.IDToken)
	if errDel != nil {
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.", errDel)
	}

	// Send Segment Event
	event.EventName = "Delete-App"
	event.Status = "Success"
	errEvent := Send(event, nil)
	if errEvent != nil {
		// Should add as log message
		//fmt.Printf("%v", errEvent)
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

func RemoveConfig() error {
	err := os.Remove(constants.CONFIGFILEPATH)
	if err != nil {
		return fmt.Errorf("Failed to remove config file")
	}
	return nil
}

// To send a segment event.
func Send(event Event, get_app map[string]interface{}) error {
	// Create a new Segment client
	client, err := segment.SegmentClient()
	if err != nil {
		return err
	}

	if get_app != nil {
		appInfo, _ := FetchAppInfo(get_app)
		event.Data = append(event.Data, *appInfo)
	}

	defer segment.Close(client)
	// Fetch the UserID and loginType
	userId, loginType, _ := FetchUserId()

	if err := segment.SendEvent(client, event.EventName, userId, event.Status, loginType, event.Data); err != nil {
		return fmt.Errorf("Failed to send segment event. Error: %v\n", err)
	}

	return nil
}

// To fetch UserID, and login type after basic validation of token.
func FetchUserId() (string, string, error) {

	// Load config, and fetch the IDToken
	config, err := LoadConfig()
	if err != nil {
		return "", "", fmt.Errorf("Failed to load config. Please login using command `appctl login`.\n")
	}
	// Get the token claims.
	claims, err := getTokenClaims(config.IDToken)
	if err != nil {
		return "", "", fmt.Errorf("%v", err)
	}

	var userId, loginType string

	// Email is empty if token is github login generated.
	if claims["email"] != nil {
		userId = fmt.Sprintf("%v", claims["email"])
		loginType = "google-auth"
	} else {
		userId = fmt.Sprintf("%v", claims["nickname"])
		loginType = "github"
	}

	return userId, loginType, nil
}

// To fetch App Info.
func FetchAppInfo(get_app map[string]interface{}) (*ListAppInfo, error) {
	// Fetch AppName, URL, Image, ReadyStatus, Creation Time from app information.
	name := fmt.Sprintf("%v", (get_app["metadata"]).(map[string]interface{})["name"])
	creationTime := fmt.Sprintf("%v", (get_app["metadata"]).(map[string]interface{})["creationTimestamp"])
	url := fmt.Sprintf("%v", (get_app["status"]).(map[string]interface{})["url"])
	template := (get_app["spec"].(map[string]interface{}))["template"].(map[string]interface{})
	detailedSpec := template["spec"].(map[string]interface{})
	containers := detailedSpec["containers"].([]interface{})[0]

	Image := fmt.Sprintf("%v", containers.(map[string]interface{})["image"])
	conditions := get_app["status"].(map[string]interface{})["conditions"]
	readyStatus := fmt.Sprintf("%s", conditions.([]interface{})[1].(map[string]interface{})["status"])

	return &ListAppInfo{name, url, Image, readyStatus, creationTime}, nil
}

// Basic token validation, and get claims.
func getTokenClaims(idToken string) (jwt.MapClaims, error) {
	// Parse the token.
	tokens, err := jwt.Parse(idToken, nil)
	if tokens == nil {
		//fmt.Printf("Empty with error :%v", err)
		return jwt.MapClaims{}, fmt.Errorf("Empty with error:%v", err)
	}

	//Fetch Claims
	claims, _ := tokens.Claims.(jwt.MapClaims)

	// Doing simple additional validation i.e if audiance == auth0 clientID
	if claims["aud"] != constants.CLIENTID {
		return jwt.MapClaims{}, fmt.Errorf("Token is invalid.")
	}

	return claims, nil
}

func checkTokenExpired(idToken string) (bool, error) {
	// Get the claims.
	claims, err := getTokenClaims(idToken)
	if err != nil {
		return true, fmt.Errorf("%v", err)
	}
	// Check if token is expired.
	if expiry, ok := claims["exp"].(float64); ok {
		expiryTime := time.Unix(int64(expiry), 0)
		if expiryTime.Before(time.Now()) {
			return true, nil
		}
	} else {
		return true, fmt.Errorf("Can't fetch token expiryAt time.\n")
	}
	return false, nil
}
