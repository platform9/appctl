package appManageAPI

import (
	"fmt"

	"github.com/platform9/pf9-appctl/pkg/appAPIs"
	"github.com/platform9/pf9-appctl/pkg/constants"
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

// To create an app.
func CreateApp(
	name string, // App name to create.
	nameSpace string, // namespace to list apps.
	image string, // Source Image to create app.
) error {
	if name == "" {
		return fmt.Errorf("App Name not specified.")
	}
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.")
	}
	if image == "" {
		return fmt.Errorf("Image not specified.")
	}
	fmt.Printf("Started Creating App\n")
	fmt.Printf("\nName is %v\nNamespace is %v\nImage is %v\n", name, nameSpace, image)
	// Fetch the running apps in given namespace.

	err := appAPIs.CreateApp(name, nameSpace, image)
	if err != nil {
		return fmt.Errorf("Failed to create app with error: %v", err)
	}
	fmt.Printf("App created with Name: %v. Run 'pf9-appctl app list' to get more information on app\n", name)
	return nil
}
