package appManageAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	isconnect "github.com/alimasyhur/is-connect"
	"github.com/briandowns/spinner"
	"github.com/golang-jwt/jwt"
	"github.com/platform9/appctl/pkg/appAPIs"
	"github.com/platform9/appctl/pkg/browser"
	"github.com/platform9/appctl/pkg/constants"
	"github.com/platform9/appctl/pkg/segment"
	"github.com/ryanuber/columnize"
)

// Config structure for configfile.
type Config struct {
	IDToken   string
	ExpiresAt time.Time
}

type Event struct {
	EventName string
	Status    string
	Data      []constants.ListAppInfo
	Error     string
}

// To list apps.
func ListAppsInfo() error {

	//Check Internet Connectivity
	if !isconnect.IsOnline() {
		return fmt.Errorf("Network unreachable. %v\n", constants.InternetConnectivity)
	}

	// Load config, and check if id_token expired
	config, err := loadConfig(constants.CONFIGFILEPATH)
	if err != nil {
		return fmt.Errorf("Failed to list apps. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)
	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To list and store output.
	var list constants.ListAppInfo
	var Output []string
	var event Event
	Output = append(Output, constants.TABLEFORMAT)

	// Fetch the running apps.
	list_apps, err := appAPIs.ListApps(config.IDToken)
	if err != nil {
		//Event is Failure.
		event.EventName = "List-Apps"
		event.Status = "Failure"
		event.Error = err.Error()
		send(event, nil)
		return fmt.Errorf("Failed to list apps with error: %v\n", err)
	}

	// Fetch the ListAppInfo for apps deployed.
	for _, items := range list_apps["items"].([]interface{}) {
		for key, appInfo := range items.(map[string]interface{}) {
			if key == "metadata" && appInfo != nil {
				list.Name = fmt.Sprintf("%v", appInfo.(map[string]interface{})["name"])
				creationTime := fmt.Sprintf("%v", appInfo.(map[string]interface{})["creationTimestamp"])
				list.CreationTime = appAge(creationTime)
			}
			// Fetch the Image name.
			if key == "spec" && appInfo != nil {
				template := (appInfo.(map[string]interface{}))["template"].(map[string]interface{})
				if template != nil {
					detailedSpec := template["spec"].(map[string]interface{})
					if detailedSpec != nil {
						containers := detailedSpec["containers"].([]interface{})[0]
						list.Image = fmt.Sprintf("%v", containers.(map[string]interface{})["image"])
					}
				}
			}

			// Fetch the URL Endpoint
			if key == "status" && appInfo != nil {
				for key, url := range appInfo.(map[string]interface{}) {
					if key == "url" {
						list.URL = fmt.Sprintf("%v", url)
					}
				}
				conditions := appInfo.(map[string]interface{})["conditions"]
				if conditions != nil {
					readyStatus := conditions.([]interface{})[1].(map[string]interface{})["status"]
					list.ReadyStatus = fmt.Sprintf("%v", readyStatus)
					list.Reason = getResponseMessage(conditions)
				}

			}
		}
		event.Data = append(event.Data, list)
		appinfo := fmt.Sprintf("%v | %v | %v | %v | %v | %v", list.Name, list.URL, list.Image, list.ReadyStatus, list.CreationTime, list.Reason)
		Output = append(Output, appinfo)
	}

	//Event is successful.
	event.EventName = "List-Apps"
	event.Status = "Success"
	send(event, nil)
	tabularAppInfo := columnize.SimpleFormat(Output)
	fmt.Println(tabularAppInfo)
	return nil
}

// To create an app.
func CreateApp(
	name string, // App name to create.
	image string, // Source Image to create app.
	username string, // User name in case of private container registry
	password string, // Password in case of private container registry
	env []string, // Environment varialbes of app.
	port string, // Port where application listens on.
) error {
	if name == "" || image == "" {
		return fmt.Errorf("Either or both of app name and image not specified.\n")
	}

	//Check Internet Connectivity
	if !isconnect.IsOnline() {
		return fmt.Errorf("Network unreachable. %v\n", constants.InternetConnectivity)
	}

	// Load config, and check if id_token expired
	config, err := loadConfig(constants.CONFIGFILEPATH)
	if err != nil {
		return fmt.Errorf("Failed to deploy app. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)
	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// To check if app with same name already exists.
	appExists, err := appAPIs.GetAppByName(name, config.IDToken)
	if err == nil && appExists != nil {
		return fmt.Errorf("App with same name already exists!! Please use different name.\n")
	}

	// Send Segment Event
	var event Event
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("red")
	s.Start()
	s.Suffix = " Deploying app.."

	errCreate := appAPIs.CreateApp(name, image, username, password, env, port, config.IDToken)
	if errCreate != nil {
		//Event is Failure.
		event.EventName = "Deploy-App"
		event.Status = "Failure"
		event.Error = errCreate.Error()
		send(event, nil)
		s.Stop()
		if errCreate.Error() == constants.MaxAppDeployLimitError {
			return fmt.Errorf("%v\n", errCreate.Error())

		}
		return fmt.Errorf("%v\n", errCreate)
	}

	time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
	// Polling to fetch URL if app is deployed.
	var count = 0
	var status, securedAppURL bool
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
			//Event is Failure.
			event.EventName = "Deploy-App"
			event.Status = "Failure"
			event.Error = invalidImage
			send(event, get_app)
			s.Stop()
			return fmt.Errorf("%v %v.\nPlease check if the application image path provided is valid, and is from a public registry.\n", invalidImage, image)
		}
		if !status {
			// Wait until stauts of app deployed is ready and true.
			time.Sleep(constants.APPDEPLOYINTERVAL * time.Second)
			continue
		}

		// URL Endpoint where the app service is available.
		url := (get_app["status"]).(map[string]interface{})["url"]
		if url != nil && status {

			// Check if app url is secured.
			securedAppURL = checkSecuredURL(url)
			if !securedAppURL {
				time.Sleep(constants.SECUREENDPOINT * time.Second)
				continue
			}

			s.Stop()
			fmt.Printf("\nApp %v is deployed and can be accessed at URL: %v\n", name, url)
			//Event is Successful.
			event.EventName = "Deploy-App"
			event.Status = "Success"
			send(event, get_app)
			return nil
		} else {
			s.Stop()
			fmt.Printf("\nApp deploy taking time. Check latest status by running command `appctl list`.\n")
			return nil
		}
	}
	if !status || !securedAppURL {
		s.Stop()
		fmt.Printf("\nApp deploy taking time. Check latest status by running command `appctl list`.\n")
	}
	return nil
}

// Check if all three status are true and ready.
func checkStatusReady(get_app map[string]interface{}) (bool, string) {
	var configurationStatus, readyStatus, routeStatus string

	if get_app["status"] != nil {
		conditions := get_app["status"].(map[string]interface{})["conditions"]
		if conditions != nil {
			configurationStatus = fmt.Sprintf("%s", conditions.([]interface{})[0].(map[string]interface{})["status"])
			readyStatus = fmt.Sprintf("%s", conditions.([]interface{})[1].(map[string]interface{})["status"])
			routeStatus = fmt.Sprintf("%s", conditions.([]interface{})[2].(map[string]interface{})["status"])

			// Check if Image given is invalid
			configurationMessage := fmt.Sprintf("%s", conditions.([]interface{})[0].(map[string]interface{})["message"])
			if strings.Contains(configurationMessage, constants.InvalidImage) {
				return false, constants.InvalidImage
			}
		}
	}

	if configurationStatus == "True" && readyStatus == "True" && routeStatus == "True" {
		return true, ""
	}

	return false, ""
}

// To get a detailed information of particular app by name.
func GetAppByNameInfo(
	name string, // app name
) error {
	if name == "" {
		return fmt.Errorf("App name not specified.\n")
	}

	//Check Internet Connectivity
	if !isconnect.IsOnline() {
		return fmt.Errorf("Network unreachable. %v\n", constants.InternetConnectivity)
	}

	// Load config, and check if id_token expired
	config, err := loadConfig(constants.CONFIGFILEPATH)
	if err != nil {
		return fmt.Errorf("Failed to get app information. Please login using command `appctl login`.\n")
	}

	// Check if Token is expired or not.
	expired, _ := checkTokenExpired(config.IDToken)

	if expired {
		return fmt.Errorf("Login expired. Please login again using command `appctl login`\n")
	}

	// Send Segment Event
	var event Event

	// Fetch the detailedapp information for given appname.
	get_app, err := appAPIs.GetAppByName(name, config.IDToken)
	if err != nil {
		//Event is Failure.
		event.EventName = "Describe-App"
		event.Status = "Failure"
		event.Error = err.Error()
		send(event, get_app)
		return fmt.Errorf("Failed to get app information with error: %v\nCheck 'appctl list' for more information on apps running.\n", err)
	}

	event.EventName = "Describe-App"
	event.Status = "Success"
	send(event, get_app)
	jsonFormatted, err := json.MarshalIndent(get_app, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(jsonFormatted))
	return nil
}

// To login using Device authentication and access appctl.
func LoginApp() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("red")

	//Check Internet Connectivity
	if !isconnect.IsOnline() {
		return fmt.Errorf("Network unreachable. %v\n", constants.InternetConnectivity)
	}

	fmt.Printf("Starting login process.\n")
	// Get the device code.
	deviceCode, err := appAPIs.GetDeviceCode()
	if err != nil {
		return fmt.Errorf("Unable to generate device code.\nError: %v\n", err)
	}

	fmt.Printf("Device verification is required to continue login.\n")
	fmt.Printf("Your Device Confirmation code is: %v\n", deviceCode.UserCode)

	// To open browser, for device verification and SSO.
	err = browser.OpenBrowser(deviceCode.VerificationUrlComplete)
	if err != nil {
		fmt.Printf("\nCouldn't open the URL, kindly do it manually: %v\n", deviceCode.VerificationUrlComplete)
	}

	var Token *appAPIs.TokenInfo

	// Wait for device verification in browser and if its success request the token.
	s.Start()
	s.Suffix = " Waiting for login to complete in browser..."

	// Send Segment Event
	var event Event
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
			//Event is Failure.
			event.EventName = "Login"
			event.Status = "Failure"
			event.Error = Token.Error
			send(event, nil)
			return fmt.Errorf("\nThe device code was expired as the app was not authorized in time!\n" +
				"Login again using `appctl login`!!\n")
		}

		// If access is Denied.
		if Token.Error == "access_denied" {
			s.Stop()
			//Event is Failure.
			event.EventName = "Login"
			event.Status = "Failure"
			event.Error = Token.Error
			send(event, nil)
			return fmt.Errorf("\nCannot login. Please try again.\n")
		}
	}
	// To create and write to config file.
	var config = Config{
		IDToken:   Token.IdToken,
		ExpiresAt: time.Now().Add(time.Duration(Token.ExpiresIn) * time.Second),
	}

	errConfig := createConfig(config, constants.CONFIGFILEPATH)
	if errConfig != nil {
		//Event is Failure.
		event.EventName = "Login"
		event.Status = "Failure"
		event.Error = errConfig.Error()
		send(event, nil)
		return fmt.Errorf("Cannot login. Please try again.\n")
	}
	// Send info to app-controller api.
	errLogin := appAPIs.Login(config.IDToken)
	if errLogin != nil {
		removeConfig(constants.CONFIGFILEPATH)
		//Event is Failure.
		event.EventName = "Login"
		event.Status = "Failure"
		event.Error = errLogin.Error()
		send(event, nil)
		return fmt.Errorf("\nCannot login!! Backend server is down. Please try later.\n")
	}

	event.EventName = "Login"
	event.Status = "Success"
	send(event, nil)
	fmt.Printf("\nSuccessfully logged in!!\n")
	return nil
}

// To get a detailed information of particular app by name.
func DeleteApp(
	name string, // app name
) error {
	if name == "" {
		return fmt.Errorf("App name not specified.\n")
	}

	//Check Internet Connectivity
	if !isconnect.IsOnline() {
		return fmt.Errorf("Network unreachable. %v\n", constants.InternetConnectivity)
	}

	// Load config, and check if id_token expired
	config, err := loadConfig(constants.CONFIGFILEPATH)
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
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.\n", errApp)
	}
	// Fetch app info prior to deletion.
	var event Event
	appInfo, _ := fetchAppInfo(get_app)
	event.Data = append(event.Data, *appInfo)

	// Fetch the detailedapp information for given appname.
	errDel := appAPIs.DeleteAppByName(name, config.IDToken)
	if errDel != nil {
		//Event is Failure.
		event.EventName = "Delete-App"
		event.Status = "Failure"
		event.Error = errDel.Error()
		send(event, nil)
		return fmt.Errorf("Failed to delete app with error: %v\nCheck 'appctl list' for more information on apps running.\n", errDel)
	}

	// Send Segment Event
	event.EventName = "Delete-App"
	event.Status = "Success"
	send(event, nil)
	return nil
}

func createConfig(config Config, configFilePath string) error {
	//Create the pf9 config directory to store configfile
	err := createDirectoryIfNotExist(constants.CONFIGDIR)
	if err != nil {
		return fmt.Errorf("Failed to create config directory!!")
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	// Write data in to file
	if err := ioutil.WriteFile(configFilePath, data, 0600); err != nil {
		return err
	}
	return nil
}

func createDirectoryIfNotExist(configPath string) error {
	// Create a pf9 directory
	var err error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if errdir := os.MkdirAll(configPath, 0700); errdir != nil {
			return errdir
		}
	}
	return err
}

func loadConfig(configFilePath string) (*Config, error) {
	config, _ := ioutil.ReadFile(configFilePath)

	readConfig := Config{}

	err := json.Unmarshal([]byte(config), &readConfig)
	if err != nil {
		return &Config{}, fmt.Errorf("Failed to parse config with error: %s", err)
	}
	return &readConfig, nil
}

func removeConfig(configFilePath string) error {
	err := os.Remove(configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to remove config file")
	}
	return nil
}

// To send a segment event.
func send(event Event, get_app map[string]interface{}) error {
	// Create a new Segment client
	client, err := segment.SegmentClient()
	if err != nil {
		return err
	}

	defer segment.Close(client)

	// Segment event for List Apps
	if event.EventName == "List-Apps" || event.EventName == "Login" {
		userId, loginType, _ := fetchUserId()
		if err := segment.SendEventList(client, event.EventName, userId, event.Status, loginType, event.Error, event.Data); err != nil {
			return fmt.Errorf("Failed to send segment event. Error: %v\n", err)
		}

	} else {
		//Segment events for Deploy, describe, delete app.
		if get_app != nil {
			appInfo, _ := fetchAppInfo(get_app)
			event.Data = append(event.Data, *appInfo)
		}
		// Fetch the UserID and loginType
		userId, loginType, _ := fetchUserId()
		if err := segment.SendEvent(client, event.EventName, userId, event.Status, loginType, event.Error, event.Data); err != nil {
			return fmt.Errorf("Failed to send segment event. Error: %v\n", err)
		}
	}

	return nil
}

// To fetch UserID, and login type after basic validation of token.
func fetchUserId() (string, string, error) {

	// Load config, and fetch the IDToken
	config, err := loadConfig(constants.CONFIGFILEPATH)
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
func fetchAppInfo(get_app map[string]interface{}) (*constants.ListAppInfo, error) {
	// Fetch AppName, URL, Image, ReadyStatus, Creation Time from app information.
	var name, creationTime, url, image, readyStatus, port, reason string

	if get_app["metadata"] != nil {
		name = fmt.Sprintf("%v", (get_app["metadata"]).(map[string]interface{})["name"])
		creationTime = fmt.Sprintf("%v", (get_app["metadata"]).(map[string]interface{})["creationTimestamp"])
	}

	if get_app["status"] != nil {
		// Fetch URL
		url = fmt.Sprintf("%v", (get_app["status"]).(map[string]interface{})["url"])

		//Fetch app status.
		conditions := get_app["status"].(map[string]interface{})["conditions"]
		if len(conditions.([]interface{})) > 2 {
			reason = getResponseMessage(conditions)
			readyStatus = fmt.Sprintf("%s", conditions.([]interface{})[1].(map[string]interface{})["status"])
		}
	}

	if get_app["spec"] != nil {
		template := (get_app["spec"].(map[string]interface{}))["template"].(map[string]interface{})
		if template != nil {
			detailedSpec := template["spec"].(map[string]interface{})
			if detailedSpec != nil {
				containers := detailedSpec["containers"].([]interface{})[0]
				// Fetch Image.
				image = fmt.Sprintf("%v", containers.(map[string]interface{})["image"])

				//Fetch Container Port.
				port = getPort(containers.(map[string]interface{}))
			}

		}
	}

	return &constants.ListAppInfo{Name: name, URL: url, Image: image, Port: port, ReadyStatus: readyStatus, CreationTime: creationTime, Reason: reason}, nil
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

// Get the container port
func getPort(container map[string]interface{}) string {
	for key, value := range container {
		if key == "ports" {
			port := value.([]interface{})[0].(map[string]interface{})["containerPort"]
			return fmt.Sprintf("%v", port)
		}
	}
	return ""
}

//appAge gives the age of app since its creation.
func appAge(appCreationTime string) string {
	appCreatedTimeParsed, err := time.Parse(constants.UTCClusterTimeStamp, appCreationTime)
	if err != nil {
		// If can't parse then return same UTC time stamp.
		return appCreationTime
	}
	currentTime := time.Now()
	appCreateTime := currentTime.Sub(appCreatedTimeParsed).String()
	// Consider only seconds, exclude micro/nano seconds.
	return fmt.Sprintf("%ss", strings.Split(appCreateTime, ".")[0])
}

//Get the response message for deployed apps.
func getResponseMessage(conditions interface{}) string {
	readyStatus := conditions.([]interface{})[1].(map[string]interface{})["status"]
	//fmt.Printf("The ready status is %v\n", readyStatus)
	if readyStatus != "True" {
		configurationMessage := fmt.Sprintf("%s", conditions.([]interface{})[0].(map[string]interface{})["message"])
		if strings.Contains(configurationMessage, constants.InvalidImage) {
			return constants.InvalidImage
		}
		readyMessage := fmt.Sprintf("%v", conditions.([]interface{})[1].(map[string]interface{})["message"])
		readyReason := fmt.Sprintf("%v", conditions.([]interface{})[1].(map[string]interface{})["reason"])
		return fmt.Sprintf("%v  %v", readyReason, readyMessage)
	}
	return "nil"
}

// Check if app url is secured.
func checkSecuredURL(url interface{}) bool {
	urlString := fmt.Sprintf("%v", url)
	return strings.Contains(urlString, constants.HTTPS)
}
