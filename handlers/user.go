package handlers

import (
	"database/sql"
	"log"
	"time"

	"diamond/models"
	"diamond/utils.go"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
)

// 用户登录
func Login(c *fiber.Ctx) error {
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Otp      string `json:"otp"`
	}
	lg := &login{}
	if err := c.BodyParser(lg); err != nil {
		return RespMsgSuccess(c, 1, err.Error())
	}
	if len(lg.Username) == 0 || len(lg.Password) == 0 {
		return RespMsgSuccess(c, 2, "用户名或密码不能为空！")
	}
	user := &models.User{}
	if err := models.DB.Where("username = ?", lg.Username).First(user); err != nil {
		return RespMsgSuccess(c, 3, "用户不存在！")
	}
	// 验证密码
	if !utils.CheckPassword(user.Password, lg.Password) {
		return RespMsgSuccess(c, 4, "密码错误！")
	}
	if !user.IsActive {
		return RespMsgSuccess(c, 5, "账号被禁用！")
	}
	if len(user.GoogleKey.String) > 0 {
		if len(lg.Otp) == 0 {
			return RespMsgSuccess(c, 6, "需要二次认证验证码！")
		}
		valid := totp.Validate(lg.Otp, user.GoogleKey.String)
		if !valid {
			return RespMsgSuccess(c, 7, "验证码错误！")
		}
	}
	// 生成token
	token := utils.J.EncodeToken(user.ID, user.Username, user.IsSuperuser)
	// 将token写入redis
	utils.SetToken(user.ID, token)
	// 更新用户登录IP和登录时间(不触发更新钩子)
	last_login_ip := sql.NullString{String: c.IPs()[0], Valid: true}
	last_login_time := sql.NullTime{Time: time.Now(), Valid: true}
	result := models.DB.Model(&user).UpdateColumns(models.User{LastLoginIP: last_login_ip, LastLoginTime: last_login_time})
	if result.Error != nil {
		return RespMsgSuccess(c, 8, result.Error.Error())
	}
	return c.JSON(fiber.Map{
		"code":    0,
		"message": "登录成功！",
		"token":   token,
	})
}

// 用户登出
func Logout(c *fiber.Ctx) error {
	uid := c.Locals("uid").(uint)
	utils.DelToken(uid)
	return RespMsgSuccess(c, 0, "注销成功！")
}

// 用户获取个人信息
func UserInfo(c *fiber.Ctx) error {
	uid := c.Locals("uid").(uint)
	username := c.Locals("username").(string)
	isSuperuser := c.Locals("is_superuser").(bool)
	// 管理员
	if isSuperuser {
		menus := &models.Menus{}
		if result := models.DB.Select("name").Find(menus); result.Error != nil {
			return RespMsgSuccess(c, 1, result.Error.Error())
		}
		menuNames := make([]string, 0, len(*menus))
		for _, v := range *menus {
			menuNames = append(menuNames, v.Name)
		}
		return c.JSON(fiber.Map{
			"code":         0,
			"is_superuser": true,
			"menus":        menuNames,
			"username":     username,
		})
	}
	// 普通用户
	user := &models.User{}
	menusMap := map[string]byte{}
	if result := models.DB.Preload("Roles.Menus").First(user, uid); result.Error != nil {
		return RespMsgSuccess(c, 1, result.Error.Error())
	}
	for _, role := range user.Roles {
		if role.IsActive {
			for _, menu := range role.Menus {
				if menu.IsActive {
					menusMap[menu.Name] = 0
				}
			}
		}
	}
	menus := make([]string, 0, len(menusMap))
	for k := range menusMap {
		menus = append(menus, k)
	}
	return c.JSON(fiber.Map{
		"code":         0,
		"is_superuser": true,
		"menus":        menus,
		"username":     username,
	})
}

// 用户重置密码
func ResetPasswd(c *fiber.Ctx) error { return nil }

// 获取用户列表
func UserListPerm(c *fiber.Ctx) error {
	hn := make([]string, 0, len(c.Route().Handlers))
	for _, val := range c.Route().Handlers {
		hn = append(hn, utils.NameOfFunction(val))
	}
	log.Println(hn)
	return nil
}

// 更新用户信息
func UpdateUserPerm(c *fiber.Ctx) error { return nil }

// 新建用户
func CreateUserPerm(c *fiber.Ctx) error { return nil }

// 删除用户
func DeleteUserPerm(c *fiber.Ctx) error { return nil }
