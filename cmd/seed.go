package cmd

import (
	"diamond/models"
	"log"

	"github.com/gobuffalo/nulls"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed the database[创建admin账户]",

	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		user := models.User{Username: username, Password: password, GoogleKey: nulls.NewString("seed"), IsSuperuser: true}
		result := models.DB.Create(&user)
		if result.Error != nil {
			log.Println(result.Error)
		} else {
			log.Printf("seed database success[%s:%s]", username, password)
		}
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
	seedCmd.Flags().StringP("username", "u", "admin", "seed username")
	seedCmd.Flags().StringP("password", "p", "12345678", "seed password")
}
