package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/seine/response"
)

// ExampleHandler 示例API处理器
type ExampleHandler struct{}

// NewExampleHandler 创建示例API处理器
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// Hello 示例Hello接口
func (h *ExampleHandler) Hello(c *fiber.Ctx) error {
	return response.Succeed(c, map[string]interface{}{
		"message": "Hello from Seine Framework!",
		"version": "1.0.0",
	}, nil)
}

// Ping 健康检查接口
func (h *ExampleHandler) Ping(c *fiber.Ctx) error {
	return response.Succeed(c, map[string]string{
		"status": "ok",
	}, nil)
}
