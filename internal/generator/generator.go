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

	// 先尝试从模板目录读取文件
	var tmplContent string
	if _, err := os.Stat(templatePath); err == nil {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("读取模板文件失败: %v", err)
		}
		tmplContent = string(content)
	} else {
		// 如果模板文件不存在，使用内置模板
		tmplContent = GetTemplateContent(templateName)
		if tmplContent == "" {
			return fmt.Errorf("模板文件 %s 不存在", templateName)
		}
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
