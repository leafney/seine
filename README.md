# Seine Framework

一个基于 Go 语言的高性能 Web 开发框架，采用三层架构和依赖注入设计模式。

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=flat)](https://gofiber.io)
[![GORM](https://img.shields.io/badge/GORM-v1.31-orange?style=flat)](https://gorm.io)
[![Wire](https://img.shields.io/badge/Wire-v0.7-blue?style=flat)](https://github.com/google/wire)

## ✨ 特性

- 🏗️ **清晰的三层架构**：API → BIZ → DAL 分层设计
- 💉 **依赖注入**：基于 Google Wire 的编译时依赖注入
- 🚀 **高性能**：基于 Fiber 框架，零内存分配路由
- 🗄️ **多数据库支持**：MySQL / SQLite
- 🔐 **安全可靠**：bcrypt 密码加密，软删除机制
- 📦 **丰富的基础库**：整合 13+ 基础工具库
- 🔧 **模块化设计**：pkg 库可独立复用
- 🛠️ **开箱即用**：完整的示例代码和文档

## 🚀 快速开始

### 安装依赖

```bash
make deps
```

### 生成 Wire 代码

```bash
make wire
```

### 编译项目

```bash
make build
```

### 运行服务

```bash
make run
```

服务器将在 `http://0.0.0.0:8080` 启动。

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/api/v1/ping

# Hello 接口
curl http://localhost:8080/api/v1/hello

# 创建用户
curl -X POST http://localhost:8080/api/v1/user/create \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"pass123","name":"Test User","email":"test@example.com"}'

# 查询用户
curl -X POST http://localhost:8080/api/v1/user/query \
  -H "Content-Type: application/json" \
  -d '{"page":1,"pageSize":10}'
```

## 📖 文档

- [快速开始指南](docs/QUICKSTART.md) - 快速上手，5 分钟搭建第一个 API
- [架构文档](docs/ARCHITECTURE.md) - 完整的架构设计和实现说明

## 🏗️ 项目结构

```
seine/
├── bin/                # 编译输出
├── config/             # 配置定义
├── core/               # 核心框架
│   ├── app.go          # 应用容器
│   ├── router.go       # 路由注册
│   └── server.go       # 服务器创建
├── data/               # 配置和数据
│   ├── config.yaml     # 主配置文件
│   └── seine.db        # SQLite 数据库
├── docs/               # 文档
├── internal/           # 业务代码
│   ├── api/            # API 层
│   ├── biz/            # BIZ 层
│   ├── dal/            # DAL 层
│   ├── model/          # 数据模型
│   └── vmodel/         # 视图模型
├── pkg/                # 基础库（13+ 工具库）
├── response/           # 响应封装
├── wire/               # Wire 依赖注入
├── main.go             # 应用入口
└── Makefile            # 编译脚本
```

## 🔧 技术栈

### 核心框架

- **Web 框架**：[Fiber v2](https://gofiber.io) - 基于 fasthttp 的高性能框架
- **ORM**：[GORM](https://gorm.io) - 功能强大的 Go ORM
- **依赖注入**：[Wire](https://github.com/google/wire) - Google 的编译时依赖注入工具
- **配置解析**：[YAML v3](https://github.com/go-yaml/yaml) - YAML 配置文件支持

### 数据库驱动

- **MySQL**：gorm.io/driver/mysql
- **SQLite**：gorm.io/driver/sqlite

### 基础库（pkg/）

| 库名 | 说明 |
|------|------|
| errc | 错误码定义 |
| errx | 错误处理封装 |
| configx | 配置接口定义 |
| xlogx | XLog 日志服务 |
| zlogx | Zap 日志封装 |
| gormx | GORM 封装 |
| badgerx | BadgerDB KV 存储 |
| cachex | 缓存服务 |
| leveldbx | LevelDB 封装 |
| memcx | Memcached 客户端 |
| redisx | Redis 封装 |
| utils | 通用工具函数 |
| versionx | 版本信息管理 |

## 📋 API 列表

### 示例接口

- `GET /api/v1/hello` - 欢迎信息
- `GET /api/v1/ping` - 健康检查

### 用户管理

- `POST /api/v1/user/create` - 创建用户
- `POST /api/v1/user/query` - 查询用户（分页）
- `POST /api/v1/user/update` - 更新用户
- `DELETE /api/v1/user/:id` - 删除用户（软删除）
- `POST /api/v1/user/login` - 用户登录

## 🎯 核心概念

### 三层架构

```
┌─────────────────────────────┐
│  API Layer                  │  ← HTTP 请求处理
├─────────────────────────────┤
│  BIZ Layer                  │  ← 业务逻辑
├─────────────────────────────┤
│  DAL Layer                  │  ← 数据访问
└─────────────────────────────┘
         ↓
    ┌────────┐
    │   DB   │
    └────────┘
```

### 配置接口化

```go
// pkg/configx - 定义接口
type LogConfig interface {
    GetLogLevel() string
}

// config/ - 实现接口
type Config struct {
    Log LogConfig `yaml:"log"`
}

// pkg/xlogx - 使用接口
func NewXLogSvc(debug, enable, level) *XLogSvc
```

### Wire 依赖注入

```go
// 定义依赖关系
wire.Build(
    InfrastructureProviderSet,
    DALProviderSet,
    BIZProviderSet,
    APIProviderSet,
    core.NewApp,
)

// Wire 自动生成代码
app, err := wire.InitializeApp(cfg, stop)
```

## 🛠️ Make 命令

```bash
make deps          # 安装依赖
make wire          # 生成 Wire 代码
make build         # 本地编译
make all           # 多平台编译
make run           # 编译并运行
make clean         # 清理编译文件
make version       # 查看版本
```

## 📦 编译平台

- Windows AMD64
- Linux AMD64
- Linux ARM64
- macOS AMD64
- macOS ARM64

编译后的二进制文件位于 `bin/` 目录。

## ⚙️ 配置说明

主配置文件：`data/config.yaml`

```yaml
server:
  host: "0.0.0.0"
  port: 8080

log:
  xdebug: true
  xenable: true
  xlevel: "debug"

database:
  driver: "sqlite"  # mysql / sqlite
  debug: true
  sqlite:
    dsn: "./data/seine.db"
  mysql:
    dsn: "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True"
```

## 🔐 安全特性

- ✅ **密码加密**：使用 bcrypt 加密存储
- ✅ **软删除**：GORM 软删除机制，数据可恢复
- ✅ **SQL 注入防护**：GORM 参数化查询
- ⏳ **JWT 认证**：待集成 jwtx 库
- ⏳ **RBAC 权限**：计划中

## 📈 性能指标

- 二进制大小：2.2MB
- 启动时间：< 100ms
- API 响应：< 10ms（本地测试）
- 并发处理：基于 Fiber 的高性能路由

## 🧪 测试验证

✅ 所有核心功能已通过测试：

- 配置加载
- 日志输出
- 数据库连接
- Wire 依赖注入
- API 路由
- CRUD 操作
- 密码加密验证
- 软删除
- 分页查询
- 错误处理

## 🗺️ 后续规划

### 短期

- [ ] 集成 JWT 认证（jwtx）
- [ ] 添加请求日志中间件
- [ ] 实现 RBAC 权限系统
- [ ] 单元测试覆盖
- [ ] Swagger API 文档

### 中期

- [ ] 定时任务支持（cronx）
- [ ] HTTP 客户端（reqx）
- [ ] Redis 缓存集成
- [ ] 文件上传功能
- [ ] 性能监控

### 长期

- [ ] 微服务支持（gRPC）
- [ ] 分布式追踪
- [ ] 消息队列（RabbitMQ）
- [ ] Docker 容器化
- [ ] CI/CD 集成

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/amazing`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 保持测试覆盖率

## 📝 版本历史

### v1.0 (2025-12-21)

- ✅ 三层架构实现
- ✅ Wire 依赖注入
- ✅ 13+ 基础库集成
- ✅ 用户 CRUD 示例
- ✅ 多数据库支持
- ✅ 完整文档

## 📄 许可证

MIT License

## 👤 作者

**leafney**

- GitHub: [@leafney](https://github.com/leafney)

## 🙏 致谢

感谢以下项目提供灵感和基础库：

- [grape](https://github.com/leafney/grape) - 提供核心基础库
- [smart-life-assistant](https://github.com/leafney/smart-life-assistant) - 提供配置系统设计
- [three](https://github.com/leafney/three) - 提供架构参考

## 📞 联系方式

- 问题反馈：[GitHub Issues](https://github.com/leafney/seine/issues)
- 邮箱：leafney@example.com

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**