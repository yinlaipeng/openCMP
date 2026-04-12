# Findings & Decisions

## Requirements
- 继续完成 openCMP 多云管理平台开发
- 重点：多云管理模块（同步策略、资源同步规则、定时同步任务）
- 重点：云厂商适配器完善（阿里云 VPC/Database/Middleware、腾讯云/AWS/Azure）

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

### 后端代码结构发现

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