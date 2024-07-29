/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:45
 * @Description:
 */

package global

import (
	rcache "github.com/leafney/rose-cache"
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/rose-notify/notify"
	rredis "github.com/leafney/rose-redis"
	rzap "github.com/leafney/rose-zap"

	"github.com/leafney/rose/xlog"
	"github.com/leafney/seine/config"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	GConfig config.Config
	GRedis  *rredis.Redis
	//GSQueue  *queue.SQueue
	GMongoDB *mongo.Database
	GCache   *rcache.Cache
	GLevelDB *rleveldb.LevelDB
	GLog     *rzap.Logger
	GXLog    *xlog.Log
	//GMQueue  *queue.MessageQueue
	GNotify notify.Notify
)
