package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ComponentGenerator 组件生成器
type ComponentGenerator struct {
	ProjectName string
	ProjectPath string
	Date        string
	ModulePath  string
	TemplateDir string
	BuildTime   string
}

// NewComponentGenerator 创建新的组件生成器
func NewComponentGenerator(name, projectPath string) *ComponentGenerator {
	// 获取模块名称
	goModPath := filepath.Join(projectPath, "go.mod")
	moduleName := ""
	
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "module ") {
					moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					break
				}
			}
		}
	}
	
	// 获取模板目录的绝对路径
	execPath, _ := os.Executable()
	templateDir := filepath.Join(filepath.Dir(execPath), "template")
	
	// 如果是开发环境，使用相对路径
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		templateDir = filepath.Join("template")
	}
	
	return &ComponentGenerator{
		ProjectName: name,
		ProjectPath: projectPath,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
		ModulePath:  moduleName,
		TemplateDir: templateDir,
		BuildTime:   time.Now().Format("2006-01-02 15:04:05"),
	}
}

// AddService 添加服务组件
func (g *ComponentGenerator) AddService(serviceName string) error {
	// 确保服务名称首字母大写
	serviceName = strings.Title(strings.ToLower(serviceName))
	
	// 服务文件路径
	servicePath := filepath.Join(g.ProjectPath, "internal/service", strings.ToLower(serviceName)+".go")
	
	// 检查文件是否已存在
	if _, err := os.Stat(servicePath); err == nil {
		return fmt.Errorf("服务 %s 已存在", serviceName)
	}
	
	// 服务模板内容
	serviceTemplate := `/**
 * @Date:        {{.Date}}
 */

package service

import (
	"context"
	"log"
)

// {{.ServiceName}} 服务
type {{.ServiceName}} struct {
	// 在这里添加依赖
	logger *log.Logger
}

// New{{.ServiceName}} 创建{{.ServiceName}}服务
func New{{.ServiceName}}() (*{{.ServiceName}}, error) {
	return &{{.ServiceName}}{
		logger: log.New(os.Stdout, "[{{.ServiceName}}] ", log.LstdFlags),
	}, nil
}

// Start 启动服务
func (s *{{.ServiceName}}) Start(ctx context.Context) error {
	s.logger.Println("{{.ServiceName}}服务启动")
	return nil
}

// Stop 停止服务
func (s *{{.ServiceName}}) Stop(ctx context.Context) error {
	s.logger.Println("{{.ServiceName}}服务停止")
	return nil
}
`
	
	// 创建模板数据
	data := struct {
		Date        string
		ServiceName string
		ModulePath  string
	}{
		Date:        g.Date,
		ServiceName: serviceName,
		ModulePath:  g.ModulePath,
	}
	
	// 解析模板
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}
	
	// 创建文件
	file, err := os.Create(servicePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	
	// 更新wire.go文件
	if err := g.updateWireFile(serviceName, "service"); err != nil {
		return err
	}
	
	return nil
}

// AddController 添加控制器组件
func (g *ComponentGenerator) AddController(controllerName string) error {
	// 确保控制器名称首字母大写
	controllerName = strings.Title(strings.ToLower(controllerName))
	
	// 创建API、BIZ和DAO文件
	if err := g.createAPIFile(controllerName); err != nil {
		return err
	}
	
	if err := g.createBIZFile(controllerName); err != nil {
		return err
	}
	
	if err := g.createDAOFile(controllerName); err != nil {
		return err
	}
	
	// 更新wire.go文件
	if err := g.updateWireFile(controllerName, "controller"); err != nil {
		return err
	}
	
	// 更新router.go文件
	if err := g.updateRouterFile(controllerName); err != nil {
		return err
	}
	
	return nil
}

// createAPIFile 创建API文件
func (g *ComponentGenerator) createAPIFile(controllerName string) error {
	// API文件路径
	apiPath := filepath.Join(g.ProjectPath, "internal/api", strings.ToLower(controllerName)+".go")
	
	// 检查文件是否已存在
	if _, err := os.Stat(apiPath); err == nil {
		return fmt.Errorf("API %s 已存在", controllerName)
	}
	
	// API模板内容
	apiTemplate := `/**
 * @Date:        {{.Date}}
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/internal/biz"
	"{{.ModulePath}}/pkg/response"
)

// {{.ControllerName}} 控制器
type {{.ControllerName}} struct {
	biz *biz.{{.ControllerName}}
}

// New{{.ControllerName}} 创建{{.ControllerName}}控制器
func New{{.ControllerName}}(biz *biz.{{.ControllerName}}) *{{.ControllerName}} {
	return &{{.ControllerName}}{
		biz: biz,
	}
}

// Get 获取{{.ControllerName}}
func (a *{{.ControllerName}}) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	
	result, err := a.biz.Get(c.Context(), id)
	if err != nil {
		return response.Fail(c, err)
	}
	
	return response.Success(c, result)
}

// Create 创建{{.ControllerName}}
func (a *{{.ControllerName}}) Create(c *fiber.Ctx) error {
	var req struct {
		Name string ` + "`json:\"name\"`" + `
		// 添加其他字段
	}
	
	if err := c.BodyParser(&req); err != nil {
		return response.FailWithMsg(c, 400, "无效的请求参数")
	}
	
	result, err := a.biz.Create(c.Context(), req.Name)
	if err != nil {
		return response.Fail(c, err)
	}
	
	return response.Success(c, result)
}
`
	
	// 创建模板数据
	data := struct {
		Date           string
		ControllerName string
		ModulePath     string
	}{
		Date:           g.Date,
		ControllerName: controllerName,
		ModulePath:     g.ModulePath,
	}
	
	// 解析模板
	tmpl, err := template.New("api").Parse(apiTemplate)
	if err != nil {
		return err
	}
	
	// 创建文件
	file, err := os.Create(apiPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	
	return nil
}

// createBIZFile 创建BIZ文件
func (g *ComponentGenerator) createBIZFile(controllerName string) error {
	// BIZ文件路径
	bizPath := filepath.Join(g.ProjectPath, "internal/biz", strings.ToLower(controllerName)+".go")
	
	// 检查文件是否已存在
	if _, err := os.Stat(bizPath); err == nil {
		return fmt.Errorf("BIZ %s 已存在", controllerName)
	}
	
	// BIZ模板内容
	bizTemplate := `/**
 * @Date:        {{.Date}}
 */

package biz

import (
	"context"
	"{{.ModulePath}}/internal/dao"
	"{{.ModulePath}}/pkg/errx"
)

// {{.ControllerName}} 业务逻辑
type {{.ControllerName}} struct {
	dao *dao.{{.ControllerName}}
}

// New{{.ControllerName}} 创建{{.ControllerName}}业务逻辑
func New{{.ControllerName}}(dao *dao.{{.ControllerName}}) *{{.ControllerName}} {
	return &{{.ControllerName}}{
		dao: dao,
	}
}

// Get 获取{{.ControllerName}}
func (b *{{.ControllerName}}) Get(ctx context.Context, id string) (interface{}, *errx.Error) {
	if id == "" {
		return nil, errx.NewWithMsg(errx.CodeParamError, "ID不能为空")
	}
	
	result, err := b.dao.Get(ctx, id)
	if err != nil {
		return nil, errx.FromError(err)
	}
	
	return result, nil
}

// Create 创建{{.ControllerName}}
func (b *{{.ControllerName}}) Create(ctx context.Context, name string) (interface{}, *errx.Error) {
	if name == "" {
		return nil, errx.NewWithMsg(errx.CodeParamError, "名称不能为空")
	}
	
	result, err := b.dao.Create(ctx, name)
	if err != nil {
		return nil, errx.FromError(err)
	}
	
	return result, nil
}
`
	
	// 创建模板数据
	data := struct {
		Date           string
		ControllerName string
		ModulePath     string
	}{
		Date:           g.Date,
		ControllerName: controllerName,
		ModulePath:     g.ModulePath,
	}
	
	// 解析模板
	tmpl, err := template.New("biz").Parse(bizTemplate)
	if err != nil {
		return err
	}
	
	// 创建文件
	file, err := os.Create(bizPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	
	return nil
}

// createDAOFile 创建DAO文件
func (g *ComponentGenerator) createDAOFile(controllerName string) error {
	// DAO文件路径
	daoPath := filepath.Join(g.ProjectPath, "internal/dao", strings.ToLower(controllerName)+".go")
	
	// 检查文件是否已存在
	if _, err := os.Stat(daoPath); err == nil {
		return fmt.Errorf("DAO %s 已存在", controllerName)
	}
	
	// DAO模板内容
	daoTemplate := `/**
 * @Date:        {{.Date}}
 */

package dao

import (
	"context"
	"fmt"
)

// {{.ControllerName}} 数据访问对象
type {{.ControllerName}} struct {
	// 在这里添加依赖，如数据库连接
}

// New{{.ControllerName}} 创建{{.ControllerName}}数据访问对象
func New{{.ControllerName}}() *{{.ControllerName}} {
	return &{{.ControllerName}}{}
}

// Get 获取{{.ControllerName}}
func (d *{{.ControllerName}}) Get(ctx context.Context, id string) (interface{}, error) {
	// 模拟数据，实际应该从数据库获取
	return map[string]interface{}{
		"id":   id,
		"name": "示例{{.ControllerName}}",
	}, nil
}

// Create 创建{{.ControllerName}}
func (d *{{.ControllerName}}) Create(ctx context.Context, name string) (interface{}, error) {
	// 模拟数据，实际应该保存到数据库
	return map[string]interface{}{
		"id":   "1",
		"name": name,
	}, nil
}
`
	
	// 创建模板数据
	data := struct {
		Date           string
		ControllerName string
		ModulePath     string
	}{
		Date:           g.Date,
		ControllerName: controllerName,
		ModulePath:     g.ModulePath,
	}
	
	// 解析模板
	tmpl, err := template.New("dao").Parse(daoTemplate)
	if err != nil {
		return err
	}
	
	// 创建文件
	file, err := os.Create(daoPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	
	return nil
}

// updateWireFile 更新wire.go文件
func (g *ComponentGenerator) updateWireFile(name, componentType string) error {
	// wire.go文件路径
	wirePath := filepath.Join(g.ProjectPath, "internal/wire.go")
	
	// 检查文件是否存在
	if _, err := os.Stat(wirePath); os.IsNotExist(err) {
		return fmt.Errorf("wire.go文件不存在")
	}
	
	// 读取文件内容
	content, err := os.ReadFile(wirePath)
	if err != nil {
		return err
	}
	
	// 文件内容
	fileContent := string(content)
	
	// 根据组件类型添加不同的内容
	var newImport, newProvider string
	
	if componentType == "service" {
		// 服务导入
		importStr := fmt.Sprintf(`	"{{.ModulePath}}/internal/service"`, g.ModulePath)
		if !strings.Contains(fileContent, importStr) {
			newImport = importStr
		}
		
		// 服务提供者
		providerStr := fmt.Sprintf(`	service.New%s,`, name)
		if !strings.Contains(fileContent, providerStr) {
			newProvider = providerStr
		}
	} else if componentType == "controller" {
		// 控制器导入
		importStr := fmt.Sprintf(`	"{{.ModulePath}}/internal/api"
	"{{.ModulePath}}/internal/biz"
	"{{.ModulePath}}/internal/dao"`, g.ModulePath)
		if !strings.Contains(fileContent, importStr) {
			newImport = importStr
		}
		
		// 控制器提供者
		providerStr := fmt.Sprintf(`	dao.New%s,
	biz.New%s,
	api.New%s,`, name, name, name)
		if !strings.Contains(fileContent, providerStr) {
			newProvider = providerStr
		}
	}
	
	// 更新导入
	if newImport != "" {
		// 查找导入块的结束位置
		importEndIndex := strings.Index(fileContent, ")")
		if importEndIndex != -1 {
			fileContent = fileContent[:importEndIndex] + newImport + "\n" + fileContent[importEndIndex:]
		}
	}
	
	// 更新提供者
	if newProvider != "" {
		// 查找提供者块的结束位置
		setIndex := strings.Index(fileContent, "var Set = wire.NewSet(")
		if setIndex != -1 {
			setEndIndex := strings.Index(fileContent[setIndex:], ")")
			if setEndIndex != -1 {
				setEndIndex += setIndex
				fileContent = fileContent[:setEndIndex] + "\n" + newProvider + "\n" + fileContent[setEndIndex:]
			}
		}
	}
	
	// 写入文件
	return os.WriteFile(wirePath, []byte(fileContent), 0644)
}

// updateRouterFile 更新router.go文件
func (g *ComponentGenerator) updateRouterFile(controllerName string) error {
	// router.go文件路径
	routerPath := filepath.Join(g.ProjectPath, "cmd/router.go")
	
	// 检查文件是否存在
	if _, err := os.Stat(routerPath); os.IsNotExist(err) {
		return fmt.Errorf("router.go文件不存在")
	}
	
	// 读取文件内容
	content, err := os.ReadFile(routerPath)
	if err != nil {
		return err
	}
	
	// 文件内容
	fileContent := string(content)
	
	// 添加控制器字段
	routerStructIndex := strings.Index(fileContent, "type DefRouter struct {")
	if routerStructIndex != -1 {
		endBrace := strings.Index(fileContent[routerStructIndex:], "}")
		if endBrace != -1 {
			endBrace += routerStructIndex
			fieldStr := fmt.Sprintf("\n\t%sApi *api.%s", controllerName, controllerName)
			if !strings.Contains(fileContent, fieldStr) {
				fileContent = fileContent[:endBrace] + fieldStr + fileContent[endBrace:]
			}
		}
	}
	
	// 添加路由
	routesIndex := strings.Index(fileContent, "func (r *DefRouter) SetupRoutes(")
	if routesIndex != -1 {
		v1GIndex := strings.Index(fileContent[routesIndex:], "v1G := app.Group(\"/api/v1\")")
		if v1GIndex != -1 {
			v1GIndex += routesIndex
			braceIndex := strings.Index(fileContent[v1GIndex:], "{")
			if braceIndex != -1 {
				braceIndex += v1GIndex
				closeBraceIndex := strings.Index(fileContent[braceIndex:], "}")
				if closeBraceIndex != -1 {
					closeBraceIndex += braceIndex
					
					// 添加路由
					routeStr := fmt.Sprintf(`
		// %s路由
		v1G.Get("/%s/:id", r.%sApi.Get)
		v1G.Post("/%s", r.%sApi.Create)`, 
						strings.ToLower(controllerName),
						strings.ToLower(controllerName),
						controllerName,
						strings.ToLower(controllerName),
						controllerName)
					
					if !strings.Contains(fileContent, routeStr) {
						fileContent = fileContent[:closeBraceIndex] + routeStr + fileContent[closeBraceIndex:]
					}
				}
			}
		}
	}
	
	// 写入文件
	return os.WriteFile(routerPath, []byte(fileContent), 0644)
}
