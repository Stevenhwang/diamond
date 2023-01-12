package cmd

import (
	"diamond/sshd"

	"github.com/spf13/cobra"
)

var sshdCmd = &cobra.Command{
	Use:   "sshd",
	Short: "start sshd server[开启sshd服务]",

	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		keyPath, _ := cmd.Flags().GetString("keyPath")
		sshd.Start(addr, keyPath)
	},
}

func init() {
	sshdCmd.Flags().StringP("addr", "a", ":2222", "sshd listen address")
	sshdCmd.Flags().StringP("keyPath", "k", "C:/Users/90hua/.ssh/id_rsa", "sshd private key path")
	RootCmd.AddCommand(sshdCmd)
}
