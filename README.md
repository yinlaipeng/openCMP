# openCMP - 多云管理平台

[![Go Report Card](https://goreportcard.com/badge/github.com/opencmp/opencmp)](https://goreportcard.com/report/github.com/opencmp/opencmp)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

openCMP 是一个基于 Go + Vue 的全栈多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

## 特性

- **统一多云管理**: 屏蔽底层云厂商差异，提供统一 API 接口
- **快速接入**: 新增云厂商只需实现标准接口并注册
- **核心资源覆盖**: 支持 VM、VPC、Subnet、SecurityGroup、Disk、EIP 等资源管理
- **云账户管理**: 支持多云账户配置和验证
- **现代化前端**: 基于 Vue 3 + Element Plus 的响应式管理界面

## 项目结构

```
openCMP/
├── backend/           # 后端 Go 代码
│   ├── cmd/          # 应用入口
│   ├── configs/      # 配置文件
│   ├── internal/     # 内部包 (handler/service/model)
│   ├── pkg/          # 公共包 (cloudprovider)
│   ├── scripts/      # 脚本文件
│   └── ...
├── frontend/         # 前端 Vue 代码
│   ├── src/
│   │   ├── api/     # API 接口
│   │   ├── layout/  # 布局组件
│   │   ├── router/  # 路由配置
│   │   ├── types/   # TypeScript 类型
│   │   ├── utils/   # 工具函数
│   │   └── views/   # 页面组件
│   └── ...
├── docs/            # 文档
└── README.md        # 项目说明
```

## 支持的云厂商

| 云厂商 | 后端实现 | 前端支持 | 支持资源 |
|--------|---------|---------|----------|
| 阿里云 | ✅ 完整实现 | ✅ | VM, VPC, Subnet, SecurityGroup, EIP, Disk, Snapshot |
| 腾讯云 | 🚧 骨架 | ✅ | 待实现 |
| AWS | 🚧 骨架 | ✅ | 待实现 |
| Azure | 🚧 骨架 | ✅ | 待实现 |

## 快速开始

### 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 初始化数据库
mysql -u root -p < scripts/init.sql

# 修改配置
vim configs/config.yaml

# 运行
go run cmd/server/main.go -config configs/config.yaml
```

后端服务运行在 `http://localhost:8080`

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev
```

前端服务运行在 `http://localhost:3000`

### Docker 部署

```bash
# 构建后端镜像
cd backend
docker build -t opencmp-backend:latest .

# 构建前端镜像
cd frontend
docker build -t opencmp-frontend:latest .

# 运行
docker run -d -p 8080:8080 opencmp-backend:latest
docker run -d -p 3000:3000 opencmp-frontend:latest
```

## 功能截图

### 云账户管理
- 添加/编辑/删除云账户
- 云账户连接验证
- 支持阿里云、腾讯云、AWS、Azure

### 计算资源管理
- 虚拟机列表查看
- 虚拟机启动/停止/重启
- 虚拟机创建（对接云厂商 API）
- 镜像管理

### 网络资源管理
- VPC 创建和管理
- 子网管理
- 安全组管理
- 弹性 IP 管理

## API 端点

### 云账户管理
- `GET /api/v1/cloud-accounts` - 列出云账户
- `POST /api/v1/cloud-accounts` - 创建云账户
- `DELETE /api/v1/cloud-accounts/:id` - 删除云账户
- `POST /api/v1/cloud-accounts/:id/verify` - 验证云账户

### 计算资源
- `GET /api/v1/compute/vms` - 列出虚拟机
- `POST /api/v1/compute/vms` - 创建虚拟机
- `POST /api/v1/compute/vms/:id/action` - 虚拟机操作
- `GET /api/v1/compute/images` - 列出镜像

### 网络资源
- `GET /api/v1/network/vpcs` - 列出 VPC
- `POST /api/v1/network/vpcs` - 创建 VPC
- `GET /api/v1/network/subnets` - 列出子网
- `GET /api/v1/network/security-groups` - 列出安全组
- `GET /api/v1/network/eips` - 列出弹性 IP

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: Gorm
- **日志**: Zap
- **架构**: 单体模块化 + 适配器模式

### 前端
- **框架**: Vue 3 + TypeScript
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios

## 接入新云厂商

### 后端接入

1. 在 `backend/pkg/cloudprovider/adapters/` 下创建新目录
2. 实现 `ICloudProvider` 接口
3. 在 `init()` 函数中注册

```go
func init() {
    cloudprovider.RegisterProvider("huawei", NewHuaweiProvider)
}
```

### 前端支持

前端已支持所有云厂商的通用接口，无需修改即可使用新接入的云厂商。

## 开发

### 后端开发

```bash
cd backend

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 构建
go build -o opencmp cmd/server/main.go
```

### 前端开发

```bash
cd frontend

# 运行测试
npm run test

# 代码格式化
npm run lint

# 构建
npm run build
```

## License

Apache License 2.0
