package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error: %s", err),
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// routers
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from diamond!",
		})
	})

	app.Use(authTokenMW())
	app.Use(checkPermMW())
	app.POST("/login", Login)
	app.POST("/logout", Logout)
	app.GET("/user_info", UserInfo)
	app.POST("/reset_pw", ResetPasswd)

	app.GET("/users", UserListPerm)
	app.POST("/users", CreateUserPerm)
	app.PUT("/users/:id", UpdateUserPerm)
	app.DELETE("/users/:id", DeleteUserPerm)

	app.GET("/menus", MenuListPerm)
	app.POST("/menus", CreateMenuPerm)
	app.PUT("/menus/:id", UpdateMenuPerm)
	app.DELETE("/menus/:id", DeleteMenuPerm)

	app.GET("/group", GroupListPerm)
	app.POST("/group", CreateGroupPerm)
	app.PUT("/group/:id", UpdateGroupPerm)
	app.DELETE("/group/:id", DeleteGroupPerm)

	app.GET("/logs", LogListPerm)

	return app
}
