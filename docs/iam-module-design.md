# IAM（身份与访问管理）模块详细设计文档

## 1. 概述

### 1.1 文档目的
本文档详细描述 openCMP 多云管理平台的身份与访问管理（IAM）模块的设计方案，包括域、项目、组、用户、角色、权限和认证源的完整设计。

### 1.2 设计目标
- 提供灵活的多层级权限管理机制
- 支持细粒度的资源访问控制
- 实现多租户隔离和资源共享
- 提供统一的认证与授权接口
- 支持多种认证源集成

### 1.3 设计原则
- **最小权限原则**: 用户仅获得完成工作所需的最小权限
- **职责分离原则**: 不同角色承担不同职责，避免权限过度集中
- **权限继承原则**: 通过层级结构实现权限的自然继承
- **安全审计原则**: 所有权限变更和访问行为可追溯

## 2. 核心概念

### 2.1 域 (Domain)
域是最高级别的组织单元，代表一个独立的管理实体，如企业或组织。每个域具有独立的用户、组、角色和权限管理体系。

**特性:**
- 域之间完全隔离，无法跨域访问资源
- 每个域有唯一的标识符和名称
- 域管理员负责域内所有资源的管理
- 支持域级别的配置和策略

### 2.2 项目 (Project)
项目是域内的逻辑分组单元，用于组织和管理相关的云资源。一个域可以包含多个项目，项目之间可以共享资源。

**特性:**
- 项目隶属于特定域
- 项目内可以创建和管理云资源
- 支持项目级别的权限控制
- 项目成员可以访问项目内的资源

### 2.3 用户 (User)
用户是系统的实际使用者，可以是个人账户或服务账户。

**特性:**
- 用户隶属于特定域
- 用户可以属于多个组
- 用户可以直接被赋予角色
- 支持本地认证和外部认证源

### 2.4 组 (Group)
组是用户的逻辑集合，用于批量管理用户权限。

**特性:**
- 组隶属于特定域
- 一个用户可以属于多个组
- 组可以被赋予角色
- 便于批量权限管理

### 2.5 角色 (Role)
角色是一组权限的集合，定义了执行特定任务的能力。

**特性:**
- 角色可以是系统预定义或用户自定义
- 角色可以分配给用户或组
- 支持角色继承和组合
- 角色具有作用范围（域级、项目级）

### 2.6 权限 (Permission)
权限是执行特定操作的许可，是最小的授权单元。

**特性:**
- 权限格式为 `resource:action`
- 支持通配符匹配
- 权限可以聚合到策略中
- 支持显式允许和拒绝

### 2.7 认证源 (Auth Source)
认证源是验证用户身份的外部系统，如LDAP、AD、OAuth2等。

**特性:**
- 支持多种认证协议
- 可以映射外部用户到本地用户
- 支持用户同步和自动创建
- 提供统一的认证接口

## 3. 权限模型

### 3.1 RBAC模型扩展
采用基于角色的访问控制（RBAC）模型，并进行以下扩展：
- 支持多层级权限（域级、项目级）
- 支持权限继承
- 支持条件权限（基于时间、IP等条件）

### 3.2 权限评估流程
```
用户请求资源
    ↓
确定请求的作用域（域/项目）
    ↓
获取用户在该作用域的所有角色
    ↓
收集所有角色的权限
    ↓
应用策略规则（Allow/Deny）
    ↓
权限检查结果
```

### 3.3 权限继承关系
```
域级角色 → 项目级角色 → 资源级权限
    ↓           ↓           ↓
  域管理员   → 项目管理员 → 资源操作员
```

## 4. 数据模型设计

### 4.1 域表 (domains)
```go
type Domain struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    OwnerID     uint      `gorm:"index" json:"owner_id"` // 域管理员用户ID
    Status      string    `gorm:"type:varchar(20)" json:"status"` // active/inactive/deleted
    IsDefault   bool      `json:"is_default"` // 是否为默认域
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 4.2 项目表 (projects)
```go
type Project struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    DomainID    uint      `gorm:"index;not null" json:"domain_id"`
    OwnerID     uint      `gorm:"index" json:"owner_id"` // 项目管理员用户ID
    Status      string    `gorm:"type:varchar(20)" json:"status"` // active/inactive/deleted
    Quota       datatypes.JSON `gorm:"type:json" json:"quota"` // 配额限制
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Domain      Domain    `gorm:"foreignKey:DomainID" json:"-"`
}
```

### 4.3 用户表 (users)
```go
type User struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    Name            string    `gorm:"uniqueIndex;not null" json:"name"`
    Email           string    `gorm:"index" json:"email"`
    DisplayName     string    `json:"display_name"`
    PasswordHash    string    `gorm:"type:text" json:"-"` // 仅本地用户需要
    DomainID        uint      `gorm:"index;not null" json:"domain_id"`
    AuthSourceID    *uint     `gorm:"index" json:"auth_source_id"` // 外部认证源ID
    ExternalID      string    `json:"external_id"` // 外部认证源中的用户ID
    Status          string    `gorm:"type:varchar(20)" json:"status"` // active/inactive/locked
    LastLoginAt     *time.Time `json:"last_login_at"`
    PasswordExpiredAt *time.Time `json:"password_expired_at"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    
    Domain          Domain    `gorm:"foreignKey:DomainID" json:"-"`
    AuthSource      *AuthSource `gorm:"foreignKey:AuthSourceID" json:"-"`
}
```

### 4.4 组表 (groups)
```go
type Group struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    DomainID    uint      `gorm:"index;not null" json:"domain_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Domain      Domain    `gorm:"foreignKey:DomainID" json:"-"`
}
```

### 4.5 用户组关联表 (user_groups)
```go
type UserGroup struct {
    UserID   uint      `gorm:"primaryKey;not null" json:"user_id"`
    GroupID  uint      `gorm:"primaryKey;not null" json:"group_id"`
    CreatedAt time.Time `json:"created_at"`
    
    User     User      `gorm:"foreignKey:UserID" json:"-"`
    Group    Group     `gorm:"foreignKey:GroupID" json:"-"`
}
```

### 4.6 角色表 (roles)
```go
type Role struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    Type        string    `gorm:"type:varchar(20)" json:"type"` // system/custom
    Scope       string    `gorm:"type:varchar(20)" json:"scope"` // domain/project/global
    DomainID    *uint     `gorm:"index" json:"domain_id"` // 仅自定义角色需要
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Domain      *Domain   `gorm:"foreignKey:DomainID" json:"-"`
}
```

### 4.7 权限表 (permissions)
```go
type Permission struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"` // resource:action格式
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    ResourceType string   `gorm:"type:varchar(50);not null" json:"resource_type"`
    Action      string    `gorm:"type:varchar(50);not null" json:"action"`
    Type        string    `gorm:"type:varchar(20)" json:"type"` // system/custom
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 4.8 角色权限关联表 (role_permissions)
```go
type RolePermission struct {
    RoleID       uint      `gorm:"primaryKey;not null" json:"role_id"`
    PermissionID uint      `gorm:"primaryKey;not null" json:"permission_id"`
    Effect       string    `gorm:"type:varchar(10);default:Allow" json:"effect"` // Allow/Deny
    CreatedAt    time.Time `json:"created_at"`
    
    Role         Role      `gorm:"foreignKey:RoleID" json:"-"`
    Permission   Permission `gorm:"foreignKey:PermissionID" json:"-"`
}
```

### 4.9 用户角色关联表 (user_roles)
```go
type UserRole struct {
    UserID    uint      `gorm:"primaryKey;not null" json:"user_id"`
    RoleID    uint      `gorm:"primaryKey;not null" json:"role_id"`
    ScopeID   *uint     `gorm:"index" json:"scope_id"` // 作用域ID（项目ID或域ID）
    ScopeType string    `gorm:"type:varchar(20)" json:"scope_type"` // domain/project
    CreatedAt time.Time `json:"created_at"`
    
    User      User      `gorm:"foreignKey:UserID" json:"-"`
    Role      Role      `gorm:"foreignKey:RoleID" json:"-"`
}
```

### 4.10 组角色关联表 (group_roles)
```go
type GroupRole struct {
    GroupID   uint      `gorm:"primaryKey;not null" json:"group_id"`
    RoleID    uint      `gorm:"primaryKey;not null" json:"role_id"`
    ScopeID   *uint     `gorm:"index" json:"scope_id"` // 作用域ID（项目ID或域ID）
    ScopeType string    `gorm:"type:varchar(20)" json:"scope_type"` // domain/project
    CreatedAt time.Time `json:"created_at"`
    
    Group     Group     `gorm:"foreignKey:GroupID" json:"-"`
    Role      Role      `gorm:"foreignKey:RoleID" json:"-"`
}
```

### 4.11 认证源表 (auth_sources)
```go
type AuthSource struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    DisplayName string    `json:"display_name"`
    Type        string    `gorm:"type:varchar(20);not null" json:"type"` // ldap/ad/oauth2/local
    Config      datatypes.JSON `gorm:"type:json" json:"config"` // 认证源配置
    Status      string    `gorm:"type:varchar(20)" json:"status"` // active/inactive
    SyncEnabled bool      `json:"sync_enabled"` // 是否启用用户同步
    SyncInterval int      `json:"sync_interval"` // 同步间隔（分钟）
    LastSyncAt  *time.Time `json:"last_sync_at"` // 最后同步时间
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## 5. 预定义角色和权限

### 5.1 系统预定义角色

#### 5.1.1 系统管理员 (system_admin)
- **作用域**: 全局
- **权限**: 拥有系统所有权限
- **描述**: 系统最高权限角色，可管理所有资源

#### 5.1.2 域管理员 (domain_admin)
- **作用域**: 域级
- **权限**: 
  - `domain:update`
  - `user:*`
  - `group:*`
  - `role:*`
  - `project:*`
  - `auth_source:list`
  - `auth_source:get`
- **描述**: 负责管理特定域内的所有资源

#### 5.1.3 项目管理员 (project_admin)
- **作用域**: 项目级
- **权限**:
  - `project:update`
  - `project:delete`
  - `vm:*`
  - `network:*`
  - `storage:*`
  - `database:*`
  - `cloud_account:list`
  - `cloud_account:get`
- **描述**: 负责管理特定项目的资源

#### 5.1.4 云资源操作员 (cloud_operator)
- **作用域**: 项目级
- **权限**:
  - `vm:create`
  - `vm:update`
  - `vm:delete`
  - `vm:action`
  - `network:*`
  - `storage:create`
  - `storage:delete`
  - `database:create`
  - `database:delete`
- **描述**: 可以创建和管理云资源的操作员

#### 5.1.5 云资源查看者 (cloud_viewer)
- **作用域**: 项目级
- **权限**:
  - `vm:list`
  - `vm:get`
  - `network:list`
  - `network:get`
  - `storage:list`
  - `storage:get`
  - `database:list`
  - `database:get`
- **描述**: 只能查看云资源的用户

### 5.2 预定义权限
参考权限系统设计文档中的预定义权限，包括：
- 云账户管理权限
- 虚拟机管理权限
- 网络资源管理权限
- 存储资源管理权限
- 数据库服务管理权限
- 用户与权限管理权限
- 认证与安全权限

## 6. API 设计

### 6.1 域管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/domains` | 创建域 |
| GET | `/api/v1/domains` | 列出域 |
| GET | `/api/v1/domains/:id` | 获取域详情 |
| PUT | `/api/v1/domains/:id` | 更新域 |
| DELETE | `/api/v1/domains/:id` | 删除域 |

### 6.2 项目管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/projects` | 创建项目 |
| GET | `/api/v1/projects` | 列出项目 |
| GET | `/api/v1/projects/:id` | 获取项目详情 |
| PUT | `/api/v1/projects/:id` | 更新项目 |
| DELETE | `/api/v1/projects/:id` | 删除项目 |

### 6.3 用户管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/users` | 创建用户 |
| GET | `/api/v1/users` | 列出用户 |
| GET | `/api/v1/users/:id` | 获取用户详情 |
| PUT | `/api/v1/users/:id` | 更新用户 |
| DELETE | `/api/v1/users/:id` | 删除用户 |
| POST | `/api/v1/users/:id/reset-password` | 重置用户密码 |
| POST | `/api/v1/users/:id/enable` | 启用用户 |
| POST | `/api/v1/users/:id/disable` | 禁用用户 |

### 6.4 组管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/groups` | 创建组 |
| GET | `/api/v1/groups` | 列出组 |
| GET | `/api/v1/groups/:id` | 获取组详情 |
| PUT | `/api/v1/groups/:id` | 更新组 |
| DELETE | `/api/v1/groups/:id` | 删除组 |
| GET | `/api/v1/groups/:id/members` | 获取组成员 |
| POST | `/api/v1/groups/:id/members` | 添加组成员 |
| DELETE | `/api/v1/groups/:id/members/:user_id` | 移除组成员 |

### 6.5 角色管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/roles` | 创建角色 |
| GET | `/api/v1/roles` | 列出角色 |
| GET | `/api/v1/roles/:id` | 获取角色详情 |
| PUT | `/api/v1/roles/:id` | 更新角色 |
| DELETE | `/api/v1/roles/:id` | 删除角色 |
| GET | `/api/v1/roles/:id/permissions` | 获取角色权限 |
| POST | `/api/v1/roles/:id/permissions` | 为角色添加权限 |
| DELETE | `/api/v1/roles/:id/permissions/:permission_id` | 从角色移除权限 |

### 6.6 权限管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/permissions` | 列出权限 |
| GET | `/api/v1/permissions/:id` | 获取权限详情 |

### 6.7 认证源管理 API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth-sources` | 创建认证源 |
| GET | `/api/v1/auth-sources` | 列出认证源 |
| GET | `/api/v1/auth-sources/:id` | 获取认证源详情 |
| PUT | `/api/v1/auth-sources/:id` | 更新认证源 |
| DELETE | `/api/v1/auth-sources/:id` | 删除认证源 |
| POST | `/api/v1/auth-sources/:id/test` | 测试认证源连接 |
| POST | `/api/v1/auth-sources/:id/sync` | 同步用户 |

### 6.8 权限分配 API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/users/:id/roles` | 获取用户的角色 |
| POST | `/api/v1/users/:id/roles` | 为用户分配角色 |
| DELETE | `/api/v1/users/:id/roles/:role_id` | 移除用户的角色 |
| GET | `/api/v1/groups/:id/roles` | 获取组的角色 |
| POST | `/api/v1/groups/:id/roles` | 为组分配角色 |
| DELETE | `/api/v1/groups/:id/roles/:role_id` | 移除组的角色 |

## 7. 认证与授权流程

### 7.1 用户认证流程
```
用户登录请求
    ↓
解析认证信息（用户名/密码 或 Token）
    ↓
确定认证源（本地 或 外部）
    ↓
执行认证验证
    ↓
生成会话Token
    ↓
返回认证结果
```

### 7.2 权限检查流程
```
用户发起资源请求
    ↓
解析用户身份和Token
    ↓
确定请求的作用域
    ↓
获取用户在该作用域的角色
    ↓
收集角色对应的权限
    ↓
检查是否满足请求权限
    ↓
允许/拒绝访问
```

### 7.3 Token管理
- 使用JWT作为认证Token格式
- Token包含用户ID、域ID、过期时间等信息
- 支持Token刷新机制
- 提供Token吊销功能

## 8. 安全考虑

### 8.1 密码安全
- 密码必须符合复杂度要求
- 使用强哈希算法（bcrypt）存储密码
- 支持密码过期策略
- 记录失败登录尝试，实施锁定策略

### 8.2 会话管理
- 设置合理的会话超时时间
- 支持并发会话控制
- 提供主动登出功能
- 定期清理过期会话

### 8.3 权限安全
- 实施最小权限原则
- 定期审计权限分配
- 记录权限变更日志
- 支持权限审批流程

### 8.4 数据保护
- 敏感数据（如密码）加密存储
- 认证源配置信息加密存储
- 实施数据访问日志记录
- 支持数据脱敏功能

## 9. 性能优化

### 9.1 缓存策略
- 缓存用户角色和权限信息
- 缓存认证源配置
- 实现缓存失效和更新机制
- 使用分布式缓存提高性能

### 9.2 数据库优化
- 为常用查询字段建立索引
- 优化权限查询SQL
- 使用连接池管理数据库连接
- 实现读写分离（如需要）

### 9.3 API优化
- 实现分页查询减少数据传输
- 提供批量操作接口
- 使用异步处理长时间操作
- 实现请求频率限制

## 10. 扩展性设计

### 10.1 插件化认证源
- 设计认证源插件接口
- 支持动态加载认证源插件
- 提供认证源SDK
- 支持自定义认证逻辑

### 10.2 条件权限
- 支持基于时间的权限控制
- 支持基于IP地址的权限控制
- 支持基于设备的权限控制
- 提供条件表达式引擎

### 10.3 审计日志扩展
- 记录所有认证和授权操作
- 支持审计日志导出
- 提供审计日志分析功能
- 支持第三方审计系统集成

## 11. 部署与配置

### 11.1 配置文件
```yaml
# configs/iam-config.yaml
iam:
  jwt:
    secret: "your-jwt-secret-key"
    expiration: 3600  # 1 hour
    refresh_expiration: 604800  # 7 days
  
  password_policy:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_number: true
    require_special_char: true
    expire_days: 90
  
  session:
    timeout: 3600  # 1 hour
    max_concurrent: 5
  
  cache:
    enabled: true
    type: redis
    address: "localhost:6379"
    db: 0
```

### 11.2 初始化脚本
系统启动时执行以下初始化操作：
- 创建默认域
- 创建系统预定义角色
- 创建初始管理员用户
- 注册认证源插件

## 12. 监控与运维

### 12.1 监控指标
- 认证成功率/失败率
- 权限检查响应时间
- 在线用户数量
- 会话创建/销毁速率

### 12.2 日志记录
- 认证事件日志
- 权限检查日志
- 权限变更日志
- 安全事件日志

### 12.3 告警机制
- 连续认证失败告警
- 异常权限访问告警
- 系统异常告警
- 性能指标告警

## 13. 测试策略

### 13.1 单元测试
- 用户管理功能测试
- 角色权限功能测试
- 认证授权流程测试
- 数据模型验证测试

### 13.2 集成测试
- API端点功能测试
- 数据库操作测试
- 认证源集成测试
- 权限继承测试

### 13.3 安全测试
- 权限绕过测试
- 认证漏洞测试
- 数据泄露测试
- 会话劫持测试

## 14. 未来扩展

### 14.1 多因素认证 (MFA)
- 支持TOTP认证
- 支持短信验证码
- 支持硬件令牌
- 支持生物识别

### 14.2 单点登录 (SSO)
- 支持SAML 2.0
- 支持OpenID Connect
- 支持企业微信/钉钉集成
- 支持自定义SSO协议

### 14.3 动态权限
- 基于风险的动态授权
- 基于上下文的权限调整
- 临时权限提升
- 基于AI的行为分析

### 14.4 合规性功能
- GDPR合规支持
- SOX合规报告
- 等保2.0支持
- 行业特定合规

## 15. 自审检查

- [x] 无 TBD/TODO 占位符
- [x] 所有核心概念明确定义
- [x] 数据模型设计完整
- [x] API设计清晰
- [x] 安全考虑充分
- [x] 扩展性设计合理
- [x] 部署配置方案明确
- [x] 监控运维策略完善
- [x] 测试策略全面
- [x] 未来扩展方向明确