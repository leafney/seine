/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:35
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/seine/internal/biz"
)

type MemoApi struct {
	memoBiz *biz.MemoBiz
}

func NewMemoApi(biz *biz.MemoBiz) *MemoApi {
	return &MemoApi{
		memoBiz: biz,
	}
}

func (a *MemoApi) Add(c *fiber.Ctx) error {

	return c.JSON("")
}

func (a *MemoApi) List(c *fiber.Ctx) error {

	return c.JSON("helloo")
}
