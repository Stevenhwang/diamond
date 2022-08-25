package actions

import (
	"database/sql/driver"
	"diamond/config"
	"diamond/middlewares"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/nulls"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var App *echo.Echo

type H map[string]interface{}

// 自定义错误处理函数
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}
	// c.Logger().Error(err)
	c.JSON(code, H{"code": code, "message": message})
}

// 自定义validator
type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func init() {
	// Echo instance
	e := echo.New()

	// custom error handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// custom validator
	validate := validator.New()
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(driver.Valuer); ok {
			val, err := valuer.Value()
			if err == nil { // nulls 包里这个err一直是nil, 所以不用继续处理
				return val
			}
			// handle the error how you want
		}
		return nil
	}, nulls.String{}, nulls.Time{}, nulls.UInt32{})
	e.Validator = &customValidator{validator: validate}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// With proxies using X-Forwarded-For header to get real client ip
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// home route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// api group route
	api := e.Group("/api")
	api.POST("/login", login)

	// jwt middleware
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.Config.GetString("jwt.secret")),
	}))
	// token middleware
	api.Use(middlewares.Token)
	// enforce middleware
	api.Use(middlewares.Enforce)

	api.POST("/logout", logout)
	api.GET("/user_info", user_info)
	api.POST("/reset_pw", resetPW)

	api.GET("/users", getUsers)
	api.POST("/users", createUser)
	api.PUT("/users/:id", updateUser)
	api.DELETE("/users/:id", deleteUser)

	api.GET("/roles", getRoles)
	api.POST("/roles", createRole)
	api.PUT("/roles/:id", updateRole)
	api.DELETE("/roles/:id", deleteRole)

	// api.GET("/permissions", getPermissions)
	// api.POST("/permissions", createPermission)
	// api.PUT("/permissions/:id", updatePermission)
	// api.DELETE("/permissions/:id", deletePermission)

	api.GET("/keys", getKeys)
	api.POST("/keys", createKey)
	api.PUT("/keys/:id", updateKey)
	api.DELETE("/keys/:id", deleteKey)

	api.GET("/servers", getServers)
	api.POST("/servers", createServer)
	api.PUT("/servers/:id", updateServer)
	api.DELETE("/servers/:id", deleteServer)

	api.GET("/groups", getGroups)
	api.POST("/groups", createGroup)
	api.PUT("/groups/:id", updateGroup)
	api.DELETE("/groups/:id", deleteGroup)

	App = e
}
