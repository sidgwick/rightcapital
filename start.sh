#!/bin/bash

echo "Starting Notification System..."

# 确保依赖已安装
echo "Installing dependencies..."
go mod tidy

# 构建项目
echo "Building project..."
go build -o bin/notification-system cmd/server/main.go

# 运行项目
echo "Running server..."
./bin/notification-system
