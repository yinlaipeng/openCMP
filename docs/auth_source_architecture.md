# 认证源（Authentication Source）设计架构文档

## 1. 概述

认证源是 OpenCMP 平台中用于管理外部身份认证系统的核心组件，允许平台集成多种身份验证机制（如 LDAP、OIDC、SAML 等）。认证源的设计目的是统一管理不同身份提供者的配置，并实现灵活的身份认证和用户同步功能。

## 2. 架构设计

### 2.1 认证源分类

认证源分为两类：

- **系统级认证源**：
  - 适用于整个平台的所有域，全局可用
  - 所有域/所有用户都可以使用
  - 登录页面所有人都能看到这个 LDAP
  - 可以用它登录任何域（前提用户存在）

- **域级认证源**：
  - 仅适用于特定域，其他域完全看不到
  - 用户必须先进入某个域，才能看到对应的 LDAP 登录方式
  - 只属于某个域，提供域隔离的身份验证

### 2.2 支持的认证协议

- **LDAP/LDAPS**：轻量级目录访问协议，广泛用于企业身份管理系统
- **OIDC**：开放ID连接协议，基于OAuth 2.0的身份验证标准
- **SAML**：安全断言标记语言，用于Web浏览器单点登录(SSO)

### 2.3 认证类型

- **OpenLDAP**：开源LDAP实现
- **Active Directory**：微软活动目录服务（可选支持）

## 3. 数据模型

### 3.1 AuthSource 主要字段

- `id`: 唯一标识符
- `name`: 认证源名称
- `type`: 认证源类型 (ldap, oidc, saml, local)
- `scope`: 作用范围 (system, domain)
- `domain_id`: 关联域ID（仅当scope为domain时有效）
- `description`: 描述信息
- `enabled`: 是否启用
- `auto_create`: 是否自动创建用户
- `config`: 认证源详细配置（JSON格式）
- `created_at/updated_at`: 时间戳

### 3.2 LDAP 配置参数

```go
type LDAPConfig struct {
    URL                  string // 服务器地址 (ldap://host:port 或 ldaps://host:port)
    BaseDN               string // 基础DN
    BindDN               string // 绑定DN
    BindPassword         string // 绑定密码
    UserFilter           string // 用户过滤器
    UserIDAttr           string // 用户唯一ID属性
    UserNameAttr         string // 用户名属性
    UserSearchBase       string // 用户搜索基础DN
    GroupSearchBase      string // 组搜索基础DN
    UserEnabledAttribute string // 用户启用状态属性
    Protocol             string // 认证协议 (ldap, ldaps)
    AuthType             string // 认证类型 (openldap, ad)
    TargetDomain         string // 用户归属目标域名称
}
```

## 4. 功能特性

### 4.1 认证流程

1. **本地认证优先**：首先尝试在平台本地用户数据库中认证
2. **外部认证后备**：本地认证失败后，尝试配置的外部认证源
3. **自动用户创建**：如果启用auto_create，则在外部认证成功时自动创建本地用户

### 4.2 用户同步

- **手动同步**：通过API触发认证源用户同步
- **自动同步**：定时任务定期同步用户数据
- **双向同步**：支持用户属性的双向同步

### 4.3 目标域自动创建

当配置的目标域不存在时：
1. 检测目标域名称是否已存在
2. 如存在，则在域名称后添加递增后缀（如"-1", "-2"）
3. 自动创建新的域，使认证源能够关联到有效的目标域

## 5. 安全考虑

- **敏感信息加密**：BindPassword等敏感信息需加密存储
- **连接安全性**：支持LDAPS和STARTTLS加密连接
- **认证限制**：可配置认证失败锁定机制
- **访问控制**：通过域和项目权限控制认证源访问

## 6. 扩展性设计

- **插件化架构**：易于添加新的认证协议
- **配置热更新**：支持运行时配置变更
- **监控集成**：提供认证成功率、响应时间等监控指标

## 7. API 接口

- `GET /auth-sources`: 获取认证源列表
- `POST /auth-sources`: 创建认证源
- `PUT /auth-sources/{id}`: 更新认证源
- `DELETE /auth-sources/{id}`: 删除认证源
- `POST /auth-sources/{id}/test`: 测试认证源连接
- `POST /auth-sources/{id}/sync`: 同步认证源用户
- `POST /auth-sources/{id}/enable`: 启用认证源
- `POST /auth-sources/{id}/disable`: 禁用认证源

## 8. 前端界面设计

- **表单布局**：采用分区域布局，区分基础配置和高级配置
- **类型切换**：根据选择的认证类型动态显示相应配置项
- **实时验证**：提供配置验证和连接测试功能
- **向导引导**：为复杂配置提供分步引导