/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 07:35
 * @Description:
 */

package redisx

import (
	"context"
	"github.com/leafney/seine/config"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisSvc struct {
	*redis.Client
}

func NewRedisSvc(cfg *config.Config, stop chan struct{}) *RedisSvc {
	cfgRedis := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfgRedis.Addr,
		Password: cfgRedis.Pwd,
		DB:       cfgRedis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
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

	return &RedisSvc{client}
}
