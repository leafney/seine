package core

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(app *fiber.App, deps *App) {
	// 创建 /api/v1 路由组
	v1 := app.Group("/api/v1")

	// 示例接口（无需认证）
	v1.Get("/hello", deps.ExampleHandler.Hello)
	v1.Get("/ping", deps.ExampleHandler.Ping)

	// 用户登录（不需要认证）
	v1.Post("/user/login", deps.UserHandler.Login)

	// 用户管理（暂不使用认证中间件）
	user := v1.Group("/user")
	{
		user.Post("/create", deps.UserHandler.Create)
		user.Post("/query", deps.UserHandler.Query)
		user.Post("/update", deps.UserHandler.Update)
		user.Delete("/:id", deps.UserHandler.Delete)
	}
}
