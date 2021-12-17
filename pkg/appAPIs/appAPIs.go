package appAPIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/platform9/appctl/pkg/color"
	"github.com/platform9/appctl/pkg/constants"
)

const (
	MaxAppDeployStatusCode = "429"
)

// Type definition for struct encapsulating app manager APIs.
type appAPI struct {
	Client  *http.Client
	BaseURL string
}

//Environmnet variable struct.
type Env struct {
	key   string
	value string
}

// To store device information fetched during device authorization.
type DeviceInfo struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	Verification_URL        string `json:"verification_uri"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VerificationUrlComplete string `json:"verification_uri_complete"`
}

// To store token information fetched during retrive token.
type TokenInfo struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IdToken          string `json:"id_token"`
	Scope            string `json:"scope"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// To fetch the information (list,create,device,token).
var (
	listAppsInfo  map[string]interface{}
	getAppInfo    map[string]interface{}
	getDeviceInfo DeviceInfo
	getTokenInfo  TokenInfo
)

// API to list/get all apps.
func (cli_api *appAPI) ListAppsAPI(token string) ([]byte, error) {

	req, err := http.NewRequest("GET", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	idToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", idToken)

	resp, err := cli_api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return data, nil
}

// To get all the apps information.
func ListApps(token string) (map[string]interface{}, error) {
	// Endpoint to list apps.
	url := fmt.Sprintf(constants.APPURL)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	list_apps, err := cli_api.ListAppsAPI(token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(list_apps), &listAppsInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return listAppsInfo, nil
}

// API to list/get all apps.
func (cli_api *appAPI) CreateAppAPI(createInfo string, token string) ([]byte, error) {
	payload := strings.NewReader(fmt.Sprintf("%s", createInfo))
	req, err := http.NewRequest("POST", cli_api.BaseURL, payload)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	idToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", idToken)

	req.Header.Add("Content-Type", "application/json")
	resp, err := cli_api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer resp.Body.Close()

	return []byte(strconv.Itoa(resp.StatusCode)), nil
}

// To get all the apps information.
func CreateApp(name string, image string, env []string, port string, token string) error {
	// Endpoint to list apps.
	url := fmt.Sprintf(constants.APPURL)
	var createInfo string
	if env != nil {
		if port != "" {
			createInfo = fmt.Sprintf(`{"name":"%s", "image":"%s", "port": "%s", "envs": %v}`, name, image, port, GenEnvSlice(env))
		} else {
			createInfo = fmt.Sprintf(`{"name":"%s", "image":"%s", "envs": %v}`, name, image, GenEnvSlice(env))
		}
	} else {
		if port != "" {
			createInfo = fmt.Sprintf(`{"name":"%s", "image":"%s", "port": "%s"}`, name, image, port)
		}
		createInfo = fmt.Sprintf(`{"name":"%s", "image":"%s"}`, name, image)
	}

	client := &http.Client{}

	cli_api := appAPI{client, url}
	data, err := cli_api.CreateAppAPI(createInfo, token)
	if err != nil {
		return err
	}
	if string(data) == MaxAppDeployStatusCode {
		return fmt.Errorf("Maximum apps deploy limit reached!!")
	}
	return nil
}

// API to get a particular app by name.
func (cli_api *appAPI) GetAppByNameAPI(token string) ([]byte, error) {

	req, err := http.NewRequest("GET", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	idToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", idToken)

	resp, err := cli_api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request processing failed with error: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response. Error: %v", err)
	}
	return data, nil
}

// To get a particular app information.
func GetAppByName(appName string, token string) (map[string]interface{}, error) {
	// Endpoint to get a particular app.
	url := fmt.Sprintf(constants.APPURL+"/%s", appName)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	get_app, err := cli_api.GetAppByNameAPI(token)
	if err != nil {
		return nil, err
	}
	// If incorrect app name is given, then empty response.
	if len(get_app) == 0 {
		return nil, fmt.Errorf("Cannot find the app " + color.Yellow(appName) + "!!")
	}
	err = json.Unmarshal([]byte(get_app), &getAppInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the response. Error: %s", err)
	}
	return getAppInfo, nil
}

// To get device code for login.
func (cli_api *appAPI) GetDeviceCodeAPI(getDevice string) ([]byte, error) {

	payload := strings.NewReader(getDevice)
	req, err := http.NewRequest("POST", cli_api.BaseURL, payload)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return body, nil
}

func GetDeviceCode() (*DeviceInfo, error) {
	// Endpoint to get device code and verification url.
	url := fmt.Sprintf("%s", constants.DEVICECODEURL)

	deviceRequest := fmt.Sprintf("%s", constants.DEVICEREQUESTPAYLOAD)
	client := &http.Client{}

	cli_api := appAPI{client, url}

	deviceInfo, err := cli_api.GetDeviceCodeAPI(deviceRequest)
	if err != nil {
		return &DeviceInfo{}, err
	}

	err = json.Unmarshal([]byte(deviceInfo), &getDeviceInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return &getDeviceInfo, nil
}

// To request token after successful device verification.
func (cli_api *appAPI) RequestTokenAPI(requestToken string) ([]byte, error) {

	payload := strings.NewReader(requestToken)

	req, err := http.NewRequest("POST", cli_api.BaseURL, payload)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return body, nil
}

func RequestToken(deviceCode string) (*TokenInfo, error) {
	// Endpoint to request for token.
	url := fmt.Sprintf("https://%s/oauth/token", constants.DOMAIN)

	deviceRequest := fmt.Sprintf("%s&device_code=%s&client_id=%s", constants.GrantType, deviceCode, constants.CLIENTID)

	client := &http.Client{}

	cli_api := appAPI{client, url}

	tokenInfo, err := cli_api.RequestTokenAPI(deviceRequest)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(tokenInfo), &getTokenInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return &getTokenInfo, nil
}

// API to delete a particular app by name.
func (cli_api *appAPI) DeleteAppByNameAPI(token string) ([]byte, error) {

	req, err := http.NewRequest("DELETE", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	idToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", idToken)

	resp, err := cli_api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response. Error: %v", err)
	}
	return data, nil
}

// To delete a particular app information.
func DeleteAppByName(appName string, token string) error {
	// Endpoint to get a particular app.
	url := fmt.Sprintf(constants.APPURL+"/%s", appName)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	_, err := cli_api.DeleteAppByNameAPI(token)
	if err != nil {
		return err
	}

	return nil
}

// Login API
func (cli_api *appAPI) LoginAPI(token string) ([]byte, error) {
	req, err := http.NewRequest("POST", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	idToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", idToken)

	req.Header.Add("Content-Type", "application/json")

	resp, err := cli_api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return data, nil
}

// Login app
func Login(token string) error {
	// Endpoint to login.
	url := fmt.Sprintf(constants.APPURL + "/login")

	client := &http.Client{}

	cli_api := appAPI{client, url}
	_, err := cli_api.LoginAPI(token)
	if err != nil {
		return err
	}
	return nil
}

// Generate environemnt slice as per create command. [{ "key":"ENV1", "value":"val1"}, { "key":"ENV2", "value":"val2"}]
func GenEnvSlice(env []string) []string {
	var envSlice []string

	if env != nil {
		for _, value := range env {
			splitEnv := strings.Split(value, "=")
			envSlice = append(envSlice, fmt.Sprintf(`{"key": "%v", "value": "%v"}`, splitEnv[0], splitEnv[1]))
		}
	}
	for count := 0; count < len(envSlice)-1; count++ {
		envSlice[count] = envSlice[count] + ","
	}
	return envSlice
}
