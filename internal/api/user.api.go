package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/seine/internal/biz"
	"github.com/leafney/seine/internal/vmodel"
	"github.com/leafney/seine/response"
)

// UserHandler 用户API处理器
type UserHandler struct {
	userBiz *biz.UserBiz
}

// NewUserHandler 创建用户API处理器
func NewUserHandler(userBiz *biz.UserBiz) *UserHandler {
	return &UserHandler{userBiz: userBiz}
}

// Create 创建用户
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req vmodel.UserCreateReq
	if err := c.BodyParser(&req); err != nil {
		return response.Failed(c, err)
	}

	err := h.userBiz.Create(c.Context(), &req)
	return response.Succeed(c, nil, err)
}

// Update 更新用户
func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req vmodel.UserUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.Failed(c, err)
	}

	err := h.userBiz.Update(c.Context(), &req)
	return response.Succeed(c, nil, err)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.Failed(c, err)
	}

	err = h.userBiz.Delete(c.Context(), uint(id))
	return response.Succeed(c, nil, err)
}

// Query 查询用户列表
func (h *UserHandler) Query(c *fiber.Ctx) error {
	var req vmodel.UserQueryReq
	if err := c.BodyParser(&req); err != nil {
		return response.Failed(c, err)
	}

	resp, err := h.userBiz.Query(c.Context(), &req)
	return response.Succeed(c, resp, err)
}

// Login 用户登录
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req vmodel.LoginReq
	if err := c.BodyParser(&req); err != nil {
		return response.Failed(c, err)
	}

	resp, err := h.userBiz.Login(c.Context(), &req)
	return response.Succeed(c, resp, err)
}
