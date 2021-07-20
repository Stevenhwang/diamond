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

	app.GET("/roles", RoleListPerm)
	app.POST("/roles", CreateRolePerm)
	app.PUT("/roles/:id", UpdateRolePerm)
	app.DELETE("/roles/:id", DeleteRolePerm)
	app.POST("/roles/:id", ResAssignPerm)

	app.GET("/menus", MenuListPerm)
	app.POST("/menus", CreateMenuPerm)
	app.PUT("/menus/:id", UpdateMenuPerm)
	app.DELETE("/menus/:id", DeleteMenuPerm)

	app.GET("/perms", PermissionListPerm)
	app.PUT("/perms/:id", UpdatePermissionPerm)

	app.GET("/groups", GroupListPerm)
	app.POST("/groups", CreateGroupPerm)
	app.PUT("/groups/:id", UpdateGroupPerm)
	app.DELETE("/groups/:id", DeleteGroupPerm)

	app.GET("/servers", ServerListPerm)
	app.POST("/servers", CreateServerPerm)
	app.PUT("/servers/:id", UpdateServerPerm)
	app.DELETE("/servers/:id", DeleteServerPerm)

	app.GET("/logs", LogListPerm)

	return app
}
