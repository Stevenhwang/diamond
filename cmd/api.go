package cmd

import (
	"diamond/config"
	"diamond/crons"
	"diamond/handlers"
	"log"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server[开启 api 服务器]",

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("开启定时任务")
		crons.C.Start()
		app := handlers.App()
		addr := config.Config.Get("api.addr").(string)
		app.Run(addr)
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)
}
