package cmd

import (
	"diamond/sshd"

	"github.com/spf13/cobra"
)

var sshdCmd = &cobra.Command{
	Use:   "sshd",
	Short: "start sshd server",

	Run: func(cmd *cobra.Command, args []string) {
		sshd.Start()
	},
}

func init() {
	RootCmd.AddCommand(sshdCmd)
}
