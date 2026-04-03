# openCMP IAM 模块 API 接口文档

## 1. 认证与授权

### 1.1 用户登录
- **POST** `/api/v1/auth/login`
- **请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```
- **响应体**:
```json
{
  "token": "string",
  "expires_in": "int"
}
```

## 2. 域管理 (Domains)

### 2.1 获取域列表
- **GET** `/api/v1/domains`
- **查询参数**:
  - `keyword`: 搜索关键字
  - `enabled`: 状态过滤
  - `limit`: 每页数量
  - `offset`: 偸移量
- **响应体**:
```json
{
  "items": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "enabled": true,
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ],
  "total": 10
}
```

### 2.2 创建域
- **POST** `/api/v1/domains`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "enabled": true
}
```

### 2.3 更新域
- **PUT** `/api/v1/domains/{id}`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "enabled": true
}
```

### 2.4 删除域
- **DELETE** `/api/v1/domains/{id}`

### 2.5 启用/禁用域
- **POST** `/api/v1/domains/{id}/enable` (启用)
- **POST** `/api/v1/domains/{id}/disable` (禁用)

## 3. 项目管理 (Projects)

### 3.1 获取项目列表
- **GET** `/api/v1/projects`
- **查询参数**:
  - `domain_id`: 域ID
  - `keyword`: 搜索关键字
  - `enabled`: 状态过滤
  - `limit`: 每页数量
  - `offset`: 偸移量

### 3.2 创建项目
- **POST** `/api/v1/projects`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "domain_id": 1
}
```

### 3.3 更新项目
- **PUT** `/api/v1/projects/{id}`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "domain_id": 1
}
```

### 3.4 删除项目
- **DELETE** `/api/v1/projects/{id}`

### 3.5 启用/禁用项目
- **POST** `/api/v1/projects/{id}/enable` (启用)
- **POST** `/api/v1/projects/{id}/disable` (禁用)

## 4. 用户管理 (Users)

### 4.1 获取用户列表
- **GET** `/api/v1/users`
- **查询参数**:
  - `domain_id`: 域ID
  - `limit`: 每页数量
  - `offset`: 偸移量

### 4.2 创建用户
- **POST** `/api/v1/users`
- **请求体**:
```json
{
  "name": "string",
  "display_name": "string",
  "email": "string",
  "phone": "string",
  "password": "string",
  "domain_id": 1
}
```

### 4.3 更新用户
- **PUT** `/api/v1/users/{id}`
- **请求体**:
```json
{
  "display_name": "string",
  "email": "string",
  "phone": "string"
}
```

### 4.4 删除用户
- **DELETE** `/api/v1/users/{id}`

### 4.5 启用/禁用用户
- **POST** `/api/v1/users/{id}/enable` (启用)
- **POST** `/api/v1/users/{id}/disable` (禁用)

### 4.6 重置用户密码
- **POST** `/api/v1/users/{id}/reset-password`
- **请求体**:
```json
{
  "password": "string"
}
```

### 4.7 获取用户角色
- **GET** `/api/v1/users/{id}/roles`

### 4.8 分配角色给用户
- **POST** `/api/v1/users/{id}/roles`
- **请求体**:
```json
{
  "role_id": 1,
  "domain_id": 1
}
```

### 4.9 撤销用户角色
- **DELETE** `/api/v1/users/{id}/roles?role_id={roleId}&domain_id={domainId}`

## 5. 用户组管理 (Groups)

### 5.1 获取用户组列表
- **GET** `/api/v1/groups`
- **查询参数**:
  - `limit`: 每页数量
  - `offset`: 偸移量

### 5.2 创建用户组
- **POST** `/api/v1/groups`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "domain_id": 1
}
```

### 5.3 更新用户组
- **PUT** `/api/v1/groups/{id}`
- **请求体**:
```json
{
  "description": "string"
}
```

### 5.4 删除用户组
- **DELETE** `/api/v1/groups/{id}`

### 5.5 获取用户组用户
- **GET** `/api/v1/groups/{id}/users`

### 5.6 添加用户到用户组
- **POST** `/api/v1/groups/{id}/users`
- **请求体**:
```json
{
  "user_id": 1
}
```

### 5.7 从用户组移除用户
- **DELETE** `/api/v1/groups/{id}/users?user_id={userId}`

## 6. 角色管理 (Roles)

### 6.1 获取角色列表
- **GET** `/api/v1/roles`
- **查询参数**:
  - `domain_id`: 域ID
  - `keyword`: 搜索关键字
  - `type`: 角色类型 (system/custom)
  - `enabled`: 状态过滤
  - `limit`: 每页数量
  - `offset`: 偸移量

### 6.2 创建角色
- **POST** `/api/v1/roles`
- **请求体**:
```json
{
  "name": "string",
  "display_name": "string",
  "description": "string",
  "domain_id": 1,
  "type": "string"
}
```

### 6.3 更新角色
- **PUT** `/api/v1/roles/{id}`
- **请求体**:
```json
{
  "name": "string",
  "display_name": "string",
  "description": "string",
  "domain_id": 1,
  "type": "string"
}
```

### 6.4 删除角色
- **DELETE** `/api/v1/roles/{id}`

### 6.5 启用/禁用角色
- **POST** `/api/v1/roles/{id}/enable` (启用)
- **POST** `/api/v1/roles/{id}/disable` (禁用)

### 6.6 公开角色
- **POST** `/api/v1/roles/{id}/public`

### 6.7 获取角色权限
- **GET** `/api/v1/roles/{id}/permissions`

### 6.8 分配权限给角色
- **POST** `/api/v1/roles/{id}/permissions`
- **请求体**:
```json
{
  "permission_id": 1
}
```

### 6.9 从角色撤销权限
- **DELETE** `/api/v1/roles/{id}/permissions?permission_id={permissionId}`

## 7. 权限管理 (Permissions)

### 7.1 获取权限列表
- **GET** `/api/v1/permissions`
- **查询参数**:
  - `resource`: 资源类型
  - `action`: 操作类型
  - `type`: 权限类型 (system/custom)
  - `keyword`: 搜索关键字
  - `limit`: 每页数量
  - `offset`: 偸移量

### 7.2 创建权限
- **POST** `/api/v1/permissions`
- **请求体**:
```json
{
  "name": "string",
  "display_name": "string",
  "resource": "string",
  "action": "string",
  "type": "string",
  "description": "string"
}
```

### 7.3 更新权限
- **PUT** `/api/v1/permissions/{id}`
- **请求体**:
```json
{
  "display_name": "string",
  "description": "string"
}
```

### 7.4 删除权限
- **DELETE** `/api/v1/permissions/{id}`

## 8. 策略管理 (Policies)

### 8.1 获取策略列表
- **GET** `/api/v1/policies`
- **查询参数**:
  - `scope`: 作用域 (system/domain/project)
  - `domain_id`: 域ID
  - `limit`: 每页数量
  - `offset`: 偸移量

### 8.2 创建策略
- **POST** `/api/v1/policies`
- **请求体**:
```json
{
  "name": "string",
  "scope": "string",
  "description": "string",
  "domain_id": "string",
  "policy": {},
  "is_system": false,
  "is_public": false
}
```

### 8.3 更新策略
- **PUT** `/api/v1/policies/{id}`
- **请求体**: 策略配置

### 8.4 删除策略
- **DELETE** `/api/v1/policies/{id}`

### 8.5 获取角色的策略列表
- **GET** `/api/v1/roles/{id}/policies`

### 8.6 分配策略给角色
- **POST** `/api/v1/roles/{id}/policies`
- **请求体**:
```json
{
  "policy_id": "string"
}
```

### 8.7 从角色撤销策略
- **DELETE** `/api/v1/roles/{id}/policies?policy_id={policyId}`

## 9. 认证源管理 (Auth Sources)

### 9.1 获取认证源列表
- **GET** `/api/v1/auth-sources`
- **查询参数**:
  - `keyword`: 搜索关键字
  - `type`: 类型 (ldap/local/sql)
  - `enabled`: 状态过滤
  - `limit`: 每页数量
  - `offset`: 偸移量

### 9.2 创建认证源
- **POST** `/api/v1/auth-sources`
- **请求体**:
```json
{
  "name": "string",
  "type": "string",
  "description": "string",
  "config": {},
  "enabled": true,
  "auto_create": true
}
```

### 9.3 更新认证源
- **PUT** `/api/v1/auth-sources/{id}`
- **请求体**:
```json
{
  "name": "string",
  "description": "string",
  "config": {},
  "enabled": true,
  "auto_create": true
}
```

### 9.4 删除认证源
- **DELETE** `/api/v1/auth-sources/{id}`

### 9.5 测试认证源
- **POST** `/api/v1/auth-sources/{id}/test`

### 9.6 启用/禁用认证源
- **POST** `/api/v1/auth-sources/{id}/enable` (启用)
- **POST** `/api/v1/auth-sources/{id}/disable` (禁用)