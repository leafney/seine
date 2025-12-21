//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/leafney/seine/config"
	"github.com/leafney/seine/core"
)

// InitializeApp 初始化应用（Wire 会自动生成实现）
func InitializeApp(cfg *config.Config, stop chan struct{}) (*core.App, error) {
	wire.Build(ProviderSet)
	return &core.App{}, nil
}
