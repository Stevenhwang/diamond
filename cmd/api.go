package cmd

import (
	"diamond/config"
	"diamond/handlers"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server[开启 api 服务器]",

	Run: func(cmd *cobra.Command, args []string) {
		app := handlers.App()
		addr := config.Config.Get("web.addr").(string)
		app.Run(addr)
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)
}