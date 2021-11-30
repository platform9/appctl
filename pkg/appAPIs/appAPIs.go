package appAPIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/platform9/pf9-appctl/pkg/constants"
)

// Type definition for struct encapsulating app manager APIs.
type appAPI struct {
	Client  *http.Client
	BaseURL string
}

// To fetch the list apps Information, create app.
var (
	listAppsInfo map[string]interface{}
	getAppInfo   map[string]interface{}
)

// API to list/get all apps.
func (cli_api *appAPI) ListAppsAPI() ([]byte, error) {

	req, err := http.NewRequest("GET", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

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
func ListApps(nameSpace string) (map[string]interface{}, error) {
	// Endpoint to list apps from a given namespace
	url := fmt.Sprintf(constants.APPURL+"/%s", nameSpace)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	list_apps, err := cli_api.ListAppsAPI()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(list_apps), &listAppsInfo)
	if err != nil {
		fmt.Printf("The error is %v", err)
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return listAppsInfo, nil
}

// API to list/get all apps.
func (cli_api *appAPI) CreateAppAPI(createInfo string) ([]byte, error) {
	payload := strings.NewReader(fmt.Sprintf("%s", createInfo))
	req, err := http.NewRequest("POST", cli_api.BaseURL, payload)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}
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

// To get all the apps information.
func CreateApp(name string, nameSpace string, image string) error {
	// Endpoint to list apps from a given namespace
	url := fmt.Sprintf(constants.APPURL)
	createInfo := fmt.Sprintf(`{"name":"%s", "space":"%s", "image":"%s"}`, name, nameSpace, image)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	_, err := cli_api.CreateAppAPI(createInfo)
	if err != nil {
		return err
	}
	return nil
}

// API to get a particular app by name.
func (cli_api *appAPI) GetAppByNameAPI() ([]byte, error) {

	req, err := http.NewRequest("GET", cli_api.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

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

// To get a particular app information.
func GetAppByName(appName string, nameSpace string) (map[string]interface{}, error) {
	// Endpoint to get a particular app from a given namespace
	url := fmt.Sprintf(constants.APPURL+"/%s/%s", nameSpace, appName)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	get_app, err := cli_api.GetAppByNameAPI()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(get_app), &getAppInfo)
	if err != nil {
		fmt.Printf("The error is %v", err)
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return getAppInfo, nil
}
