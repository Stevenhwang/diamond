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

	app.Use(authTokenMW)
	app.Get("/users", UserListPerm)

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
