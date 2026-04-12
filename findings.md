# Findings & Decisions

## Requirements
- 实现 openCMP 项目完整落地，前端与后端真实对接
- 不再使用模拟接口，实现真实的云厂商 SDK 调用
- 完成从云账号添加 → 资源同步 → 资源管理的完整业务流程

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