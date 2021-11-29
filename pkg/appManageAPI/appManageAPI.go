package appManageAPI

import (
	"fmt"

	"github.com/platform9/pf9-appctl/pkg/appAPIs"
)

// To list apps.
func ListAppsInfo(
	nameSpace string, // namespace to list apps.
) error {
	if nameSpace == "" {
		return fmt.Errorf("Namespace not specified.")
	}

	// Fetch the running apps in given namespace.
	list_apps, err := appAPIs.ListApps(nameSpace)
	if err != nil {
		return fmt.Errorf("Failed to list apps with error: %v", err)
	}

	fmt.Printf("*************")
	fmt.Printf("\nThe full apps are \n\n%+v\n", list_apps)

	/* To add output formating.
	for _, apps := range list_apps.(map[string]interface{}) {
		fmt.Printf("The apps are %v", apps)
	}*/
	return nil
}
