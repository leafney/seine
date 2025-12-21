/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     grape
 * @Date:        2024-11-15 01:45
 * @Description:
 */

package leveldbx

import (
	"github.com/leafney/rose"
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/seine/pkg/configx"
	"github.com/leafney/seine/pkg/xlogx"
)

type LevelDBSvc struct {
	*rleveldb.LevelDB
}

func NewLevelDBSvc(cfg configx.BadgerDBConfig, log *xlogx.XLogSvc, stop chan struct{}) *LevelDBSvc {

	dbPath := cfg.GetBadgerDBPath()

	if err := rose.DirExistsEnsure(dbPath); err != nil {
		log.Fatalf("[Leveldb] dbPath exist [%v] error [%v]", dbPath, err)
	}

	db, err := rleveldb.NewLevelDB(dbPath)
	if err != nil {
		log.Fatalf("[Leveldb] OpenFile [%v] error [%v]", dbPath, err)
	}

	go func() {
		// 等待停止信号
		<-stop
		if err := db.Close(); err != nil {
			log.Errorf("[Leveldb] Closed error [%v]", err)
		} else {
			log.Infoln("[Leveldb] Exit successful")
		}
	}()

	log.Infoln("[Leveldb] Load successful")

	return &LevelDBSvc{db}
}
