# IAM 模块实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建完整的IAM（身份与访问管理）模块，包括域、项目、用户、组、角色、权限管理功能

**Architecture:** 采用分层架构，包含数据模型层(model)、业务逻辑层(service)、API接口层(handler)和中间件层(middleware)

**Tech Stack:** Go, Gin, Gorm, PostgreSQL/MySQL

---

## 已完成的任务

### 1. 数据模型层
- [x] 创建域数据模型 (internal/model/domain.go)
- [x] 创建项目数据模型 (internal/model/project.go)
- [x] 创建用户数据模型 (internal/model/user.go)
- [x] 创建组数据模型 (internal/model/group.go)
- [x] 创建角色数据模型 (internal/model/role.go)
- [x] 创建权限数据模型 (internal/model/permission.go)

### 2. 业务逻辑层
- [x] 创建域服务 (internal/service/domain_service.go)
- [x] 创建项目服务 (internal/service/project_service.go)
- [x] 创建用户服务 (internal/service/user_service.go)
- [x] 创建组服务 (internal/service/group_service.go)
- [x] 创建角色服务 (internal/service/role_service.go)
- [x] 创建权限服务 (internal/service/permission_service.go)

### 3. API接口层
- [x] 创建域API处理器 (internal/handler/domain_handler.go)
- [x] 创建项目API处理器 (internal/handler/project_handler.go)
- [x] 创建用户API处理器 (internal/handler/user_handler.go)
- [x] 创建组API处理器 (internal/handler/group_handler.go)
- [x] 创建角色和权限API处理器 (internal/handler/role_permission_handler.go)

### 4. 单元测试
- [x] 创建域和项目单元测试 (internal/service/domain_service_test.go, internal/service/project_service_test.go)
- [x] 创建用户和组单元测试 (internal/service/user_service_test.go, internal/service/group_service_test.go)
- [x] 创建角色和权限单元测试 (internal/service/role_service_test.go, internal/service/permission_service_test.go)
- [x] 创建API处理器单元测试 (internal/handler/domain_handler_test.go, internal/handler/project_handler_test.go, internal/handler/user_handler_test.go, internal/handler/group_handler_test.go)

## 下一步计划

### 5. 完善认证和授权中间件
- [ ] 创建JWT认证中间件
- [ ] 创建权限检查中间件
- [ ] 实现认证和授权流程

### 6. 完善主服务器路由
- [ ] 更新main.go以包含所有IAM相关路由
- [ ] 确保所有API端点正确注册

### 7. 完善数据模型关系
- [ ] 定义用户-角色关联模型
- [ ] 定义用户-组关联模型
- [ ] 定义角色-权限关联模型

### 8. 完善业务逻辑
- [ ] 实现用户角色分配功能
- [ ] 实现用户权限检查功能
- [ ] 实现基于角色的访问控制(RBAC)

## 详细实施步骤

### Task 1: 创建用户-角色关联模型
**Files:**
- Create: `internal/model/user_role.go`

```go
package model

import (
	"time"
)

// UserRole represents the relationship between a user and a role
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	RoleID    uint      `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user"`
	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}
```

### Task 2: 创建角色-权限关联模型
**Files:**
- Create: `internal/model/role_permission.go`

```go
package model

import (
	"time"
)

// RolePermission represents the relationship between a role and a permission
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	
	// Relationships
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}
```

### Task 3: 更新主服务器文件
**Files:**
- Modify: `cmd/server/main.go`
- Add routes for new handlers

### Task 4: 创建认证中间件
**Files:**
- Create: `internal/middleware/auth.go`

### Task 5: 创建权限检查中间件
**Files:**
- Create: `internal/middleware/permission.go`

## 验证和测试

### Task 6: 集成测试
- [ ] 创建端到端测试场景
- [ ] 测试用户认证流程
- [ ] 测试权限检查流程

### Task 7: 安全测试
- [ ] 验证认证机制的安全性
- [ ] 测试权限绕过漏洞
- [ ] 验证敏感信息保护

## 部署和监控

### Task 8: 配置管理
- [ ] 添加配置选项
- [ ] 实现环境变量支持

### Task 9: 日志和监控
- [ ] 添加审计日志
- [ ] 实现操作跟踪
- [ ] 添加性能指标

计划完成，IAM模块的基础功能已实现，下一步将完善高级功能和安全特性。