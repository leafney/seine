/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:44
 * @Description:
 */

package core

import (
	"errors"
	rzap "github.com/leafney/rose-zap"
	"github.com/leafney/seine/global"
	"syscall"
)

func InitZLog(stop chan struct{}) {
	// TODO 日志配置
	//cfgLog:=global.GConfig

	cfg := rzap.NewConfig()
	// 自定义设置
	cfg.
		OutMultiFile(true).
		ShowCaller(true).
		ShowStacktrace(false).
		SetLevel("debug")

	logger := rzap.NewLogger(cfg)

	// 停止
	go func() {
		<-stop
		if err := logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			global.GXLog.Errorf("[Zap] Sync error [%v]", err)
		} else {
			global.GXLog.Infoln("[Zap] Exit successful")
		}
	}()

	global.GLog = logger

	global.GXLog.Infoln("[Zap] Load successful")
}
