package core

import (
	"github.com/leafney/seine/internal/api"
)

// App 应用依赖集合
type App struct {
	UserHandler    *api.UserHandler
	ExampleHandler *api.ExampleHandler
}

// NewApp 创建应用依赖集合
func NewApp(
	userHandler *api.UserHandler,
	exampleHandler *api.ExampleHandler,
) *App {
	return &App{
		UserHandler:    userHandler,
		ExampleHandler: exampleHandler,
	}
}
