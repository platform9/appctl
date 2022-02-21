package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// usage example
var loginExample = `
  # Login using Google account/Github account to use appctl.
  appctl login
 `

// loginCmd represents "Login and use appctl".
var (
	loginCmd = &cobra.Command{
		Use:     "login",
		Short:   "Login using Google account/Github account to use appctl",
		Example: loginExample,
		Long:    `Login using Google account/Github account to use appctl`,
		Run:     loginCmdRun,
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

// To login.
func loginCmdRun(cmd *cobra.Command, args []string) {
	errapi := appManageAPI.LoginApp()
	if errapi != nil {
		fmt.Printf("%v", errapi)
	}
}
