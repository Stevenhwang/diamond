package cmd

import (
	"diamond/models"
	"log"

	"github.com/gobuffalo/nulls"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Initialize the database[创建admin账户]",

	Run: func(cmd *cobra.Command, args []string) {
		user := models.User{Username: "admin", GoogleKey: nulls.NewString("seed"), IsSuperuser: true}
		result := models.DB.Create(&user)
		if result.Error != nil {
			log.Println(result.Error)
		} else {
			log.Println("Initialize database success")
		}
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
}
