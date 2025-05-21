package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// updateConfigForDB 更新配置文件以支持数据库
func (g *ComponentGenerator) updateConfigForDB(dbType string) error {
	// 配置文件路径
	configFilePath := filepath.Join(g.ProjectPath, "config/config.go")
	configDefaultPath := filepath.Join(g.ProjectPath, "config/config.toml.default")

	// 检查文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("config.go文件不存在")
	}

	if _, err := os.Stat(configDefaultPath); os.IsNotExist(err) {
		return fmt.Errorf("config.toml.default文件不存在")
	}

	// 读取配置文件内容
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	defaultConfigContent, err := os.ReadFile(configDefaultPath)
	if err != nil {
		return err
	}

	// 文件内容
	configFileContent := string(configContent)
	defaultConfigFileContent := string(defaultConfigContent)

	// 根据数据库类型添加不同的配置
	var configStructAdd, configTypeAdd, configDefaultAdd string

	switch dbType {
	case "mysql":
		// 检查是否已存在MySQL配置
		if strings.Contains(configFileContent, "MySQL") {
			return nil
		}

		configStructAdd = "\tMySQL   MySQL\n"
		configTypeAdd = `
type MySQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
	Debug    bool ` + "`default:\"true\"`" + `
}`

		configDefaultAdd = `
# MySQL配置
[MySQL]
Host = "localhost"
Port = "3306"
User = "root"
Password = "password"
DBName = "test"
DSN = ""
Debug = true`

	case "sqlite":
		// 检查是否已存在SQLite配置
		if strings.Contains(configFileContent, "SQLite") {
			return nil
		}

		configStructAdd = "\tSQLite  SQLite\n"
		configTypeAdd = `
type SQLite struct {
	DSN   string ` + "`default:\"data/app.db\"`" + `
	Debug bool   ` + "`default:\"true\"`" + `
}`

		configDefaultAdd = `
# SQLite配置
[SQLite]
DSN = "data/app.db"
Debug = true`

	case "postgresql":
		// 检查是否已存在PostgreSQL配置
		if strings.Contains(configFileContent, "PostgreSQL") {
			return nil
		}

		configStructAdd = "\tPostgreSQL PostgreSQL\n"
		configTypeAdd = `
type PostgreSQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
	Debug    bool ` + "`default:\"true\"`" + `
}`

		configDefaultAdd = `
# PostgreSQL配置
[PostgreSQL]
Host = "localhost"
Port = "5432"
User = "postgres"
Password = "password"
DBName = "test"
DSN = ""
Debug = true`

	case "mongodb":
		// 检查是否已存在MongoDB配置
		if strings.Contains(configFileContent, "MongoDB") {
			return nil
		}

		configStructAdd = "\tMongoDB MongoDB\n"
		configTypeAdd = `
type MongoDB struct {
	URI      string
	Database string
}`

		configDefaultAdd = `
# MongoDB配置
[MongoDB]
URI = "mongodb://localhost:27017"
Database = "test"`
	}

	// 更新Config结构体
	configStructIndex := strings.Index(configFileContent, "type Config struct {")
	if configStructIndex != -1 {
		endBrace := strings.Index(configFileContent[configStructIndex:], "}")
		if endBrace != -1 {
			endBrace += configStructIndex
			configFileContent = configFileContent[:endBrace] + configStructAdd + configFileContent[endBrace:]
		}
	}

	// 添加数据库配置类型
	typeIndex := strings.Index(configFileContent, "type (")
	if typeIndex != -1 {
		endParen := strings.Index(configFileContent[typeIndex:], ")")
		if endParen != -1 {
			endParen += typeIndex
			configFileContent = configFileContent[:endParen] + configTypeAdd + "\n" + configFileContent[endParen:]
		}
	}

	// 更新配置文件
	if err := os.WriteFile(configFilePath, []byte(configFileContent), 0644); err != nil {
		return err
	}

	// 更新默认配置文件
	defaultConfigFileContent += configDefaultAdd
	if err := os.WriteFile(configDefaultPath, []byte(defaultConfigFileContent), 0644); err != nil {
		return err
	}

	return nil
}

// updateWireFileForDB 更新wire.go文件以支持数据库
func (g *ComponentGenerator) updateWireFileForDB(dbType string) error {
	// wire.go文件路径
	internalWirePath := filepath.Join(g.ProjectPath, "internal/wire.go")
	cmdWirePath := filepath.Join(g.ProjectPath, "cmd/wire.go")
	cmdInjectorPath := filepath.Join(g.ProjectPath, "cmd/injector.go")

	// 检查文件是否存在
	if _, err := os.Stat(internalWirePath); os.IsNotExist(err) {
		return fmt.Errorf("internal/wire.go文件不存在")
	}

	if _, err := os.Stat(cmdWirePath); os.IsNotExist(err) {
		return fmt.Errorf("cmd/wire.go文件不存在")
	}

	if _, err := os.Stat(cmdInjectorPath); os.IsNotExist(err) {
		return fmt.Errorf("cmd/injector.go文件不存在")
	}

	// 读取文件内容
	internalWireContent, err := os.ReadFile(internalWirePath)
	if err != nil {
		return err
	}

	cmdInjectorContent, err := os.ReadFile(cmdInjectorPath)
	if err != nil {
		return err
	}

	// 文件内容
	internalWireFileContent := string(internalWireContent)
	cmdInjectorFileContent := string(cmdInjectorContent)

	// 根据数据库类型添加不同的导入和提供者
	var importAdd, providerAdd, injectorFieldAdd, injectorImportAdd string

	switch dbType {
	case "mysql", "sqlite", "postgresql":
		// 检查是否已存在Gorm导入
		if strings.Contains(internalWireFileContent, "gormx") {
			return nil
		}

		importAdd = fmt.Sprintf(`	"%s/pkg/gormx"`, g.ModulePath)
		providerAdd = "	gormx.NewGormDBSvc,"
		injectorFieldAdd = "	DB  *gormx.GormDBSvc\n"
		injectorImportAdd = fmt.Sprintf(`	"%s/pkg/gormx"`, g.ModulePath)

	case "mongodb":
		// 检查是否已存在MongoDB导入
		if strings.Contains(internalWireFileContent, "mongox") {
			return nil
		}

		importAdd = fmt.Sprintf(`	"%s/pkg/mongox"`, g.ModulePath)
		providerAdd = "	mongox.NewMongoSvc,"
		injectorFieldAdd = "	DB  *mongox.MongoSvc\n"
		injectorImportAdd = fmt.Sprintf(`	"%s/pkg/mongox"`, g.ModulePath)
	}

	// 更新internal/wire.go文件
	// 添加导入
	importIndex := strings.Index(internalWireFileContent, "import (")
	if importIndex != -1 {
		endImport := strings.Index(internalWireFileContent[importIndex:], ")")
		if endImport != -1 {
			endImport += importIndex
			internalWireFileContent = internalWireFileContent[:endImport] + "\n" + importAdd + internalWireFileContent[endImport:]
		}
	}

	// 添加提供者
	setIndex := strings.Index(internalWireFileContent, "var Set = wire.NewSet(")
	if setIndex != -1 {
		endSet := strings.Index(internalWireFileContent[setIndex:], ")")
		if endSet != -1 {
			endSet += setIndex
			internalWireFileContent = internalWireFileContent[:endSet] + "\n" + providerAdd + internalWireFileContent[endSet:]
		}
	}

	// 更新internal/wire.go文件
	if err := os.WriteFile(internalWirePath, []byte(internalWireFileContent), 0644); err != nil {
		return err
	}

	// 更新cmd/injector.go文件
	// 添加导入
	importIndex = strings.Index(cmdInjectorFileContent, "import (")
	if importIndex != -1 {
		endImport := strings.Index(cmdInjectorFileContent[importIndex:], ")")
		if endImport != -1 {
			endImport += importIndex
			cmdInjectorFileContent = cmdInjectorFileContent[:endImport] + "\n" + injectorImportAdd + cmdInjectorFileContent[endImport:]
		}
	}

	// 添加字段
	structIndex := strings.Index(cmdInjectorFileContent, "type Injector struct {")
	if structIndex != -1 {
		endStruct := strings.Index(cmdInjectorFileContent[structIndex:], "}")
		if endStruct != -1 {
			endStruct += structIndex
			cmdInjectorFileContent = cmdInjectorFileContent[:endStruct] + injectorFieldAdd + cmdInjectorFileContent[endStruct:]
		}
	}

	// 更新cmd/injector.go文件
	if err := os.WriteFile(cmdInjectorPath, []byte(cmdInjectorFileContent), 0644); err != nil {
		return err
	}

	return nil
}
