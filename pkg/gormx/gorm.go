/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     grape
 * @Date:        2024-10-20 15:04
 * @Description:
 */

package gormx

import (
	"github.com/leafney/seine/pkg/xlogx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDBSvc struct {
	*gorm.DB
}

// NewGormDBSvc 创建 GORM 数据库服务
// driver: 数据库驱动（"mysql" 或 "sqlite"）
// dsn: 数据库连接字符串
// debug: 是否开启调试模式
// log: xlog 实例用于输出启动信息
// stop: 停止信号通道
func NewGormDBSvc(driver string, dsn string, debug bool, log *xlogx.XLogSvc, stop chan struct{}) *GormDBSvc {
	var dialector gorm.Dialector

	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		log.Fatalf("[Gorm] Unsupported database driver: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("[Gorm] connect error [%v]", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("[Gorm] Get sql.DB error [%v]", err)
	}

	// 一些默认配置项
	//sqlDb.SetMaxOpenConns()

	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("[Gorm] Ping error [%v]", err)
	}

	go func() {
		<-stop
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("[Gorm] Get sql.DB error [%v]", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Fatalf("[Gorm] Closed error [%v]", err)
		} else {
			log.Infoln("[Gorm] Exit successful")
		}
	}()

	log.Infoln("[Gorm] Load successful")

	if debug {
		// 启用 gorm 调试模式，打印 sql 调试信息
		db = db.Debug()
		log.Infoln("[Gorm] Debug mode enabled")
	} else {
		// 禁用 gorm 调试模式，不打印 sql 调试信息
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	return &GormDBSvc{DB: db}
}
