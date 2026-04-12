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

## Phase 4 实现结果 (2026-04-12)

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