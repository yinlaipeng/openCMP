# 云账户管理增强 - 设计方案

> 创建日期: 2026-04-14
> 设计版本: v1.0

## 一、功能概述

本次增强包含两个主要功能模块：

1. **更新云账号弹窗** - 支持编辑备注信息、密钥凭证，并提供测试连接验证
2. **属性设置-设置自动同步弹窗** - 包含8个子页面，提供云账户的完整属性管理

## 二、功能 1: 更新云账号弹窗

### 2.1 组件设计

**组件名称**: `EditAccountDialog.vue`

**组件位置**: `frontend/src/views/cloud-accounts/components/EditAccountDialog.vue`

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  更新云账号                                                          │
├────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  基本信息                                                            │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 名称:        [云账户名称] (只读)                              │  │
│  │ 平台:        [阿里云] (只读)                                  │  │
│  │ 状态:        [已连接] (只读)                                  │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  编辑信息                                                            │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 备注信息:    [________________] (textarea, 可编辑)            │  │
│  │ 密钥ID:      [________________] (input, 可编辑)               │  │
│  │              (Access Key ID)                                  │  │
│  │ 密钥密码:    [________________] (password, 可编辑)            │  │
│  │              (Access Key Secret)                              │  │
│  │                                                              │  │
│  │ [测试连接]  ✓ 连接成功，18个区域可用                          │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│                                              [取消]  [保存更改]      │
└────────────────────────────────────────────────────────────────────┘
```

### 2.2 字段规格

| 字段 | 类型 | 是否必填 | 是否可编辑 | 说明 |
|------|------|---------|-----------|------|
| 名称 | text | - | 否 | 只读显示，不可修改 |
| 平台 | tag | - | 否 | 只读显示云厂商类型 |
| 状态 | tag | - | 否 | 只读显示连接状态 |
| 备注信息 | textarea | 否 | 是 | 最大500字符 |
| 密钥ID | input | 是 | 是 | Access Key ID |
| 密钥密码 | password | 是 | 是 | Access Key Secret，显示时隐藏 |

### 2.3 测试连接按钮行为

**触发条件**: 用户输入密钥ID和密钥密码后点击

**API调用**: `POST /cloud-accounts/:id/test-connection-with-credentials`

**请求参数**:
```json
{
  "access_key_id": "LTAI...",
  "access_key_secret": "..."
}
```

**响应格式**:
```json
{
  "connected": true,
  "message": "连接成功，18个区域可用",
  "regions": ["cn-hangzhou", "cn-shanghai", ...]
}
```

**状态显示**:
- 测试中: 按钮 loading 状态，文字显示"测试中..."
- 成功: 绿色文字"✓ 连接成功，XX个区域可用"
- 失败: 红色文字"✗ 连接失败: [错误信息]"

### 2.4 保存行为

**验证规则**:
- 密钥ID: 必填，长度 >= 10
- 密钥密码: 必填，长度 >= 10
- 备注信息: 可选，最大500字符

**API调用**: `PUT /cloud-accounts/:id`

**请求参数**:
```json
{
  "description": "备注信息",
  "credentials": {
    "access_key_id": "LTAI...",
    "access_key_secret": "...",
    "region_id": "cn-hangzhou"
  }
}
```

---

## 三、功能 2: 属性设置-设置自动同步弹窗

### 3.1 主弹窗框架设计

**组件名称**: `CloudAccountDetailDialog.vue`

**组件位置**: `frontend/src/views/cloud-accounts/components/CloudAccountDetailDialog.vue`

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  云账户属性设置 - [云账户名称]                                      │
├────────────────────────────────────────────────────────────────────┤
│  ┌─────────┬─────────┬─────────┬─────────┬─────────┬─────────┬────┐│
│  │ 详情    │ 资源统计 │ 订阅    │ 云用户  │ 云用户组│ 云上项目│... ││
│  └─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴────┘│
│                                                                      │
│  [当前选中 Tab 的内容区域]                                           │
│                                                                      │
│                                              [关闭]                  │
└────────────────────────────────────────────────────────────────────┘
```

**Tabs 顺序**:
1. 详情 (Detail)
2. 资源统计 (Resource Stats)
3. 订阅 (Subscriptions)
4. 云用户 (Cloud Users)
5. 云用户组 (Cloud User Groups)
6. 云上项目 (Cloud Projects)
7. 定时任务 (Scheduled Tasks)
8. 操作日志 (Operation Logs)

---

## 四、子页面详细设计

### 4.1 详情子页面

**Tab 名称**: 详情

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  基本信息                                                            │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 名称:          云账户名称                                     │  │
│  │ 平台:          阿里云                                         │  │
│  │ 状态:          已连接                                         │  │
│  │ 启用状态:      启用                                           │  │
│  │ 健康状态:      正常                                           │  │
│  │ 创建时间:      2026-04-14 10:30:00                            │  │
│  │ 上次同步:      2026-04-14 09:00:00                            │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  账号信息                                                            │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 账号ID:        [阿里云账号ID]                                  │  │
│  │ 账户余额:      ¥ 1,234.56                                     │  │
│  │ 默认区域:      cn-hangzhou                                    │  │
│  │ 所属域:        默认域                                         │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  云平台权限                                                          │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ ┌────────────────┬────────────────┬────────────┐             │  │
│  │ │ 权限名称       │ 权限描述       │ 状态       │             │  │
│  │ ├────────────────┼────────────────┼────────────┤             │  │
│  │ │ AliyunECSFull  │ ECS完全访问   │ ✓ 已授权   │             │  │
│  │ │ AliyunVPCFull  │ VPC完全访问   │ ✓ 已授权   │             │  │
│  │ │ AliyunRDSRead  │ RDS只读      │ ✗ 未授权   │             │  │
│  │ └────────────────┴────────────────┴────────────┘             │  │
│  │ [刷新权限]                                                   │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────┘
```

**API 端点**:
- `GET /cloud-accounts/:id` - 获取基本信息
- `GET /cloud-accounts/:id/permissions` - 获取云平台权限列表

---

### 4.2 资源统计子页面

**Tab 名称**: 资源统计

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  资源概览 (统计卡片网格)                                             │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┬─────────┐│
│  │ 虚拟机   │ RDS实例  │ Redis实例│ 存储桶   │ EIP      │ VPC     ││
│  │    12    │     3    │     2    │     5    │     8    │    4   ││
│  ├──────────┼──────────┼──────────┼──────────┼──────────┼─────────┤│
│  │ 公网IP   │ 快照     │ IP子网   │ IP总量   │          │         ││
│  │     6    │    15    │    10    │   256    │          │         ││
│  └──────────┴──────────┴──────────┴──────────┴──────────┴─────────┘│
│                                                                      │
│  使用率统计                                                          │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 虚拟机开机率:  ████████████░░░░░░░░  60%  (12台/20台)        │  │
│  │ 磁盘挂载率:    ████████████████░░░░  80%  (16块/20块)        │  │
│  │ EIP使用率:     ██████████████░░░░░░  70%  (7个/10个)         │  │
│  │ IP使用率:      ████████░░░░░░░░░░░░  40%  (102/256)          │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  [刷新统计]                                                          │
└────────────────────────────────────────────────────────────────────┘
```

**统计指标**:
| 指标 | 数据来源 | 计算方式 |
|------|---------|---------|
| 虚拟机数量 | ListVMs API | count |
| RDS实例 | (待实现) | count |
| Redis实例 | (待实现) | count |
| 存储桶数量 | (待实现) | count |
| EIP数量 | ListEIPs API | count |
| 公网IP数量 | ListEIPs API (public类型) | count |
| 快照数量 | (待实现) | count |
| VPC数量 | ListVPCs API | count |
| IP子网数量 | ListSubnets API | count |
| IP总量 | ListSubnets CIDR计算 | sum |
| 虚拟机开机率 | GetVMStatus | running/total |
| 磁盘挂载率 | ListDisks | attached/total |
| EIP使用率 | ListEIPs | bound/total |
| IP使用率 | Subnet CIDR | used/total |

**API 端点**:
- `GET /cloud-accounts/:id/resource-stats` - 获取资源统计数据

---

### 4.3 订阅子页面

**Tab 名称**: 订阅

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  [新建订阅]                                                          │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 名称 │ 订阅ID │ 启用 │ 状态 │ 同步时间 │ 耗时 │ 同步状态 │ 域 ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ 订阅1 │ sub-001 │ ✓ │ 正常 │ 2026-04-14 │ 12s │ 同步完成 │ 默认域 ││
│  │ 订阅2 │ sub-002 │ ✓ │ 正常 │ 2026-04-13 │ 8s  │ 同步失败 │ 默认域 ││
│  └────────────────────────────────────────────────────────────────┘│
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 所属域: [默认域] │ 默认项目: [下拉选择] │ 操作               ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ 默认域  │ 项目A      │ [更改项目] [同步策略] [更多▼]          ││
│  │ 默认域  │ 项目B      │ [更改项目] [同步策略] [更多▼]          ││
│  └────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  更多下拉菜单:                                                       │
│  ├─ 同步策略设置                                                     │
│  ├─ 启用                                                             │
│  ├─ 禁用                                                             │
│  └─ 删除                                                             │
└────────────────────────────────────────────────────────────────────┘
```

**数据模型** (新增):
```go
type CloudSubscription struct {
    ID                  uint           `gorm:"primaryKey" json:"id"`
    CloudAccountID      uint           `gorm:"index;not null" json:"cloud_account_id"`
    Name                string         `gorm:"size:200;not null" json:"name"`
    SubscriptionID      string         `gorm:"size:100;not null" json:"subscription_id"`
    Enabled             bool           `gorm:"default:true" json:"enabled"`
    Status              string         `gorm:"size:20;default:'normal'" json:"status"`
    SyncTime            *time.Time     `json:"sync_time"`
    SyncDuration        int            `gorm:"default:0" json:"sync_duration"` // 秒
    SyncStatus          string         `gorm:"size:20;default:'completed'" json:"sync_status"`
    DomainID            uint           `gorm:"index;default:1" json:"domain_id"`
    DefaultProjectID    *uint          `json:"default_project_id"`
    SyncPolicyID        *uint          `json:"sync_policy_id"`
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**API 端点**:
- `GET /cloud-accounts/:id/subscriptions` - 订阅列表
- `POST /cloud-accounts/:id/subscriptions` - 创建订阅
- `PUT /cloud-accounts/:id/subscriptions/:sid` - 更新订阅
- `DELETE /cloud-accounts/:id/subscriptions/:sid` - 删除订阅
- `POST /cloud-accounts/:id/subscriptions/:sid/sync` - 同步订阅
- `PUT /cloud-accounts/:id/subscriptions/:sid/project` - 更改项目
- `PUT /cloud-accounts/:id/subscriptions/:sid/policy` - 设置同步策略

---

### 4.4 云用户子页面

**Tab 名称**: 云用户

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  [新建云用户]                                                        │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 用户名 │ 控制台登录 │ 状态 │ 密码 │ 登录地址 │ 关联本地用户 │ 操作 ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ admin  │ ✓         │ 正常 │ ●●●  │ https://.. │ 用户A       │ [编辑][删除] ││
│  │ ops    │ ✗         │ 正常 │ ●●●  │ -          │ -          │ [编辑][删除] ││
│  └────────────────────────────────────────────────────────────────┘│
└────────────────────────────────────────────────────────────────────┘
```

**数据模型** (新增):
```go
type CloudUser struct {
    ID                uint           `gorm:"primaryKey" json:"id"`
    CloudAccountID    uint           `gorm:"index;not null" json:"cloud_account_id"`
    Username          string         `gorm:"size:100;not null" json:"username"`
    ConsoleLogin      bool           `gorm:"default:false" json:"console_login"`
    Status            string         `gorm:"size:20;default:'normal'" json:"status"`
    Password          string         `gorm:"size:255" json:"-"` // 加密存储，不返回前端
    LoginURL          string         `gorm:"size:255" json:"login_url"`
    LocalUserID       *uint          `json:"local_user_id"`
    Platform          string         `gorm:"size:20" json:"platform"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
    DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**API 端点**:
- `GET /cloud-accounts/:id/cloud-users` - 云用户列表
- `POST /cloud-accounts/:id/cloud-users` - 创建云用户
- `PUT /cloud-accounts/:id/cloud-users/:uid` - 更新云用户
- `DELETE /cloud-accounts/:id/cloud-users/:uid` - 删除云用户

---

### 4.5 云用户组子页面

**Tab 名称**: 云用户组

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  [新建云用户组]                                                      │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 名称   │ 状态 │ 权限     │ 平台   │ 所属云账号 │ 所属域 │ 操作 ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ Admins │ 正常 │ 管理员   │ 阿里云 │ 账户A     │ 默认域 │ [编辑][删除] ││
│  │ DevOps │ 正常 │ 运维     │ 阿里云 │ 账户A     │ 默认域 │ [编辑][删除] ││
│  └────────────────────────────────────────────────────────────────┘│
└────────────────────────────────────────────────────────────────────┘
```

**数据模型** (新增):
```go
type CloudUserGroup struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`
    Name            string         `gorm:"size:200;not null" json:"name"`
    Status          string         `gorm:"size:20;default:'normal'" json:"status"`
    Permissions     datatypes.JSON `gorm:"type:json" json:"permissions"`
    Platform        string         `gorm:"size:20" json:"platform"`
    DomainID        uint           `gorm:"index;default:1" json:"domain_id"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**API 端点**:
- `GET /cloud-accounts/:id/cloud-user-groups` - 云用户组列表
- `POST /cloud-accounts/:id/cloud-user-groups` - 创建云用户组
- `PUT /cloud-accounts/:id/cloud-user-groups/:gid` - 更新云用户组
- `DELETE /cloud-accounts/:id/cloud-user-groups/:gid` - 删除云用户组

---

### 4.6 云上项目子页面

**Tab 名称**: 云上项目

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  [同步云上项目]                                                      │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 云上项目 │ 订阅 │ 状态 │ 标签 │ 所属域 │ 本地项目 │ 优先级 │ 操作 ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ proj-001 │ sub-1 │ 正常 │ env:prod │ 默认域 │ 项目A │ 1 │ [编辑][映射] ││
│  │ proj-002 │ sub-1 │ 正常 │ env:dev  │ 默认域 │ 项目B │ 2 │ [编辑][映射] ││
│  └────────────────────────────────────────────────────────────────┘│
└────────────────────────────────────────────────────────────────────┘
```

**数据模型** (新增):
```go
type CloudProject struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`
    Name            string         `gorm:"size:200;not null" json:"name"`
    SubscriptionID  *uint          `json:"subscription_id"`
    Status          string         `gorm:"size:20;default:'normal'" json:"status"`
    Tags            datatypes.JSON `gorm:"type:json" json:"tags"`
    DomainID        uint           `gorm:"index;default:1" json:"domain_id"`
    LocalProjectID  *uint          `json:"local_project_id"`
    Priority        int            `gorm:"default:0" json:"priority"`
    SyncTime        *time.Time     `json:"sync_time"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**API 端点**:
- `GET /cloud-accounts/:id/cloud-projects` - 云上项目列表
- `POST /cloud-accounts/:id/cloud-projects/sync` - 从云平台同步项目
- `PUT /cloud-accounts/:id/cloud-projects/:pid` - 更新云上项目
- `PUT /cloud-accounts/:id/cloud-projects/:pid/map` - 映射到本地项目
- `DELETE /cloud-accounts/:id/cloud-projects/:pid` - 删除云上项目

---

### 4.7 定时任务子页面

**Tab 名称**: 定时任务

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  [新建定时任务]                                                      │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ 名称   │ 状态 │ 启用状态 │ 操作动作 │ 策略详情 │ 操作         ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ 每日同步 │ 正常 │ ✓      │ 同步资源 │ 每日02:00 │ [编辑][执行][更多] ││
│  │ 每周统计 │ 正常 │ ✓      │ 生成报告 │ 每周一    │ [编辑][执行][更多] ││
│  └────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  更多下拉菜单:                                                       │
│  ├─ 启用                                                             │
│  ├─ 禁用                                                             │
│  └─ 删除                                                             │
└────────────────────────────────────────────────────────────────────┘
```

**使用现有模型**: ScheduledTask（已有 cloud_account_id 字段）

**API 端点** (复用现有):
- `GET /scheduled-tasks?cloud_account_id=:id` - 定时任务列表
- `POST /scheduled-tasks` - 创建定时任务
- `PUT /scheduled-tasks/:tid` - 更新定时任务
- `POST /scheduled-tasks/:tid/execute` - 立即执行
- `DELETE /scheduled-tasks/:tid` - 删除定时任务

---

### 4.8 操作日志子页面

**Tab 名称**: 操作日志

**UI 结构**:
```
┌────────────────────────────────────────────────────────────────────┐
│  筛选: [操作类型▼] [结果▼] [时间范围]                               │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │ #ID │ 操作时间 │ 资源名称 │ 资源类型 │ 操作类型 │ 结果 │ 发起人 │ 操作 ││
│  ├────────────────────────────────────────────────────────────────┤│
│  │ 001 │ 10:30:00 │ VM-001  │ 虚拟机  │ 创建    │ 成功 │ admin │ [查看] ││
│  │ 002 │ 10:31:00 │ VPC-001 │ 网络    │ 删除    │ 成功 │ ops   │ [查看] ││
│  └────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  分页: [◀ 1 2 3 ... 10 ▶]                                           │
│                                                                      │
│  日志详情弹窗:                                                       │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │ 操作日志详情                                                 │  │
│  ├──────────────────────────────────────────────────────────────┤  │
│  │ #ID:            001                                          │  │
│  │ 操作时间:       2026-04-14 10:30:00                          │  │
│  │ 资源名称:       VM-001                                       │  │
│  │ 资源类型:       虚拟机                                       │  │
│  │ 操作类型:       创建                                         │  │
│  │ 服务类型:       compute                                      │  │
│  │ 风险级别:       低                                           │  │
│  │ 事件类型:       API调用                                      │  │
│  │ 结果:           成功                                         │  │
│  │ 发起人:         admin                                        │  │
│  │ 所属项目:       项目A                                        │  │
│  │ 请求详情:       {...}                                        │  │
│  │ 响应详情:       {...}                                        │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────┘
```

**扩展现有模型**: OperationLog
```go
// 新增字段
CloudAccountID *uint  `json:"cloud_account_id,omitempty" gorm:"column:cloud_account_id;index"`
EventType      string `json:"event_type" gorm:"column:event_type;size:50;default:'api_call'"`
RequestDetail  datatypes.JSON `json:"request_detail,omitempty" gorm:"type:json"`
ResponseDetail datatypes.JSON `json:"response_detail,omitempty" gorm:"type:json"`
```

**API 端点**:
- `GET /cloud-accounts/:id/operation-logs` - 操作日志列表
- `GET /operation-logs/:lid` - 操作日志详情

---

## 五、后端 API 规范

### 5.1 API 端点汇总

| API | 方法 | 说明 | 需新增 |
|------|------|------|--------|
| `/cloud-accounts/:id` | GET | 云账户详情 | 否 |
| `/cloud-accounts/:id` | PUT | 更新云账户 | 否 |
| `/cloud-accounts/:id/test-connection` | POST | 测试连接 | 否 |
| `/cloud-accounts/:id/test-connection-with-credentials` | POST | 带凭证测试连接 | **是** |
| `/cloud-accounts/:id/resource-stats` | GET | 资源统计 | **是** |
| `/cloud-accounts/:id/permissions` | GET | 权限列表 | **是** |
| `/cloud-accounts/:id/subscriptions` | GET/POST | 订阅 CRUD | **是** |
| `/cloud-accounts/:id/subscriptions/:sid` | PUT/DELETE | 订阅 CRUD | **是** |
| `/cloud-accounts/:id/cloud-users` | GET/POST | 云用户 CRUD | **是** |
| `/cloud-accounts/:id/cloud-users/:uid` | PUT/DELETE | 云用户 CRUD | **是** |
| `/cloud-accounts/:id/cloud-user-groups` | GET/POST | 云用户组 CRUD | **是** |
| `/cloud-accounts/:id/cloud-user-groups/:gid` | PUT/DELETE | 云用户组 CRUD | **是** |
| `/cloud-accounts/:id/cloud-projects` | GET | 云上项目列表 | **是** |
| `/cloud-accounts/:id/cloud-projects/sync` | POST | 同步云上项目 | **是** |
| `/cloud-accounts/:id/cloud-projects/:pid` | PUT/DELETE | 云上项目 CRUD | **是** |
| `/cloud-accounts/:id/operation-logs` | GET | 操作日志列表 | **是** |

### 5.2 API 请求/响应规范

**测试连接（带凭证）**
```typescript
// Request
POST /cloud-accounts/:id/test-connection-with-credentials
{
  "access_key_id": string,
  "access_key_secret": string
}

// Response
{
  "connected": boolean,
  "message": string,
  "regions": string[]  // 可选，成功时返回区域列表
}
```

**资源统计**
```typescript
// Response
{
  "resources": {
    "vms": number,
    "rds": number,
    "redis": number,
    "buckets": number,
    "eips": number,
    "public_ips": number,
    "snapshots": number,
    "vpcs": number,
    "subnets": number,
    "total_ips": number
  },
  "usage_rates": {
    "vm_running_rate": number,  // 0-100
    "disk_mounted_rate": number,
    "eip_bound_rate": number,
    "ip_used_rate": number
  }
}
```

---

## 六、前端组件文件结构

```
frontend/src/
├── views/cloud-accounts/
│   ├── index.vue                           # 云账户列表页（修改）
│   └── components/
│       ├── EditAccountDialog.vue           # 更新云账号弹窗（新建）
│       ├── CloudAccountDetailDialog.vue    # 属性设置主弹窗（新建）
│       └── tabs/
│           ├── DetailTab.vue               # 详情子页面（新建）
│           ├── ResourceStatsTab.vue        # 资源统计（新建）
│           ├── SubscriptionTab.vue         # 订阅（新建）
│           ├── CloudUserTab.vue            # 云用户（新建）
│           ├── CloudUserGroupTab.vue       # 云用户组（新建）
│           ├── CloudProjectTab.vue         # 云上项目（新建）
│           ├── ScheduledTaskTab.vue        # 定时任务（新建）
│           └── OperationLogTab.vue         # 操作日志（新建）
├── api/
│   ├── cloud-account.ts                    # 修改：新增API函数
│   ├── cloud-subscription.ts               # 新建
│   ├── cloud-user.ts                       # 新建
│   └── cloud-user-group.ts                 # 新建
├── types/
│   ├── cloud-account.ts                    # 修改：新增类型
│   ├── cloud-subscription.ts               # 新建
│   ├── cloud-user.ts                       # 新建
│   └── cloud-user-group.ts                 # 新建
```

---

## 七、后端文件结构

```
backend/
├── internal/
│   ├── model/
│   │   ├── cloud_account.go                # 修改：扩展字段
│   │   ├── cloud_subscription.go           # 新建
│   │   ├── cloud_user.go                   # 新建
│   │   ├── cloud_user_group.go             # 新建
│   │   ├── cloud_project.go                # 新建
│   │   ├── operation_log.go                # 修改：扩展字段
│   │   └── scheduled_task.go               # 已有，确认字段
│   ├── handler/
│   │   ├── cloud_account.go                # 修改：新增端点
│   │   ├── cloud_subscription.go           # 新建
│   │   ├── cloud_user.go                   # 新建
│   │   ├── cloud_user_group.go             # 新建
│   │   ├── cloud_project.go                # 新建
│   │   └── operation_log.go                # 修改：新增端点
│   ├── service/
│   │   ├── cloud_account.go                # 修改：新增方法
│   │   ├── cloud_subscription.go           # 新建
│   │   ├── cloud_user.go                   # 新建
│   │   ├── cloud_user_group.go             # 新建
│   │   ├── cloud_project.go                # 新建
│   │   └── operation_log.go                # 修改：新增方法
├── cmd/server/
│   └── main.go                             # 修改：注册新路由和模型
```

---

## 八、实施优先级

| 阶段 | 任务 | 复杂度 | 依赖 |
|------|------|--------|------|
| **阶段1** | 更新云账号弹窗 | 低 | 无 |
| **阶段1** | 属性设置主弹窗框架 | 低 | 无 |
| **阶段1** | 详情子页面 | 中 | 权限API |
| **阶段1** | 资源统计子页面 | 中 | 统计API |
| **阶段2** | 操作日志子页面 | 低 | 扩展字段 |
| **阶段2** | 定时任务子页面 | 低 | 已有模型 |
| **阶段3** | 订阅子页面 | 高 | 新模型 |
| **阶段3** | 云用户子页面 | 高 | 新模型 |
| **阶段3** | 云用户组子页面 | 高 | 新模型 |
| **阶段3** | 云上项目子页面 | 高 | 新模型 |

---

## 九、风险与注意事项

1. **凭证安全**: Access Key Secret 需加密存储，前端传输后需验证
2. **权限API**: 云厂商权限查询API可能有限制或需特殊权限
3. **资源统计性能**: 统计API需调用多个云厂商API，可能耗时较长
4. **数据一致性**: 新增模型需考虑与现有IAM模型的关系
5. **向后兼容**: 扩展字段需确保不影响现有功能

---

## 十、确认清单

请确认以下设计点：

- [ ] 更新云账号弹窗的字段是否满足需求？
- [ ] 属性设置弹窗的8个子页面是否完整？
- [ ] 资源统计指标是否需要调整？
- [ ] 新增数据模型字段是否合理？
- [ ] API端点设计是否满足前端需求？
- [ ] 实施优先级是否合适？
- [ ] 是否有其他特殊需求需要补充？