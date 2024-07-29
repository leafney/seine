/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:42
 * @Description:
 */

package core

import (
	"github.com/leafney/rose"
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/seine/global"
	"github.com/leafney/seine/global/vars"
)

func InitLevelDB(stop chan struct{}) {
	dbPath := global.GConfig.LevelDB
	if rose.StrIsEmpty(dbPath) {
		dbPath = vars.DefLEVDBName
	}

	// 保证路径存在
	if err := rose.DEnsurePathExist(dbPath); err != nil {
		global.GXLog.Fatalf("[Leveldb] dbPath exist [%v] error [%v]", dbPath, err)
	}

	db, err := rleveldb.NewLevelDB(dbPath)
	if err != nil {
		global.GXLog.Fatalf("[Leveldb] OpenFile [%v] error [%v]", dbPath, err)
	}

	go func() {
		// 等待停止信号
		<-stop
		if err := db.Close(); err != nil {
			global.GXLog.Errorf("[Leveldb] Closed error [%v]", err)
		} else {
			global.GXLog.Infoln("[Leveldb] Exit successful")
		}
	}()

	global.GLevelDB = db

	global.GXLog.Infoln("[Leveldb] Load successful")
}
