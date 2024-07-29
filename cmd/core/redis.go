/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-29 16:56
 * @Description:
 */

package core

import (
	rredis "github.com/leafney/rose-redis"
	"github.com/leafney/seine/global"
)

func InitRedis(stop chan struct{}) {
	redisCfg := global.GConfig.Redis
	client, err := rredis.NewRedis(redisCfg.Addr, &rredis.Option{
		Pass: redisCfg.Pwd,
		Db:   redisCfg.Db,
		Type: rredis.TypeNode,
	})
	if err != nil {
		global.GXLog.Fatalf("[Redis] connect error [%v]", err)
	}

	ping := client.Ping()
	if !ping {
		global.GXLog.Fatalln("[Redis] ping error")
	}

	go func() {
		<-stop // 等待停止信号
		if err := client.Close(); err != nil {
			global.GXLog.Errorf("[Redis] disconnect error [%v]", err)
		} else {
			global.GXLog.Infoln("[Redis] Exit successful")
		}
	}()

	global.GRedis = client

	////
	//que := queue.NewSQueue(client)
	//global.GSQueue = que

	global.GXLog.Infoln("[Redis] Load successful")
}
