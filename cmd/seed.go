package cmd

import (
	"diamond/misc"
	"diamond/models"
	"fmt"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed admin account[创建admin账户]",

	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		admin := models.User{Username: "admin", Password: password, IsActive: true}
		res := models.DB.Create(&admin)
		if res.Error != nil {
			misc.Logger.Error().Err(res.Error).Msg("")
		} else {
			fmt.Printf("seed admin success with password: %s", password)
		}
	},
}

func init() {
	seedCmd.Flags().StringP("password", "p", "12345678", "seed admin password")
	RootCmd.AddCommand(seedCmd)
}
