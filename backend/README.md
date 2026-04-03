# openCMP - 开放多云管理平台

openCMP 是一个基于 Go + Gin + Gorm 的单体架构多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

## 架构概览

```
┌──────────────────────────────────────────────────────────────────────┐
│                         API Layer (Gin)                               │
│                    RESTful API / HTTP Handlers                        │
├──────────────────────────────────────────────────────────────────────┤
│                        Service Layer                                  │
│           业务逻辑层 (资源管理/账户管理/任务编排/权限控制)              │
├──────────────────────────────────────────────────────────────────────┤
│                     Cloud Provider Layer                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ Alibaba  │  │ Tencent  │  │   AWS    │  │  Azure   │             │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │             │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘             │
├──────────────────────────────────────────────────────────────────────┤
│                  Cloud Interface Layer (分层标准接口)                   │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐     │
│  │   计算资源  │  │   网络资源  │  │   存储资源  │  │   数据服务  │     │
│  │ ICompute   │  │ INetwork   │  │ IStorage   │  │ IDatabase  │     │
│  └────────────┘  └────────────┘  └────────────┘  └────────────┘     │
├──────────────────────────────────────────────────────────────────────┤
│                      Data Layer (Gorm)                                │
│              MySQL/PostgreSQL / 云账户配置 / 资源元数据存储            │
└──────────────────────────────────────────────────────────────────────┘
```

## IAM 模块功能

我们已经完成了IAM（身份与访问管理）模块的开发，包括：

### 1. 核心实体
- **域 (Domain)** - 租户隔离顶层单位
- **项目 (Project)** - 域内资源分组
- **用户 (User)** - 系统用户管理
- **组 (Group)** - 用户批量管理
- **角色 (Role)** - 权限集合模板
- **权限 (Permission)** - 最小授权单元
- **策略 (Policy)** - 权限集合和条件控制

### 2. 权限管理
- 基于RBAC模型的权限控制
- 支持用户-角色、组-角色、角色-权限关联
- 支持策略(Policy)和策略语句(Statement)
- 细粒度的资源和操作权限控制

### 3. 认证功能
- JWT Token认证
- 支持多种认证源（本地、LDAP、OIDC等）
- 密码安全存储（bcrypt）

## 快速开始

### 环境要求
- Go 1.19+
- MySQL 5.7+ 或 PostgreSQL 9.6+

### 安装步骤

1. 克隆项目
```bash
git clone https://github.com/opencmp/opencmp.git
cd opencmp
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库
修改 `configs/config.yaml` 中的数据库连接信息

4. 运行服务
```bash
go run cmd/server/main.go
```

服务将启动在 `http://localhost:8080`

## API 文档

API 文档可以通过 Swagger 访问（如果已集成）：
- `/swagger/index.html` - Swagger UI
- `/api/v1/swagger/doc.json` - API 规范

## 支持的云厂商

- 阿里云 (Alibaba Cloud)
- 腾讯云 (Tencent Cloud)
- AWS (Amazon Web Services)
- Azure (Microsoft Azure)

## 项目结构

```
opencmp/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                 # HTTP Handler (Gin)
│   ├── service/                 # 业务逻辑层
│   ├── model/                   # 数据模型 (Gorm)
│   ├── middleware/              # Gin 中间件
│   └── migration/               # 数据库迁移
├── pkg/
│   └── cloudprovider/           # 云适配器层
├── configs/                     # 配置文件
├── docs/                        # 文档
└── scripts/                     # 脚本
```

## 开发

### 添加新的云厂商适配器

1. 创建适配器目录：`pkg/cloudprovider/adapters/newprovider/`
2. 实现云提供商接口
3. 在 `init()` 函数中注册适配器

### 添加新的资源类型

1. 在 `pkg/cloudprovider/interfaces_*.go` 中定义接口
2. 在各个适配器中实现接口
3. 在后端服务层添加业务逻辑
4. 在API层添加路由和处理器

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进 openCMP。

## 许可证

MIT License