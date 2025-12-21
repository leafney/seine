package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leafney/seine/config"
	"github.com/leafney/seine/core"
	"github.com/leafney/seine/pkg/versionx"
	"github.com/leafney/seine/wire"
)

var (
	configFile = flag.String("f", "data/config.yaml", "配置文件路径")
	version    = flag.Bool("v", false, "显示版本信息")

	// 编译时注入（通过 Makefile）
	Version   = "dev"
	GitBranch = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	flag.Parse()

	// 注入版本信息
	versionx.VersionInfo.Version = Version
	versionx.VersionInfo.GitCommit = GitCommit
	versionx.VersionInfo.BuildTime = BuildTime

	// 显示版本信息
	if *version {
		fmt.Printf("Seine Framework\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Git Branch: %s\n", GitBranch)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
		fmt.Printf("  Build Time: %s\n", BuildTime)
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
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
		fmt.Printf("Config:     %s\n", *configFile)
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
