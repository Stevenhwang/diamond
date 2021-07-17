package handlers

import (
	"log"

	"diamond/utils.go"

	"github.com/gofiber/fiber/v2"
)

// 用户登录
func Login(c *fiber.Ctx) error { return nil }

// 用户登出
func Logout(c *fiber.Ctx) error { return nil }

// 用户获取个人信息
func UserInfo(c *fiber.Ctx) error { return nil }

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
