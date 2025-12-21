package response

import (
	"github.com/gofiber/fiber/v2"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Succeed 成功响应
func Succeed(c *fiber.Ctx, resp interface{}, err error) error {
	var body Body
	if err != nil {
		body.Code = -1
		body.Msg = err.Error()
	} else {
		body.Code = 0
		body.Msg = "OK"
		body.Data = resp
	}
	return c.Status(fiber.StatusOK).JSON(body)
}

// Failed 失败响应
func Failed(c *fiber.Ctx, err error) error {
	var body Body
	body.Code = -1
	if err != nil {
		body.Msg = err.Error()
	} else {
		body.Msg = "请求失败"
	}
	body.Data = nil
	return c.Status(fiber.StatusOK).JSON(body)
}

// CusFailed 自定义错误响应
func CusFailed(c *fiber.Ctx, code int, err error) error {
	var body Body
	body.Code = code
	if err != nil {
		body.Msg = err.Error()
	} else {
		body.Msg = "请求失败"
	}
	body.Data = nil
	return c.Status(fiber.StatusOK).JSON(body)
}

