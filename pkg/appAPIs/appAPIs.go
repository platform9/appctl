package appAPIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/platform9/pf9-appctl/pkg/constants"
)

// Type definition for struct encapsulating app manager APIs.
type appAPI struct {
	Client  *http.Client
	BaseURL string
}

// To fetch the list apps Information.
var listappsInfo map[string]interface{}

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
	url := fmt.Sprintf(constants.APPURL+"%s", nameSpace)

	client := &http.Client{}

	cli_api := appAPI{client, url}
	list_apps, err := cli_api.ListAppsAPI()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(list_apps), &listappsInfo)
	if err != nil {
		fmt.Printf("The error is %v", err)
		return nil, fmt.Errorf("Failed to Unmarshal with error: %s", err)
	}
	return listappsInfo, nil
}
