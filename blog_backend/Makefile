.PHONY: build run clean test lint migrate

# 应用名称
APP_NAME := blog_backend
# 构建输出目录
BUILD_DIR := ./bin
# 主入口
MAIN_PATH := ./cmd/server/main.go

# 默认目标
all: build

# 构建应用
build:
	@echo "正在构建应用..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(APP_NAME)"

# 运行应用
run:
	@echo "正在启动应用..."
	@go run $(MAIN_PATH)

# 清理构建产物
clean:
	@echo "正在清理..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

# 代码检查
lint:
	@echo "正在运行代码检查..."
	@go vet ./...

# 格式化代码
fmt:
	@echo "正在格式化代码..."
	@gofmt -w .

# 整理依赖
tidy:
	@echo "正在整理依赖..."
	@go mod tidy

# 数据库迁移（示例）
migrate:
	@echo "正在执行数据库迁移..."
	@go run $(MAIN_PATH) migrate

# 开发模式运行
dev:
	@echo "正在以开发模式启动..."
	@APP_MODE=debug go run $(MAIN_PATH)

# 生产模式构建
prod:
	@echo "正在以生产模式构建..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "生产模式构建完成"

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build    - 构建应用"
	@echo "  make run      - 运行应用"
	@echo "  make clean    - 清理构建产物"
	@echo "  make test     - 运行测试"
	@echo "  make lint     - 代码检查"
	@echo "  make fmt      - 格式化代码"
	@echo "  make tidy     - 整理依赖"
	@echo "  make dev      - 开发模式运行"
	@echo "  make prod     - 生产模式构建"
	@echo "  make help     - 显示帮助信息"
