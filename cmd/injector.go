/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 10:11
 * @Description:
 */

package cmd

import (
	"github.com/google/wire"
	"github.com/leafney/seine/config"
	"github.com/leafney/seine/internal"
	"github.com/leafney/seine/pkg/gormx"
	"github.com/leafney/seine/pkg/redisx"
	"github.com/leafney/seine/pkg/xlogx"
)

var AppSet = wire.NewSet(
	config.NewConfig,
	gormx.NewGormDBSvc,
	redisx.NewRRedisSvc,
	xlogx.NewXLogSvc,
	internal.Set,
)

type Injector struct {
	L *xlogx.XLogSvc
	R DefRouter
}
