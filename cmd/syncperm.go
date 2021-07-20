package cmd

import (
	"diamond/handlers"
	"diamond/models"
	"diamond/utils.go"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var syncPermCmd = &cobra.Command{
	Use:   "syncperm",
	Short: "sync permissions info[同步权限信息]",

	Run: func(cmd *cobra.Command, args []string) {
		// 获取所有权限
		routes := handlers.App().Routes()
		allPerms := []string{}
		for _, v := range routes {
			if strings.HasSuffix(v.Handler, "Perm") {
				allPerms = append(allPerms, v.Handler)
			}
		}
		// doc文档
		docs := make(map[string]string)
		fset := token.NewFileSet() // positions are relative to fset
		d, err := parser.ParseDir(fset, "./handlers", nil, parser.ParseComments)
		if err != nil {
			log.Fatalln(err)
		}
		for _, f := range d {
			p := doc.New(f, "./", 2)
			// 获取所有func doc
			for _, f := range p.Funcs {
				if strings.HasSuffix(f.Name, "Perm") {
					funcName := fmt.Sprintf("diamond/handlers.%v", f.Name)
					docs[funcName] = f.Doc
				}
			}
		}
		// 已有权限
		existPerms := &models.Permissions{}
		if result := models.DB.Find(existPerms); result.Error != nil {
			log.Fatalln(result.Error)
		}
		existPermsList := []string{}
		// 删除不需要的权限
		for _, v := range *existPerms {
			existPermsList = append(existPermsList, v.Name)
			if !utils.FindValInSlice(allPerms, v.Name) {
				if result := models.DB.Delete(&v); result.Error != nil {
					log.Fatalln(result.Error)
				}
				log.Println("Delete permission " + v.Name)
			}
		}
		// 新增需要的权限
		for _, v := range allPerms {
			if !utils.FindValInSlice(existPermsList, v) {
				perm := &models.Permission{}
				perm.Name = v
				perm.Remark = docs[v]
				perm.IsActive = true
				if result := models.DB.Create(perm); result.Error != nil {
					log.Fatalln(result.Error)
				}
				log.Println("Add permission " + v)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(syncPermCmd)
}
