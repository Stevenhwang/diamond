package actions

import (
	"diamond/cache"
	"diamond/misc"
	"diamond/models"
	"diamond/policy"
	"strings"
	"time"

	"github.com/Stevenhwang/gommon/nulls"
	"github.com/Stevenhwang/gommon/slice"
	"github.com/Stevenhwang/gommon/tools"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type auth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func login(c echo.Context) error {
	au := auth{}
	if err := c.Bind(&au); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&au); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	user := models.User{}
	if result := models.DB.Where("username = ?", au.Username).First(&user); result.Error != nil {
		cache.Ban(c.RealIP())
		return echo.NewHTTPError(400, result.Error.Error())
	}
	// 验证密码
	if !tools.CheckPassword(user.Password, au.Password) {
		cache.Ban(c.RealIP())
		return echo.NewHTTPError(400, "password invalid")
	}
	if !user.IsActive {
		return echo.NewHTTPError(400, "accound forbiden")
	}
	// 生成token
	claims := jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"iss":      "diamond",
		"uid":      user.ID,
		"username": user.Username,
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(misc.Config.GetString("jwt.secret")))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 更新用户登录IP和登录时间(不触发更新钩子)
	last_login_ip := nulls.NewString(c.RealIP())
	last_login_time := nulls.NewTime(time.Now())
	result := models.DB.Model(&user).UpdateColumns(models.User{LastLoginIP: last_login_ip, LastLoginTime: last_login_time})
	if result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "token": t})
}

func userInfo(c echo.Context) error {
	username := c.Get("username").(string)
	menus := c.Get("menus").(datatypes.JSON)
	return c.JSON(200, echo.Map{"success": true, "name": username, "menus": menus})
}

type resetPw struct {
	OldPw  string `json:"old_pw" validate:"required"`
	NewPw1 string `json:"new_pw1" validate:"required"`
	NewPw2 string `json:"new_pw2" validate:"required"`
}

// 重置密码
func resetPW(c echo.Context) error {
	rpw := resetPw{}
	if err := c.Bind(&rpw); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&rpw); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if rpw.NewPw1 != rpw.NewPw2 {
		return echo.NewHTTPError(400, "the two passwords not the same")
	}
	uid := c.Get("uid").(uint)
	user := models.User{}
	if result := models.DB.Find(&user, uid); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if !tools.CheckPassword(user.Password, rpw.OldPw) {
		return echo.NewHTTPError(400, "the old password is not right")
	}
	user.Password = rpw.NewPw1
	models.DB.Updates(&user)
	return c.JSON(200, echo.Map{"success": true})
}

func getUsers(c echo.Context) error {
	var total int64
	users := models.Users{}
	if res := models.DB.Model(&models.User{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Omit("password").Find(&users); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": users, "total": total})
}

func createUser(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&user); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Create(&user); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateUser(c echo.Context) error {
	user := models.User{}
	if result := models.DB.First(&user, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&user); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if user.Username == "admin" {
		return echo.NewHTTPError(400, "can not update admin")
	}
	// 处理password更新
	excludeColumns := []string{}
	if len(user.Password) == 0 {
		excludeColumns = append(excludeColumns, "password")
	}
	if result := models.DB.Select("*").Omit(excludeColumns...).Updates(&user); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteUser(c echo.Context) error {
	user := models.User{}
	if result := models.DB.Find(&user, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if user.Username == "admin" {
		return echo.NewHTTPError(400, "can not delete admin")
	}
	if result := models.DB.Delete(&user); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

type userPerm struct {
	Username string   `json:"username" validate:"required"`
	Perms    []string `json:"perms" validate:"required"` // 权限名称列表，不用id因为可能变化，名称固定
}

// 给用户分配权限
func assignUserPerm(c echo.Context) error {
	up := userPerm{}
	if err := c.Bind(&up); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&up); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 先清除所有权限
	policies := policy.Enforcer.GetFilteredPolicy(0, up.Username)
	if _, err := policy.Enforcer.RemovePolicies(policies); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 再更新用户权限
	var rules [][]string
	permissions := models.Permissions{}
	if res := models.DB.Where("name IN ?", up.Perms).Find(&permissions); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	for _, p := range permissions {
		rules = append(rules, []string{up.Username, p.URL, p.Method})
	}
	if _, err := policy.Enforcer.AddPolicies(rules); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

type checkedPerm struct {
	models.Permission
	Check bool `json:"check"`
}

// 根据用户id查找授予的权限
func getUserPerms(c echo.Context) error {
	id := c.QueryParam("id")
	if len(id) == 0 {
		return echo.NewHTTPError(400, "need id")
	}
	user := models.User{}
	if res := models.DB.First(&user, id); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	// 查找系统权限
	permissions := models.Permissions{}
	if res := models.DB.Order("created_at desc").Find(&permissions); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	var requests [][]interface{}
	for _, p := range permissions {
		requests = append(requests, []interface{}{user.Username, p.URL, p.Method})
	}
	// 批量检查用户权限
	checks, err := policy.Enforcer.BatchEnforce(requests)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	checkPerms := []checkedPerm{}
	for i, p := range permissions {
		checkPerms = append(checkPerms, checkedPerm{p, checks[i]})
	}
	return c.JSON(200, echo.Map{"success": true, "data": checkPerms})
}

type checkedServer struct {
	models.Server
	Check bool `json:"check"`
}

// 根据用户id查找授予的服务器
func getUserServers(c echo.Context) error {
	id := c.QueryParam("id")
	if len(id) == 0 {
		return echo.NewHTTPError(400, "need id")
	}
	user := models.User{}
	if res := models.DB.Preload("Servers", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&user, id); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	var sids []int
	for _, s := range user.Servers {
		sids = append(sids, int(s.ID))
	}
	// 查找所有服务器
	servers := models.Servers{}
	if res := models.DB.Order("created_at desc").Find(&servers); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	checkServers := []checkedServer{}
	for _, s := range servers {
		if slice.FindValInIntSlice(sids, int(s.ID)) {
			checkServers = append(checkServers, checkedServer{s, true})
		} else {
			checkServers = append(checkServers, checkedServer{s, false})
		}
	}
	return c.JSON(200, echo.Map{"success": true, "data": checkServers})
}

type userServer struct {
	Username string `json:"username" validate:"required"`
	Servers  []int  `json:"servers" validate:"required"` // 服务器ID列表
}

// 给用户分配服务器
func assignUserServer(c echo.Context) error {
	us := userServer{}
	if err := c.Bind(&us); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&us); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 先清除所有服务器
	user := models.User{}
	if res := models.DB.Where("username = ?", us.Username).First(&user); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if err := models.DB.Model(&user).Association("Servers").Clear(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 再更新用户服务器，列表为空就没必要继续了
	if len(us.Servers) > 0 {
		servers := models.Servers{}
		if res := models.DB.Where("id IN ?", us.Servers).Find(&servers); res.Error != nil {
			return echo.NewHTTPError(400, res.Error.Error())
		}
		if err := models.DB.Model(&user).Association("Servers").Append(servers); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
	}
	return c.JSON(200, echo.Map{"success": true})
}

// 同步系统权限
func syncPerms(c echo.Context) error {
	// sync permissions，先清空表，再更新
	models.DB.Exec("TRUNCATE TABLE permissions")
	routes := App.Routes()
	perms := models.Permissions{}
	prefix := "diamond/actions."
	for _, r := range routes {
		if strings.HasPrefix(r.Name, prefix) {
			n := strings.ReplaceAll(r.Name, prefix, "")
			if !slice.FindValInStringSlice([]string{"login", "terminal", "userInfo", "resetPW"}, n) { // 白名单
				perms = append(perms, models.Permission{Name: n, Method: r.Method, URL: r.Path})
			}
		}
	}
	models.DB.Create(&perms)
	return c.JSON(200, echo.Map{"success": true})
}

func getBanIPs(c echo.Context) error {
	keyword := c.QueryParam("ip")
	ips, err := cache.FilterBanIPs(keyword)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": ips})
}

type delBan struct {
	IP string `json:"ip" validate:"required"`
}

func delBanIP(c echo.Context) error {
	db := delBan{}
	if err := c.Bind(&db); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&db); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := cache.DelBanIP(db.IP); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}
