# 权限系统设计文档

## 概述

本文档描述了 openCMP 多云管理平台的权限系统设计，参考了主流云管理平台的权限管理最佳实践。

## 核心概念

### 1. 权限 (Permission)

权限是最小的授权单元，格式为 `resource:action`。

**结构:**
```
权限标识：resource:action
显示名称：人类可读的权限名称
资源类型：cloud_account, vm, network, user, role 等
操作类型：list, get, create, update, delete, action, grant, revoke
类型：system(系统) / custom(自定义)
```

**示例:**
- `cloud_account:list` - 查看云账户列表
- `vm:create` - 创建虚拟机
- `user:delete` - 删除用户

### 2. 策略 (Policy)

策略是权限的集合，用于批量授权。

**结构:**
```json
{
  "id": 1,
  "name": "VMManager",
  "description": "虚拟机管理员",
  "type": "custom",
  "statements": [
    {
      "effect": "Allow",
      "resource": "vm:*",
      "action": ["list", "get", "create", "update", "delete", "action"]
    }
  ]
}
```

**策略语句 (Statement):**
- **Effect**: Allow (允许) / Deny (拒绝)
- **Resource**: 资源类型，支持通配符 `*`
- **Action**: 操作列表

### 3. 角色 (Role)

角色是权限/策略的载体，用户通过扮演角色获得权限。

**系统角色:**
- `admin` - 系统管理员，拥有所有权限
- `user` - 普通用户，拥有基础权限

**自定义角色:**
- 用户可创建角色并分配权限/策略

### 4. 用户 (User)

用户通过以下方式获得权限：
1. 直接分配角色
2. 加入用户组，通过用户组获得角色

## 权限层次结构

```
用户 (User)
  ├── 直接角色 (UserRole)
  │     └── 权限/策略
  └── 用户组 (Group)
        └── 角色 (GroupRole)
              └── 权限/策略
```

## 预定义权限

### 云账户管理
| 权限标识 | 显示名称 | 说明 |
|---------|---------|------|
| cloud_account:list | 查看云账户 | 查看云账户列表 |
| cloud_account:get | 查看云账户详情 | 查看单个云账户 |
| cloud_account:create | 创建云账户 | 创建新的云账户 |
| cloud_account:update | 更新云账户 | 更新云账户信息 |
| cloud_account:delete | 删除云账户 | 删除云账户 |
| cloud_account:verify | 验证云账户 | 验证云账户连接 |

### 虚拟机管理
| 权限标识 | 显示名称 | 说明 |
|---------|---------|------|
| vm:list | 查看虚拟机 | 查看虚拟机列表 |
| vm:get | 查看虚拟机详情 | 查看单个虚拟机 |
| vm:create | 创建虚拟机 | 创建新的虚拟机 |
| vm:update | 更新虚拟机 | 更新虚拟机配置 |
| vm:delete | 删除虚拟机 | 删除虚拟机 |
| vm:action | 操作虚拟机 | 启动/停止/重启虚拟机 |

### 网络资源管理
| 权限标识 | 显示名称 | 说明 |
|---------|---------|------|
| network:vpc:list | 查看 VPC | 查看 VPC 列表 |
| network:vpc:create | 创建 VPC | 创建新的 VPC |
| network:vpc:delete | 删除 VPC | 删除 VPC |
| network:subnet:list | 查看子网 | 查看子网列表 |
| network:subnet:create | 创建子网 | 创建新的子网 |
| network:security-group:list | 查看安全组 | 查看安全组列表 |
| network:eip:list | 查看弹性 IP | 查看弹性 IP 列表 |

### 用户与权限管理
| 权限标识 | 显示名称 | 说明 |
|---------|---------|------|
| user:list | 查看用户 | 查看用户列表 |
| user:create | 创建用户 | 创建新用户 |
| user:update | 更新用户 | 更新用户信息 |
| user:delete | 删除用户 | 删除用户 |
| role:list | 查看角色 | 查看角色列表 |
| role:create | 创建角色 | 创建新角色 |
| role:grant | 角色授权 | 为角色分配权限 |
| permission:list | 查看权限 | 查看权限列表 |

### 认证与安全
| 权限标识 | 显示名称 | 说明 |
|---------|---------|------|
| auth_source:list | 查看认证源 | 查看认证源列表 |
| auth_source:create | 创建认证源 | 创建新的认证源 |
| auth_source:test | 测试认证源 | 测试认证源连接 |
| message:list | 查看消息 | 查看消息列表 |
| alert:list | 查看告警 | 查看安全告警 |
| alert:resolve | 处理告警 | 处理安全告警 |

## 预定义策略

### AdministratorAccess (系统管理员)
```json
{
  "name": "AdministratorAccess",
  "statement": [
    {
      "effect": "Allow",
      "resource": "*:*",
      "action": ["*"]
    }
  ]
}
```

### CloudAccountReadOnly (云账户只读)
```json
{
  "name": "CloudAccountReadOnly",
  "statement": [
    {
      "effect": "Allow",
      "resource": "cloud_account:*",
      "action": ["list", "get"]
    }
  ]
}
```

### VMManager (虚拟机管理员)
```json
{
  "name": "VMManager",
  "statement": [
    {
      "effect": "Allow",
      "resource": "vm:*",
      "action": ["list", "get", "create", "update", "delete", "action"]
    }
  ]
}
```

### UserViewer (用户查看者)
```json
{
  "name": "UserViewer",
  "statement": [
    {
      "effect": "Allow",
      "resource": "user:*",
      "action": ["list", "get"]
    }
  ]
}
```

## 权限检查流程

```
用户请求资源
    ↓
获取用户所有角色
    ↓
收集所有角色的权限/策略
    ↓
检查是否有 Deny 策略 → 有 → 拒绝访问
    ↓ 无
检查是否有 Allow 策略 → 有 → 允许访问
    ↓ 无
拒绝访问
```

## 最佳实践

1. **最小权限原则**: 只授予完成工作所需的最小权限
2. **使用策略而非直接授权**: 通过策略管理权限，便于复用和维护
3. **定期审计权限**: 定期检查用户权限，移除不必要的权限
4. **使用用户组**: 通过用户组管理批量用户权限
5. **系统权限不可修改**: 系统预定义权限和策略不可删除或修改

## 未来扩展

1. **资源级权限**: 支持对特定资源实例的权限控制
2. **条件策略**: 支持基于时间、IP 等条件的权限控制
3. **权限审计日志**: 记录所有权限相关操作
4. **临时授权**: 支持有时效性的临时权限授权
5. **跨域授权**: 支持跨域的资源访问授权
