package cmd

import (
	"diamond/actions"

	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "start app server[开启app服务]",

	Run: func(cmd *cobra.Command, args []string) {
		actions.CronStart()
		addr, _ := cmd.Flags().GetString("addr")
		app := actions.App
		app.Logger.Fatal(app.Start(addr))
	},
}

func init() {
	appCmd.Flags().StringP("addr", "a", ":8000", "app listen address")
	RootCmd.AddCommand(appCmd)
}
