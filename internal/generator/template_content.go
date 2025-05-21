/**
 * @Project:     seine
 * @Date:        2025-05-21
 */

package generator

// GetTemplateContent 获取模板文件内容
// 这个函数替代了原来的 getEmbeddedTemplate 函数
// 所有模板内容都集中在这里管理，方便后续修改
func GetTemplateContent(templateName string) string {
	// 根据模板名称返回对应的模板内容
	if content, ok := templateContents[templateName]; ok {
		return content
	}
	return ""
}

// templateContents 存储所有模板内容的映射
var templateContents = map[string]string{
	// 这里只列出几个示例，实际需要包含所有模板文件
	"main.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"{{.ModulePath}}/cmd"
	"{{.ModulePath}}/config/vars"
	"github.com/spf13/pflag"
)

var (
	v         bool
	h         bool
	d         bool
	p         string
	Version   = "v0.1.0"
	GitBranch = ""
	GitCommit = ""
	BuildTime = "{{.BuildTime}}"
)

func init() {
	pflag.BoolVarP(&v, "version", "v", false, "显示版本信息")
	pflag.BoolVarP(&h, "help", "h", false, "显示帮助信息")
	pflag.BoolVarP(&d, "daemon", "d", false, "以守护进程方式运行")
	pflag.StringVarP(&p, "port", "p", "", "指定端口")
}

func main() {
	pflag.Parse()

	if v {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git Branch: %s\n", GitBranch)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	if h {
		pflag.Usage()
		return
	}

	// 设置端口
	if p != "" {
		vars.HttpPort = p
	}

	// 创建上下文
	ctx := context.Background()

	// 启动应用
	if err := cmd.Start(ctx, d); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}`,

	"pkg/zlogx/zlogx.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package zlogx

import (
	"os"

	"{{.ModulePath}}/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZLogSvc Zap日志服务
type ZLogSvc struct {
	*zap.Logger
	debug bool
}

// NewZLogSvc 创建Zap日志服务
func NewZLogSvc(cfg *config.Config) (*ZLogSvc, error) {
	// 确保日志目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	// 获取配置
	logCfg := cfg.Log
	if !logCfg.ZEnable {
		return &ZLogSvc{
			Logger: zap.NewNop(),
			debug:  false,
		}, nil
	}

	// 设置日志级别
	var level zapcore.Level
	switch logCfg.ZLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 创建日志记录器
	var logger *zap.Logger
	if logCfg.ZCaller {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}

	return &ZLogSvc{
		Logger: logger,
		debug:  logCfg.Debug,
	}, nil
}

// 全局日志实例
var (
	Logger *ZLogSvc
)

// InitLogger 初始化全局日志实例
func InitLogger(cfg *config.Config) error {
	logger, err := NewZLogSvc(cfg)
	if err != nil {
		return err
	}
	Logger = logger
	return nil
}

// GetLogger 获取全局日志实例
func GetLogger() *ZLogSvc {
	if Logger == nil {
		// 如果未初始化，返回一个空日志实例
		return &ZLogSvc{
			Logger: zap.NewNop(),
			debug:  false,
		}
	}
	return Logger
}`,

	"pkg/versionx/versionx.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package versionx

// InfoSvc 版本信息结构体
type InfoSvc struct {
	Version   string
	GitCommit string
	BuildTime string
}

var (
	// VersionInfo 全局版本信息实例
	VersionInfo = &InfoSvc{}
)

// NewInfoSvc 提供一个构造函数用于 wire 注入
func NewInfoSvc() *InfoSvc {
	return VersionInfo
}

// GetVersion 获取版本信息
func (v *InfoSvc) GetVersion() map[string]string {
	return map[string]string{
		"version":    v.Version,
		"git_commit": v.GitCommit,
		"build_time": v.BuildTime,
	}
}`,

	// 可以继续添加其他模板文件...
}
