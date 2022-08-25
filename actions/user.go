package actions

import (
	"diamond/config"
	"diamond/models"
	"diamond/policy"
	"diamond/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type auth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Otp      string `json:"otp"`
}

func login(c echo.Context) error {
	au := &auth{}
	if err := c.Bind(au); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(au); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	user := &models.User{}
	if result := models.DB.Where("username = ?", au.Username).First(user); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	// 验证密码
	if !utils.CheckPassword(user.Password, au.Password) {
		return c.JSON(http.StatusOK, H{"code": 4, "message": "password invalid"})
	}
	if !user.IsActive {
		return c.JSON(http.StatusOK, H{"code": 5, "message": "user is forbidden"})
	}
	if len(user.GoogleKey.String) > 0 {
		if len(au.Otp) == 0 {
			return c.JSON(http.StatusOK, H{"code": 6, "message": "need otp password"})
		}
		valid := totp.Validate(au.Otp, user.GoogleKey.String)
		if !valid {
			return c.JSON(http.StatusOK, H{"code": 7, "message": "invalid otp password"})
		}
	}
	// 生成token
	claims := jwt.MapClaims{
		"iat":          time.Now().Unix(),
		"iss":          "diamond",
		"uid":          user.ID,
		"username":     user.Username,
		"is_superuser": user.IsSuperuser,
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.GetString("jwt.secret")))
	if err != nil {
		return c.JSON(http.StatusOK, H{"code": 8, "message": err})
	}
	// 将token写入redis
	utils.SetToken(user.ID, t)
	// 更新用户登录IP和登录时间(不触发更新钩子)
	last_login_ip := nulls.NewString(c.RealIP())
	last_login_time := nulls.NewTime(time.Now())
	result := models.DB.Model(&user).UpdateColumns(models.User{LastLoginIP: last_login_ip, LastLoginTime: last_login_time})
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 9, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "token": t, "message": "login success"})
}

func logout(c echo.Context) error {
	uid := c.Get("uid").(uint)
	utils.DelToken(uid)
	return c.JSON(http.StatusOK, H{"code": 0, "message": "logout success"})
}

func user_info(c echo.Context) error {
	is_superuser := c.Get("is_superuser").(bool)
	name := c.Get("username").(string)
	uid := c.Get("uid").(int)
	var menus []string
	if is_superuser {
		menus = []string{}
	} else {
		sub := fmt.Sprintf("user::%d", uid)
		perms, _ := policy.Enforcer.GetNamedImplicitPermissionsForUser("p", sub)
		for _, perm := range perms {
			if utils.FindValInSlice(perm, "menu") {
				menus = append(menus, perm[1])
			}
		}
	}
	return c.JSON(http.StatusOK, H{"code": 0, "is_superuser": is_superuser, "name": name, "menus": menus})
}

type resetPw struct {
	OldPw  string `json:"old_pw" validate:"required"`
	NewPw1 string `json:"new_pw1" validate:"required"`
	NewPw2 string `json:"new_pw2" validate:"required"`
}

func resetPW(c echo.Context) error {
	rpw := &resetPw{}
	if err := c.Bind(rpw); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(rpw); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if rpw.NewPw1 != rpw.NewPw2 {
		return c.JSON(http.StatusOK, H{"code": 3, "message": "the two passwords not the same"})
	}
	uid := c.Get("uid").(uint)
	user := &models.User{}
	if result := models.DB.Find(user, uid); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	if !utils.CheckPassword(user.Password, rpw.OldPw) {
		return c.JSON(http.StatusOK, H{"code": 5, "message": "the old password is not right"})
	}
	models.DB.Model(&user).Update("password", rpw.NewPw1)
	return c.JSON(http.StatusOK, H{"code": 0, "message": "reset password success"})
}

func getUsers(c echo.Context) error {
	users := models.Users{}
	var total int64
	var result *gorm.DB
	// 使用role查找的时候不用分页，也不用filter
	if role := c.QueryParam("role"); len(role) > 0 {
		r := &models.Role{}
		models.DB.Where("name = ?", role).First(r)
		sub := fmt.Sprintf("role::%d", r.ID)
		usersliceraw, _ := policy.Enforcer.GetUsersForRole(sub)
		var userslice []string
		for _, u := range usersliceraw {
			userslice = append(userslice, strings.ReplaceAll(u, "user::", ""))
		}
		result = models.DB.Where("name IN ?", userslice).Omit("password").Find(&users)
		total = int64(len(users))
	} else {
		models.DB.Model(&models.User{}).Scopes(models.Filter(models.User{}, c)).Count(&total)
		result = models.DB.Scopes(models.Filter(models.User{}, c), models.Paginate(c)).Omit("password").Find(&users)
	}
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "data": users, "total": total})
}

func createUser(c echo.Context) error {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if result := models.DB.Create(user); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "create user success"})
}

func updateUser(c echo.Context) error {
	user := &models.User{}
	if result := models.DB.Find(user, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
	}
	// 处理password和otp_key更新
	excludeColumns := []string{}
	if len(user.Password) == 0 {
		excludeColumns = append(excludeColumns, "password")
	}
	otpKey, _ := strconv.Atoi(user.GoogleKey.String)
	if otpKey == 1 {
		excludeColumns = append(excludeColumns, "google_key")
	}
	if result := models.DB.Select("*").Omit(excludeColumns...).Updates(user); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "update user success"})
}

func deleteUser(c echo.Context) error {
	user := &models.User{}
	if result := models.DB.Find(user, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if user.IsSuperuser || user.ID == 1 {
		return c.JSON(http.StatusOK, H{"code": 2, "message": "can not delete a superuser"})
	}
	if result := models.DB.Delete(user); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete user success"})
}
