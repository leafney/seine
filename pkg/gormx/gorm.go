/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     grape
 * @Date:        2024-10-20 15:04
 * @Description:
 */

package gormx

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/leafney/seine/pkg/xlogx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDBSvc struct {
	*gorm.DB
}

const (
	// gormConnectMaxAttempts 最大连接尝试次数（含首次尝试）；达到次数后仍失败则直接退出进程
	gormConnectMaxAttempts = 10
	// gormConnectInitialBackoff 首次失败后的初始等待时间（之后按指数退避翻倍）
	gormConnectInitialBackoff = 1 * time.Second
	// gormConnectMaxBackoff 指数退避等待的上限（避免等待无限增长）
	gormConnectMaxBackoff = 30 * time.Second
	// gormConnectJitterMax 每次等待叠加的随机抖动上限（避免多实例同时重连导致“惊群”）
	gormConnectJitterMax = 100 * time.Millisecond
)

func openWithRetry(dialector gorm.Dialector, log *xlogx.XLogSvc, stop <-chan struct{}) (*gorm.DB, *sql.DB) {
	backoff := gormConnectInitialBackoff
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var lastErr error
	for attempt := 1; attempt <= gormConnectMaxAttempts; attempt++ {
		select {
		case <-stop:
			log.Fatalln("[Gorm] connect canceled")
		default:
		}

		db, err := gorm.Open(dialector, &gorm.Config{})
		if err == nil {
			sqlDb, err := db.DB()
			if err == nil {
				if err := sqlDb.Ping(); err == nil {
					return db, sqlDb
				} else {
					lastErr = err
					_ = sqlDb.Close()
				}
			} else {
				lastErr = err
			}
		} else {
			lastErr = err
		}

		if attempt == gormConnectMaxAttempts {
			log.Fatalf("[Gorm] connect failed after %d attempts [%v]", gormConnectMaxAttempts, lastErr)
		}

		sleep := backoff
		if gormConnectJitterMax > 0 {
			sleep += time.Duration(rng.Int63n(int64(gormConnectJitterMax) + 1))
		}

		log.Errorf("[Gorm] connect attempt %d/%d failed [%v], retrying in %s", attempt, gormConnectMaxAttempts, lastErr, sleep)

		timer := time.NewTimer(sleep)
		select {
		case <-stop:
			if !timer.Stop() {
				<-timer.C
			}
			log.Fatalln("[Gorm] connect canceled")
		case <-timer.C:
		}

		backoff *= 2
		if backoff > gormConnectMaxBackoff {
			backoff = gormConnectMaxBackoff
		}
	}

	log.Fatalf("[Gorm] connect failed [%v]", lastErr)
	return nil, nil
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

	db, sqlDb := openWithRetry(dialector, log, stop)

	go func() {
		<-stop
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
