package actions

import (
	"database/sql/driver"
	"diamond/middlewares"
	"diamond/misc"
	"fmt"
	"io/fs"
	"net/http"
	"reflect"

	"diamond/frontend"

	"github.com/Stevenhwang/gommon/nulls"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(frontend.Index, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

var App *echo.Echo

// 自定义错误处理函数
func customHTTPErrorHandler(err error, c echo.Context) {
	code := 500
	message := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}
	// c.Logger().Error(err)
	c.JSON(200, echo.Map{"success": false, "code": code, "reason": message})
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
	}, nulls.String{}, nulls.Time{}, nulls.UInt{})
	e.Validator = &customValidator{validator: validate}

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogStatus:    true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			username := c.Get("username")
			var un string
			if username == nil {
				un = ""
			} else {
				un = username.(string)
			}
			misc.Logger.Info().
				Str("from", "app").
				Str("user", un).
				Str("remote_ip", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Int("status", v.Status).
				Str("uri", v.URI).
				Str("latency", v.Latency.String()).
				Str("user_agent", v.UserAgent).
				Msg("")
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// With proxies using X-Forwarded-For header to get real client ip
	// e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.IPExtractor = echo.ExtractIPDirect()

	// home route
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(200, "Hello from diamond!")
	// })

	// banip middleware
	e.Use(middlewares.BanIP)

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(assetHandler))
	e.GET("*.png", echo.WrapHandler(assetHandler))
	// e.GET("/favicon.ico", echo.WrapHandler(assetHandler))

	// ssh records 目录
	e.Static("/records", "./records")

	// api group route
	api := e.Group("/api")
	api.POST("/login", login)
	api.GET("/terminal", terminal)
	// jwt middleware
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(misc.Config.GetString("jwt.secret")),
	}))
	// token middleware
	api.Use(middlewares.Token)
	api.GET("/user_info", userInfo)
	api.POST("/reset_pw", resetPW)
	// enforce middleware
	api.Use(middlewares.Enforce)

	api.GET("/users", getUsers)
	api.POST("/users", createUser)
	api.PUT("/users/:id", updateUser)
	api.DELETE("/users/:id", deleteUser)

	api.GET("/banips", getBanIPs)
	api.POST("/delbanip", delBanIP)

	api.POST("/syncPerms", syncPerms)

	api.GET("/userPerms", getUserPerms)
	api.POST("/userPerms", assignUserPerm)

	api.GET("/userServers", getUserServers)
	api.POST("/userServers", assignUserServer)

	api.GET("/servers", getServers)
	api.POST("/servers", createServer)
	api.PUT("/servers/:id", updateServer)
	api.DELETE("/servers/:id", deleteServer)

	api.GET("/records", getRecords)

	api.GET("/credentials", getCredentials)
	api.POST("/credentials", createCredential)
	api.PUT("/credentials/:id", updateCredential)
	api.DELETE("/credentials/:id", deleteCredential)

	api.GET("/scripts", getScripts)
	api.POST("/scripts", createScript)
	api.PUT("/scripts/:id", updateScript)
	api.DELETE("/scripts/:id", deleteScript)

	api.GET("/tasks", getTasks)
	api.POST("/tasks", createTask)
	api.PUT("/tasks/:id", updateTask)
	api.DELETE("/tasks/:id", deleteTask)
	api.POST("/tasks/:id", invokeTask)
	api.GET("/taskhist", getTasksHist)
	api.GET("/taskhist/:id", getTasksHistDetail)

	api.GET("/crons", getCrons)
	api.POST("/crons", createCron)
	api.PUT("/crons/:id", updateCron)
	api.DELETE("/crons/:id", deleteCron)

	App = e
}
