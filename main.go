package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/leafney/seine/config"
	"github.com/leafney/seine/core"
	"github.com/leafney/seine/pkg/versionx"
	"github.com/leafney/seine/wire"
	"github.com/spf13/pflag"
)

var (
	v bool
	h bool
	c string

	// 编译时注入（通过 Makefile）
	Version   = "dev"
	GitBranch = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	pflag.StringVarP(&c, "config", "f", "data/config.yaml", "配置文件路径")
	pflag.BoolVarP(&v, "version", "v", false, "显示版本信息")
	pflag.BoolVarP(&h, "help", "h", false, "显示帮助信息")
	pflag.Parse()

	// 注入版本信息
	versionx.VersionInfo.Version = Version
	versionx.VersionInfo.GitCommit = GitCommit
	versionx.VersionInfo.BuildTime = BuildTime

	if h {
		pflag.PrintDefaults()
	} else if v {
		// 显示版本信息
		fmt.Printf("Seine Framework\n")
		fmt.Println("Version:      " + Version)
		fmt.Println("Git branch:   " + GitBranch)
		fmt.Println("Git commit:   " + GitCommit)
		fmt.Println("Built time:   " + BuildTime)
		fmt.Println("Go version:   " + runtime.Version())
		fmt.Println("OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH)
	} else {
		// 加载配置
		cfg, err := config.LoadConfig(c)
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}

		// 创建停止信号通道
		stop := make(chan struct{})

		// Wire 初始化应用
		app, err := wire.InitializeApp(cfg, stop)
		if err != nil {
			log.Fatalf("初始化应用失败: %v", err)
		}

		// 创建 Fiber 服务器
		fiberApp := core.NewServer(app, cfg)

		// 服务器地址
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

		// 启动服务器（在 goroutine 中）
		go func() {
			fmt.Printf("\n========================================\n")
			fmt.Printf("Seine Framework Started\n")
			fmt.Printf("========================================\n")
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Config:     %s\n", c)
			fmt.Printf("Address:    http://%s\n", addr)
			fmt.Printf("========================================\n\n")

			if err := fiberApp.Listen(addr); err != nil {
				log.Fatalf("服务器启动失败: %v", err)
			}
		}()

		// 监听系统信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("\n正在关闭服务器...")

		// 关闭停止信号通道
		close(stop)

		// 优雅关闭 Fiber 服务器
		if err := fiberApp.Shutdown(); err != nil {
			log.Printf("服务器关闭失败: %v", err)
		}

		fmt.Println("服务器已关闭")
	}
}