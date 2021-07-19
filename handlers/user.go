package handlers

import (
	"database/sql"
	"strconv"
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
	if result := models.DB.Where("username = ?", lg.Username).First(user); result.Error != nil {
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
	last_login_ip := sql.NullString{String: c.Locals("user_ip").(string), Valid: true}
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
func ResetPasswd(c *fiber.Ctx) error {
	type resetPw struct {
		OldPw  string `json:"old_pw"`
		NewPw1 string `json:"new_pw1"`
		NewPw2 string `json:"new_pw2"`
	}
	rpw := &resetPw{}
	if err := c.BodyParser(rpw); err != nil {
		return RespMsgSuccess(c, 1, err.Error())
	}
	if len(rpw.OldPw) == 0 || len(rpw.NewPw1) == 0 || len(rpw.NewPw2) == 0 {
		return RespMsgSuccess(c, 2, "关键参数不能为空！")
	}
	if rpw.NewPw1 != rpw.NewPw2 {
		return RespMsgSuccess(c, 3, "两个新密码输入不一致！")
	}
	uid := c.Locals("uid").(uint)
	user := &models.User{}
	if result := models.DB.Find(user, uid); result.Error != nil {
		return RespMsgSuccess(c, 4, result.Error.Error())
	}
	if !utils.CheckPassword(user.Password, rpw.OldPw) {
		return RespMsgSuccess(c, 5, "原始密码错误！")
	}
	models.DB.Model(&user).Update("password", rpw.NewPw1)
	return RespMsgSuccess(c, 0, "修改密码成功！")
}

// 获取用户列表
func UserListPerm(c *fiber.Ctx) error {
	users, total, err := models.GetUserList(c)
	if err != nil {
		return RespMsgSuccess(c, 1, err.Error())
	}
	return RespDataSuccess(c, 0, users, total)
}

// 更新用户信息
func UpdateUserPerm(c *fiber.Ctx) error {
	user := &models.User{}
	if result := models.DB.Find(user, c.Params("id")); result.Error != nil {
		return RespMsgSuccess(c, 1, result.Error.Error())
	}
	if err := c.BodyParser(user); err != nil {
		return RespMsgSuccess(c, 2, err.Error())
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
	if result := models.DB.Omit(excludeColumns...).Updates(user); result.Error != nil {
		return RespMsgSuccess(c, 3, result.Error.Error())
	}
	return RespMsgSuccess(c, 0, "更新成功！")
}

// 新建用户
func CreateUserPerm(c *fiber.Ctx) error {
	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return RespMsgSuccess(c, 1, err.Error())
	}
	if result := models.DB.Create(user); result.Error != nil {
		return RespMsgSuccess(c, 2, result.Error.Error())
	}
	return RespMsgSuccess(c, 0, "创建成功！")
}

// 删除用户
func DeleteUserPerm(c *fiber.Ctx) error {
	user := &models.User{}
	if result := models.DB.Find(user, c.Params("id")); result.Error != nil {
		return RespMsgSuccess(c, 1, result.Error.Error())
	}
	if user.IsSuperuser {
		return RespMsgSuccess(c, 2, "超级管理员不可被删除！")
	}
	if result := models.DB.Delete(user); result.Error != nil {
		return RespMsgSuccess(c, 3, result.Error.Error())
	}
	return RespMsgSuccess(c, 0, "删除成功！")
}
