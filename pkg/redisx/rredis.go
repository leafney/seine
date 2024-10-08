/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 07:36
 * @Description:
 */

package redisx

import (
	rredis "github.com/leafney/rose-redis"
	"github.com/leafney/seine/config"
	"log"
)

type RRedisSvc struct {
	*rredis.Redis
}

func NewRRedisSvc(cfg *config.Config, stop chan struct{}) *RRedisSvc {
	cfgRedis := cfg.Redis
	client, err := rredis.NewRedis(cfgRedis.Addr, &rredis.Option{
		Pass: cfgRedis.Pwd,
		Db:   cfgRedis.DB,
		Type: rredis.TypeNode,
	})
	if err != nil {
		log.Fatalf("[RRedis] connect error [%v]", err)
	}

	go func() {
		<-stop // 等待停止信号
		if err := client.Close(); err != nil {
			log.Fatalf("[RRedis] disconnect error [%v]", err)
		} else {
			log.Println("[RRedis] Exit successful")
		}
	}()

	log.Println("[RRedis] Load successful")

	return &RRedisSvc{client}
}
