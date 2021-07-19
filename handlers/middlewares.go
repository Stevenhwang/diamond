package handlers

import (
	"database/sql"
	"diamond/models"
	"diamond/utils.go"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func authTokenMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 登录接口跳过
		if c.Request.URL.Path == "/login" {
			c.Next()
		}
		// 检查token
		xToken := c.Request.Header.Get("X-Token")
		if len(xToken) == 0 {
			c.JSON(401, gin.H{
				"code":    1,
				"message": "请登录!",
			})
			c.Abort()
			return
		}
		uid, username, isSuperuser := utils.J.DecodeToken(xToken)
		if uid == 0 {
			c.JSON(401, gin.H{
				"code":    2,
				"message": "Token非法!",
			})
			c.Abort()
			return
		}
		rdToken := utils.GetToken(uid)
		if len(rdToken) == 0 || rdToken != xToken {
			c.JSON(401, gin.H{
				"code":    3,
				"message": "Token失效!",
			})
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Set("username", username)
		c.Set("is_superuser", isSuperuser)
		// 记录访问日志(除get请求)
		cCp := c.Copy()
		if c.Request.Method != "GET" {
			log := &models.Log{}
			log.Username = username
			log.IP = cCp.ClientIP()
			log.Method = cCp.Request.Method
			log.URL = cCp.Request.URL.Path
			b, _ := ioutil.ReadAll(cCp.Request.Body)
			if len(b) > 0 {
				log.Data = sql.NullString{String: string(b), Valid: true}
			}
			go func(log *models.Log) {
				models.DB.Create(log)
			}(log)
		}
		c.Next()
	}
}
