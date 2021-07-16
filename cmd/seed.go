package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Initialize the database[初始化数据库]",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("initialize the database")
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
}
