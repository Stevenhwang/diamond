package cmd

import (
	"diamond/models"
	"fmt"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed user account[创建用户账户]",

	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		user := models.User{Username: username, Password: password, IsActive: true}
		res := models.DB.Create(&user)
		if res.Error != nil {
			fmt.Printf("seed user error: %s\n", res.Error.Error())
		} else {
			fmt.Printf("seed user success: user=> %s password=> %s\n", username, password)
		}
	},
}

func init() {
	seedCmd.Flags().StringP("username", "u", "admin", "seed username")
	seedCmd.Flags().StringP("password", "p", "12345678", "seed password")
	RootCmd.AddCommand(seedCmd)
}
