package cmd

import (
	"diamond/models"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "auto migrate[运行自动迁移]",
	Long: `AutoMigrate 会创建表、缺失的外键、约束、列和索引。 
	    如果大小、精度、是否为空可以更改，则 AutoMigrate 会改变列的类型。 
	    出于保护您数据的目的，它不会删除未使用的列`,

	Run: func(cmd *cobra.Command, args []string) {
		models.DB.AutoMigrate(&models.User{}, &models.Server{}, &models.Role{},
			&models.Permission{}, &models.Menu{}, &models.Log{},
			&models.Group{}, &models.Record{})
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
