# openCMP 多 Agent 持续开发设计文档

**日期**: 2026-04-01
**状态**: 已批准
**适用项目**: openCMP 多云管理平台

---

## 1. 概述

本文档定义 openCMP 项目的多 Agent 并行开发体系，采用**文档驱动流水线**模式，将产品设计、任务分发、后端开发、前端开发、测试验收五个角色解耦，实现两条业务线并行推进。

---

## 2. Agent 角色与职责

| Agent | 角色 | 输入 | 输出 |
|-------|------|------|------|
| **产品设计 agent** | 负责模块设计，输出可执行规范 | 用户需求、参考截图、现有代码 | `docs/specs/<模块>-design.md` |
| **任务分发 agent** | 读取设计文档，拆解为后端/前端可执行任务 | spec 文档 | `docs/tasks/<模块>-tasks.md` |
| **后端开发 agent** | 实现 Go 后端（handler/service/model/provider） | tasks 文档 + 现有代码 | Go 代码，更新任务状态 |
| **前端开发 agent** | 实现 Vue 3 前端页面和 API 调用 | tasks 文档 + 现有代码 | Vue/TS 代码，更新任务状态 |
| **测试 agent** | 验收后端接口和前端功能 | 完成的任务 + 代码 | 测试报告，问题反馈 |

### 文档目录结构

```
docs/
├── superpowers/specs/        # 产品设计 agent 输出
│   ├── iam-design.md
│   ├── message-center-design.md
│   └── multicloud-sync-design.md
├── tasks/                    # 任务分发 agent 输出
│   ├── iam-tasks.md
│   ├── message-center-tasks.md
│   └── multicloud-sync-tasks.md
└── PROGRESS.md               # 全局进度看板
```

### 流水线节奏

```
产品设计 agent 输出 spec
        ↓
任务分发 agent 拆解任务（后端 + 前端分离）
        ↓（同步开始）
后端 agent                    前端 agent
  └ model/migration             └ views/api
  └ service                     └ store/router
  └ handler/router              └ components
        ↓（后端先完成接口，前端联调）
测试 agent 验收
        ↓
下一个模块
```

---

## 3. 两条业务线并行

### 线 1：IAM 核心模块

```
阶段 1（基础）:  域(Domain) → 项目(Project) → 组(Group)
阶段 2（身份）:  用户(User) → 角色(Role) → 权限(Permission)
阶段 3（接入）:  完善认证源(Auth Source) → 基于项目的资源划分
```

### 线 2：平台功能模块

```
消息中心:   站内信 → 消息订阅 → 通知渠道 → 接受人管理 → 机器人管理
多云管理:   云账号增强 → 同步策略 → 资源同步规则(标签→项目) → 定时同步任务
```

### 云厂商适配器补全（后端 agent 间隙推进）

- 腾讯云适配器
- AWS 适配器
- Azure 适配器

---

## 4. 云厂商资源统一分类体系

### 4.1 主机

| 子类 | 资源项 |
|------|--------|
| 虚拟机 | 虚拟机、主机模版、弹性伸缩组 |
| 镜像 | 系统镜像 |
| 密钥 | 密钥对 |
| 主机存储 | 硬盘、硬盘快照、主机快照 |

### 4.2 网络

| 子类 | 资源项 |
|------|--------|
| 地域 | 区域、可用区 |
| 基础网络 | VPC互联、VPC对等链接、全局VPC、路由表、二层网络、IP子网 |
| 网络服务 | 弹性公网IP、NAT网关、DNS解析、IPv6网关 |
| 负载均衡 | 实例、访问控制、证书 |
| 内容分发网络 | CDN域名 |

### 4.3 存储

| 子类 | 资源项 |
|------|--------|
| 块存储 | 块存储 |
| 对象存储 | 存储桶 |
| 文件存储 | 文件系统、NAS权限组 |

### 4.4 数据库

| 子类 | 资源项 |
|------|--------|
| 关系型 | RDS实例 |
| 缓存 | Redis实例 |
| 文档型 | MongoDB实例 |

### 4.5 中间件

| 子类 | 资源项 |
|------|--------|
| 消息队列 | Kafka |
| 数据分析 | Elasticsearch |

---

## 5. IAM 模块设计

### 5.1 数据模型层级关系

```
域 (Domain)
  └── 项目 (Project)
        └── 资源归属（云资源按项目划分）
  └── 用户 (User)
        └── 用户组 (Group)
  └── 角色 (Role)
        └── 权限集合 (Permissions)
  └── 认证源 (Auth Source) → LDAP / OIDC / SAML
```

### 5.2 各模块职责

| 模块 | 核心功能 |
|------|---------|
| **域** | 租户隔离顶层单位；域管理员、域配置、域间切换 |
| **项目** | 域内资源分组；项目成员、项目配额、项目资源视图 |
| **用户** | 创建/启停/删除；归属域、归属组、绑定角色 |
| **组** | 用户批量管理；组绑定角色、组归属项目 |
| **角色** | 权限集合模板；系统角色 + 自定义角色 |
| **权限** | 细粒度资源操作权限；格式：`<模块>:<资源>:<动作>` |
| **认证源** | 完善 LDAP/OIDC/SAML 配置；连接测试；用户同步 |

### 5.3 权限格式

```
compute:vm:list           # 查看虚拟机列表
compute:vm:create         # 创建虚拟机
compute:vm:delete         # 删除虚拟机
network:vpc:list          # 查看VPC列表
network:vpc:create        # 创建VPC
network:vpc:delete        # 删除VPC
storage:disk:list         # 查看硬盘
iam:user:create           # 创建用户
iam:user:disable          # 禁用用户
iam:role:assign           # 分配角色
message:notify:send       # 发送通知
```

### 5.4 基于项目的资源划分

- 每个云资源归属一个项目
- 项目成员只能查看/操作本项目资源
- 项目管理员可配置项目配额（VM数量上限、磁盘容量上限等）
- 云资源同步时可通过标签自动归入对应项目

---

## 6. 消息中心模块设计

### 6.1 模块结构

```
消息中心
  ├── 站内信        → 系统通知/告警/任务完成推送给用户
  ├── 消息订阅      → 用户自选关注的事件类型
  ├── 通知渠道      → 邮件 / 企业微信 / 钉钉 / Webhook
  ├── 接受人管理    → 通知目标（用户/用户组/角色）
  └── 机器人管理    → 企业微信机器人 / 钉钉机器人 Token 配置
```

### 6.2 核心数据流

```
系统事件触发（VM创建完成、同步失败、登录告警等）
        ↓
消息路由（查询该事件的订阅规则）
        ↓
匹配接受人（用户/组/角色）
        ↓
按通知渠道分发
    ├── 站内信 → 存 DB，前端轮询/WebSocket
    ├── 邮件   → SMTP 发送
    ├── 企微/钉钉机器人 → Webhook POST
    └── 自定义 Webhook → HTTP POST
```

### 6.3 关键数据模型

| 实体 | 关键字段 |
|------|---------|
| 消息(Message) | id, title, content, type, level(info/warn/error), read_at, user_id, created_at |
| 订阅规则(Subscription) | event_type, receiver_type(user/group/role), receiver_id, channels, enabled |
| 通知渠道(NotifyChannel) | name, type(email/wechat/dingtalk/webhook), config(JSON), enabled |
| 机器人(Robot) | name, type(wechat/dingtalk), webhook_url, secret, enabled |

---

## 7. 多云管理模块设计

### 7.1 云账号增强

- 支持多云账号管理（阿里云/腾讯云/AWS/Azure）
- 账号连接状态监控
- 账号级别权限隔离（归属域/项目）

### 7.2 同步策略

```
同步策略配置
  ├── 同步范围：按资源类型勾选（主机/网络/存储/数据库/中间件）
  ├── 同步规则：云上资源标签 → 映射到 openCMP 项目
  │       例：tag: project=frontend → 归入 "frontend" 项目
  │       例：tag: env=prod → 归入 "production" 项目
  └── 同步时机
        ├── 手动触发（立即同步）
        └── 定时任务（cron 表达式，如 0 * * * * 每小时同步）
```

### 7.3 定时同步任务

- 基于 cron 表达式配置同步周期
- 支持按云账号、按资源类型独立配置
- 同步日志记录（成功/失败/变更详情）
- 同步失败触发消息中心告警

---

## 8. 全局进度看板（PROGRESS.md 格式）

每个模块完成后更新 `docs/PROGRESS.md`：

```markdown
| 模块 | 设计 | 后端 | 前端 | 测试 | 状态 |
|------|------|------|------|------|------|
| 域   | ✅   | 🚧   | ⬜   | ⬜   | 开发中 |
| 项目 | ✅   | ⬜   | ⬜   | ⬜   | 待开发 |
...
```

---

## 9. 开发优先级

### 第一批（IAM 基础）
1. 域 (Domain)
2. 项目 (Project)
3. 组 (Group)
4. 用户 (User)
5. 角色 (Role)
6. 权限 (Permission)

### 第二批（平台功能）
7. 消息中心（站内信 → 通知渠道 → 订阅 → 机器人）
8. 完善认证源（LDAP/OIDC/SAML）
9. 基于项目的资源划分

### 第三批（多云增强）
10. 云账号增强
11. 同步策略 + 资源同步规则
12. 定时同步任务

### 云厂商适配器（持续补全）
- 腾讯云：主机/网络/存储/数据库/中间件
- AWS：主机/网络/存储/数据库/中间件
- Azure：主机/网络/存储/数据库/中间件
