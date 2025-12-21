# Seine Framework Makefile

# 定义可执行文件的名称
EXECUTABLE = seine
# 定义不同平台的可执行文件名称
WINDOWS = $(EXECUTABLE)_windows_amd64.exe
LINUX = $(EXECUTABLE)_linux_amd64
DARWIN = $(EXECUTABLE)_darwin_amd64
DARWIN_ARM64 = $(EXECUTABLE)_darwin_arm64
LINUX_ARM64 = $(EXECUTABLE)_linux_arm64

# 获取当前版本为最近一次提交的短哈希值
DEFAULT_VERSION = $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
# 支持用户自定义版本，如果未提供则使用默认版本
VERSION ?= $(DEFAULT_VERSION)

GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")

# Go 编译参数
LDFLAGS = -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.GitBranch=$(GIT_BRANCH)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(BUILD_TIME)'

# 默认目标
.DEFAULT_GOAL := help

# ========== 编译目标 ==========

# 编译所有平台
all: windows linux darwin darwin-arm64 linux-arm64

# Windows 平台
windows: $(WINDOWS)
$(WINDOWS):
	@echo "Building for Windows AMD64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/$(WINDOWS) -ldflags="$(LDFLAGS)" ./main.go

# Linux 平台
linux: $(LINUX)
$(LINUX):
	@echo "Building for Linux AMD64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/$(LINUX) -ldflags="$(LDFLAGS)" ./main.go

# macOS 平台 (Intel)
darwin: $(DARWIN)
$(DARWIN):
	@echo "Building for macOS AMD64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/$(DARWIN) -ldflags="$(LDFLAGS)" ./main.go

# macOS 平台 (ARM64)
darwin-arm64: $(DARWIN_ARM64)
$(DARWIN_ARM64):
	@echo "Building for macOS ARM64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/$(DARWIN_ARM64) -ldflags="$(LDFLAGS)" ./main.go

# Linux ARM64 平台
linux-arm64: $(LINUX_ARM64)
$(LINUX_ARM64):
	@echo "Building for Linux ARM64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/$(LINUX_ARM64) -ldflags="$(LDFLAGS)" ./main.go

# 本地编译
build:
	@echo "Building for local platform..."
	go build -o ./bin/$(EXECUTABLE) -ldflags="$(LDFLAGS)" ./main.go

# ========== 开发工具命令 ==========

# Wire 生成
wire:
	@echo "Running Wire..."
	cd wire && wire

# 运行
run: build
	@echo "Running $(EXECUTABLE)..."
	./bin/$(EXECUTABLE) -f data/config.yaml

# 查看版本
version: build
	@./bin/$(EXECUTABLE) -v

# 初始化项目（创建必要的目录）
init:
	@echo "Initializing project directories..."
	mkdir -p bin data logs data/.cache
	@echo "Done."

# ========== 清理目标 ==========

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf ./bin/*
	rm -f wire/wire_gen.go

# 深度清理（包括日志和缓存）
clean-all: clean
	@echo "Cleaning logs and cache..."
	rm -rf ./logs/*
	rm -rf ./data/.cache/*

# ========== 依赖管理 ==========

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 更新依赖
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# ========== 测试 ==========

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 运行测试（带覆盖率）
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ========== 帮助信息 ==========

help:
	@echo "Seine Framework - Makefile Commands"
	@echo ""
	@echo "编译命令:"
	@echo "  make build          - 编译当前平台的可执行文件"
	@echo "  make all            - 编译所有平台的可执行文件"
	@echo "  make windows        - 编译 Windows 平台"
	@echo "  make linux          - 编译 Linux AMD64 平台"
	@echo "  make darwin         - 编译 macOS Intel 平台"
	@echo "  make darwin-arm64   - 编译 macOS ARM64 平台"
	@echo "  make linux-arm64    - 编译 Linux ARM64 平台"
	@echo ""
	@echo "开发命令:"
	@echo "  make wire           - 运行 Wire 生成依赖注入代码"
	@echo "  make run            - 编译并运行项目"
	@echo "  make version        - 显示版本信息"
	@echo "  make init           - 初始化项目目录"
	@echo ""
	@echo "清理命令:"
	@echo "  make clean          - 清理编译文件"
	@echo "  make clean-all      - 深度清理（包括日志和缓存）"
	@echo ""
	@echo "依赖管理:"
	@echo "  make deps           - 安装依赖"
	@echo "  make update-deps    - 更新依赖"
	@echo ""
	@echo "测试命令:"
	@echo "  make test           - 运行测试"
	@echo "  make test-coverage  - 运行测试并生成覆盖率报告"
	@echo ""
	@echo "自定义版本:"
	@echo "  make VERSION=v1.0.0 build - 使用自定义版本构建"

.PHONY: all windows linux darwin darwin-arm64 linux-arm64 build \
        wire run version init clean clean-all deps update-deps \
        test test-coverage help
