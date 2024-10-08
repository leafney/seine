/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 12:24
 * @Description:
 */

package zlogx

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	rzap "github.com/leafney/rose-zap"
	"github.com/leafney/seine/config"
	"syscall"
)

type ZLogSvc struct {
	*rzap.Logger
}

func NewZLogSvc(cfg *config.Config, stop chan struct{}) *ZLogSvc {

	cfgLog := cfg.Log

	zc := rzap.NewConfig()

	zc.
		SetEnable(cfgLog.ZEnable).
		SetLevel(cfgLog.ZLevel).
		OutMultiFile(true).
		ShowCaller(false).
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
