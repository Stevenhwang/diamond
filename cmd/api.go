package cmd

import (
	"diamond/actions"
	"diamond/config"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server[开启 api 服务器]",

	Run: func(cmd *cobra.Command, args []string) {
		app := actions.App
		addr := config.Config.Get("api.addr").(string)
		app.Logger.Fatal(app.Start(addr))
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)
}
