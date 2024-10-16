/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-16 01:15
 * @Description:
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartHttpServer(injector *Injector, quit chan struct{}) {

	f := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// default middlewares
	//f.Use()

	// middlewares
	injector.R.UseMiddlewares(f)

	// routers
	injector.R.SetupRoutes(f)

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
		close(quit)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := f.ShutdownWithContext(ctx); err != nil {
			injector.L.Errorf("[Server] Shutdown error [%v]", err)
		} else {
			injector.L.Info("[Server] Shutdown successful")
		}
		time.Sleep(100 * time.Millisecond)
	}()

	// start
	injector.L.Info("[Server] Load successful")
	if err := f.Listen(fmt.Sprintf(":%s", "8080")); err != nil {
		injector.L.Errorf("[Server] Listen error [%v]", err)
	}

	wg.Wait()
	injector.L.Info("[Server] Exit successful")
}
