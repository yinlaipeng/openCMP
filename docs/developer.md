# openCMP 开发者文档

## 项目架构

openCMP 采用分层单体架构：

```
┌──────────────────────────────────────────────────────────────┐
│                    API Layer (Gin)                            │
│                  RESTful API / HTTP Handlers                  │
├──────────────────────────────────────────────────────────────┤
│                    Service Layer                              │
│        业务逻辑层 (资源管理/账户管理/任务编排/权限控制)          │
├──────────────────────────────────────────────────────────────┤
│                  Cloud Provider Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │ Alibaba  │  │ Tencent  │  │   AWS    │  │  Azure   │      │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
├──────────────────────────────────────────────────────────────┤
│                Cloud Interface Layer (标准化接口)              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│  │ ICompute   │  │ INetwork   │  │ IStorage   │              │
│  └────────────┘  └────────────┘  └────────────┘              │
├──────────────────────────────────────────────────────────────┤
│                    Data Layer (Gorm)                          │
│            MySQL/PostgreSQL / 云账户配置 / 资源元数据           │
└──────────────────────────────────────────────────────────────┘
```

## 目录结构

### 后端结构

```
backend/
├── cmd/server/main.go      # 入口文件
├── configs/config.yaml     # 配置文件
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # 业务逻辑
│   ├── model/              # 数据模型
│   ├── middleware/         # 中间件
│   └── utils/              # 内部工具
├── pkg/
│   ├── cloudprovider/      # 云厂商适配器
│   │   ├── interfaces_*.go # 标准接口
│   │   └── adapters/       # 各厂商实现
│   │       └── alibaba/    # 阿里云适配器
│   ├── scheduler/          # 任务调度
│   └── utils/              # 公共工具
├── scripts/                # 初始化脚本
└── tests/                  # 测试文件
```

### 前端结构

```
frontend/
├── src/
│   ├── views/              # 页面组件
│   │   ├── compute/        # 计算资源
│   │   ├── network/        # 网络资源
│   │   ├── storage/        # 存储资源
│   │   ├── database/       # 数据库资源
│   │   ├── iam/            # IAM模块
│   │   └── ...
│   ├── api/                # API客户端
│   ├── components/         # 公共组件
│   ├── layout/             # 布局组件
│   ├── router.ts           # 路由配置
│   ├── utils/              # 工具函数
│   └── types/              # 类型定义
├── public/                 # 静态资源
└── package.json            # 依赖配置
```

## 开发环境

### 后端开发

```bash
# 安装依赖
cd backend
go mod tidy

# 运行开发服务
go run cmd/server/main.go

# 运行测试
go test -v ./...

# 构建
go build -o opencmp cmd/server/main.go
```

### 前端开发

```bash
# 安装依赖
cd frontend
npm install

# 运行开发服务
npm run dev

# 构建
npm run build
```

## 核心设计模式

### 云厂商适配器模式

通过标准接口（ICompute/INetwork/IStorage）统一不同云厂商的操作：

```go
// pkg/cloudprovider/interfaces_compute.go
type ICompute interface {
    ListVMs(ctx context.Context, accountID uint, opts ListVMOptions) ([]CloudVM, error)
    CreateVM(ctx context.Context, accountID uint, opts CreateVMOptions) (*CloudVM, error)
    DeleteVM(ctx context.Context, accountID uint, vmID string) error
    StartVM(ctx context.Context, accountID uint, vmID string) error
    StopVM(ctx context.Context, accountID uint, vmID string) error
}
```

### 适配器注册机制

每个适配器在初始化时自动注册：

```go
// pkg/cloudprovider/adapters/alibaba/provider.go
func init() {
    registry.Register("alibaba", &AlibabaProvider{})
}
```

### 数据来源规则

- **列表查询**: 从本地数据库获取（CloudVM/CloudVPC等）
- **详情查询**: 实时调用云厂商API
- **创建/操作**: 直接调用云厂商SDK，成功后写入本地数据库
- **状态变更**: 记录到resource_state_logs表

## 添加新的云厂商适配器

### 1. 创建适配器目录

```bash
mkdir -p pkg/cloudprovider/adapters/newcloud
```

### 2. 实现标准接口

```go
// pkg/cloudprovider/adapters/newcloud/provider.go
package newcloud

import "github.com/opencmp/opencmp/pkg/cloudprovider"

type NewCloudProvider struct{}

func init() {
    cloudprovider.Register("newcloud", &NewCloudProvider{})
}

func (p *NewCloudProvider) Name() string {
    return "newcloud"
}

// 实现ICompute接口
func (p *NewCloudProvider) ListVMs(...) ([]cloudprovider.CloudVM, error) {
    // 调用新云厂商SDK
}

// 实现其他接口...
```

### 3. 引入适配器

在 `cmd/server/main.go` 中添加引入：

```go
import _ "github.com/opencmp/opencmp/pkg/cloudprovider/adapters/newcloud"
```

## 添加新的API端点

### 1. 创建Handler

```go
// internal/handler/new_feature.go
package handler

func NewNewFeatureHandler(db *gorm.DB, logger *zap.Logger) *NewFeatureHandler {
    return &NewFeatureHandler{db: db, logger: logger}
}

func (h *NewFeatureHandler) List(c *gin.Context) {
    // 实现列表逻辑
}

func (h *NewFeatureHandler) Create(c *gin.Context) {
    // 实现创建逻辑
}
```

### 2. 创建Service

```go
// internal/service/new_feature.go
package service

type NewFeatureService struct {
    db *gorm.DB
}

func (s *NewFeatureService) List(ctx context.Context) ([]model.NewFeature, error) {
    // 实现业务逻辑
}
```

### 3. 注册路由

在 `cmd/server/main.go` 中添加路由：

```go
newFeatureHandler := handler.NewNewFeatureHandler(db, logger)
newFeatureGroup := v1.Group("/new-features")
{
    newFeatureGroup.GET("", newFeatureHandler.List)
    newFeatureGroup.POST("", newFeatureHandler.Create)
}
```

## 权限系统

### 权限格式

权限使用格式：`<module>:<resource>:<action>`

### 权限中间件

```go
// internal/middleware/auth.go
func PermissionMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取用户权限
        // 检查是否有对应权限
        // 无权限返回403
    }
}
```

### 项目隔离中间件

```go
func ProjectIsolationMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 根据项目上下文过滤资源
    }
}
```

## 数据模型

### 使用GORM

```go
// internal/model/new_model.go
package model

type NewModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
    
    Name      string         `gorm:"size:100;not null" json:"name"`
    Status    string         `gorm:"size:20;default:'active'" json:"status"`
}

func (NewModel) TableName() string {
    return "new_models"
}
```

### 自动迁移

在 `main.go` 中添加：

```go
db.AutoMigrate(&model.NewModel{})
```

## 测试

### 单元测试

```go
// internal/service/new_feature_test.go
package service

import "testing"

func TestNewFeatureService_List(t *testing.T) {
    // 设置测试环境
    // 执行测试
    // 验证结果
}
```

### 使用Mock数据库

```go
import "github.com/DATA-DOG/go-sqlmock"

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        panic(err)
    }
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: db,
    }), &gorm.Config{})
    return gormDB, mock
}
```

## 前端开发规范

### 页面结构

```vue
<template>
  <div class="xxx-container">
    <div class="page-header">
      <h2>页面标题</h2>
    </div>
    <div class="filter-card">
      <!-- 筛选条件 -->
    </div>
    <el-table :data="tableData" row-key="id">
      <!-- 表格列 -->
    </el-table>
    <div class="pagination">
      <!-- 分页 -->
    </div>
  </div>
</template>

<style scoped>
.xxx-container {
  padding: 20px;
}
</style>
```

### API调用

```typescript
// src/api/new_feature.ts
import request from '@/utils/request'

export function getList(params: any) {
  return request.get('/api/v1/new-features', { params })
}

export function create(data: any) {
  return request.post('/api/v1/new-features', data)
}
```

### 使用CloudAccountSelector

选择云账户时使用统一组件：

```vue
<CloudAccountSelector v-model="cloudAccountId" />
```

## 部署

### Docker部署

```bash
# 使用docker-compose
docker-compose up -d
```

### 环境变量

关键配置通过环境变量传递：

```bash
DB_HOST=mysql
DB_PORT=3306
DB_USER=opencmp
DB_PASSWORD=your-password
JWT_SECRET=your-secret-key
```

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码（遵循conventional commits）
4. 运行测试确保通过
5. 提交Pull Request