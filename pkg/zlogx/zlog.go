/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     grape
 * @Date:        2024-10-22 14:26
 * @Description:
 */

package zlogx

import (
	"errors"
	"syscall"

	"github.com/leafney/seine/pkg/xlogx"
	rzap "github.com/leafney/rose-zap"
)

type ZLogSvc struct {
	*rzap.Logger
}

// NewZLogSvc 创建 Zap 日志服务
// enable: 是否启用
// level: 日志级别
// caller: 是否显示调用者信息
// log: xlog 实例用于输出启动信息
// stop: 停止信号通道
func NewZLogSvc(enable bool, level string, caller bool, log *xlogx.XLogSvc, stop chan struct{}) *ZLogSvc {
	zc := rzap.NewConfig()

	zc.
		SetEnable(enable).
		SetLevel(level).
		OutMultiFile(true).
		ShowCaller(caller).
		ShowStacktrace(false)

	logger := rzap.NewLogger(zc)

	go func() {
		<-stop
		if err := logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Fatalf("[Zap] Sync error [%v]", err)
		} else {
			log.Info("[Zap] Exit successful")
		}
	}()

	log.Info("[Zap] Load successful")

	return &ZLogSvc{logger}
}
