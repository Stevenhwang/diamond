package cmd

import (
	"diamond/config"
	"diamond/handlers"
	"log"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "start web api server[开启 web api 服务器]",

	Run: func(cmd *cobra.Command, args []string) {
		app := handlers.App()
		addr := config.Config.Get("web.addr").(string)
		log.Fatal(app.Listen(addr))
	},
}

func init() {
	RootCmd.AddCommand(webCmd)
}
