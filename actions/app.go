package actions

import (
	"database/sql/driver"
	"diamond/middlewares"
	"diamond/misc"
	"diamond/models"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"reflect"

	"diamond/frontend"

	"github.com/Stevenhwang/gommon/nulls"
	"github.com/Stevenhwang/gommon/tools"

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
			misc.Logger.Info().
				Str("from", "app").
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
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// home route
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(200, "Hello from diamond!")
	// })

	assetHandler := http.FileServer(getFileSystem())
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(assetHandler))
	e.GET("*.png", echo.WrapHandler(assetHandler))
	// e.GET("/favicon.ico", echo.WrapHandler(assetHandler))

	// ssh records 目录
	e.Static("/records", "./records")

	// 反向代理 navicat http 隧道，密码保护
	naviURL, _ := url.Parse(misc.Config.GetString("navicate.url"))
	navi := e.Group("/ntunnel_mysql.php", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		user := models.User{}
		if result := models.DB.Where("username = ?", username).First(&user); result.Error != nil {
			return false, nil
		}
		// 验证密码
		if !tools.CheckPassword(user.Password, password) {
			return false, nil
		}
		if !user.IsActive {
			return false, nil
		}
		return true, nil
	}))
	navi.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: naviURL,
		},
	})))

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

	// // sync permissions，先清空表，再更新
	// models.DB.Exec("TRUNCATE TABLE permissions")
	// routes := e.Routes()
	// perms := models.Permissions{}
	// prefix := "diamond/actions."
	// for _, r := range routes {
	// 	if strings.HasPrefix(r.Name, prefix) {
	// 		n := strings.ReplaceAll(r.Name, prefix, "")
	// 		if n != "login" && n != "terminal" { // 白名单
	// 			perms = append(perms, models.Permission{Name: n, Method: r.Method, URL: r.Path})
	// 		}
	// 	}
	// }
	// models.DB.Create(&perms)

	App = e
}
