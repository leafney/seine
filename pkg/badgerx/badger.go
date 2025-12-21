/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     smart-life-assistant
 * @Date:        2025-05-30 10:33
 * @Description:
 */

package badgerx

import (
	"github.com/leafney/rose"
	rbadger "github.com/leafney/rose-badger"
	"github.com/leafney/seine/pkg/configx"
	"github.com/leafney/seine/pkg/xlogx"
)

type BadgerSvc struct {
	*rbadger.BadgerDB
}

func NewBadgerSvc(cfg configx.BadgerDBConfig, log *xlogx.XLogSvc, stop chan struct{}) *BadgerSvc {

	dbPath := cfg.GetBadgerDBPath()
	if err := rose.DirExistsEnsure(dbPath); err != nil {
		log.Fatalf("[BadgerDB] dbPath exist [%v] error [%v]", dbPath, err)
	}

	db, err := rbadger.NewBadgerDB(dbPath)
	if err != nil {
		log.Fatalf("[BadgerDB] OpenFile [%v] error [%v]", dbPath, err)
	}

	go func() {
		<-stop
		if err := db.Close(); err != nil {
			log.Errorf("[BadgerDB] Close error [%v]", err)
		} else {
			log.Infoln("[BadgerDB] Exit successful")
		}
	}()

	log.Infoln("[BadgerDB] Load successful")

	return &BadgerSvc{db}
}
