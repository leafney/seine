package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddRoute 添加路由
func (g *ComponentGenerator) AddRoute(routeName string) error {
	// 检查路由名称
	if routeName == "" {
		return fmt.Errorf("路由名称不能为空")
	}

	// 创建路由文件
	routerDir := filepath.Join(g.ProjectPath, "internal", "router")
	if err := os.MkdirAll(routerDir, 0755); err != nil {
		return fmt.Errorf("创建路由目录失败: %v", err)
	}

	// 路由文件名
	routerFile := filepath.Join(routerDir, fmt.Sprintf("%s.go", strings.ToLower(routeName)))

	// 检查文件是否已存在
	if _, err := os.Stat(routerFile); err == nil {
		return fmt.Errorf("路由文件 %s 已存在", routerFile)
	}

	// 创建路由文件
	file, err := os.Create(routerFile)
	if err != nil {
		return fmt.Errorf("创建路由文件失败: %v", err)
	}
	defer file.Close()

	// 生成路由名称（首字母大写）
	routeNameUpper := strings.ToUpper(routeName[:1]) + routeName[1:]

	// 写入路由文件内容
	routerContent := fmt.Sprintf(`/**
 * @Project:     %s
 * @Date:        %s
 */

package router

import (
	"github.com/gofiber/fiber/v2"
	"%s/pkg/response"
)

// %sRouter %s路由
func %sRouter(app *fiber.App) {
	// 在这里添加%s相关的路由
	%sGroup := app.Group("/%s")
	{
		// 获取%s列表
		%sGroup.Get("/", Get%sList)
	}
}

// Get%sList 获取%s列表
func Get%sList(c *fiber.Ctx) error {
	// 这里是获取%s列表的示例实现
	return response.OkWithData(c, map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"id": 1,
				"name": "示例%s 1",
				"created_at": "2025-01-01 12:00:00",
			},
			{
				"id": 2,
				"name": "示例%s 2",
				"created_at": "2025-01-02 12:00:00",
			},
		},
		"total": 2,
	})
}
`,
		g.ProjectName, g.Date, g.ModulePath,
		routeNameUpper, routeName, routeNameUpper, routeName,
		strings.ToLower(routeName), strings.ToLower(routeName),
		routeName, strings.ToLower(routeName), routeNameUpper,
		routeNameUpper, routeName, routeNameUpper,
		routeName, routeName, routeName)

	if _, err := file.WriteString(routerContent); err != nil {
		return fmt.Errorf("写入路由文件失败: %v", err)
	}

	// 更新路由注册
	return g.updateRouterRegistration(routeName)
}

// updateRouterRegistration 更新路由注册
func (g *ComponentGenerator) updateRouterRegistration(routeName string) error {
	// 主路由文件
	routerFile := filepath.Join(g.ProjectPath, "internal", "router", "router.go")

	// 检查主路由文件是否存在
	if _, err := os.Stat(routerFile); os.IsNotExist(err) {
		// 创建主路由文件
		file, err := os.Create(routerFile)
		if err != nil {
			return fmt.Errorf("创建主路由文件失败: %v", err)
		}
		defer file.Close()

		// 生成路由名称（首字母大写）
		routeNameUpper := strings.ToUpper(routeName[:1]) + routeName[1:]

		// 写入主路由文件内容
		routerContent := fmt.Sprintf(`/**
 * @Project:     %s
 * @Date:        %s
 */

package router

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"%s/internal/api"
	"%s/internal/service"
)

// Router 路由管理器
type Router struct {
	HomeApi *api.Home

	// 服务列表
	CronSvc   *service.Cron
	NotifySvc *service.Notify
}

// Init 初始化
func (r *Router) Init(ctx context.Context) error {
	// 自动数据库迁移
	if err := r.AutoMigrate(); err != nil {
		return err
	}

	// 自动启动服务
	if err := r.AutoStart(ctx); err != nil {
		return err
	}

	return nil
}

// AutoMigrate 自动迁移数据库
func (r *Router) AutoMigrate() error {
	// 在这里添加数据库迁移代码
	return nil
}

// AutoStart 自动启动服务
func (r *Router) AutoStart(ctx context.Context) error {
	// 启动定时任务服务
	if r.CronSvc != nil {
		if err := r.CronSvc.LoadJobs(ctx); err != nil {
			return err
		}
		// 延迟5秒启动定时任务
		time.AfterFunc(5*time.Second, func() {
			r.CronSvc.Start()
		})
	}

	// 启动通知服务
	if r.NotifySvc != nil {
		if err := r.NotifySvc.Start(); err != nil {
			return err
		}
	}

	return nil
}

// UseMiddlewares 使用中间件
func (r *Router) UseMiddlewares(app *fiber.App, injector interface{}) {
	app.Use(cors.New())
	app.Use(logger.New())
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes(app *fiber.App, injector interface{}) {
	// 健康检查路由
	app.Get("/health", r.Health)

	// API路由组
	v1 := app.Group("/api/v1")
	{
		// 版本信息
		v1.Get("/version", r.HomeApi.Version)

		// 在这里添加更多API路由
	}

	// 注册路由
	%sRouter(app)
}

// Health 健康检查
func (r *Router) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
`,
			g.ProjectName, g.Date, g.ModulePath, g.ModulePath, routeNameUpper)

		if _, err := file.WriteString(routerContent); err != nil {
			return fmt.Errorf("写入主路由文件失败: %v", err)
		}

		return nil
	}

	// 读取主路由文件内容
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return fmt.Errorf("读取主路由文件失败: %v", err)
	}

	// 生成路由名称（首字母大写）
	routeNameUpper := strings.ToUpper(routeName[:1]) + routeName[1:]

	// 检查路由是否已注册
	if strings.Contains(string(content), fmt.Sprintf("%sRouter(app)", routeNameUpper)) {
		return nil
	}

	// 查找注册路由的位置
	lines := strings.Split(string(content), "\n")
	insertIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "// 注册路由") {
			insertIndex = i + 1
			break
		}
	}

	if insertIndex == -1 {
		return fmt.Errorf("未找到路由注册位置")
	}

	// 插入新路由注册
	lines = append(lines[:insertIndex+1], lines[insertIndex:]...)
	lines[insertIndex] = fmt.Sprintf("\t%sRouter(app)", routeNameUpper)

	// 写入更新后的内容
	if err := os.WriteFile(routerFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("更新主路由文件失败: %v", err)
	}

	return nil
}
