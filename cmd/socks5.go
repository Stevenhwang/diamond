package cmd

import (
	"diamond/socks5"

	"github.com/spf13/cobra"
)

var socks5Cmd = &cobra.Command{
	Use:   "socks5",
	Short: "start socks5 server[开启socks5服务]",

	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		socks5.Start(addr)
	},
}

func init() {
	socks5Cmd.Flags().StringP("addr", "a", ":8538", "socks5 listen address")
	RootCmd.AddCommand(socks5Cmd)
}
