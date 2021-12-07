package cmd

import (
	"fmt"

	"github.com/platform9/appctl/pkg/appManageAPI"
	"github.com/spf13/cobra"
)

// loginCmd represents "Login and use appctl".
var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to use appctl",
		Long:  `Login to use appctl`,
		Run:   loginCmdRun,
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

// To login.
func loginCmdRun(cmd *cobra.Command, args []string) {
	errapi := appManageAPI.LoginApp()
	if errapi != nil {
		fmt.Printf("Not able to login: Error %v", errapi)
	}
}
