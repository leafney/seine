# Seine - Go Web项目基础框架

Seine是一个用于快速创建Go Web项目的基础框架生成工具，它提供了一套完整的项目结构和基础组件，帮助开发者快速搭建Go Web应用。

## 功能特点

- 项目结构生成：自动创建标准的Go Web项目目录结构
- 组件添加：支持添加服务、控制器、数据库和路由等组件
- 依赖注入：基于Wire的依赖注入系统
- 配置管理：基于Koanf的配置管理
- Web服务：基于Fiber的Web服务
- 数据库支持：支持MySQL、SQLite、PostgreSQL和MongoDB

## 安装

```bash
go install github.com/leafney/seine@latest
```

## 使用方法

### 创建新项目

```bash
seine create myproject
```

这将创建一个名为`myproject`的新项目，包含完整的目录结构和基础代码。

### 添加服务组件

```bash
cd myproject
seine add service cron
```

这将在`internal/service`目录下创建一个名为`cron.go`的服务组件，并自动更新依赖注入配置。

### 添加控制器组件

```bash
seine add controller user
```

这将创建以下文件：
- `internal/api/user.go`：API处理器
- `internal/biz/user.go`：业务逻辑
- `internal/dao/user.go`：数据访问对象

并自动更新依赖注入配置。

### 添加数据库支持

```bash
seine add db mysql     # 添加MySQL支持
seine add db sqlite    # 添加SQLite支持
seine add db postgresql # 添加PostgreSQL支持
seine add db mongodb   # 添加MongoDB支持
```

这将在`pkg`目录下创建相应的数据库支持组件，并自动更新配置文件和依赖注入配置。

### 添加API路由

```bash
seine add route product
```

这将在`pkg/router`目录下创建一个名为`product.go`的路由处理器，包含基本的CRUD操作，并自动注册到路由系统中。

## 项目结构

```
myproject/
├── cmd/                # 命令行相关代码
│   ├── injector.go     # 依赖注入
│   ├── router.go       # 路由设置
│   ├── start.go        # 服务启动
│   └── wire.go         # Wire依赖注入
├── config/             # 配置相关代码
│   ├── cache/          # 缓存配置
│   ├── vars/           # 常量定义
│   ├── config.go       # 配置管理
│   └── config.toml.default # 默认配置模板
├── internal/           # 内部代码
│   ├── api/            # API处理器
│   ├── biz/            # 业务逻辑
│   ├── dao/            # 数据访问
│   ├── model/          # 数据模型
│   ├── service/        # 服务组件
│   ├── vmodel/         # 视图模型
│   └── wire.go         # 内部依赖注入
├── pkg/                # 公共包
│   ├── errx/           # 错误处理
│   ├── gormx/          # Gorm数据库支持
│   ├── middlewarex/    # 中间件
│   ├── mongox/         # MongoDB支持
│   ├── response/       # 响应处理
│   ├── router/         # 路由管理
│   ├── utils/          # 工具函数
│   └── xlogx/          # 日志处理
├── web/                # Web资源
├── data/               # 数据文件
├── logs/               # 日志文件
├── main.go             # 主入口
├── go.mod              # 依赖管理
├── Makefile            # 构建脚本
└── README.md           # 说明文档
```

## 自定义扩展

### 添加新的pkg组件

如果需要添加新的pkg组件，可以按照以下步骤操作：

1. 在`pkg`目录下创建新的组件目录，如`pkg/cachex`
2. 实现组件功能
3. 在`cmd/injector.go`中添加组件导入和依赖注入
4. 在`internal/wire.go`中添加组件提供者

例如，添加一个缓存组件：

```go
// pkg/cachex/cache.go
package cachex

import (
    "time"
    "{{.ModulePath}}/config"
)

type CacheSvc struct {
    // 缓存实现
}

func NewCacheSvc(cfg *config.Config) *CacheSvc {
    return &CacheSvc{}
}
```

然后更新依赖注入配置：

```go
// cmd/injector.go
import (
    "{{.ModulePath}}/pkg/cachex"
)

// 添加到AppSet
var AppSet = wire.NewSet(
    // ...
    cachex.NewCacheSvc,
)

// 添加到Injector结构体
type Injector struct {
    // ...
    Cache *cachex.CacheSvc
}
```

### 自定义路由

Seine框架使用集中式的路由管理，所有路由都定义在`pkg/router`目录下。如果需要添加新的路由，可以使用`seine add route`命令，或者手动创建路由文件。

路由文件的基本结构如下：

```go
package router

import (
    "github.com/gofiber/fiber/v2"
    "{{.ModulePath}}/pkg/response"
)

// XxxRoutes 注册Xxx相关路由
func XxxRoutes(router fiber.Router) {
    r := router.Group("/xxx")
    
    r.Get("/", getXxxList)
    r.Get("/:id", getXxxById)
    // 其他路由...
}

// 路由处理函数...
```

然后在`pkg/router/router.go`中注册路由：

```go
func RegisterRoutes(app *fiber.App) {
    api := app.Group("/api/v1")
    
    // 注册路由
    XxxRoutes(api)
}
```

## 许可证

MIT
