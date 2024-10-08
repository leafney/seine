/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 12:12
 * @Description:
 */

package xlogx

import (
	"github.com/leafney/rose/xlog"
	"github.com/leafney/seine/config"
)

type XLogSvc struct {
	*xlog.Log
}

func NewXLogSvc(cfg *config.Config) *XLogSvc {
	cfgLog := cfg.Log

	x := xlog.NewXLog()
	// debug model
	x.SetDebug(cfgLog.XDebug)
	// enable
	x.SetEnable(cfgLog.XEnable)
	// default log level
	x.SetLevelStr(cfgLog.XLevel)

	x.Infoln("[XLog] Load successful")

	return &XLogSvc{x}
}
