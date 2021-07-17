package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncPermCmd = &cobra.Command{
	Use:   "syncPerm",
	Short: "sync permissions info[同步权限信息]",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync permissions info")
	},
}

func init() {
	RootCmd.AddCommand(syncPermCmd)
}
