package redisx

import (
	"context"
	"time"

	"github.com/leafney/seine/pkg/xlogx"
	"github.com/redis/go-redis/v9"
)

type RedisSvc struct {
	*redis.Client
	Addr string
}

func NewRedisSvc(addr, password string, db int, log *xlogx.XLogSvc, stop chan struct{}) (*RedisSvc, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	go func() {
		<-stop
		if err := client.Close(); err != nil {
			log.Errorf("[Redis] Close error: %v", err)
		} else {
			log.Info("[Redis] Exit successful")
		}
	}()

	log.Infof("[Redis] Connected successfully to: %s", addr)
	return &RedisSvc{Client: client, Addr: addr}, nil
}
