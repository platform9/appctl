package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// appCmd represents the app commands can be run
var (
	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Create or get or revise app",
		Long:  `Create or get or revise app, app.run controller`,
	}

	appCmdList = &cobra.Command{
		Use:   "list",
		Short: "List apps running",
		Long:  `List apps running`,
		Run:   appCmdListRun,
	}
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appCmdList)
	//appCmdList.Flags().StringVarP(&nameSpace, "namespace", "n", "", "set namespace")
}

func appCmdListRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Into the applist command, logic to be implemeted")
}
