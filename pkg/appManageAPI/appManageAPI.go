package appManageAPI

import (
	"fmt"

	"github.com/platform9/pf9-appctl/pkg/appAPIs"
	"github.com/platform9/pf9-appctl/pkg/constants"
	"github.com/ryanuber/columnize"
	//"k8s.io/cli-runtime/pkg/printers"
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
		return fmt.Errorf("Namespace not specified.")
	}

	// To list and store output.
	var list ListAppInfo
	var Output []string
	Output = append(Output, constants.TABLEFORMAT)

	// Fetch the running apps in given namespace.
	list_apps, err := appAPIs.ListApps(nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to list apps with error: %v", err)
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
