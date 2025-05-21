/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 02:08
 * @Description:
 */

package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewMiddlewares() []*fiber.App {
	return []*fiber.App{
		//Middleware1,
		//Middleware2,
		//Middleware3,
		//Cors(),
	}
}

func Cors() fiber.Handler {
	// cors
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	})
}
