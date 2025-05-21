package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// AddDatabase 添加数据库支持
func (g *ComponentGenerator) AddDatabase(dbType string) error {
	// 根据数据库类型选择不同的实现
	switch dbType {
	case "mysql", "sqlite", "postgresql":
		return g.addGormDatabase(dbType)
	case "mongodb":
		return g.addMongoDatabase()
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbType)
	}
}

// addGormDatabase 添加Gorm数据库支持
func (g *ComponentGenerator) addGormDatabase(dbType string) error {
	// 创建pkg/gormx目录
	gormxDir := filepath.Join(g.ProjectPath, "pkg/gormx")
	if err := os.MkdirAll(gormxDir, 0755); err != nil {
		return err
	}

	// 创建gorm.go文件
	gormFilePath := filepath.Join(gormxDir, "gorm.go")

	// 检查文件是否已存在
	if _, err := os.Stat(gormFilePath); err == nil {
		return fmt.Errorf("gorm.go文件已存在")
	}

	// Gorm模板内容
	var gormTemplate string
	switch dbType {
	case "mysql":
		gormTemplate = `/**
 * @Date:        {{.Date}}
 */

package gormx

import (
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/pkg/xlogx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDBSvc struct {
	*gorm.DB
}

func NewGormDBSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *GormDBSvc {
	dsn := cfg.MySQL.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MySQL.User,
			cfg.MySQL.Password,
			cfg.MySQL.Host,
			cfg.MySQL.DBName,
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Gorm] 连接错误 [%v]", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
	}

	// 设置连接池
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)

	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("[Gorm] Ping错误 [%v]", err)
	}

	go func() {
		<-stop
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Fatalf("[Gorm] 关闭错误 [%v]", err)
		} else {
			log.Infoln("[Gorm] 退出成功")
		}
	}()

	log.Infoln("[Gorm] 加载成功")

	// 打印SQL调试信息
	if cfg.MySQL.Debug {
		db = db.Debug()
	}

	return &GormDBSvc{DB: db}
}`
	case "sqlite":
		gormTemplate = `/**
 * @Date:        {{.Date}}
 */

package gormx

import (
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/pkg/xlogx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormDBSvc struct {
	*gorm.DB
}

func NewGormDBSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *GormDBSvc {
	db, err := gorm.Open(sqlite.Open(cfg.SQLite.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Gorm] 连接错误 [%v]", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
	}

	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("[Gorm] Ping错误 [%v]", err)
	}

	go func() {
		<-stop
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Fatalf("[Gorm] 关闭错误 [%v]", err)
		} else {
			log.Infoln("[Gorm] 退出成功")
		}
	}()

	log.Infoln("[Gorm] 加载成功")

	// 打印SQL调试信息
	if cfg.SQLite.Debug {
		db = db.Debug()
	}

	return &GormDBSvc{DB: db}
}`
	case "postgresql":
		gormTemplate = `/**
 * @Date:        {{.Date}}
 */

package gormx

import (
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/pkg/xlogx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDBSvc struct {
	*gorm.DB
}

func NewGormDBSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *GormDBSvc {
	dsn := cfg.PostgreSQL.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.PostgreSQL.Host,
			cfg.PostgreSQL.User,
			cfg.PostgreSQL.Password,
			cfg.PostgreSQL.DBName,
			cfg.PostgreSQL.Port,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Gorm] 连接错误 [%v]", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
	}

	// 设置连接池
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)

	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("[Gorm] Ping错误 [%v]", err)
	}

	go func() {
		<-stop
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("[Gorm] 获取sql.DB错误 [%v]", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Fatalf("[Gorm] 关闭错误 [%v]", err)
		} else {
			log.Infoln("[Gorm] 退出成功")
		}
	}()

	log.Infoln("[Gorm] 加载成功")

	// 打印SQL调试信息
	if cfg.PostgreSQL.Debug {
		db = db.Debug()
	}

	return &GormDBSvc{DB: db}
}`
	}

	// 创建模板数据
	data := struct {
		Date       string
		ModulePath string
	}{
		Date:       g.Date,
		ModulePath: g.ModulePath,
	}

	// 解析模板
	tmpl, err := template.New("gorm").Parse(gormTemplate)
	if err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(gormFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	// 更新配置文件
	if err := g.updateConfigForDB(dbType); err != nil {
		return err
	}

	// 更新wire.go文件
	if err := g.updateWireFileForDB(dbType); err != nil {
		return err
	}

	return nil
}
