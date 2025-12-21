package wire

import (
	"github.com/google/wire"

	"github.com/leafney/seine/config"
	"github.com/leafney/seine/core"
	"github.com/leafney/seine/internal/api"
	"github.com/leafney/seine/internal/biz"
	"github.com/leafney/seine/internal/dal"
	"github.com/leafney/seine/internal/model"
	"github.com/leafney/seine/pkg/gormx"
	"github.com/leafney/seine/pkg/xlogx"

	"gorm.io/gorm"
)

// ProviderSet Wire Provider 集合
var ProviderSet = wire.NewSet(
	// 基础设施层
	InfrastructureProviderSet,

	// DAL 层
	DALProviderSet,

	// BIZ 层
	BIZProviderSet,

	// API 层
	APIProviderSet,

	// 组装 App
	core.NewApp,
)

// InfrastructureProviderSet 基础设施层
var InfrastructureProviderSet = wire.NewSet(
	provideXLogSvc,
	provideGormDB,
)

// DALProviderSet DAL 层
var DALProviderSet = wire.NewSet(
	dal.NewUserDal,
)

// BIZProviderSet BIZ 层
var BIZProviderSet = wire.NewSet(
	biz.NewUserBiz,
)

// APIProviderSet API 层
var APIProviderSet = wire.NewSet(
	api.NewUserHandler,
	api.NewExampleHandler,
)

// ========== Provider 函数 ==========

// provideXLogSvc 提供日志服务
func provideXLogSvc(cfg *config.Config, stop chan struct{}) *xlogx.XLogSvc {
	return xlogx.NewXLogSvc(cfg.Log.XDebug, cfg.Log.XEnable, cfg.Log.XLevel)
}

// provideGormDB 提供数据库连接
func provideGormDB(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *gorm.DB {
	driver := cfg.Database.Driver
	var dsn string

	switch driver {
	case "mysql":
		dsn = cfg.Database.MySQL.DSN
	case "sqlite":
		dsn = cfg.Database.SQLite.DSN
	default:
		log.Fatalf("[Wire] Unsupported database driver: %s", driver)
	}

	dbSvc := gormx.NewGormDBSvc(driver, dsn, cfg.Database.Debug, log, stop)

	// 自动迁移数据表
	if err := dbSvc.DB.AutoMigrate(&model.User{}); err != nil {
		log.Errorf("[Wire] AutoMigrate failed: %v", err)
	} else {
		log.Infoln("[Wire] AutoMigrate completed successfully")
	}

	return dbSvc.DB
}
