//go:build wireinject
// +build wireinject

/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:00
 * @Description:
 */

package cmd

import (
	"github.com/google/wire"
)

func InitApp(stop chan struct{}) (*Injector, func(), error) {
	wire.Build(
		AppSet,
		wire.Struct(new(DefRouter), "*"),
		wire.NewSet(wire.Struct(new(Injector), "*")),
	)
	return &Injector{}, nil, nil
}
