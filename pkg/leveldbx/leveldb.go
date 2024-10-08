/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 18:40
 * @Description:
 */

package leveldbx

import (
	"github.com/leafney/rose"
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/seine/config"
	"log"
)

type LevelDBSvc struct {
	*rleveldb.LevelDB
}

func NewLevelDBSvc(cfg *config.Config, stop chan struct{}) *LevelDBSvc {

	cfgLD := cfg.LevelDB
	dbPath := cfgLD.Path

	// 保证路径存在
	if err := rose.DEnsurePathExist(dbPath); err != nil {
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
			log.Fatalf("[Leveldb] Closed error [%v]", err)
		} else {
			log.Println("[Leveldb] Exit successful")
		}
	}()

	log.Println("[Leveldb] Load successful")
	return &LevelDBSvc{db}
}
