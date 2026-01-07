# 项目结构说明

## 目录结构

```
.
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/                    # 私有应用代码
│   ├── api/
│   │   └── handler.go           # HTTP 处理器
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── deliver/
│   │   └── deliver.go           # 投递执行器
│   ├── model/
│   │   └── task.go              # 数据模型
│   ├── platform/
│   │   ├── platform.go          # 平台适配器接口
│   │   └── implementations.go   # 平台实现
│   ├── repository/
│   │   └── repository.go        # 数据存储接口
│   ├── semantic/
│   │   └── semantic.go          # 投递语义处理
│   └── service/
│       └── service.go           # 业务逻辑层
├── pkg/                         # 公共库代码
│   └── utils/
│       └── ptr.go               # 工具函数
├── Makefile                     # 构建脚本
├── go.mod                       # Go 模块定义
├── README.md                    # 项目说明
└── API_EXAMPLES.md              # API 使用示例
```

## 核心模块说明

### 1. Model (internal/model)
定义核心数据模型:
- `NotificationTask`: 通知任务
- `TaskStatus`: 任务状态
- `DeliverySemantic`: 投递语义
- `CreateTaskRequest`: 创建任务请求
- `TaskQueryResponse`: 任务查询响应
- `CallbackRequest`: 回调请求

### 2. Repository (internal/repository)
数据存储抽象，支持不同存储实现:
- 当前实现内存存储
- 可扩展为 MySQL、PostgreSQL 等

### 3. Platform (internal/platform)
平台适配器，实现不同通知平台的投递:
- 接口: `Platform`
- 实现: HTTP、Email、SMS

### 4. Semantic (internal/semantic)
投递语义处理器:
- `AtLeastOnceHandler`: 至少一次，失败后重试
- `AtMostOnceHandler`: 至多一次，失败后不重试
- `ExactlyOnceHandler`: 有且只有一次，通过事务性更新任务状态实现幂等性（DB 持久化）

### 5. Deliver (internal/deliver)
投递执行器:
- 定时扫描待处理任务
- 根据语义进行投递
- 处理重试逻辑
- 发送回调通知

### 6. Service (internal/service)
业务逻辑层:
- 任务创建
- 任务查询
- 任务取消
- 回调处理

### 7. API (internal/api)
HTTP API 层:
- `POST /tasks`: 创建任务
- `GET /tasks/:id`: 查询任务
- `DELETE /tasks/:id`: 取消任务

## 运行说明

### 构建
```bash
make build
```

### 运行
```bash
make run
```

### 测试
```bash
make test
```

## 扩展指南

### 添加新平台
1. 在 `internal/platform/implementations.go` 实现新平台
2. 在 `main.go` 中注册新平台

### 添加新投递语义
1. 在 `internal/semantic/semantic.go` 实现新语义处理器
2. 在 `main.go` 中注册新语义

### 切换存储
实现新的 `Repository` 接口，替换内存存储
