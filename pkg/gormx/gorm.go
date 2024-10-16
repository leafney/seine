/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:03
 * @Description:
 */

package gormx

import (
	"github.com/leafney/seine/config"
	"github.com/leafney/seine/pkg/xlogx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormDBSvc struct {
	*gorm.DB
}

func NewGormDBSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *GormDBSvc {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Sqlite] connect error [%v]", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("[Sqlite] Get sql.DB error [%v]", err)
	}

	// 一些默认配置项
	//sqlDb.SetMaxOpenConns()

	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("[Sqlite] Ping error [%v]", err)
	}

	go func() {
		<-stop
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("[Sqlite] Get sql.DB error [%v]", err)
		}
		if err := sqlDb.Close(); err != nil {
			log.Fatalf("[Sqlite] Closed error [%v]", err)
		} else {
			log.Infoln("[Sqlite] Exit successful")
		}
	}()

	log.Infoln("[Sqlite] Load successful")

	return &GormDBSvc{DB: db}
}
