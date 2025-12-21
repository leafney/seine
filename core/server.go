package core

import (
	"github.com/leafney/seine/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// NewServer 创建并启动 Fiber 服务器
func NewServer(app *App, cfg *config.Config) *fiber.App {
	// 创建 Fiber 实例
	fiberApp := fiber.New(fiber.Config{
		AppName: "Seine API",
	})

	// 添加内置中间件
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	// 注册路由
	RegisterRoutes(fiberApp, app)

	return fiberApp
}

