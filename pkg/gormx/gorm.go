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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type DBService struct {
	*gorm.DB
}

func NewDBService(cfg *config.Config, stop chan struct{}) *DBService {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
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
			log.Printf("[Sqlite] Closed error [%v]", err)
		} else {
			log.Println("[Sqlite] Exit successful")
		}
	}()

	log.Println("[Sqlite] Load successful")

	return &DBService{
		DB: db,
	}

}

//func (ds *DBService) Stop() error {
//	sqlDb, err := ds.DB.DB()
//	if err != nil {
//		return err
//	}
//	return sqlDb.Close()
//}
