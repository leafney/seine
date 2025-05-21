/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:54
 * @Description:
 */

package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/seine/internal/api"
	"github.com/leafney/seine/internal/model"
	"github.com/leafney/seine/pkg/gormx"
)

type DefRouter struct {
	DB      *gormx.GormDBSvc
	memoApi *api.MemoApi
}

func (r *DefRouter) AutoMigrate() error {
	return r.DB.AutoMigrate(
		new(model.Memo),
	)
}

func (r *DefRouter) Init() error {
	if err := r.AutoMigrate(); err != nil {
		return err
	}
	return nil
}

func (r *DefRouter) UseMiddlewares(app *fiber.App) {
	app.Use(Cors())
}

func (r *DefRouter) SetupRoutes(app *fiber.App) {

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})

	app.Get("/aaa", r.memoApi.List)

}
