# Seine Framework 快速开始

## 快速安装

### 1. 克隆项目

```bash
cd /Users/jason/xyz/projects/seine-frame/seine
```

### 2. 安装依赖

```bash
make deps
```

### 3. 生成 Wire 代码

```bash
make wire
```

### 4. 编译项目

```bash
make build
```

---

## 运行服务

### 启动服务器

```bash
./bin/seine -f ./data/config.yaml
```

或使用 Makefile：

```bash
make run
```

服务器将在 `http://0.0.0.0:8080` 启动。

### 查看版本信息

```bash
./bin/seine -v
```

输出示例：
```
Seine Framework
  Version:    c2abf23
  Git Branch: feat/51220
  Git Commit: c2abf23
  Build Time: 2025-12-21 16:28:02
```

---

## 快速测试

### 1. 健康检查

```bash
curl http://localhost:8080/api/v1/ping
```

**响应**：
```json
{"code":0,"msg":"OK","data":{"status":"ok"}}
```

### 2. Hello 接口

```bash
curl http://localhost:8080/api/v1/hello
```

**响应**：
```json
{
  "code":0,
  "msg":"OK",
  "data":{
    "message":"Hello from Seine Framework!",
    "version":"1.0.0"
  }
}
```

---

## 用户管理示例

### 1. 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/user/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure123",
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

**响应**：
```json
{"code":0,"msg":"OK"}
```

### 2. 查询用户列表

```bash
curl -X POST http://localhost:8080/api/v1/user/query \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "pageSize": 10,
    "username": "",
    "email": ""
  }'
```

**响应**：
```json
{
  "code":0,
  "msg":"OK",
  "data":{
    "list":[
      {
        "id":1,
        "username":"john_doe",
        "name":"John Doe",
        "email":"john@example.com",
        "status":1
      }
    ],
    "total":1,
    "page":1,
    "size":10
  }
}
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure123"
  }'
```

**成功响应**：
```json
{
  "code":0,
  "msg":"OK",
  "data":{
    "access_token":"mock-token-john_doe",
    "expires_at":0
  }
}
```

**失败响应**（密码错误）：
```json
{
  "code":-1,
  "msg":"error: code = 4008 desc = 用户名或密码错误"
}
```

### 4. 更新用户信息

```bash
curl -X POST http://localhost:8080/api/v1/user/update \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "John Smith",
    "email": "john.smith@example.com",
    "status": 1
  }'
```

**响应**：
```json
{"code":0,"msg":"OK"}
```

### 5. 删除用户（软删除）

```bash
curl -X DELETE http://localhost:8080/api/v1/user/1
```

**响应**：
```json
{"code":0,"msg":"OK"}
```

---

## 配置说明

### 主配置文件：`data/config.yaml`

```yaml
# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080

# 日志配置
log:
  xdebug: true       # 开启调试模式
  xenable: true      # 启用日志
  xlevel: "debug"    # 日志级别：debug/info/warn/error

# 数据库配置
database:
  driver: "sqlite"   # 数据库驱动：mysql/sqlite
  debug: true        # 打印 SQL 语句

  # SQLite 配置
  sqlite:
    dsn: "./data/seine.db"

  # MySQL 配置（切换 driver 为 mysql 时使用）
  mysql:
    dsn: "root:password@tcp(localhost:3306)/seine?charset=utf8mb4&parseTime=True&loc=Local"

# Redis 配置（预留）
redis:
  enable: false
  addr: "localhost:6379"
  password: ""
  db: 0

# 缓存配置（预留）
cache:
  enable: false
  ttl: 3600

# JWT 配置（预留）
jwt:
  signing_key: "your-secret-key"
  expire_hours: 24

# 加密配置（预留）
crypto:
  aes_key: "your-32-byte-aes-256-key-here"
```

---

## 开发流程

### 1. 添加新的业务模块

以添加"产品管理"为例：

#### Step 1: 创建数据模型

```bash
# 创建文件：internal/model/product.go
```

```go
package model

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string  `gorm:"type:varchar(128);not null"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"type:decimal(10,2);not null"`
    Stock       int     `gorm:"type:int;default:0"`
    Status      int     `gorm:"type:tinyint;default:1"`
}
```

#### Step 2: 创建视图模型

```bash
# 添加到：internal/vmodel/types.go
```

```go
type ProductCreateReq struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type ProductQueryResp struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    Status      int     `json:"status"`
}
```

#### Step 3: 创建 DAL 层

```bash
# 创建文件：internal/dal/product.dal.go
```

```go
package dal

import (
    "context"
    "github.com/leafney/seine/internal/model"
    "gorm.io/gorm"
)

type ProductDal struct {
    db *gorm.DB
}

func NewProductDal(db *gorm.DB) *ProductDal {
    return &ProductDal{db: db}
}

func (d *ProductDal) Create(ctx context.Context, product *model.Product) error {
    return d.db.WithContext(ctx).Create(product).Error
}

// ... 其他方法
```

#### Step 4: 创建 BIZ 层

```bash
# 创建文件：internal/biz/product.biz.go
```

```go
package biz

import (
    "context"
    "github.com/leafney/seine/internal/dal"
    "github.com/leafney/seine/internal/model"
    "github.com/leafney/seine/internal/vmodel"
)

type ProductBiz struct {
    productDal *dal.ProductDal
}

func NewProductBiz(productDal *dal.ProductDal) *ProductBiz {
    return &ProductBiz{productDal: productDal}
}

func (b *ProductBiz) Create(ctx context.Context, req *vmodel.ProductCreateReq) error {
    product := &model.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
    }
    return b.productDal.Create(ctx, product)
}

// ... 其他方法
```

#### Step 5: 创建 API 层

```bash
# 创建文件：internal/api/product.api.go
```

```go
package api

import (
    "github.com/gofiber/fiber/v2"
    "github.com/leafney/seine/internal/biz"
    "github.com/leafney/seine/internal/vmodel"
    "github.com/leafney/seine/response"
)

type ProductHandler struct {
    productBiz *biz.ProductBiz
}

func NewProductHandler(productBiz *biz.ProductBiz) *ProductHandler {
    return &ProductHandler{productBiz: productBiz}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
    var req vmodel.ProductCreateReq
    if err := c.BodyParser(&req); err != nil {
        return response.Failed(c, err)
    }
    err := h.productBiz.Create(c.Context(), &req)
    return response.Succeed(c, nil, err)
}

// ... 其他方法
```

#### Step 6: 添加 Wire Provider

```bash
# 编辑：wire/provider.go
```

```go
// 在 DALProviderSet 添加
var DALProviderSet = wire.NewSet(
    dal.NewUserDal,
    dal.NewProductDal,  // 新增
)

// 在 BIZProviderSet 添加
var BIZProviderSet = wire.NewSet(
    biz.NewUserBiz,
    biz.NewProductBiz,  // 新增
)

// 在 APIProviderSet 添加
var APIProviderSet = wire.NewSet(
    api.NewUserHandler,
    api.NewExampleHandler,
    api.NewProductHandler,  // 新增
)
```

#### Step 7: 更新 App 和路由

```bash
# 编辑：core/app.go
```

```go
type App struct {
    UserHandler    *api.UserHandler
    ExampleHandler *api.ExampleHandler
    ProductHandler *api.ProductHandler  // 新增
}

func NewApp(
    userHandler *api.UserHandler,
    exampleHandler *api.ExampleHandler,
    productHandler *api.ProductHandler,  // 新增
) *App {
    return &App{
        UserHandler:    userHandler,
        ExampleHandler: exampleHandler,
        ProductHandler: productHandler,  // 新增
    }
}
```

```bash
# 编辑：core/router.go
```

```go
func RegisterRoutes(app *fiber.App, deps *App) {
    v1 := app.Group("/api/v1")

    // ... 其他路由

    // 产品管理
    product := v1.Group("/product")
    {
        product.Post("/create", deps.ProductHandler.Create)
        product.Post("/query", deps.ProductHandler.Query)
        product.Post("/update", deps.ProductHandler.Update)
        product.Delete("/:id", deps.ProductHandler.Delete)
    }
}
```

#### Step 8: 添加自动迁移

```bash
# 编辑：wire/provider.go 中的 provideGormDB
```

```go
// 自动迁移数据表
if err := dbSvc.DB.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
    log.Errorf("[Wire] AutoMigrate failed: %v", err)
} else {
    log.Infoln("[Wire] AutoMigrate completed successfully")
}
```

#### Step 9: 重新生成 Wire 代码并编译

```bash
make wire
make build
make run
```

---

## 调试技巧

### 1. 开启详细日志

```yaml
# data/config.yaml
log:
  xdebug: true
  xlevel: "debug"

database:
  debug: true  # 打印所有 SQL 语句
```

### 2. 查看数据库内容

```bash
# SQLite
sqlite3 data/seine.db "SELECT * FROM users;"

# MySQL
mysql -u root -p -e "USE seine; SELECT * FROM users;"
```

### 3. 使用 curl 详细模式

```bash
curl -v http://localhost:8080/api/v1/user/query
```

### 4. 日志输出位置

- 标准输出：控制台
- 文件日志：可在配置中指定

---

## 常用命令

### 编译相关

```bash
make deps          # 安装依赖
make wire          # 生成 Wire 代码
make build         # 本地编译
make all           # 多平台编译
make clean         # 清理编译文件
```

### 运行相关

```bash
make run           # 编译并运行
./bin/seine -v     # 查看版本
./bin/seine -f ./data/config.yaml  # 指定配置文件运行
```

### 数据库相关

```bash
# 删除数据库重新开始
rm data/seine.db
make run  # 自动重建表
```

---

## 故障排查

### 问题：端口被占用

**错误信息**：
```
listen tcp 0.0.0.0:8080: bind: address already in use
```

**解决方法**：
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀掉进程
kill -9 <PID>

# 或修改配置文件端口
# data/config.yaml
server:
  port: 8081
```

### 问题：数据库连接失败

**错误信息**：
```
[Gorm] Failed to connect to database
```

**解决方法**：
1. 检查配置文件中的 DSN
2. 确认数据库服务已启动
3. 验证用户名密码正确

### 问题：Wire 生成失败

**错误信息**：
```
wire: could not import ...
```

**解决方法**：
```bash
# 确保依赖已安装
go mod tidy

# 重新生成
make wire
```

---

## 性能优化建议

### 1. 数据库连接池

```yaml
# data/config.yaml
database:
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

### 2. 启用 Fiber 压缩

```go
// core/server.go
import "github.com/gofiber/fiber/v2/middleware/compress"

fiberApp.Use(compress.New())
```

### 3. 生产环境关闭调试

```yaml
log:
  xdebug: false
  xlevel: "info"

database:
  debug: false
```

---

## 下一步

- 📖 阅读完整的 [架构文档](ARCHITECTURE.md)
- 🔐 集成 JWT 认证
- 📝 编写单元测试
- 🚀 部署到生产环境

---

**最后更新**：2025-12-21
