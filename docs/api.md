# openCMP API 文档

## 概述

openCMP 提供完整的 RESTful API，用于多云资源管理、IAM 身份管理、消息中心等功能。

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Token
- **请求格式**: JSON
- **响应格式**: JSON

## 认证

### 登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "admin",
    "display_name": "管理员"
  }
}
```

### 使用 Token

所有 API 请求需要在 Header 中携带 JWT Token：

```http
Authorization: Bearer <token>
```

## 云账户管理

### 列出云账户

```http
GET /api/v1/cloud-accounts
```

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "阿里云生产账号",
      "provider_type": "alibaba",
      "status": "active",
      "health_status": "healthy"
    }
  ],
  "total": 1
}
```

### 创建云账户

```http
POST /api/v1/cloud-accounts
Content-Type: application/json

{
  "name": "阿里云测试账号",
  "provider_type": "alibaba",
  "credentials": {
    "access_key_id": "AKIAXXXXX",
    "access_key_secret": "XXXXX"
  },
  "description": "阿里云测试环境账号"
}
```

### 验证云账户凭证

```http
POST /api/v1/cloud-accounts/:id/verify
```

### 同步云账户资源

```http
POST /api/v1/cloud-accounts/:id/sync
Content-Type: application/json

{
  "resource_types": ["vm", "vpc", "subnet", "disk"]
}
```

## 计算资源

### 虚拟机列表

```http
GET /api/v1/compute/vms?cloud_account_id=1&region=cn-hangzhou
```

### 创建虚拟机

```http
POST /api/v1/compute/vms
Content-Type: application/json

{
  "cloud_account_id": 1,
  "region": "cn-hangzhou",
  "zone": "cn-hangzhou-a",
  "name": "test-vm-01",
  "instance_type": "ecs.g6.large",
  "image_id": "ubuntu_22_04",
  "vpc_id": "vpc-xxx",
  "subnet_id": "subnet-xxx",
  "security_group_ids": ["sg-xxx"],
  "password": "YourPassword123"
}
```

### 虚拟机操作

```http
POST /api/v1/compute/vms/:id/action
Content-Type: application/json

{
  "action": "start"  // start, stop, restart, delete
}
```

## 网络资源

### VPC 列表

```http
GET /api/v1/network/vpcs?cloud_account_id=1
```

### 创建 VPC

```http
POST /api/v1/network/vpcs
Content-Type: application/json

{
  "cloud_account_id": 1,
  "region": "cn-hangzhou",
  "name": "test-vpc",
  "cidr": "10.0.0.0/16",
  "description": "测试VPC"
}
```

### 子网列表

```http
GET /api/v1/network/subnets?vpc_id=vpc-xxx
```

### 创建子网

```http
POST /api/v1/network/subnets
Content-Type: application/json

{
  "cloud_account_id": 1,
  "vpc_id": "vpc-xxx",
  "zone": "cn-hangzhou-a",
  "name": "subnet-a",
  "cidr": "10.0.1.0/24"
}
```

### 弹性 IP

```http
GET /api/v1/network/eips
POST /api/v1/network/eips
POST /api/v1/network/eips/:id/bind
POST /api/v1/network/eips/:id/unbind
```

## 存储资源

### 云硬盘列表

```http
GET /api/v1/storage/cloud-disks?cloud_account_id=1
```

### 创建云硬盘

```http
POST /api/v1/storage/cloud-disks
Content-Type: application/json

{
  "cloud_account_id": 1,
  "region": "cn-hangzhou",
  "zone": "cn-hangzhou-a",
  "name": "disk-01",
  "size": 100,
  "type": "cloud_ssd"
}
```

### 挂载/卸载硬盘

```http
POST /api/v1/storage/cloud-disks/:id/attach
{
  "instance_id": "vm-xxx"
}

POST /api/v1/storage/cloud-disks/:id/detach
```

## 数据库资源

### RDS 实例

```http
GET /api/v1/database/rds
POST /api/v1/database/rds
POST /api/v1/database/rds/:id/action
```

### Redis 实例

```http
GET /api/v1/database/cache
POST /api/v1/database/cache
POST /api/v1/database/cache/:id/action
```

## IAM 管理

### 域管理

```http
GET /api/v1/domains
POST /api/v1/domains
PUT /api/v1/domains/:id
DELETE /api/v1/domains/:id
```

### 项目管理

```http
GET /api/v1/projects
POST /api/v1/projects
PUT /api/v1/projects/:id
POST /api/v1/projects/:id/join
GET /api/v1/projects/:id/users
```

### 用户管理

```http
GET /api/v1/users
POST /api/v1/users
PUT /api/v1/users/:id
POST /api/v1/users/:id/reset-password
GET /api/v1/users/:id/roles
POST /api/v1/users/:id/roles
```

### 角色管理

```http
GET /api/v1/roles
POST /api/v1/roles
PUT /api/v1/roles/:id
DELETE /api/v1/roles/:id
```

### 权限管理

```http
GET /api/v1/permissions
POST /api/v1/permissions
POST /api/v1/permissions/role/:role_id/assign/:permission_id
```

## 消息中心

### 消息列表

```http
GET /api/v1/messages
GET /api/v1/messages/unread-count
PUT /api/v1/messages/:id/read
POST /api/v1/messages/mark-all-read
```

### 通知渠道

```http
GET /api/v1/notification-channels
POST /api/v1/notification-channels
POST /api/v1/notification-channels/:id/test
```

### 机器人管理

```http
GET /api/v1/robots
POST /api/v1/robots
POST /api/v1/robots/:id/test
```

## 同步管理

### 同步策略

```http
GET /api/v1/sync-policies
POST /api/v1/sync-policies
PUT /api/v1/sync-policies/:id
DELETE /api/v1/sync-policies/:id
```

### 定时任务

```http
GET /api/v1/scheduled-tasks
POST /api/v1/scheduled-tasks
PUT /api/v1/scheduled-tasks/:id
POST /api/v1/scheduled-tasks/:id/execute
```

### 同步日志

```http
GET /api/v1/sync-logs
GET /api/v1/sync-logs/statistics
GET /api/v1/sync-logs/latest
```

## 监控告警

### 告警策略

```http
GET /api/v1/monitor/alert-policies
POST /api/v1/monitor/alert-policies
PUT /api/v1/monitor/alert-policies/:id
DELETE /api/v1/monitor/alert-policies/:id
POST /api/v1/monitor/alert-policies/:id/toggle
```

### 告警历史

```http
GET /api/v1/monitor/alert-history
POST /api/v1/monitor/alert-history/:id/resolve
POST /api/v1/monitor/alert-history/:id/ignore
```

### 监控资源

```http
GET /api/v1/monitor/resources
GET /api/v1/monitor/resources/:id/metrics
```

## 费用中心

### 账单

```http
GET /api/v1/finance/bills
POST /api/v1/finance/bills/sync
POST /api/v1/finance/bills/export
```

### 成本分析

```http
GET /api/v1/finance/cost/analysis
GET /api/v1/finance/cost/reports
```

### 预算管理

```http
GET /api/v1/finance/budgets
POST /api/v1/finance/budgets
PUT /api/v1/finance/budgets/:id
DELETE /api/v1/finance/budgets/:id
```

## 健康检查

```http
GET /health
```

**响应**:
```json
{
  "status": "ok"
}
```

## 错误响应

所有错误响应遵循统一格式：

```json
{
  "error": "错误信息",
  "code": "ERROR_CODE"
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 参数验证失败 |
| 401 | 认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 分页

列表接口支持分页参数：

```http
GET /api/v1/compute/vms?page=1&page_size=20
```

响应包含分页信息：

```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

## 权限格式

权限使用格式：`<module>:<resource>:<action>`

示例：
- `compute:vm:list` - 查看虚拟机列表
- `compute:vm:create` - 创建虚拟机
- `iam:user:create` - 创建用户
- `finance:bill:view` - 查看账单