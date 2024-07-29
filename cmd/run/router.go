/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-29 17:07
 * @Description:
 */

package run

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/leafney/seine/global"
	"github.com/leafney/seine/web"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

func bindRouter(app *fiber.App) {

	// 非登录请求
	//homeG := app.Group("/api/v1")
	{
		//homeG.Post("/captcha", handler.SendCaptcha) // 发送验证码
		//homeG.Post("/login", handler.Login)         // 登录
	}

	// api 后台请求
	//v1G := app.Group("/api/v1", middleware.JWTAuth())
	//v1.Get("/user", func(c *fiber.Ctx) error {
	//	return c.JSON(fiber.Map{
	//		"say": "hello world",
	//	})
	//})
	{
		//v1G.Get("/test/:id", handler.Test1)

	}

	// webui
	uiDist, err := fs.Sub(web.UiStatic, "dist")
	if err != nil {
		global.GLog.Fatal("static dir load error [%v]", zap.Error(err))
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(uiDist),
	}))

}
