package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Generator 项目生成器
type Generator struct {
	ProjectName string
	ProjectPath string
	TemplateDir string
	ModulePath  string
	Date        string
	BuildTime   string
}

// NewGenerator 创建新的项目生成器
func NewGenerator(name, path string) *Generator {
	// 获取模板目录的绝对路径
	execPath, _ := os.Executable()
	templateDir := filepath.Join(filepath.Dir(execPath), "template")

	// 如果是开发环境，使用相对路径
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		templateDir = filepath.Join("template")
	}

	// 处理项目名称和模块路径
	var modulePath string

	// 检查是否包含域名和路径分隔符
	if strings.Contains(name, "/") || strings.Contains(name, ".") {
		// 完整的模块路径形式，如 github.com/user/repo
		modulePath = name
	} else {
		// 单个单词形式，如 helloWorld
		// 添加默认前缀
		modulePath = fmt.Sprintf("github.com/%s", name)
	}

	return &Generator{
		ProjectName: name,
		ProjectPath: path,
		TemplateDir: templateDir,
		ModulePath:  modulePath,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
		BuildTime:   time.Now().Format("2006-01-02 15:04:05"),
	}
}

// Generate 生成项目
func (g *Generator) Generate() error {
	// 创建项目根目录
	if err := os.MkdirAll(g.ProjectPath, 0755); err != nil {
		return err
	}

	// 创建基础目录结构
	dirs := []string{
		"cmd",
		"config",
		"config/cache",
		"config/vars",
		"internal",
		"internal/api",
		"internal/biz",
		"internal/dao",
		"internal/model",
		"internal/service",
		"internal/vmodel",
		"pkg",
		"pkg/errx",
		"pkg/middlewarex",
		"pkg/response",
		"pkg/utils",
		"pkg/xlogx",
		"web",
		"data",
		"logs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(g.ProjectPath, dir), 0755); err != nil {
			return err
		}
	}

	// 生成模板文件
	if err := g.generateFiles(); err != nil {
		return err
	}

	return nil
}

// generateFiles 生成项目文件
func (g *Generator) generateFiles() error {
	// 定义需要生成的文件
	files := map[string]string{
		"main.go":                    "main.go.tmpl",
		"go.mod":                     "go.mod.tmpl",
		"Makefile":                   "Makefile.tmpl",
		".gitignore":                 "gitignore.tmpl",
		"cmd/injector.go":            "cmd/injector.go.tmpl",
		"cmd/router.go":              "cmd/router.go.tmpl",
		"cmd/start.go":               "cmd/start.go.tmpl",
		"cmd/wire.go":                "cmd/wire.go.tmpl",
		"config/config.go":           "config/config.go.tmpl",
		"config/config.toml.default": "config/config.toml.default.tmpl",
		"config/vars/vars.go":        "config/vars/vars.go.tmpl",
		"internal/wire.go":           "internal/wire.go.tmpl",
		"internal/router/router.go":  "internal/router/router.go.tmpl",
		"internal/api/home.api.go":   "internal/api/home.api.go.tmpl",

		// pkg 目录下的基础库
		"pkg/errx/errx.go":          "pkg/errx/errx.go.tmpl",
		"pkg/middlewarex/logger.go": "pkg/middlewarex/logger.go.tmpl",
		"pkg/response/response.go":  "pkg/response/response.go.tmpl",
		"pkg/xlogx/xlogx.go":        "pkg/xlogx/xlogx.go.tmpl",

		// 添加更多 pkg 目录下的库
		"pkg/versionx/version.go":  "pkg/versionx/versionx.go.tmpl",
		"pkg/rmqx/rmqx.go":         "pkg/rmqx/rmqx.go.tmpl",
		"pkg/redisx/redisx.go":     "pkg/redisx/redisx.go.tmpl",
		"pkg/parsex/parsex.go":     "pkg/parsex/parsex.go.tmpl",
		"pkg/notifyx/notifyx.go":   "pkg/notifyx/notifyx.go.tmpl",
		"pkg/leveldbx/leveldbx.go": "pkg/leveldbx/leveldbx.go.tmpl",
		"pkg/errc/errc.go":         "pkg/errc/errc.go.tmpl",
		"pkg/cronx/cronx.go":       "pkg/cronx/cronx.go.tmpl",
		"pkg/cachex/cachex.go":     "pkg/cachex/cachex.go.tmpl",
		"pkg/zlogx/zlogx.go":       "pkg/zlogx/zlogx.go.tmpl",

		"README.md": "README.md.tmpl",
	}

	for filePath, templateName := range files {
		if err := g.generateFile(filePath, templateName); err != nil {
			return err
		}
	}

	return nil
}

// generateFile 生成单个文件
func (g *Generator) generateFile(filePath, templateName string) error {
	// 读取模板文件
	templatePath := filepath.Join(g.TemplateDir, templateName)

	// 如果是开发环境，使用嵌入的模板
	var tmplContent string
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		tmplContent = getEmbeddedTemplate(templateName)
	} else {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return err
		}
		tmplContent = string(content)
	}

	// 解析模板
	tmpl, err := template.New(templateName).Parse(tmplContent)
	if err != nil {
		return err
	}

	// 创建目标文件
	targetPath := filepath.Join(g.ProjectPath, filePath)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, g); err != nil {
		return err
	}

	return nil
}

// 临时函数，用于获取嵌入的模板
// 在实际实现中，应该使用go:embed功能
func getEmbeddedTemplate(name string) string {
	// 这里只是示例，实际应该使用go:embed
	templates := map[string]string{
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

func main() {
	pflag.BoolVarP(&h, "help", "h", false, "help")
	pflag.BoolVarP(&v, "version", "v", false, "version")
	pflag.BoolVarP(&d, "debug", "d", false, "debug mode")
	pflag.StringVarP(&p, "port", "p", vars.WebDefaultPort, "server listen port")
	pflag.Parse()

	if h {
		pflag.PrintDefaults()
	} else if v {
		fmt.Println("Version:      " + Version)
		fmt.Println("Git branch:   " + GitBranch)
		fmt.Println("Git commit:   " + GitCommit)
		fmt.Println("Built time:   " + BuildTime)
		fmt.Println("Go version:   " + runtime.Version())
		fmt.Println("OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH)
	} else {
		quitChan := make(chan struct{})
		injector, callback, err := cmd.BuildInjector(quitChan)
		if err != nil {
			log.Fatalln(err)
		}
		defer callback()

		ctx := context.Background()
		if err := injector.R.Init(ctx); err != nil {
			log.Fatalf("初始化异常 %v", err)
		}

		cmd.StartServer(injector, p, quitChan)
	}
}`,
		"go.mod.tmpl": `module {{.ModulePath}}

go 1.20

require (
	github.com/gofiber/fiber/v2 v2.48.0
	github.com/google/wire v0.5.0
	github.com/knadh/koanf/parsers/toml/v2 v2.0.1
	github.com/knadh/koanf/providers/file v0.1.0
	github.com/knadh/koanf/v2 v2.0.1
	github.com/creasty/defaults v1.7.0
	github.com/spf13/pflag v1.0.5
)`,
		"Makefile.tmpl": `# Makefile for {{.ProjectName}}

.PHONY: build clean run

APP_NAME={{.ProjectName}}
BUILD_DIR=./dist
BUILD_TIME=$(shell date "+%Y-%m-%d %H:%M:%S")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(shell git symbolic-ref --short -q HEAD 2>/dev/null || echo "unknown")
VERSION=v0.1.0

LDFLAGS=-ldflags "-X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -X main.GitBranch=${GIT_BRANCH} -X main.Version=${VERSION}"

build:
	@echo "Building ${APP_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME} .
	@echo "Build complete: ${BUILD_DIR}/${APP_NAME}"

clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@echo "Clean complete"

run:
	@go run ${LDFLAGS} main.go

dev:
	@air -c .air.toml`,
		"gitignore.tmpl": `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# IDE files
.idea/
.vscode/

# Build output
dist/

# Config files
data/config.toml

# Log files
logs/

# Mac OS
.DS_Store`,
		"cmd/injector.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package cmd

import (
	"github.com/google/wire"
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/internal"
	"{{.ModulePath}}/pkg/xlogx"
)

var AppSet = wire.NewSet(
	config.NewConfig,
	xlogx.NewXLogSvc,
	internal.Set,
)

type Injector struct {
	L *xlogx.XLogSvc
	R DefRouter
	C *config.Config
}`,
		"cmd/router.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package cmd

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"{{.ModulePath}}/pkg/middlewarex"
)

type DefRouter struct {
	// 在这里添加API和Service引用
}

func (r *DefRouter) Init(ctx context.Context) error {
	return nil
}

func (r *DefRouter) UseMiddlewares(app *fiber.App, inj *Injector) {
	app.Use(cors.New())
	app.Use(middlewarex.RequestLogger(inj.L))
}

func (r *DefRouter) SetupRoutes(app *fiber.App, inj *Injector) {
	// API路由设置
	v1G := app.Group("/api/v1")
	{
		// 在这里添加API路由
		v1G.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "ok",
			})
		})
	}
}`,
		"cmd/start.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/config/vars"
)

func StartServer(injector *Injector, port string, quit chan struct{}) {
	f := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// 使用中间件
	injector.R.UseMiddlewares(f, injector)

	// 设置路由
	injector.R.SetupRoutes(f, injector)

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
		time.Sleep(200 * time.Millisecond)
	}()

	// 启动服务
	injector.L.Info("[Server] Load successful")

	// 监听端口，默认 > 配置文件 > 命令行
	defPort := vars.WebDefaultPort
	if port != vars.WebDefaultPort {
		defPort = port
	}

	if err := f.Listen(fmt.Sprintf(":%s", defPort)); err != nil {
		injector.L.Fatalf("[Server] Listen error [%v]", err)
	}

	wg.Wait()
	injector.L.Info("[Server] Exit successful")
}`,
		"cmd/wire.go.tmpl": `//go:build wireinject
// +build wireinject

/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package cmd

import (
	"github.com/google/wire"
)

func BuildInjector(quit chan struct{}) (*Injector, func(), error) {
	wire.Build(AppSet, wire.Struct(new(Injector), "*"))
	return &Injector{}, nil, nil
}`,
		"config/config.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package config

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"github.com/creasty/defaults"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

//go:embed config.toml.default
var DefaultConfig embed.FS

type Config struct {
	Port  string
	Log   Log
	Redis Redis
}

type (
	Log struct {
		Enable bool   ` + "`default:\"true\"`" + `
		Debug  bool   ` + "`default:\"true\"`" + `
		Level  string
	}

	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
)

var k = koanf.New(".")

func NewConfig() (*Config, error) {
	cfg := new(Config)
	path := "data/config.toml"

	if exist := loadDefaultConfig(path); !exist {
		log.Printf("[Koanf] Create default config file [%v]", path)
		os.Exit(0)
	}

	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		log.Fatalf("[Koanf] Load config file [%v] error [%v]", path, err)
	}

	// 先设置默认值
	if err := defaults.Set(cfg); err != nil {
		log.Fatalf("[Koanf] Set default value error [%v]", err)
	}
	// 再解析配置参数
	if err := k.Unmarshal("", cfg); err != nil {
		log.Fatalf("[Koanf] Unmarshal config error [%v]", err)
	}

	log.Println("[Koanf] Load successful")
	return cfg, nil
}

func loadDefaultConfig(path string) (exist bool) {
	exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
		// 保证配置文件所在目录存在
		if err := os.MkdirAll("data", 0755); err != nil {
			log.Fatalf("[Koanf] Failed to create config directory: %v", err)
		}

		data, err := fs.ReadFile(DefaultConfig, "config.toml.default")
		if err != nil {
			log.Fatalf("[Koanf] Failed to read embedded config: %v", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			log.Fatalf("[Koanf] Failed to write default config file: %v", err)
		}

		log.Println("[Koanf] Default config file created")
	}
	return exist
}`,
		"config/config.toml.default.tmpl": `# {{.ProjectName}} 配置文件

# 服务端口
Port = "8080"

# 日志配置
[Log]
Enable = true
Debug = true
Level = "info"

# Redis配置
[Redis]
Addr = "localhost:6379"
Pwd = ""
DB = 0`,
		"config/vars/vars.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package vars

const (
	WebDefaultPort = "8080"
)`,
		"internal/wire.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package internal

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	// 在这里添加内部组件的依赖注入
)`,
		"pkg/errx/errx.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package errx

import (
	"errors"
	"fmt"
)

// 错误码定义
const (
	CodeSuccess      = 0
	CodeParamError   = 1001
	CodeInternalError = 1002
	CodeNotFound     = 1004
	CodeUnauthorized = 1005
)

// 错误信息定义
var msgMap = map[int]string{
	CodeSuccess:      "成功",
	CodeParamError:   "参数错误",
	CodeInternalError: "内部错误",
	CodeNotFound:     "资源不存在",
	CodeUnauthorized: "未授权",
}

// Error 自定义错误类型
type Error struct {
	Code int
	Msg  string
}

// Error 实现error接口
func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Msg)
}

// New 创建新的错误
func New(code int) *Error {
	msg, ok := msgMap[code]
	if !ok {
		msg = "未知错误"
	}
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// NewWithMsg 创建带自定义消息的错误
func NewWithMsg(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// FromError 从标准error转换为自定义Error
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	
	return &Error{
		Code: CodeInternalError,
		Msg:  err.Error(),
	}
}`,
		"pkg/middlewarex/logger.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package middlewarex

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/pkg/xlogx"
)

// RequestLogger 请求日志中间件
func RequestLogger(logger *xlogx.XLogSvc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// 处理请求
		err := c.Next()
		
		// 记录请求信息
		latency := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		
		logger.Infof("[%d] %s %s %s %s", status, method, path, ip, latency)
		
		return err
	}
}`,
		"pkg/response/response.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package response

import (
	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/pkg/errx"
)

// Response 统一响应结构
type Response struct {
	Code int         ` + "`json:\"code\"`" + `
	Msg  string      ` + "`json:\"msg\"`" + `
	Data interface{} ` + "`json:\"data,omitempty\"`" + `
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code: errx.CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *fiber.Ctx, err *errx.Error) error {
	return c.JSON(Response{
		Code: err.Code,
		Msg:  err.Msg,
	})
}

// FailWithMsg 带自定义消息的失败响应
func FailWithMsg(c *fiber.Ctx, code int, msg string) error {
	return c.JSON(Response{
		Code: code,
		Msg:  msg,
	})
}`,
		"pkg/xlogx/xlogx.go.tmpl": `/**
 * @Project:     {{.ProjectName}}
 * @Date:        {{.Date}}
 */

package xlogx

import (
	"log"
	"os"
	"path/filepath"

	"{{.ModulePath}}/config"
)

// XLogSvc 日志服务
type XLogSvc struct {
	logger *log.Logger
	debug  bool
}

// NewXLogSvc 创建日志服务
func NewXLogSvc(cfg *config.Config) (*XLogSvc, error) {
	// 确保日志目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	// 创建或打开日志文件
	logFile, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// 创建日志记录器
	logger := log.New(logFile, "", log.LstdFlags)

	return &XLogSvc{
		logger: logger,
		debug:  cfg.Log.Debug,
	}, nil
}

// Info 记录信息日志
func (l *XLogSvc) Info(v ...interface{}) {
	l.logger.Println(v...)
}

// Infof 记录格式化信息日志
func (l *XLogSvc) Infof(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Error 记录错误日志
func (l *XLogSvc) Error(v ...interface{}) {
	l.logger.Println("ERROR:", v)
}

// Errorf 记录格式化错误日志
func (l *XLogSvc) Errorf(format string, v ...interface{}) {
	l.logger.Printf("ERROR: "+format, v...)
}

// Debug 记录调试日志
func (l *XLogSvc) Debug(v ...interface{}) {
	if l.debug {
		l.logger.Println("DEBUG:", v)
	}
}

// Debugf 记录格式化调试日志
func (l *XLogSvc) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.logger.Printf("DEBUG: "+format, v...)
	}
}

// Fatal 记录致命错误并退出
func (l *XLogSvc) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

// Fatalf 记录格式化致命错误并退出
func (l *XLogSvc) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}`,
		"README.md.tmpl": "# {{.ProjectName}}\n\n## 简介\n\n这是一个使用Seine框架生成的Go Web项目。\n\n## 目录结构\n\n- cmd: 命令行相关代码\n- config: 配置相关代码\n- internal: 内部代码\n  - api: API处理器\n  - biz: 业务逻辑\n  - dao: 数据访问\n  - model: 数据模型\n  - service: 服务组件\n- pkg: 公共包\n- web: Web资源\n\n## 快速开始\n\n### 运行\n\n```bash\ngo run main.go\n```\n\n### 构建\n\n```bash\nmake build\n```\n\n### 开发模式\n\n```bash\nmake dev\n```\n\n## 配置\n\n配置文件位于 `data/config.toml`，首次运行会自动生成默认配置文件。",
	}

	return templates[name]
}
