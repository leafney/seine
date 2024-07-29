/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:40
 * @Description:
 */

package core

import (
	"context"
	rcache "github.com/leafney/rose-cache"
	"github.com/leafney/seine/global"
	"github.com/leafney/seine/global/vars"
)

// InitCache 内存级缓存
func InitCache(stop chan struct{}) {
	ctx := context.Background()
	// 设置缓存过期时间
	cc, err := rcache.NewCache(ctx, vars.DefCacheMinutes)
	if err != nil {
		global.GXLog.Fatalf("[Cache] NewCache error [%v]", err)
	}

	go func() {
		// 等待停止信号
		<-stop
		cc.Close()
		global.GXLog.Infoln("[Cache] Exit successful")
	}()

	global.GCache = cc
	global.GXLog.Infoln("[Cache] Load successful")
}
