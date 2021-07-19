package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func App() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// routers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from diamond!",
		})
	})

	app.Use(realIPMW)
	app.Use(authTokenMW)
	app.Post("/login", Login)
	app.Post("/logout", Logout)
	app.Get("/user_info", UserInfo)
	app.Post("/reset_pw", ResetPasswd)

	app.Get("/users", UserListPerm)
	app.Post("/users", CreateUserPerm)
	app.Put("/users/:id", UpdateUserPerm)
	app.Delete("/users/:id", DeleteUserPerm)

	app.Get("/logs", LogListPerm)

	// rl := app.Stack()
	// hn := []string{}
	// for _, r := range rl {
	// 	for _, k := range r {
	// 		for _, val := range k.Handlers {
	// 			hn = append(hn, utils.NameOfFunction(val))
	// 		}
	// 	}
	// }
	// hn = utils.RemoveDupInSlice(hn)
	// log.Println(hn)

	return app
}
