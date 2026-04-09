# openCMP 产品设计 Agent 规范指南

**日期**: 2026-04-05  
**状态**: 已启用  
**角色**: 产品设计 Agent（Product Design Agent）

---

## 1. 角色定位

产品设计 Agent 是多 Agent 开发流水线的**起点和质量把关节点**，负责：

- 把产品需求转化为可执行的技术规范
- 审核各模块的功能完整性与用户体验一致性
- 维护产品路线图和功能优先级
- 输出标准化 spec 文档供后端/前端 Agent 消费

---

## 2. 产品核心功能体系

### 2.1 认证与安全（IAM）

#### 认证源设计原则
| 认证源类型 | 适用范围 | 说明 |
|-----------|---------|------|
| 本地（local/sql） | system | 系统内置，不可删除，供系统管理员使用 |
| LDAP/OpenLDAP | system 或 domain | 可绑定到指定域，该域用户通过 LDAP 登录 |
| OIDC | system 或 domain | OAuth2/OpenID Connect 单点登录 |
| SAML | system 或 domain | 企业级 SSO 集成 |

#### 域隔离模型
```
系统（system）
  └── 域 A（Domain A）
  │     ├── 认证源：LDAP-A（domain scope, domain_id=A）
  │     ├── 用户：只属于域 A
  │     └── 资源：只展示域 A 下的云资源
  └── 域 B（Domain B）
        ├── 认证源：LDAP-B 或 本地
        └── 用户、资源同上
```

**关键约束：**
- `scope=domain` 的认证源登录的用户，JWT 中携带 `domain_id`
- 后续请求自动按 `domain_id` 过滤资源
- 跨域查看需要系统管理员权限

### 2.2 云资源管理

#### 资源类型矩阵
| 分类 | 子类 | 核心资源项 |
|------|------|-----------|
| 计算 | 虚拟机 | VM、镜像、密钥对、弹性伸缩组 |
| 网络 | 基础网络 | VPC、子网、路由表、安全组 |
| 网络 | 外网服务 | EIP、NAT 网关、负载均衡、CDN |
| 存储 | 块/对象/文件 | 云盘、快照、OSS 桶、NAS |
| 数据库 | 关系/缓存/文档 | RDS、Redis、MongoDB |
| 中间件 | 消息/搜索 | Kafka、Elasticsearch |

#### 多云适配器优先级
1. **阿里云**（Alibaba Cloud）- 主要云，优先完整
2. **腾讯云**（Tencent Cloud）- 次优先
3. **AWS** - 国际化扩展
4. **Azure** - 企业客户需求

### 2.3 消息中心

```
事件触发 → 路由规则 → 接收人匹配 → 渠道分发
                                    ├── 站内信（WebSocket/轮询）
                                    ├── 邮件（SMTP）
                                    ├── 企业微信机器人
                                    ├── 钉钉机器人
                                    └── 自定义 Webhook
```

### 2.4 多云同步

```
云账号配置
  └── 同步策略（按资源类型 + 时间计划）
        └── 资源同步规则（云标签 → openCMP 项目映射）
              └── 同步任务（cron 调度 + 日志记录）
```

---

## 3. 页面设计规范

### 3.1 列表页标准结构

```
[卡片头部：标题 + 新建按钮]
[筛选条件栏]
[数据表格]
  - 操作列：超过 3 个操作时使用"更多"下拉菜单
  - 危险操作（删除）用红色文字 + 分隔线 + Popconfirm 二次确认
  - 不可操作项用 disabled + tooltip 说明原因
[分页]
```

### 3.2 表单对话框规范

- 宽度：500px（简单）/ 700px（中等）/ 900px（复杂，如 LDAP 配置）
- 必填项：`prop` + `rules` 验证
- 类型切换时重置相关配置字段
- 提交按钮带 `:loading` 状态
- 取消时不需要确认

### 3.3 状态展示规范

| 状态类型 | 颜色方案 |
|---------|---------|
| 启用/成功 | `success`（绿色）|
| 禁用/停止 | `info`（灰色）|
| 警告/系统级 | `warning`（橙色）|
| 错误/危险 | `danger`（红色）|
| 普通标签 | `primary`（蓝色）|

---

## 4. 后端 API 设计规范

### 4.1 统一响应格式

```json
// 列表
{ "items": [...], "total": 100 }

// 单对象
{ "id": 1, "name": "...", ... }

// 操作成功
{ "message": "操作描述" }

// 错误
{ "error": "错误描述", "code": "ERROR_CODE" }
```

### 4.2 分页参数

- 列表接口统一使用 `page` + `page_size`（或 `limit` + `offset`）
- 默认 page_size: 20，最大: 100

### 4.3 认证中间件约定

- JWT 过期 → 401 + `{"code": "TOKEN_EXPIRED"}`
- JWT 无效 → 401 + `{"code": "TOKEN_INVALID"}`
- 权限不足 → 403

### 4.4 domain_id 上下文传播

通过 LDAP/域级认证源登录的用户，JWT 中携带 `domain_id` 字段。后端中间件解析后注入 Gin context：

```go
c.Set("domain_id", domainID)  // 用于资源查询过滤
c.Set("auth_source_id", authSourceID)
```

---

## 5. 认证源功能完成标准

| 功能点 | 验收标准 |
|--------|---------|
| 列表展示 | 域名正确显示（非 domain_id），更多操作下拉 |
| 创建 LDAP | 表单验证完整，config JSON 正确保存 |
| 连接测试 | TCP 连通测试，超时 5s 反馈 |
| 用户同步 | 安装 go-ldap/ldap/v3 后实现完整 bind+search |
| 域绑定登录 | LDAP 用户登录后 JWT 携带 domain_id，资源按域过滤 |
| 自动创建用户 | auto_create=true 时首次登录自动创建本地用户 |
| 默认认证源保护 | 本地认证源（type=local, scope=system）不可删除、不可禁用 |

---

## 6. 当前已知 TODO 项

| 模块 | TODO | 依赖 |
|------|------|------|
| 认证源 | LDAP 完整 bind 认证 | `go get github.com/go-ldap/ldap/v3` |
| 认证源 | LDAP 用户批量同步 | 同上 |
| 认证源 | OIDC 登录流程 | OAuth2 redirect flow |
| JWT | domain_id 注入中间件 | auth middleware 扩展 |
| 资源 | 按 domain_id 过滤查询 | 各资源 service 加 domain 过滤 |
| 多云同步 | 同步策略 + 定时任务 | cron 库 |

---

## 7. 下一阶段开发计划

### 第一优先级（当前迭代）
1. 安装 `go-ldap/ldap/v3`，实现完整 LDAP bind 认证
2. JWT 中间件扩展：解析 domain_id，注入 context
3. 资源查询加入 domain_id 过滤

### 第二优先级
4. 多云同步模块：同步策略设计 + 后端实现
5. 腾讯云适配器补全（主机 + 网络）

### 第三优先级
6. OIDC 认证源
7. AWS/Azure 适配器
8. 审计日志系统

---

## 8. 产品设计 Agent 使用说明

当用户说"启动产品设计 Agent"或"生成 XXX 模块的规范"时，该 Agent 应：

1. **阅读** `docs/PROGRESS.md` 了解当前进度
2. **阅读** 相关已有代码（model/service/handler）
3. **输出** 标准 spec 文档到 `docs/superpowers/specs/YYYY-MM-DD-<module>-design.md`，包含：
   - 功能描述与边界
   - 数据模型（字段表）
   - API 接口列表（路径、方法、请求/响应格式）
   - 前端页面结构（组件、交互、表单字段）
   - 验收标准
4. **更新** `docs/PROGRESS.md` 中对应模块的设计状态为 ✅
