# openCMP - 多云管理平台

[![Go Report Card](https://goreportcard.com/badge/github.com/opencmp/opencmp)](https://goreportcard.com/report/github.com/opencmp/opencmp)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

openCMP 是一个基于 Go + Gin + Gorm 的单体架构多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

## 特性

- **统一多云管理**: 屏蔽底层云厂商差异，提供统一 API 接口
- **快速接入**: 新增云厂商只需实现标准接口并注册
- **核心资源覆盖**: 支持 VM、VPC、Subnet、SecurityGroup、Disk、EIP 等资源管理
- **云账户管理**: 支持多云账户配置和验证

## 支持的云厂商

| 云厂商 | 状态 | 支持资源 |
|--------|------|----------|
| 阿里云 | ✅ 完整实现 | VM, VPC, Subnet, SecurityGroup, EIP, Disk, Snapshot |
| 腾讯云 | 🚧 骨架 | 待实现 |
| AWS | 🚧 骨架 | 待实现 |
| Azure | 🚧 骨架 | 待实现 |

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+ 或 PostgreSQL 12+
- Docker (可选)

### 安装

```bash
# 克隆代码
git clone https://github.com/opencmp/opencmp.git
cd opencmp

# 安装依赖
make deps

# 构建
make build

# 运行
make run
```

### 配置

编辑 `configs/config.yaml`:

```yaml
server:
  port: 8080
  mode: debug

database:
  driver: mysql
  dsn: "root:password@tcp(127.0.0.1:3306)/opencmp?charset=utf8mb4&parseTime=True&loc=Local"
```

### 初始化数据库

```bash
mysql -u root -p < scripts/init.sql
```

### Docker 部署

```bash
docker build -t opencmp:latest .
docker run -d -p 8080:8080 opencmp:latest
```

## API 使用示例

### 添加云账户

```bash
curl -X POST http://localhost:8080/api/v1/cloud-accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "name": "阿里云测试账号",
    "provider_type": "alibaba",
    "credentials": {
      "access_key_id": "your-access-key-id",
      "access_key_secret": "your-access-key-secret"
    },
    "description": "测试环境"
  }'
```

### 创建虚拟机

```bash
curl -X POST http://localhost:8080/api/v1/compute/vms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "account_id": 1,
    "name": "test-vm",
    "instance_type": "ecs.t5-lc1m1.small",
    "image_id": "m-bp1xxxxxxxx",
    "vpc_id": "vpc-bp1xxxxxxxx",
    "subnet_id": "vsw-bp1xxxxxxxx",
    "security_groups": ["sg-bp1xxxxxxxx"],
    "disk_size": 40
  }'
```

### 列出 VPC

```bash
curl -X GET "http://localhost:8080/api/v1/network/vpcs?account_id=1" \
  -H "Authorization: Bearer your-token"
```

### 虚拟机操作

```bash
# 启动虚拟机
curl -X POST http://localhost:8080/api/v1/compute/vms/vm-xxx/action \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"action": "start"}'

# 停止虚拟机
curl -X POST http://localhost:8080/api/v1/compute/vms/vm-xxx/action \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"action": "stop"}'
```

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Layer (Gin)                          │
├─────────────────────────────────────────────────────────────────┤
│                        Service Layer                             │
├─────────────────────────────────────────────────────────────────┤
│                     Cloud Provider Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ Alibaba  │  │ Tencent  │  │   AWS    │  │  Azure   │        │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
├─────────────────────────────────────────────────────────────────┤
│                  Cloud Interface Layer                           │
│         ICompute | INetwork | IStorage | IDatabase              │
├─────────────────────────────────────────────────────────────────┤
│                      Data Layer (Gorm)                           │
└─────────────────────────────────────────────────────────────────┘
```

## API 端点

### 云账户管理
- `POST /api/v1/cloud-accounts` - 创建云账户
- `GET /api/v1/cloud-accounts` - 列出云账户
- `GET /api/v1/cloud-accounts/:id` - 获取云账户详情
- `PUT /api/v1/cloud-accounts/:id` - 更新云账户
- `DELETE /api/v1/cloud-accounts/:id` - 删除云账户
- `POST /api/v1/cloud-accounts/:id/verify` - 验证云账户

### 计算资源
- `POST /api/v1/compute/vms` - 创建虚拟机
- `GET /api/v1/compute/vms` - 列出虚拟机
- `GET /api/v1/compute/vms/:id` - 获取虚拟机详情
- `DELETE /api/v1/compute/vms/:id` - 删除虚拟机
- `POST /api/v1/compute/vms/:id/action` - 虚拟机操作 (start/stop/reboot)
- `GET /api/v1/compute/images` - 列出镜像

### 网络资源
- `POST /api/v1/network/vpcs` - 创建 VPC
- `GET /api/v1/network/vpcs` - 列出 VPC
- `DELETE /api/v1/network/vpcs/:id` - 删除 VPC
- `POST /api/v1/network/subnets` - 创建子网
- `GET /api/v1/network/subnets` - 列出子网
- `POST /api/v1/network/security-groups` - 创建安全组
- `GET /api/v1/network/security-groups` - 列出安全组
- `POST /api/v1/network/eips` - 创建弹性 IP
- `GET /api/v1/network/eips` - 列出弹性 IP

## 接入新云厂商

1. 在 `pkg/cloudprovider/adapters/` 下创建新目录
2. 实现 `ICloudProvider` 接口
3. 在 `init()` 函数中注册适配器

```go
func init() {
    cloudprovider.RegisterProvider("huawei", NewHuaweiProvider)
}
```

## 开发

```bash
# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint
```

## License

Apache License 2.0
