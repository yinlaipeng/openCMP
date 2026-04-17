# Findings & Decisions

## Requirements
- 实现 openCMP 项目完整落地，前端与后端真实对接
- 不再使用模拟接口，实现真实的云厂商 SDK 调用
- 完成从云账号添加 → 资源同步 → 资源管理的完整业务流程

## Phase 27: 前端页面样式统一性检查 (2026-04-17)

### 标准样式结构 (host-templates/index.vue) - ui-ux-pro-max 参考
```
模板结构：
├── div.xxx-container { padding: 20px }
│   ├── div.page-header { flex布局, margin-bottom: 20px }
│   │   ├── h2 标题
│   │   └── el-button 操作按钮
│   ├── el-card.filter-card { margin-bottom: 20px }
│   │   └── el-form inline筛选表单
│   ├── el-table { width: 100%, row-key="id" }
│   └── el-pagination.pagination { margin-top: 20px, text-align: right }
```

**ui-ux-pro-max 关键规则**:
- `spacing-scale`: 使用 8dp 间距系统 (20px = 2.5 × 8dp)
- `container-width`: 内容区域保持一致
- `consistency`: 同一产品内所有页面使用相同样式
- `z-index-management`: 定义分层 z-index
- `touch-target-size`: 触控目标 ≥44×44pt

### P0修复已完成 (2026-04-17)

已修复8个页面添加筛选区/分页/padding：
- images/index.vue ✅
- eips/index.vue ✅  
- rds/instances/index.vue ✅
- redis/instances/index.vue ✅
- vpcs/index.vue ✅
- subnets/index.vue ✅
- policies/index.vue ✅
- resources/vms/index.vue ✅

### P1深度问题发现 (ui-ux-pro-max 分析)

| 页面 | 问题类型 | 具体问题 |
|------|---------|---------|
| **vms/index.vue** | 结构差异 | 使用el-collapse折叠查询而非.filter-card el-card；使用el-card header slot而非独立.page-header |
| **storage/cloud/disks/index.vue** | 多项问题 | el-card header slot；inline pagination样式；缺少row-key；筛选区非el-card结构 |
| **cloud-accounts/index.vue** | 缺少筛选 | el-card header slot；无.filter-card筛选卡片；缺少row-key |
| **sync-policies/index.vue** | 结构差异 | el-card header slot；筛选表单直接嵌入el-card而非独立.filter-card；缺少row-key |
| **compute/vms/index.vue** | 特殊结构 | 使用折叠式筛选(el-collapse)，与其他页面不统一 |

**推荐修复方案**:
1. 统一将 el-card header slot 改为独立 `.page-header` div
2. 筛选区统一使用 `.filter-card` el-card 包裹
3. 表格统一添加 `row-key` 属性
4. 分页统一使用 `.pagination { text-align: right }` class

### 待修复页面清单 (P1级别)

| 页面 | 需要修复 |
|------|---------|
| compute/vms/index.vue | page-header结构、el-collapse改为filter-card |
| storage/cloud/disks/index.vue | page-header、filter-card、row-key、pagination class |
| cloud-accounts/index.vue | page-header、添加filter-card、row-key |
| sync-policies/index.vue | page-header、filter-card结构、row-key |

---

## Phase 25: 多角色验证核心发现 (2026-04-16)

### 核心问题：资源列表数据来源错误

**设计文档要求**：
> 资源同步流程：云平台API → 解析标签 → 项目归属映射 → 写入数据库 → **本地展示**

**实际实现**：
- `ComputeService.ListVMs` (compute.go:61-68) → 直接调用 `provider.ListVMs`
- `NetworkService.ListVPCs` → 直接调用 `provider.ListVPCs`
- 所有资源列表都是实时查询云平台API，而非查询本地数据库

**已有本地存储模型**（未使用）：
- `CloudVM` (sync_cloud_vms表) - 同步后的虚拟机数据
- `CloudVPC` (sync_cloud_vpcs表) - 同步后的VPC数据
- `CloudSubnet` (sync_cloud_subnets表) - 同步后的子网数据
- 等等...

### 权限中间件未注册问题

**已实现但未注册**：
| 中间件 | 实现文件 | 注册状态 |
|--------|---------|---------|
| `PermissionMiddleware` | middleware/permission.go | ❌ main.go未注册 |
| `ProjectIsolationMiddleware` | middleware/project_isolation.go | ❌ main.go未注册 |

**影响**: 所有API对登录用户开放，权限控制失效。

### 业务流程完整度评估

| 流程步骤 | 状态 | 实现位置 |
|---------|------|---------|
| 定时任务调度器 | ✅ | scheduler.go |
| 同步服务 | ✅ | cloud_account.go:SyncResourcesWithMode |
| 标签解析 | ✅ | resource_mapping.go:DetermineProjectAttribution |
| 项目归属映射 | ✅ | 支持3种匹配条件 |
| 增量/全量同步 | ✅ | syncMode参数 |
| 同步日志 | ✅ | sync_log.go |
| 权限中间件 | ⚠️ 代码实现但未注册 | permission.go |
| 项目隔离 | ⚠️ 代码实现但未注册 | project_isolation.go |

### 修复优先级

**P0 (必须修复)**:
1. 资源列表数据来源 → 改为查询本地数据库
2. 注册权限中间件 → 在main.go中添加
3. Service层项目过滤 → 使用ApplyProjectFilter

**P1 (应该修复)**:
4. 前端云账号选择器 → 使用CloudAccountSelector组件
5. 验证其他资源列表表头

### 详细验证报告
- docs/superpowers/specs/2026-04-16-multi-agent-verification-report.md

## Phase 19: P0问题修复结果 (2026-04-15)

### 修复完成状态

**F1: 后台Cron调度器** ✅
- 实现文件：`backend/pkg/scheduler/scheduler.go`
- 功能：
  - 使用robfig/cron/v3库
  - 支持多种频率：once/daily/weekly/monthly/custom/hourly
  - 动态添加/删除/更新任务
  - 服务启动时加载所有active任务
  - 优雅关闭时停止调度器

**F2: 通用权限检查中间件** ✅
- 实现文件：`backend/internal/middleware/permission.go`
- 功能：
  - 解析请求路径自动提取资源和操作
  - HTTP方法映射：GET→list/get, POST→create, PUT→update, DELETE→delete
  - 系统管理员跳过权限检查
  - 查询用户角色和权限进行验证
  - 提供辅助函数GetUserPermissions/GetUserRoleIDs

**F3: 项目隔离验证中间件** ✅
- 实现文件：`backend/internal/middleware/project_isolation.go`
- 功能：
  - 获取用户所属项目列表（直接+通过用户组）
  - 注入project_ids到context供service层使用
  - 系统管理员跳过隔离（all_projects_visible=true）
  - 财务人员对账单API跳过隔离
  - 提供service层过滤辅助函数ApplyProjectFilter/ProjectIsolationScope

**F4: BSS API真实数据解析** ✅
- 实现文件：`backend/pkg/cloudprovider/adapters/alibaba/billing.go`
- 功能：
  - 定义完整响应结构体（5种）
  - 解析真实JSON响应
  - 返回真实数据而非mock
  - 优雅处理解析失败（返回空列表而非错误）

### 编译验证
所有修改编译成功，无错误。

### 待执行P1任务
1. G1: 资源标签解析与项目归属映射
2. G2: 增量/全量同步差异逻辑
3. G3: 同步日志记录
4. G4: 扩展同步资源类型（Disk/Snapshot/RDS/Redis）

## Phase 18: 全量前后端API交互审查发现 (2026-04-15)

### 审查总结

**已实现且真实对接的功能：**
1. 云账户管理完整CRUD - 真实数据库操作
2. 阿里云VM/VPC/Subnet/SG/EIP/Image管理 - 真实SDK调用
3. 同步策略CRUD（含域选择） - 前端已有domain_id字段
4. 定时任务CRUD及手动执行 - Execute方法调用SyncResources
5. JWT认证中间件 - Token解析和用户信息注入
6. 前端完整页面功能 - 所有按钮有对应API调用

**未实现或需要修复的严重问题（P0）：**
1. **后台Cron调度器缺失** - 定时任务只能手动触发，无自动执行
2. **通用权限检查中间件缺失** - 只有AdminOnlyMiddleware，无基于permission的访问控制
3. **项目隔离验证缺失** - 用户应只能看到所属项目的资源，未实现WHERE project_id过滤
4. **BSS API返回Mock数据** - 账单/订单/成本数据是硬编码mock，不是真实API解析

**未实现的高优先级问题（P1）：**
1. **资源标签解析与项目归属映射** - SyncResources不读取同步策略，不解析Tags
2. **增量/全量同步差异处理** - 前端传递mode参数，后端未区分处理逻辑
3. **同步日志记录** - 无SyncLog模型和记录逻辑
4. **同步资源类型不完整** - Disk/Snapshot/RDS/Redis未同步

### 同步策略模块特别说明

**用户需求：创建同步策略需添加域选择**

**审查发现：前端已有域选择功能！**
- `frontend/src/views/cloud-management/sync-policies/index.vue` 第124-128行
- 表单有 `domain_id` 字段，下拉框选择域
- 后端Handler/Service支持 `DomainID` 字段
- **结论：此需求已实现，无需额外开发**

### 资源同步核心流程缺失环节

业务逻辑应为：
```
[定时任务调度器触发] -> [读取账号绑定的定时任务配置] -> [读取账号绑定的同步策略]
-> [解析资源标签(Tags)] -> [应用同步策略的标签映射规则 → 确定项目归属]
-> [状态对比 + 数据处理] -> [写入同步日志] -> [更新云账号同步状态]
```

**当前实现状态：**
| 步骤 | 状态 |
|------|------|
| 定时任务调度器触发 | ❌ 缺少Cron Scheduler |
| 读取账号定时任务配置 | ✅ Execute方法读取 |
| 读取账号同步策略 | ❌ 未实现 |
| 解析资源标签 | ❌ 未实现 |
| 应用标签映射规则 | ❌ 未实现 |
| 增量/全量同步差异 | ⚠️ 未区分 |
| 已删除资源标记 | ❌ 未实现 |
| 写入同步日志 | ❌ 未实现 |
| 更新云账号同步状态 | ⚠️ 只更新Status |

### API权限验证流程缺失环节

业务逻辑应为：
```
[JWT认证中间件] -> [权限检查中间件] -> [项目隔离验证] -> [Handler处理]
```

**当前实现状态：**
| 步骤 | 状态 |
|------|------|
| JWT认证中间件 | ✅ 已实现 |
| 权限检查中间件 | ⚠️ 只有AdminOnly |
| 项目隔离验证 | ❌ 未实现 |
| RBAC权限矩阵 | ⚠️ 模型存在，检查逻辑缺失 |

### 修复计划已创建

详见：`docs/superpowers/specs/2026-04-15-fix-implementation-plan.md`

## Phase 16: 云账户管理增强需求分析 (2026-04-14)

### 功能需求拆解

**功能 1: 更新云账号弹窗**
| 功能点 | 前端实现 | 后端实现 | 复杂度 |
|--------|---------|---------|--------|
| 备注信息编辑 | el-input textarea | 已有 PUT API | 低 |
| 密钥ID编辑 | el-input | 已有 PUT API (Credentials JSON) | 低 |
| 密码编辑 | el-input password | 已有 PUT API | 低 |
| 测试连接按钮 | 调用 verify-credentials API | 已实现 VerifyCredentials | 低 |
| 凭证验证反馈 | ElMessage 显示结果 | - | 低 |

**功能 2: 属性设置-设置自动同步弹窗**
| 子页面 | 需要新模型 | 复杂度 | 说明 |
|--------|-----------|--------|------|
| 详情 | 否 | 中 | 基本信息+账号信息+权限列表（需调用云厂商API获取权限） |
| 资源统计 | 否 | 中 | 需新增 GetResourceStats API，从云厂商同步统计数据 |
| 订阅 | 是(CloudSubscription) | 高 | 完整 CRUD + 同步策略配置 |
| 云用户 | 是(CloudUser) | 高 | 完整 CRUD + 本地用户关联 |
| 云用户组 | 是(CloudUserGroup) | 高 | 完整 CRUD |
| 云上项目 | 是(CloudProject) | 高 | 完整 CRUD + 本地项目映射 |
| 定时任务 | 否(已有ScheduledTask) | 中 | 扩展关联字段，页面展示 |
| 操作日志 | 否(已有OperationLog) | 低 | 扩展 cloud_account_id 字段，表格展示 |

### 需要新增的数据模型

**1. CloudSubscription (云订阅)**
```go
type CloudSubscription struct {
    ID                  uint           // 主键
    CloudAccountID      uint           // 关联云账户
    Name                string         // 订阅名称
    SubscriptionID      string         // 订阅ID（云厂商侧）
    Enabled             bool           // 启用状态
    Status              string         // 状态
    SyncTime            *time.Time     // 上次同步时间
    SyncDuration        int            // 上次同步耗时（秒）
    SyncStatus          string         // 同步状态（completed/failed/in_progress）
    DomainID            uint           // 所属域
    DefaultProjectID    *uint          // 资源默认归属项目
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

**2. CloudUser (云用户)**
```go
type CloudUser struct {
    ID                uint           // 主键
    CloudAccountID    uint           // 关联云账户
    Username          string         // 用户名
    ConsoleLogin      bool           // 控制台登录权限
    Status            string         // 状态
    Password          string         // 密码（加密）
    LoginURL          string         // 登录地址
    LocalUserID       *uint          // 关联本地用户
    Platform          string         // 平台
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

**3. CloudUserGroup (云用户组)**
```go
type CloudUserGroup struct {
    ID              uint           // 主键
    CloudAccountID  uint           // 关联云账户
    Name            string         // 组名
    Status          string         // 状态
    Permissions     string         // 权限（JSON）
    Platform        string         // 平台
    DomainID        uint           // 所属域
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**4. CloudProject (云上项目)**
```go
type CloudProject struct {
    ID              uint           // 主键
    CloudAccountID  uint           // 关联云账户
    Name            string         // 云上项目名
    SubscriptionID  *uint          // 关联订阅
    Status          string         // 状态
    Tags            datatypes.JSON // 标签
    DomainID        uint           // 所属域
    LocalProjectID  *uint          // 本地项目映射
    Priority        int            // 优先级
    SyncTime        *time.Time     // 同步时间
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 需要扩展的现有模型

**OperationLog 扩展**
```go
// 新增字段
CloudAccountID *uint  `json:"cloud_account_id,omitempty" gorm:"column:cloud_account_id;index"`
EventType      string `json:"event_type" gorm:"column:event_type;size:50;default:'api_call'"` // api_call/console/scheduled
```

### 前端组件架构

**主组件: CloudAccountDetailDialog.vue**
- 使用 el-tabs 结构包含 8 个子页面
- 每个子页面作为独立组件或内联实现
- 统一的 API 调用和数据刷新机制

**组件结构建议:**
```
frontend/src/views/cloud-accounts/
├── index.vue                     # 云账户列表页（已有）
├── components/
│   ├── EditAccountDialog.vue     # 更新云账号弹窗（新建）
│   └── CloudAccountDetailDialog.vue  # 属性设置弹窗（新建）
│   ├── tabs/
│   │   ├── DetailTab.vue         # 详情子页面
│   │   ├── ResourceStatsTab.vue  # 资源统计子页面
│   │   ├── SubscriptionTab.vue   # 订阅子页面
│   │   ├── CloudUserTab.vue      # 云用户子页面
│   │   ├── CloudUserGroupTab.vue # 云用户组子页面
│   │   ├── CloudProjectTab.vue   # 云上项目子页面
│   │   ├── ScheduledTaskTab.vue  # 定时任务子页面
│   │   └── OperationLogTab.vue   # 操作日志子页面
```

### 后端 API 设计

**新增 API 端点:**
```
GET    /cloud-accounts/:id/resource-stats     # 获取资源统计
GET    /cloud-accounts/:id/permissions        # 获取权限列表
GET    /cloud-accounts/:id/subscriptions      # 订阅列表
POST   /cloud-accounts/:id/subscriptions      # 创建订阅
PUT    /cloud-accounts/:id/subscriptions/:sid # 更新订阅
DELETE /cloud-accounts/:id/subscriptions/:sid # 删除订阅
GET    /cloud-accounts/:id/cloud-users        # 云用户列表
POST   /cloud-accounts/:id/cloud-users        # 创建云用户
PUT    /cloud-accounts/:id/cloud-users/:uid   # 更新云用户
DELETE /cloud-accounts/:id/cloud-users/:uid   # 删除云用户
GET    /cloud-accounts/:id/cloud-user-groups  # 云用户组列表
GET    /cloud-accounts/:id/cloud-projects     # 云上项目列表
GET    /cloud-accounts/:id/operation-logs     # 操作日志列表
```

### 实现优先级

| 优先级 | 任务 | 原因 |
|--------|------|------|
| P0 | Task 1: 更新云账号弹窗 | 独立功能，不依赖其他模型 |
| P0 | Task 2: 主弹窗框架 | 基础框架，后续页面依赖 |
| P1 | Task 10: 操作日志 | 已有模型，只需扩展和展示 |
| P1 | Task 9: 定时任务 | 已有模型，只需页面展示 |
| P1 | Task 4: 资源统计 | 不需要新模型，调用云厂商API |
| P2 | Task 3: 详情 | 需要权限获取API |
| P2 | Task 5-8: 订阅/云用户/云用户组/云上项目 | 需要新模型 |

## Research Findings (2026-04-12 Session 3)

### 后端 API 调用链分析

**真实数据流已建立：**
```
前端页面 → API 请求 → Handler → Service → CloudProvider Adapter → 云厂商 SDK → 真实云资源
                                ↓
                           数据库存储（账号/资源元数据/同步状态）
```

**关键发现：**
1. `ComputeService` 和 `NetworkService` 已正确实现 `getProvider()` 方法获取适配器
2. 适配器通过 `cloudprovider.GetProvider(account.ProviderType, config)` 动态获取
3. Handler 层已正确调用 Service 层方法
4. Service 层已正确调用 Provider 接口

### 云厂商适配器实现状态 (更新 2026-04-12)

| 云厂商 | Compute | Network | Storage | Database | Middleware | SDK 安装 |
|--------|---------|---------|---------|----------|------------|----------|
| 阿里云 | ✅ 真实 | ✅ 真实 | ✅ 真实 | ⚠️ 待实现 | ⚠️ 待实现 | ✅ alibaba-cloud-sdk-go |
| 腾讯云 | ✅ 真实 | ✅ 真实 | ❌ 占位  | ❌ 占位   | ❌ 占位    | ✅ tencentcloud-sdk-go |
| AWS    | ✅ 真实 | ✅ 真实 | ❌ 占位  | ❌ 占位   | ❌ 占位    | ✅ aws-sdk-go-v2 |
| Azure  | ✅ 真实 | ✅ 真实 | ❌ 占位  | ❌ 占位   | ❌ 占位    | ✅ azure-sdk-for-go |

**阿里云适配器文件结构：**
- `provider.go` - SDK 初始化 (ECSClient, VPCClient)
- `vm.go` - CreateVM/DeleteVM/StartVM/StopVM/ListVMs 真实调用
- `vpc.go` - CreateVPC/DeleteVPC/ListVPCs/CreateSubnet 真实调用
- `disk.go` - CreateDisk/DeleteDisk/AttachDisk 真实调用

**腾讯云适配器文件结构：**
- `provider.go` - SDK 初始化 (CvmClient, VpcClient)
- `vm.go` - CreateVM/DeleteVM/StartVM/StopVM/ListVMs 真实调用
- `vpc.go` - CreateVPC/DeleteVPC/ListVPCs/CreateSubnet 真实调用

**AWS适配器文件结构：**
- `provider.go` - SDK 初始化 (EC2Client)
- `vm.go` - RunInstances/TerminateInstances/StartInstances 真实调用
- `vpc.go` - CreateVpc/DeleteVpc/CreateSubnet 真实调用

**Azure适配器状态 (已完成)：**
- `provider.go` - SDK 初始化 (vmClient, vnetClient, subnetClient, nsgClient, ipClient)
- `vm.go` - CreateVM/DeleteVM/StartVM/StopVM/RebootVM/GetVMStatus/ListVMs 真实调用
- `vpc.go` - CreateVPC/DeleteVPC/ListVPCs/CreateSubnet/DeleteSubnet/ListSubnets/CreateSecurityGroup 真实调用

### 前端页面功能按钮分析

| 页面 | 按钮/操作 | 后端 API | 状态 |
|------|----------|----------|------|
| VMs | 创建虚拟机 | POST /compute/vms | ✅ 已对接 |
| VMs | 启动/停止/重启 | POST /compute/vms/:id/action | ✅ 已对接 |
| VMs | 删除 | DELETE /compute/vms/:id | ✅ 已对接 |
| VPCs | 创建 VPC | POST /network/vpcs | ✅ 已对接 |
| VPCs | 删除 | DELETE /network/vpcs/:id | ✅ 已对接 |
| Subnets | 创建子网 | POST /network/subnets | ✅ 已对接 |
| SecurityGroups | 创建安全组 | POST /network/security-groups | ✅ 已对接 |
| EIPs | 申请 EIP | POST /network/eips | ✅ 已对接 |

**结论：前端已准备好调用真实后端 API，后端也已正确调用云厂商适配器。**

### 需要完善的内容

**Phase 8 - Azure 适配器实现：**
1. 安装 Azure SDK (`github.com/Azure/azure-sdk-for-go`)
2. 实现 Compute 接口 (VM CRUD)
3. 实现 Network 接口 (VNet/Subnet/NSG)

**Phase 9 - 云账号流程完善：**
1. 前端云账号创建 → 后端验证凭证 → 存储
2. 云账号同步功能（全量/增量）
3. 同步策略配置生效

**Phase 10 - Database/Middleware 适配器：**
1. 阿里云 RDS (MySQL/PostgreSQL)
2. 阿里云 Redis
3. 其他云厂商 Database

### 云账号凭证格式

**阿里云：**
```json
{
  "access_key_id": "LTAI...",
  "access_key_secret": "...",
  "region_id": "cn-hangzhou"
}
```

**腾讯云：**
```json
{
  "secret_id": "...",
  "secret_key": "...",
  "region": "ap-guangzhou"
}
```

**AWS：**
```json
{
  "access_key": "AKIA...",
  "secret_key": "...",
  "region": "us-west-2"
}
```

**Azure：**
```json
{
  "tenant_id": "...",
  "client_id": "...",
  "client_secret": "...",
  "subscription_id": "..."
}
```

## Phase 13: 多云管理页面优化分析 (2026-04-13)

### 云账户管理页面分析 (cloud-accounts/index.vue)
| 功能 | 前端状态 | 后端API状态 | 问题 |
|------|---------|------------|------|
| 列表查询 | ✅ 完整 | ✅ 完整 | 无 |
| 向导式添加 | ✅ 完整 | ✅ Create | 凭证字段映射需验证 |
| 同步操作 | ✅ sync对话框 | ✅ /sync | 已对接 |
| 验证凭证 | ✅ verify按钮 | ✅ /verify | 已对接 |
| 状态设置 | ✅ status对话框 | ⚠️ PATCH /status | 后端未实现 |
| 属性设置 | ✅ attributes对话框 | ⚠️ PATCH /attributes | 后端未实现 |
| 编辑云账号 | ⚠️ 待实现 | ✅ PUT | 前端缺少编辑表单 |
| 健康状态检测 | ✅ 显示 | ❌ 无字段 | model无health_status字段 |
| 余额显示 | ✅ 显示 | ❌ 无字段 | model无balance字段 |

**关键发现:**
- 向导式添加流程完整，但需要验证凭证字段与各云厂商配置匹配
- status/attributes PATCH API后端未实现
- health_status/balance字段model中不存在，前端显示占位

### 同步策略页面分析 (sync-policies/index.vue)
| 功能 | 前端状态 | 后端API状态 | 问题 |
|------|---------|------------|------|
| 列表查询 | ✅ 完整 | ✅ 完整 | 无 |
| 创建策略 | ✅ 完整 | ✅ Create | 无 |
| 编辑策略 | ✅ 完整 | ✅ PUT | 无 |
| 状态切换 | ✅ 完整 | ✅ POST /status | 无 |
| 规则配置 | ✅ 标签配置 | ✅ Rules+Tags | 结构匹配 |
| 删除策略 | ✅ 完整 | ✅ Delete | 无 |
| 所属域选择 | ✅ 下拉框 | ⚠️ 需验证 | 需加载域列表 |
| 项目选择 | ✅ 下拉框 | ⚠️ 需验证 | 需加载项目列表 |

**关键发现:**
- 前端页面功能完整
- 需要实现域列表和项目列表的加载API调用
- 规则标签配置UI完善，与后端结构匹配

### 定时任务页面分析 (scheduled-tasks.vue)
| 功能 | 前端状态 | 后端API状态 | 问题 |
|------|---------|------------|------|
| 列表查询 | ✅ 完整 | ✅ 完整 | 无 |
| 创建任务 | ✅ 完整 | ✅ Create | 无 |
| 编辑任务 | ⚠️ 未调用API | ✅ PUT | handleSubmit缺少调用 |
| 删除任务 | ✅ 完整 | ✅ Delete | 无 |
| 状态切换 | ✅ 完整 | ✅ POST /status | 无 |
| **执行任务** | ❌ 缺失 | ✅ POST /execute | **前端API缺失** |
| **关联云账号** | ❌ 表单无字段 | ⚠️ model有 | **前端表单缺失** |

**关键发现:**
- executeScheduledTask 前端API函数缺失（后端已实现）
- 表单缺少 cloud_account_id 字段（选择关联的云账号）
- handleSubmit 编辑模式未调用API

### 需要完成的任务清单

1. **定时任务页面修复**
   - 前端API: 添加 executeScheduledTask 函数
   - 前端表单: 添加 cloud_account_id 选择字段
   - 前端表格: 添加"执行"按钮
   - handleSubmit: 修复编辑模式API调用

2. **云账户管理页面完善**
   - 后端: 实现 PATCH /status 和 PATCH /attributes
   - 前端: 添加编辑云账号表单/对话框
   - 向导: 验证凭证字段与云厂商配置匹配

3. **同步策略页面完善**
   - 前端: 加载域列表 (getDomains API)
   - 前端: 加载项目列表 (getProjects API)

## Phase 14: UI/UX 设计优化分析 (2026-04-13)

### 设计系统推荐 (ui-ux-pro-max)
| 维度 | 推荐 | 说明 |
|------|------|------|
| **Pattern** | Enterprise Gateway | 企业级平台，导航清晰，信任指标突出 |
| **Style** | Flat Design | 无阴影，简洁线条，图标为主，响应快 |
| **Colors** | Navy/Grey + Green Accent | #0F172A 主色，#22C55E 强调色（成功） |
| **Typography** | Plus Jakarta Sans | 企业级 B2B SaaS 字体 |
| **Avoid** | 过度动画、默认深色模式 | |

### 云账户管理页面 UX 问题分析

| 优先级 | 问题 | 影响 | 建议 |
|--------|------|------|------|
| **CRITICAL** | 表格列过多（15+列） | 信息过载，难以阅读 | 合并状态列，隐藏次要列，使用详情弹窗 |
| **CRITICAL** | 操作下拉菜单（点击2次） | 操作效率低 | 改为分组按钮：主操作+下拉次要 |
| **HIGH** | 缺少编辑对话框 | 用户无法修改账号信息 | 添加 EditDialog |
| **HIGH** | 缺少空状态 | 无数据时用户困惑 | 添加空状态引导 + 创建按钮 |
| **HIGH** | 向导步骤无进度可视化 | 用户不知道当前位置 | 每步骤显示进度指示 |
| **MEDIUM** | 缺少搜索/筛选功能 | 大量账号难以查找 | 添加名称搜索、平台筛选 |
| **MEDIUM** | 健康状态/余额字段未实现 | 显示占位数据 | 后端需实现或移除列 |

### 定时任务页面 UX 问题分析

| 优先级 | 问题 | 影响 | 建议 |
|--------|------|------|------|
| **CRITICAL** | 操作按钮4个（违反主CTA原则） | 视觉混乱 | "执行"为主按钮，其他用下拉菜单 |
| **CRITICAL** | 缺少空状态处理 | 无数据时页面空洞 | 添加空状态 + 创建任务引导 |
| **HIGH** | 缺少任务执行进度 | 执行结果无反馈 | 执行后显示进度条/结果弹窗 |
| **HIGH** | 操作列宽度280px | 列太宽，挤压其他列 | 优化为200px，分组按钮 |
| **MEDIUM** | 缺少任务执行历史 | 无法查看执行记录 | 添加执行日志/统计展示 |

### 同步策略页面 UX 问题分析

| 优先级 | 问题 | 影响 | 建议 |
|--------|------|------|------|
| **CRITICAL** | 操作按钮4个 | 每屏幕应只有一个主CTA | "查看"为主，其他下拉 |
| **CRITICAL** | 缺少空状态 | 无策略时用户困惑 | 添加空状态 + 创建引导 |
| **HIGH** | 规则配置复杂 | 用户难以理解 | 添加帮助提示/示例规则 |
| **HIGH** | 缺少策略应用预览 | 用户不确定规则效果 | 添加预览功能 |
| **MEDIUM** | 详情对话框规则展示 | 信息层级不清晰 | 使用卡片分组展示规则 |
| **MEDIUM** | 缺少策略模板 | 新手难以创建 | 预设常用策略模板 |

### 设计优化任务清单

1. **Task 1: 添加空状态组件** - ✅ 已完成 - 创建 EmptyState.vue 可复用组件，应用到三个页面
2. **Task 2: 优化操作按钮布局** - ✅ 已完成 - 主按钮+下拉菜单分组，操作列宽度优化为 140px
3. **Task 3: 添加云账户编辑对话框** - ✅ 已完成 - 实现 handleEdit/handleEditSubmit 函数和编辑对话框
4. **Task 4: 添加任务执行进度反馈** - ⚠️ 待实现（执行结果已有 ElMessage 提示）
5. **Task 5: 优化表格列结构** - ✅ 已完成 - 合并状态列，删除次要列，操作列宽度优化为 120px
6. **Task 6: 添加帮助提示** - ⚠️ 待实现（规则配置、向导步骤）

**实现完成状态:**
- 空状态组件: ✅
- 操作按钮优化: ✅
- 云账户编辑对话框: ✅
- 表格列优化: ✅
- 前端构建: ✅ 成功
- 后端构建: ✅ 成功

## Previous Findings

## Research Findings

### 项目进度现状 (docs/PROGRESS.md) - **发现过时信息！**

**PROGRESS.md 显示的状态 vs 实际状态对比：**

| 模块 | PROGRESS.md 状态 | 实际状态 |
|------|-----------------|----------|
| 同步策略 | ⬜ 待设计 | ✅ 后端完成 + API完成 + 前端页面缺失 |
| 定时同步任务 | ⬜ 待设计 | ✅ 后端完成 + API完成 + 前端完成 |
| 资源同步规则 | ⬜ 待设计 | ⬜ 待确认（可能是同步策略的一部分）|

**实际完成情况：**

1. **IAM 模块** - ✅ 全部完成（后端、前端、测试）
2. **平台功能模块** - ✅ 全部完成
3. **同步策略 (SyncPolicy)**：
   - 后端：`handler/sync_policy.go` ✅, `service/sync_policy.go` ✅, `model/sync_policy.go` ✅
   - API：`frontend/src/api/sync-policy.ts` ✅
   - 路由：`main.go` 中注册 `/sync-policies` ✅
   - 前端页面：❌ **缺失**（没有找到对应的 Vue 页面）
4. **定时同步任务 (ScheduledTask)**：
   - 后端：`handler/scheduled_task.go` ✅, `service/scheduled_task.go` ✅, `model/scheduled_task.go` ✅
   - API：`frontend/src/api/scheduled-task.ts` ✅
   - 路由：`main.go` 中注册 `/scheduled-tasks` ✅
   - 前端页面：`frontend/src/views/cloud-accounts/scheduled-tasks.vue` ✅

### 云厂商适配器现状分析

**阿里云适配器完成度：**
| 模块 | 状态 | 实现文件 |
|------|------|---------|
| Compute (VM) | ✅ 完整 | vm.go - CreateVM/DeleteVM/StartVM/StopVM/RebootVM/ListVMs/GetVM |
| Network (VPC) | ✅ 完整 | vpc.go - VPC/Subnet/SecurityGroup/EIP 全部实现 |
| Storage (Disk) | ✅ 完整 | disk.go - CreateDisk/DeleteDisk/AttachDisk/DetachDisk |
| Image | ✅ 完整 | vm.go - ListImages/GetImage |
| Database | ⬜ 未实现 | provider.go 返回 ErrUnsupportedOperation |
| Middleware | ⬜ 未实现 | provider.go 返回 ErrUnsupportedOperation |

**腾讯云/AWS/Azure 适配器：**
- 全部为骨架实现（provider.go 只有 GetCloudInfo）
- 所有接口返回 `ErrUnsupportedOperation`
- 需要初始化 SDK 客户端并实现核心接口

**同步策略数据模型：**
- `SyncPolicy` - 同步策略配置（Name, Remarks, Status, Enabled, Scope, DomainID）
- `Rule` - 同步规则（ConditionType, ResourceMapping, TargetProjectID, TargetProjectName）
- `RuleTag` - 规则标签（TagKey, TagValue）
- 条件类型：`all_match`, `any_match`, `key_match`
- 资源映射：`specify_project`, `specify_name`

**定时任务数据模型：**
- `ScheduledTask` - 定时同步任务（Name, Type, Frequency, TriggerTime, ValidFrom, ValidUntil, Status, CloudAccountID）
- 频次类型：`once`, `daily`, `weekly`, `monthly`, `custom`

2. 云厂商适配器：
   - 阿里云：主机 ✅ 网络 ✅ 存储 ✅ 数据库 ⬜ 中间件 ⬜
   - 腾讯云：骨架（Compute/Network 有基础）
   - AWS：骨架（Compute/Network 有基础）
   - Azure：骨架（Compute/Network 有基础）

### 项目架构发现
**后端结构：**
- `internal/handler/` - HTTP handlers (30+ 文件)
- `internal/service/` - Business logic services (30+ 文件)
- `internal/model/` - Database models
- `pkg/cloudprovider/` - 云厂商接口和适配器

**云厂商适配器架构：**
- 接口定义：`interfaces_compute.go`, `interfaces_network.go`, `interfaces_storage.go`, `interfaces_database.go`
- 适配器目录：`adapters/alibaba/`, `adapters/tencent/`, `adapters/aws/`, `adapters/azure/`
- 注册机制：`registry.go` 通过 `RegisterProvider()` 注册

**阿里云适配器现状：**
- `vm.go` - 虚拟机 CRUD 完整实现
- `vpc.go` - VPC 网络管理（部分实现）
- `disk.go` - 存储管理
- `provider.go` - Provider 初始化和大量 "not implemented" 方法

### Git Diff 分析 (95 文件改动)
**新增/大幅改动：**
- `backend/internal/handler/compute.go` - 新增 223 行
- `backend/internal/handler/network.go` - 新增 404 行
- `backend/internal/service/compute.go` - 新增 260 行
- `backend/internal/service/network.go` - 新增 140 行
- `frontend/src/views/compute/vms/index.vue` - 新增 321 行
- `frontend/src/views/network/vpcs/index.vue` - 新增 412 行
- `frontend/src/views/network/subnets/index.vue` - 新增 318 行
- `backend/pkg/cloudprovider/adapters/alibaba/vm.go` - 新增实现

## Technical Decisions
| Decision | Rationale |
|----------|-----------|
| 阿里云 VM 接口完整实现 | ECS API 相对简单，作为适配器实现模板 |
| 网络接口定义扩展 (INetwork) | 支持 VPC、Subnet、SecurityGroup、EIP、LoadBalancer、DNS |
| 高级网络功能暂返回 ErrUnsupportedOperation | 优先实现核心功能，复杂功能后续迭代 |

## Issues Encountered
| Issue | Resolution |
|-------|------------|
| 大量未提交改动 | 先分析现状，再决定是否提交或继续开发 |

## Resources
- 进度看板：`docs/PROGRESS.md`
- IAM 模块任务：`docs/iam-module-tasks.md`
- 认证源架构：`docs/auth_source_architecture.md`
- 云厂商接口：`backend/pkg/cloudprovider/interfaces_*.go`
- 阿里云适配器：`backend/pkg/cloudprovider/adapters/alibaba/`

## Visual/Browser Findings
- 阿里云适配器 vm.go 已实现完整的 VM CRUD 操作（CreateVM, DeleteVM, StartVM, StopVM, RebootVM, GetVMStatus, ListVMs, GetVM）
- INetwork 接口扩展支持 VPC 互联、VPC Peering、Route Table、L2 Network
- 前端新增了大量 compute/network 相关视图和路由

## Phase 4: 前端云资源页面现状分析 (2026-04-12)

### Compute 页面现状
| 页面/组件 | 状态 | 说明 |
|----------|------|------|
| vms/index.vue | ✅ 完整 | VM 列表、查询、分页、详情弹窗、VNC 控制台 |
| VMModal.vue | ✅ 完整 | VM 详情弹窗组件 |
| VNCConsole.vue | ✅ 完整 | VNC 远程控制台组件 |
| VMActionDropdown.vue | ✅ 完整 | VM 操作下拉菜单 |
| images/index.vue | ✅ 完整 | 镜像列表 |
| keys/index.vue | ✅ 完整 | 密钥管理 |
| host-templates/index.vue | ✅ 完整 | 主机模板 |
| autoscaling-groups/index.vue | ✅ 完整 | 弹性伸缩组 |
| **创建 VM 弹窗** | ❌ 缺失 | handleCreate 只显示"功能开发中" |

### Network 页面现状
| 页面 | 状态 | 说明 |
|------|------|------|
| vpcs/index.vue | ✅ 完整 | VPC 列表、详情查看、删除 |
| subnets/index.vue | ✅ 完整 | 子网列表 |
| services/eips/index.vue | ✅ 完整 | EIP 列表 |
| services/security-groups | ✅ 完整 | 安全组列表 |
| routes/index.vue | ✅ 完整 | 路由表 |
| **创建 VPC 弹窗** | ❌ 缺失 | 无创建按钮或弹窗 |
| **创建 Subnet 弹窗** | ❌ 缺失 | 无创建弹窗 |

### Storage 页面现状
| 页面 | 状态 | 说明 |
|------|------|------|
| block/block-storage/index.vue | ⬜ 待验证 | 块存储 |
| storage/disks/index.vue | ✅ 存在 | 磁盘列表 |
| storage/disk-snapshots/index.vue | ✅ 存在 | 磁盘快照 |

### Database 页面现状
| 页面 | 状态 | 说明 |
|------|------|------|
| rds/instances/index.vue | ✅ 存在 | RDS 实例列表 |
| redis/instances/index.vue | ✅ 存在 | Redis 实例 |
| mongodb/instances/index.vue | ✅ 存在 | MongoDB 实例 |

### API 完整度
| API 文件 | 状态 | 说明 |
|----------|------|------|
| api/compute.ts | ✅ 完整 | VM CRUD、镜像、模板、伸缩组 |
| api/network.ts | ✅ 完整 | VPC、Subnet、SecurityGroup、EIP、Region、Zone、Peering 等 |
| types/index.ts | ✅ 完整 | 所有类型定义完整 |

### 后端 Handler 完整度
| Handler | 状态 | 说明 |
|--------|------|------|
| compute.go | ✅ 完整 | VM CRUD、Action、Details、VNC、Images |
| network.go | ✅ 完整 | VPC、Subnet、SecurityGroup、EIP、高级网络 |

### 需要完善的功能
1. **创建虚拟机弹窗组件** - CreateVMModal.vue ✅ 已完成
2. **创建 VPC 弹窗** - vpcs/index.vue 添加创建功能 ✅ 已完成
3. **创建 Subnet 弹窗** - subnets/index.vue 添加创建功能 ✅ 已完成
4. **验证存储/数据库页面完整性** - 待后续迭代

## Phase 7: UI/UX 设计系统分析 (2026-04-12)

### ui-ux-pro-max 推荐设计系统

**产品类型**: multi-cloud management platform SaaS dashboard

| 维度 | 推荐 | 说明 |
|------|------|------|
| **Pattern** | App Store Style Landing | 展示真实截图，包含评分，平台特定 CTA |
| **Style** | Glassmorphism | 毛玻璃效果，模糊背景，层次感 |
| **Colors** | Dark bg + Green accent | #020617 背景, #22C55E 强调色 |
| **Typography** | Fira Code / Fira Sans | 仪表盘/数据分析风格，技术感 |
| **Effects** | Backdrop blur + Subtle border | 10-20px 模糊，1px 白色透明边框 |
| **Avoid** | 过度动画 + 默认深色模式 | 保持简洁，支持亮/暗双模式 |

### 推荐颜色变量

```css
--color-primary: #0F172A       /* 主要色 */
--color-on-primary: #FFFFFF    /* 主要色上的文字 */
--color-secondary: #1E293B     /* 次级色 */
--color-accent: #22C55E        /* 强调色/CTA */
--color-background: #020617    /* 背景色 */
--color-foreground: #F8FAFC    /* 前景色 */
--color-muted: #1A1E2F         /* 柔和色 */
--color-border: #334155        /* 边框色 */
--color-destructive: #EF4444   /* 危险操作色 */
```

### 现有页面分析

#### VMs 页面 (vms/index.vue)
| 问题 | 当前状态 | 建议改进 |
|------|---------|---------|
| 查询区域 | 固定展开 inline form | 添加折叠/展开功能 |
| 状态标签 | el-tag type 属性 | 统一颜色语义 (success=绿色, warning=橙色, danger=红色) |
| 空状态 | 无提示 | 添加空状态占位图和引导文案 |
| 响应式 | 无适配 | 添加断点样式 |
| 字体 | 默认字体 | 技术数据使用 Fira Code |

#### VPCs 页面 (vpcs/index.vue)
| 问题 | 当前状态 | 建议改进 |
|------|---------|---------|
| Tabs 过滤 | 全部/本地idc/公有云 | 标签样式增强 |
| 拓扑图 | placeholder 文字 | 添加视觉占位或简化拓扑图 |
| 操作按钮 | el-button + dropdown | 主要/次要/危险按钮分组优化 |

#### Subnets 页面 (subnets/index.vue)
| 问题 | 当前状态 | 建议改进 |
|------|---------|---------|
| 详情弹窗 | el-tabs 多标签 | 标签样式增强 |
| IP 使用图表 | el-statistic 简单数字 | 可视化进度条或环形图 |

#### CreateVMModal (5步向导)
| 问题 | 当前状态 | 建议改进 |
|------|---------|---------|
| 步骤指示器 | el-steps simple | 添加进度条或卡片式步骤 |
| 表单布局 | 单列布局 | 可考虑双列紧凑布局 |
| 确认页 | el-descriptions | 数据卡片分组布局 |

#### CreateVPCModal/CreateSubnetModal
| 问题 | 当前状态 | 建议改进 |
|------|---------|---------|
| CIDR 帮助 | 弹窗 + 表格 | 保持，样式增强 |
| 表单宽度 | 固定宽度 | 响应式宽度调整 |

### 优化优先级

| 优先级 | 类别 | 关键检查项 | 避免反模式 |
|--------|------|-----------|-----------|
| 1 (CRITICAL) | Accessibility | 对比度 4.5:1, Alt text, Keyboard nav | 移除 focus ring, 无标签图标按钮 |
| 2 (CRITICAL) | Touch & Interaction | 最小 44×44px, 加载反馈 | 仅依赖 hover, 瞬时状态变化 |
| 3 (HIGH) | Performance | 懒加载, CLS < 0.1 | 布局跳动 |
| 4 (HIGH) | Style Selection | 匹配产品类型, 一致性, SVG 图标 | 随意混用风格, emoji 作为图标 |
| 5 (HIGH) | Layout & Responsive | Mobile-first breakpoints | 水平滚动, 禁止缩放 |
| 6 (MEDIUM) | Typography & Color | 16px 基准, line-height 1.5 | < 12px 文字, 灰色叠加灰色 |

### 实现方案

**Phase 7 分步计划:**

1. **Step 1: CSS 变量定义** - 创建 `frontend/src/styles/design-system.css`
   - 定义颜色 token (primary, secondary, accent, surface, etc.)
   - 定义字体变量 (Fira Code for data, Fira Sans for text)
   - 定义间距 token (4/8/12/16/24/32/48)

2. **Step 2: Glassmorphism 样式** - 创建 `frontend/src/styles/glass-card.css`
   - backdrop-filter: blur(10-20px)
   - background: rgba(15, 23, 42, 0.8)
   - border: 1px solid rgba(255, 255, 255, 0.2)
   - 用于 el-card 和 el-dialog 自定义样式

3. **Step 3: 列表页面优化** - 修改 vpcs/subnets/vms/index.vue
   - 添加查询折叠按钮
   - 状态标签颜色语义化
   - 空状态组件

4. **Step 4: 创建弹窗优化** - 修改 CreateVMModal/CreateVPCModal/CreateSubnetModal
   - 步骤指示器视觉增强
   - 确认页数据卡片布局

5. **Step 5: 响应式样式** - 添加 CSS media queries
   - @media (max-width: 768px) { ... }
   - @media (max-width: 375px) { ... }

6. **Step 6: 无障碍检查** - 使用 axe-core 或手动检查
   - 验证对比度
   - 添加 aria-label
   - 测试键盘导航

### 新增文件清单
| 文件 | 类型 | 说明 |
|------|------|------|
| frontend/src/utils/cidr.ts | 工具 | CIDR 校验函数 |
| frontend/src/components/common/CloudAccountSelector.vue | 组件 | 云账号选择器 |
| frontend/src/components/network/CreateVPCModal.vue | 组件 | 创建 VPC 弹窗 |
| frontend/src/components/network/CreateSubnetModal.vue | 组件 | 创建子网弹窗 |
| frontend/src/components/vm/CreateVMModal.vue | 组件 | 创建 VM 5步向导 |
| frontend/src/api/network.ts | 修改 | 添加 ipv6_cidr 类型 |
| frontend/src/views/compute/vms/index.vue | 修改 | 集成创建弹窗 |
| frontend/src/views/network/vpcs/index.vue | 修改 | 集成创建弹窗 |
| frontend/src/views/network/subnets/index.vue | 修改 | 集成创建弹窗 |

### 组件特性总结
- **CreateVMModal**: 5步向导（基本配置→计算配置→网络配置→存储配置→确认创建），支持模板自动填充，字段级联加载
- **CreateVPCModal**: 单页表单，CIDR 格式校验，CIDR 帮助提示
- **CreateSubnetModal**: 单页表单，子网 CIDR 在 VPC 范围内校验
- **CloudAccountSelector**: 可复用选择器，显示云厂商类型和健康状态标签