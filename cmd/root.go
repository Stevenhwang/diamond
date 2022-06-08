package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "help",
	Short: "diamond devops",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to diamond!")
		// policy.Enforcer.AddNamedPolicy("p", "data_admin", "route_group", "POST")
		// policy.Enforcer.AddNamedGroupingPolicy("g", "alice", "data_admin")
		// policy.Enforcer.AddNamedGroupingPolicy("g2", "/api/users/:id", "route_group")
		// pass, err := policy.Enforcer.Enforce("alice", "/api/users/1", "POST")
		// log.Println("pass: ", pass)
		// log.Println("error: ", err)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
