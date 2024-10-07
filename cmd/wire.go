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

func InitApp(stop chan struct{}) (*App, func(), error) {
	wire.Build(AppSet)
	return &App{}, nil, nil
}
