package handlers

import (
	"diamond/models"
	"diamond/utils.go"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

func authTokenMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 登录接口跳过
		if c.Request.URL.Path == "/login" {
			c.Next()
			return
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
				log.Data = nulls.NewString(string(b))
			}
			go func(log *models.Log) {
				models.DB.Create(log)
			}(log)
		}
		c.Next()
	}
}

func checkPermMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 管理员跳过
		if isSuperuser := c.GetBool("is_superuser"); isSuperuser {
			c.Next()
			return
		}
		// 不需要权限的handler跳过
		if !strings.HasSuffix(c.HandlerName(), "Perm") {
			c.Next()
			return
		}
		// 检查权限
		user := &models.User{}
		if result := models.DB.Preload("Roles.Permissions").First(user, c.GetUint("user_id")); result.Error != nil {
			c.JSON(401, gin.H{
				"code":    1,
				"message": result.Error.Error(),
			})
			c.Abort()
			return
		}
		for _, role := range user.Roles {
			if role.IsActive {
				for _, permission := range role.Permissions {
					if permission.IsActive && permission.Name == c.HandlerName() {
						c.Next()
						return
					}
				}
			}
		}
		c.JSON(401, gin.H{
			"code":    2,
			"message": "无权限！",
		})
		c.Abort()
	}
}
