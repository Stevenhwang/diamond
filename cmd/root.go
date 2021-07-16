package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	// Use: "help",
	// Short: "diamond devops",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Welcome to diamond!")
	// },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
