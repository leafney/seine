/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:52
 * @Description:
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/leafney/seine/internal/api"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	*fiber.App
	port string
}

func NewApp(api *api.MemoApi) *App {
	app := fiber.New()

	// TODO router middlewares
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})

	app.Get("/memo", api.List)

	return &App{App: app, port: "8090"}
}

func (a *App) Start(stop chan struct{}) {

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// 添加信号监听，支持优雅退出
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		// 等待停止信号
		<-signalChan

		// 关闭通道
		close(stop)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.ShutdownWithContext(ctx); err != nil {
			log.Errorf("[Server] Shutdown error [%v]", err)
		} else {
			log.Info("[Server] Shutdown successful")
		}
		time.Sleep(100 * time.Millisecond)
	}()

	// start
	log.Info("[Server] Load successful")
	if err := a.Listen(fmt.Sprintf(":%s", a.port)); err != nil {
		log.Errorf("[Server] Listen error [%v]", err)
	}

	wg.Wait()
	log.Info("[Server] Exit successful")
}
