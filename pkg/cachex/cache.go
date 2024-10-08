/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 18:31
 * @Description:
 */

package cachex

import (
	rcache "github.com/leafney/rose-cache"
	"github.com/leafney/seine/config"
	"log"
)

type CacheSvc struct {
	*rcache.Cache
}

func NewCacheSvc(cfg *config.Config, stop chan struct{}) *CacheSvc {
	cfgCache := cfg.Cache

	cc, err := rcache.NewCache(cfgCache.Minutes)
	if err != nil {
		log.Fatalf("[Cache] NewCache error [%v]", err)
	}

	go func() {
		// 等待停止信号
		<-stop
		cc.Close()
		log.Println("[Cache] Exit successful")
	}()

	log.Println("[Cache] Load successful")

	return &CacheSvc{cc}
}
