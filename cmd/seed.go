package cmd

import (
	"database/sql"
	"diamond/models"
	"log"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Initialize the database[初始化数据库]",

	Run: func(cmd *cobra.Command, args []string) {
		user := models.User{Username: "admin", GoogleKey: sql.NullString{String: "seed", Valid: true}, IsSuperuser: true}
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
