# openCMP - 开放多云管理平台

openCMP 是一个基于 Go + Gin + Gorm 的单体架构多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

## 🚀 功能特性

### IAM 模块 (身份与访问管理)
- **域管理**: 租户隔离顶层单位，支持多域管理
- **项目管理**: 域内资源分组，支持项目级别的资源划分
- **用户管理**: 完整的用户生命周期管理，支持密码策略和MFA
- **组管理**: 用户批量管理，支持组内权限分配
- **角色管理**: 基于RBAC模型的角色管理，支持系统角色和自定义角色
- **权限管理**: 细粒度的权限控制，支持资源级权限
- **策略管理**: 灵活的策略引擎，支持条件策略和策略语句
- **认证源管理**: 支持多种认证源（LDAP、本地、SQL等）

### 云资源管理
- **计算资源**: 虚拟机、镜像、密钥对、实例模板管理
- **网络资源**: VPC、子网、安全组、弹性IP、负载均衡管理
- **存储资源**: 云磁盘、快照、对象存储管理
- **数据库服务**: RDS、Redis、MongoDB 等数据库管理
- **中间件服务**: 消息队列、数据分析等中间件管理

### 消息中心
- **站内信**: 系统通知和消息推送
- **消息订阅**: 用户自选关注的事件类型
- **通知渠道**: 支持邮件、企业微信、钉钉、Webhook
- **接受人管理**: 通知目标管理（用户/用户组/角色）
- **机器人管理**: 企业微信机器人、钉钉机器人配置

### 多云管理
- **云账号管理**: 支持多云账号统一管理
- **同步策略**: 按资源类型配置同步范围
- **资源同步规则**: 通过标签自动映射到项目
- **定时同步任务**: 支持cron表达式的定时同步

## 🏗️ 架构设计

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

## 🛠️ 技术栈

### 后端技术栈
- **语言**: Go 1.19+
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: PostgreSQL 12+ / MySQL 8+
- **认证**: JWT
- **密码加密**: bcrypt
- **日志**: Zap
- **配置**: Viper

### 前端技术栈
- **框架**: Vue 3 + TypeScript
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **构建工具**: Vite
- **测试**: Vitest + Vue Test Utils

## 📦 快速开始

### 环境要求
- Go 1.19+
- Node.js 18+
- PostgreSQL 12+ 或 MySQL 8+
- Docker (可选)

### 本地开发

#### 1. 克装依赖
```bash
# 后端
cd backend
go mod tidy

# 前端
cd frontend
npm install
```

#### 2. 配置数据库
```bash
# 创建 PostgreSQL 数据库
CREATE DATABASE opencmp;
CREATE USER opencmp WITH PASSWORD 'opencmp123';
GRANT ALL PRIVILEGES ON DATABASE opencmp TO opencmp;
```

#### 3. 配置文件
```bash
# backend/configs/config.yaml
server:
  port: 8080
  mode: debug

database:
  driver: postgres
  dsn: "host=localhost user=opencmp password=opencmp123 dbname=opencmp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
  max_idle_conns: 10
  max_open_conns: 100

auth:
  jwt_secret: "your-super-secret-jwt-key-change-this-in-production"
  token_expire_hours: 24
```

#### 4. 启动服务
```bash
# 启动后端
cd backend
go run cmd/server/main.go

# 启动前端
cd frontend
npm run dev
```

### Docker 部署（推荐）

#### 快速部署
```bash
# 1. 克隆项目
git clone https://github.com/opencmp/opencmp.git
cd opencmp

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，修改生产环境配置

# 3. 一键启动所有服务
docker-compose up -d

# 4. 查看服务状态
docker-compose ps
```

#### 服务访问
- 前端界面: http://localhost:80
- API服务: http://localhost:8080/api/v1
- 健康检查: http://localhost:8080/health

#### 默认账户
- 用户名: `admin`
- 密码: `admin123`

#### 手动构建镜像
```bash
# 构建后端镜像
cd backend
docker build -t opencmp/backend .

# 构建前端镜像
cd frontend
docker build -t opencmp/frontend .
```

### 本地开发部署

## 🔐 安全特性

### 认证与授权
- JWT Token 认证
- 基于角色的访问控制 (RBAC)
- 细粒度权限控制
- 多因子认证 (MFA) 支持

### 数据安全
- 密码 bcrypt 加密存储
- API 请求签名验证
- SQL 注入防护
- XSS 和 CSRF 防护

## 📚 文档

- [API 文档](./docs/api/iam-api-spec.md)
- [前端开发指南](./docs/frontend/iam-frontend-guide.md)
- [后端开发指南](./docs/backend/iam-backend-guide.md)
- [部署指南](./docs/deployment/guide.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进 openCMP。

### 开发规范
- Go 代码遵循 Go 官方规范
- TypeScript 代码遵循 ESLint + Prettier 规范
- Git 提交信息遵循 conventional commits 规范
- 代码提交前需通过所有测试

## 📄 许可证

MIT License - 详见 [LICENSE](./LICENSE) 文件

## 🆘 支持

如有问题或建议，请通过以下方式联系我们：
- 提交 Issue
- 发送邮件至 [support@opencmp.org](mailto:support@opencmp.org)

---

感谢您使用 openCMP！如果您觉得这个项目有用，请给我们一个 Star ⭐。