/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:46
 * @Description:
 */

package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/seine/global/respcode"
)

type Response struct {
	Code  int         `json:"code"`  // 自定义错误码
	Msg   string      `json:"msg"`   // 给用户看的错误信息
	Data  interface{} `json:"data"`  // 返回的数据
	Error string      `json:"error"` // 给开发者看的详细错误信息
}

func jsonResult(c *fiber.Ctx, statusCode int, errCode int, msg string, data interface{}, err string) error {
	return c.Status(statusCode).JSON(Response{
		Code:  errCode,
		Msg:   msg,
		Data:  data,
		Error: err,
	})
}

func Ok(c *fiber.Ctx) error {
	msg := respcode.GetErrMsg(respcode.Success, "")
	return jsonResult(c, fiber.StatusOK, respcode.Success, msg, map[string]interface{}{}, "")
}

func OkWithData(c *fiber.Ctx, data interface{}) error {
	msg := respcode.GetErrMsg(respcode.Success, "")
	return jsonResult(c, fiber.StatusOK, respcode.Success, msg, data, "")
}

func OkWithMsg(c *fiber.Ctx, msg string) error {
	return jsonResult(c, fiber.StatusOK, respcode.Success, msg, map[string]interface{}{}, "")
}

func OkWithDetailed(c *fiber.Ctx, msg string, data interface{}) error {
	return jsonResult(c, fiber.StatusOK, respcode.Success, msg, data, "")
}

// *******************

// Fail 操作失败
func Fail(c *fiber.Ctx, err string) error {
	msg := respcode.GetErrMsg(respcode.Failed, "")
	return jsonResult(c, fiber.StatusOK, respcode.Failed, msg, map[string]interface{}{}, err)
}

// FailWithCode 使用统一错误消息
func FailWithCode(c *fiber.Ctx, errCode int, err string) error {
	msg := respcode.GetErrMsg(errCode, "")
	return jsonResult(c, fiber.StatusOK, errCode, msg, map[string]interface{}{}, err)
}

// FailWithMsg 自定义错误消息
func FailWithMsg(c *fiber.Ctx, msg, err string) error {
	return jsonResult(c, fiber.StatusOK, respcode.Failed, msg, map[string]interface{}{}, err)
}

// FailWith400Msg client error
func FailWith400Msg(c *fiber.Ctx, err string) error {
	msg := respcode.GetErrMsg(respcode.ErrClient, "")
	return jsonResult(c, fiber.StatusBadRequest, respcode.ErrClient, msg, map[string]interface{}{}, err)
}

// FailWith401Msg Unauthorized 认证相关
func FailWith401Msg(c *fiber.Ctx, err string) error {
	msg := respcode.GetErrMsg(respcode.ErrUnAuthorized, "")
	return jsonResult(c, fiber.StatusUnauthorized, respcode.ErrUnAuthorized, msg, map[string]interface{}{}, err)
}

// FailWith403Msg Forbidden 权限相关
func FailWith403Msg(c *fiber.Ctx, err string) error {
	msg := respcode.GetErrMsg(respcode.ErrForbidden, "")
	return jsonResult(c, fiber.StatusUnauthorized, respcode.ErrForbidden, msg, map[string]interface{}{}, err)
}

// FailWith500Msg server error
func FailWith500Msg(c *fiber.Ctx, err string) error {
	msg := respcode.GetErrMsg(respcode.ErrServer, "")
	return jsonResult(c, fiber.StatusInternalServerError, respcode.ErrServer, msg, map[string]interface{}{}, err)
}
