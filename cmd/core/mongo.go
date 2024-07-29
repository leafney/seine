/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:42
 * @Description:
 */

package core

import (
	"context"
	"github.com/leafney/seine/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(stop chan struct{}) {
	cfgMongo := global.GConfig.Mongo

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfgMongo.Addr))
	if err != nil {
		global.GXLog.Fatalf("[Mongo] connect error [%v]", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		global.GXLog.Fatalf("[Mongo] ping error [%v]", err)
	}

	go func() {
		<-stop
		if err := client.Disconnect(ctx); err != nil {
			global.GXLog.Errorf("[Mongo] disconnect error [%v]", err)
		} else {
			global.GXLog.Infoln("[Mongo] Exit successful")
		}
	}()

	global.GMongoDB = client.Database(cfgMongo.Db)

	global.GXLog.Infoln("[Mongo] Load successful")

}
