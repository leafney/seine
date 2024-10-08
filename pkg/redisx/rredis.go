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

type RRedisService struct {
	*rredis.Redis
}

func NewRRedisService(cfg *config.Config, stop chan struct{}) *RRedisService {
	cfgRedis := cfg.Redis
	client, err := rredis.NewRedis(cfgRedis.Addr, &rredis.Option{
		Pass: cfgRedis.Pwd,
		Db:   cfgRedis.DB,
		Type: rredis.TypeNode,
	})
	if err != nil {
		log.Fatalf("[Redis] connect error [%v]", err)
	}

	ping := client.Ping()
	if !ping {
		log.Fatalln("[Redis] ping error")
	}

	go func() {
		<-stop // 等待停止信号
		if err := client.Close(); err != nil {
			log.Fatalf("[Redis] disconnect error [%v]", err)
		} else {
			log.Println("[Redis] Exit successful")
		}
	}()

	log.Println("[Redis] Load successful")

	return &RRedisService{
		client,
	}

}
