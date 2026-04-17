# Task Plan: openCMP 项目落地实现

## Goal
实现 openCMP 多云管理平台的完整功能落地，前端页面与后端 API 真实对接，云厂商适配器调用真实 SDK，实现从云账号添加 → 资源同步 → 资源管理的完整业务流程。

## Current Phase
Phase 27 - 前端页面样式统一检查与修复 (Complete)
Next: 执行specs/目录中的开发Specs

## Specs (Ralph Wiggum)
开发Specs基于设计文档创建，按优先级执行：
- [specs/001-phase10-resource-management/spec.md](specs/001-phase10-resource-management/spec.md) - High Priority
- [specs/002-phase11-integration-verification/spec.md](specs/002-phase11-integration-verification/spec.md) - High Priority
- [specs/003-phase12-delivery-preparation/spec.md](specs/003-phase12-delivery-preparation/spec.md) - Medium Priority
- [specs/004-resource-project-mapping/spec.md](specs/004-resource-project-mapping/spec.md) - Medium Priority
- [specs/005-billing-cost-analysis/spec.md](specs/005-billing-cost-analysis/spec.md) - Medium Priority
- [specs/006-monitoring-alerts/spec.md](specs/006-monitoring-alerts/spec.md) - Low Priority

## Phases

### Phase 27: 前端页面样式统一检查与修复
- **目标**: 检查所有资源列表页面的样式统一性，参考主机模版页面的样式进行统一
- **参考页面**: frontend/src/views/compute/host-templates/index.vue
- **参考技能**: ui-ux-pro-max (spacing-scale, consistency, touch-target-size)
- **Status:** complete ✅
- [x] **检查完成** ✅ - 4个Agent并行检查，发现汇总完成
- [x] **P0修复: 补全筛选区和分页** ✅
  - images/index.vue: 添加筛选区、分页、padding ✅
  - eips/index.vue: 添加筛选区、分页、padding ✅
  - rds/instances/index.vue: 添加分页、修改padding ✅
  - redis/instances/index.vue: 添加分页、修改padding ✅
  - vpcs/index.vue: 添加分页、修改padding ✅
  - subnets/index.vue: 添加分页、修改padding ✅
  - policies/index.vue: 添加分页、修改padding ✅
  - resources/vms/index.vue: 添加分页、修改padding ✅
- [x] **P1修复: 统一页头和卡片结构** ✅
  - compute/vms/index.vue: el-card header改为.page-header，el-collapse改为.filter-card ✅
  - storage/cloud/disks/index.vue: page-header、filter-card、row-key、pagination class ✅
  - cloud-accounts/index.vue: page-header、添加filter-card筛选区、row-key ✅
  - sync-policies/index.vue: page-header、筛选表单独立为.filter-card、row-key ✅
- **修复成果**:
  - 12个页面统一为标准样式结构
  - 容器统一使用 `padding: 20px`
  - 页头统一使用独立 `.page-header` div (不再使用 el-card header slot)
  - 筛选区统一使用 `el-card.filter-card` 包裹
  - 分页统一使用 `.pagination` class + `text-align: right`
  - 表格统一添加 `row-key="id"`

## Phases

### Phase 25: 多角色Agent验证与设计审查
- **目标**: 创建不同角色的Agent智能体，从多个角度验证项目完成程度，审查设计文档，检查资源展示问题
- **Status:** complete
- [x] **V1: 启动并行验证Agent** ✅（手动验证完成）
- [x] **V2: 收集验证报告并汇总** ✅
- [x] **V3: 根据验证结果制定修复计划** ✅
- [ ] **V4: 执行关键问题修复** 🔄
- **验证报告**: docs/superpowers/specs/2026-04-16-multi-agent-verification-report.md
- **设计文档**: openCMP综合设计文档.md
- **关键发现**:
  - **核心问题**: 资源列表（VM/VPC等）直接调用云平台API，应改为查询本地数据库（CloudVM/CloudVPC等表）
  - **权限问题**: PermissionMiddleware和ProjectIsolationMiddleware已实现但未在main.go中注册
  - **前端问题**: 云账号选择器应使用CloudAccountSelector组件而非手动输入ID

### Phase 26: 关键问题修复实施
- **目标**: 修复验证发现的P0级别问题
- **Status:** complete
- [x] **P0-1: 修改资源列表数据来源** ✅ - 从本地数据库获取
  - ComputeService.ListVMs → 查询CloudVM表
  - NetworkService.ListVPCs → 查询CloudVPC表
  - NetworkService.ListSubnets → 查询CloudSubnet表
  - NetworkService.ListSecurityGroups → 查询CloudSecurityGroup表
  - NetworkService.ListEIPs → 查询CloudEIP表
- [x] **P0-2: 注册权限中间件** ✅ - 在main.go中添加
  - PermissionMiddleware已注册
  - ProjectIsolationMiddleware已注册
- [x] **P0-3: Service层应用项目过滤** ✅ - 使用projectIDs参数
  - ComputeService.ListVMs添加projectIDs参数
  - NetworkService.ListVPCs添加projectIDs参数
  - NetworkService.ListSubnets添加projectIDs参数
  - Handler层使用GetProjectFilter提取project_ids并传递给Service
- [x] **P1-1: 前端云账号选择器** ✅ - 使用CloudAccountSelector组件
  - compute/vms/index.vue已改用CloudAccountSelector组件
  - 其他页面（VPC/Subnet/RDS/Redis）已使用该组件

### Phase 24: 监控与费用模块完善
- **目标**: 完善监控告警和费用中心模块
- **Status:** complete
- [x] **M1: 监控模型** ✅
  - 创建backend/internal/model/monitor.go
  - AlertPolicy/AlertHistory/MonitorResource模型
- [x] **M2: 监控Handler/Service** ✅
  - 创建backend/internal/handler/monitor.go
  - 创建backend/internal/service/monitor.go
  - 告警策略CRUD/告警历史/监控资源同步API
- [x] **M3: 监控路由注册** ✅
  - 在main.go添加/monitor路由
- [x] **M4: 前端监控API** ✅
  - 创建frontend/src/api/monitor.ts
- [x] **M5: 前端告警策略页面更新** ✅
  - 使用真实API替换Mock数据
  - 添加创建/编辑/删除/启用/禁用功能
- [x] **M6: 前端VM监控页面更新** ✅
  - 使用真实API替换Mock数据
  - 添加云账号选择器、同步、查看指标功能
- [x] **B1: 费用中心验证** ✅
  - 阿里云billing adapter已实现
  - 前端账单/续费页面已使用真实API
  - 后端finance service正确调用billing provider

### Phase 22: 数据库资源管理实现
- **目标**: 完善数据库资源管理模块（RDS/Redis）的业务流程
- **Status:** complete
- [x] **D1: 阿里云RDS适配器** ✅
  - 创建backend/pkg/cloudprovider/adapters/alibaba/rds.go
  - 实现IRDS接口：ListRDSInstances/CreateRDSInstance/DeleteRDSInstance等
  - 使用阿里云RDS SDK真实调用
- [x] **D2: 阿里云Redis适配器** ✅
  - 创建backend/pkg/cloudprovider/adapters/alibaba/redis.go
  - 实现IElasticCache接口：ListCacheInstances/CreateCacheInstance等
  - 使用阿里云RDS SDK（Redis通过相似API模式管理）
- [x] **D3: 数据库Handler/Service** ✅
  - 创建backend/internal/handler/database.go
  - 创建backend/internal/service/database.go
  - 实现RDS/Redis完整CRUD和操作API
- [x] **D4: 数据库路由注册** ✅
  - 在main.go添加/database/rds和/database/cache路由
- [x] **D5: 前端数据库API** ✅
  - 创建frontend/src/api/database.ts
  - 实现listRDS/createRDS/deleteRDS等API函数
- [x] **D6: 前端RDS页面更新** ✅
  - 添加CloudAccountSelector组件
  - 使用真实API调用替换Mock数据
  - 添加创建/删除/操作/调整规格/备份功能
- [x] **D7: 前端Redis页面更新** ✅
  - 添加CloudAccountSelector组件
  - 使用真实API调用替换Mock数据
  - 添加创建/删除/重启/调整规格/备份功能

### Phase 21: 资源管理业务流程实现
- **目标**: 完善核心资源管理模块（VM/网络/存储）的业务流程
- **Status:** complete
- [x] **P1: VM生命周期管理** ✅
  - 前端VM列表页面使用真实API
  - VM创建/启动/停止/重启/删除完整流程
  - VM详情模态框（基本信息/安全组/网络/磁盘/快照/日志）
  - 阿里云adapter使用真实SDK调用
- [x] **P2: 网络资源管理** ✅
  - VPC页面更新：添加云账号选择器，调用真实API
  - Subnet页面更新：添加云账号选择器，调用真实API
  - 后端VPC/Subnet/SecurityGroup/EIP handler完整实现
  - 阿里云VPC adapter使用真实SDK调用
- [x] **P3: 存储资源管理** ✅
  - 云硬盘管理页面完整实现
  - 支持创建/删除/挂载/卸载/扩容/快照
  - 支持同步云平台数据
  - 阿里云Disk adapter使用真实SDK调用

### Phase 20: 前后端集成与验证
- **目标**: 完善前端与后端API的集成，验证核心业务流程
- **Status:** complete
- [x] **H1: 数据库迁移集成** ✅
  - 更新migration.go添加所有新模型
  - 包含SyncLog/Rule/RuleTag/CloudVM等同步资源模型
- [x] **H2: 同步API集成** ✅
  - 更新cloud_account.go handler支持sync_mode参数
  - 前端已有SyncCloudAccountOptions接口
  - 后端正确解析mode并调用SyncResourcesWithMode
- [x] **H3: 同步日志API** ✅
  - 创建backend/internal/handler/sync_log.go
  - 创建frontend/src/api/sync-log.ts
  - 添加sync-logs路由到main.go
- [x] **H4: 前端同步日志页面** ✅
  - 创建frontend/src/views/cloud-accounts/sync-logs.vue
  - 添加路由和侧边栏菜单
  - 实现统计卡片、日志列表、详情弹窗

### Phase 19: 核心流程修复实施

### Phase 18: 全量前后端API交互与业务逻辑审查
- **目标**: 全面扫描前后端代码，检查所有API交互是否真实对接，验证业务逻辑实现完整性
- **Status:** complete
- [x] **Task 0: 创建审查计划** - 规划审查范围和方法
- [x] **Task 1: 云账户管理模块审查**
  - 前端按钮/API对应关系 ✅
  - 后端接口实现状态 ✅
  - 云账号添加流程验证 ✅
  - 同步功能真实性（调用SDK）✅
  - 测试连接功能 ✅
  - 同步资源类型完整性 ⚠️（缺少Disk/Snapshot/RDS/Redis）
- [x] **Task 2: 同步策略模块审查**
  - 创建同步策略逻辑 ✅ 已有域选择
  - 规则配置与业务流程匹配 ✅
  - 增量/全量同步逻辑实现 ⚠️ 后端未区分处理
  - 资源标签映射规则实现 ❌ 未实现
- [x] **Task 3: 定时任务模块审查**
  - 任务调度器实现状态 ❌ 缺少后台Cron调度器
  - 任务执行与云账号同步关联 ✅
  - 执行状态跟踪 ⚠️ 缺少日志
- [x] **Task 4: 资源同步核心流程审查**
  - 标签解析与项目归属映射 ❌ 未实现
  - 增量同步 vs 全量同步实现 ⚠️ 前端传递参数，后端未区分
  - 同步日志记录 ❌ 未实现
  - 云账号同步状态更新 ⚠️ 只更新Status
- [x] **Task 5: API权限验证流程审查**
  - JWT认证中间件 ✅ 已实现
  - 权限检查中间件 ⚠️ 只有AdminOnly，缺少通用权限检查
  - 项目隔离验证 ❌ 未实现
  - RBAC权限矩阵实现 ⚠️ 模型存在，检查逻辑缺失
- [x] **Task 6: IAM模块审查**
  - 域管理功能 ✅
  - 项目管理功能 ✅
  - 用户/组/角色管理 ✅
  - 权限策略配置 ✅
- [x] **Task 7: 费用中心模块审查**
  - 阿里云BSS API对接 ⚠️ 发送请求但返回Mock数据
  - 账单/订单/续费同步 ⚠️ Mock数据
  - 成本分析实现 ⚠️ Mock数据
- [x] **Task 8: 消息中心模块审查**（未深度审查）
- [x] **Task 9: 云厂商适配器完整性审查**
  - 阿里云：Compute ✅ Network ✅ Storage ⚠️ Database ❌ Middleware ❌
  - 腾讯云：Compute ⚠️ Network ⚠️ Storage ❌ Database ❌
  - AWS：Compute ⚠️ Network ⚠️ Storage ❌ Database ❌
  - Azure：Compute ⚠️ Network ⚠️ Storage ❌ Database ❌
- [x] **Task 10: 创建修复计划** ✅
- **详细报告**: docs/superpowers/specs/2026-04-15-api-interaction-review.md
- **修复计划**: docs/superpowers/specs/2026-04-15-fix-implementation-plan.md

### Phase 19: 核心流程修复实施
- **目标**: 修复Phase 18发现的严重问题
- **Status:** complete
- [x] **F1: 创建后台Cron调度器** (P0) ✅
  - 创建backend/pkg/scheduler/scheduler.go和task_runner.go
  - 使用robfig/cron/v3库
  - 支持动态添加/删除任务
  - 在main.go中启动调度器，服务启动时加载所有active任务
- [x] **F2: 实现通用权限检查中间件** (P0) ✅
  - 创建backend/internal/middleware/permission.go
  - 解析请求路径提取资源和操作
  - 查询用户角色和权限进行验证
  - 系统管理员跳过权限检查
- [x] **F3: 实现项目隔离验证** (P0) ✅
  - 创建backend/internal/middleware/project_isolation.go
  - 获取用户所属项目列表（直接+通过组）
  - 注入project_ids到context供service层使用
  - 系统管理员和财务人员（账单API）跳过隔离
- [x] **F4: 修复BSS API返回真实数据** (P0) ✅
  - 修改billing.go解析真实JSON响应
  - 定义响应结构体（BillListResponse等）
  - 实现parseXxxResponse解析真实数据
  - 添加错误处理
- [x] **G1: 实现资源标签解析与项目归属映射** (P1) ✅
  - 创建backend/internal/service/resource_mapping.go
  - 实现DetermineProjectAttribution函数
  - 支持all_match/any_match/key_match条件类型
  - 支持正则表达式和精确匹配
  - 添加BatchDetermineProjectAttribution批量处理
- [x] **G2: 实现增量/全量同步差异逻辑** (P1) ✅
  - 创建backend/internal/service/cloud_account.go SyncResourcesWithMode
  - 增量同步：新资源INSERT，已存在资源SKIP
  - 全量同步：新资源INSERT，已存在UPDATE，云平台已删除标记terminated
  - 添加convertTagsToJSON辅助函数
- [x] **G3: 实现同步日志记录** (P1) ✅
  - 创建backend/internal/model/sync_log.go
  - 创建backend/internal/service/sync_log.go
  - 实现StartSyncLog/CompleteSyncLog/UpdateSyncLogProgress
  - 支持详细日志记录和统计
- [x] **G4: 扩展同步资源类型** (P1) ✅
  - 扩展storage.go模型：CloudDisk/CloudSnapshot添加ProjectID/Tags/RegionID
  - 创建cloud_resources_sync.go模型：CloudVM/CloudVPC/CloudSubnet/CloudSecurityGroup/CloudEIP/CloudImage/CloudRDS/CloudRedis
  - 添加syncDisks/syncSnapshots/syncRDS/syncRedis同步函数
  - 扩展ScheduledTask模型：添加Enabled/SyncPolicyID字段

### Phase 17: 项目全面完善
- **目标**: 根据全面审查结果，系统性完善项目各模块功能
- [x] **Task 0: 全面项目审查** - 完成项目代码审查，创建完善计划
- [x] **Phase A: 云账户详情完善**（已完成）
  - [x] A1: SubscriptionTab 完整 CRUD - 后端添加 PUT/DELETE/Toggle/Sync API，前端完整实现
  - [x] A2: ScheduledTaskTab 调用真实 API - 支持 cloud_account_id 筛选，完整 CRUD
  - [x] A3: CloudUser/CloudUserGroup/CloudProject CRUD API + 前端完整实现
  - [x] A4: SubscriptionTab 更改项目/同步策略/启用禁用功能实现
- [x] **Phase B: 费用中心真实API对接**（已完成）
  - [x] B1: 阿里云 BSS SDK 集成 - 创建 IBilling 接口和 Alibaba billing adapter
  - [x] B2: 账单/订单/成本 API 实现 - SyncBills/SyncOrders/SyncRenewals/GetAccountBalance
  - [x] B3: 前端对接真实数据 - 续费管理同步功能调用真实API
- [x] **Phase C: 云厂商适配器完善**（已完成）
  - [x] C1: Storage 接口实现 - 阿里云 Disk/Snapshot adapter 已存在，创建 Storage Handler
  - [x] C2: 云硬盘管理页面 - 前端实现创建/删除/挂载/卸载/扩容/快照功能
  - [x] C3: 存储路由注册 - /storage/cloud-disks, /storage/cloud-snapshots
- [x] **Phase D: 网络模块功能完善**
  - [x] D1: 扩展网络接口 - UpdateSubnet/AddSecurityGroupRule/DeleteSecurityGroupRule/BindEIP/UnbindEIP
  - [x] D2: Alibaba provider 实现 - vpc.go 新增方法实现
  - [x] D3: 其他 provider stub 实现 - Tencent/Azure/AWS 添加 stub 方法
  - [x] D4: Handler 方法实现 - network.go 扩展方法
  - [x] D5: 路由注册 - main.go 新增 subnet/sg/eip 扩展路由
  - [x] D6: 前端 API - network.ts 扩展 API 方法
- [x] **Phase E: 监控/用户中心完善**
  - [x] E1: 监控数据导出 - monitoring/query/index.vue 实现 CSV 导出
  - [x] E2: 告警策略配置 - monitoring/resources/vms/index.vue 实现新增/编辑策略
  - [x] E3: 修改密码功能 - layout/index.vue + auth API 实现
  - [x] E4: 个人信息编辑 - layout/index.vue + auth API 实现
- **Status:** Phase C complete
- **详细计划**: docs/superpowers/specs/2026-04-14-project-improvement-plan.md

### Phase 16: 云账户管理增强
- **目标**: 完善云账户管理功能，实现更新云账号弹窗、属性设置自动同步弹窗及其8个子页面
- [x] **Task 0: 现状分析与需求梳理** ✅
- [x] **Task 1: 更新云账号弹窗** ✅
- [x] **Task 2: 属性设置主弹窗框架** ✅
- [x] **Task 3: 详情子页面完善** ✅（骨架完成）
- [x] **Task 4: 资源统计子页面完善** ✅（骨架完成）
- [x] **Task 5: 订阅子页面完善** ✅（骨架+后端API）
- [x] **Task 6: 云用户子页面完善** ✅（骨架+后端API）
- [x] **Task 7: 云用户组子页面完善** ✅（骨架+后端API）
- [x] **Task 8: 云上项目子页面完善** ✅（骨架+后端API）
- [x] **Task 9: 定时任务子页面完善** ✅（骨架完成）
- [x] **Task 10: 操作日志子页面完善** ✅（骨架+后端API）
- [x] **Task 11: 数据库迁移与编译验证** ✅
  - 新增模型：CloudSubscription、CloudUser、CloudUserGroup、CloudProject
  - 扩展模型：OperationLog 添加 cloud_account_id 字段
  - 后端编译成功
  - 前端编译成功
- **Status:** complete（框架和基础功能完成，后续迭代可完善细节）
- [ ] **Task 2: 属性设置主弹窗框架**
  - 前端：创建 CloudAccountDetailDialog.vue 组件
  - 结构：el-tabs 包含 8 个子页面
  - API：新增获取云账户详情统计的 API
- [ ] **Task 3: 详情子页面**
  - 基本信息：名称、平台、状态、创建时间等
  - 账号信息：账号ID、余额、上次同步时间
  - 权限列表：云账号在云平台中的权限（需要同步获取）
- [ ] **Task 4: 资源统计子页面**
  - 统计卡片：虚拟机、RDS、Redis、存储桶、EIP、VPC、子网等数量
  - 使用率指标：虚拟机开机率、磁盘挂载率、EIP使用率、IP使用率
  - 后端：新增 GetResourceStats API
- [ ] **Task 5: 订阅子页面**（需要新模型）
  - 数据模型：CloudSubscription（订阅ID、启用状态、同步状态、所属域、默认项目）
  - 前端表格：名称、订阅ID、状态、同步时间、所属域、操作
  - 操作功能：更改项目、同步策略设置、启用/禁用/删除
- [ ] **Task 6: 云用户子页面**（需要新模型）
  - 数据模型：CloudUser（用户名、控制台登录、状态、密码、登录地址、关联本地用户）
  - 前端表格：用户名、状态、平台、所属云账号、操作
- [ ] **Task 7: 云用户组子页面**（需要新模型）
  - 数据模型：CloudUserGroup（名称、状态、权限、平台、所属云账号、所属域）
  - 前端表格：名称、状态、权限、平台、所属域、操作
- [ ] **Task 8: 云上项目子页面**（需要新模型）
  - 数据模型：CloudProject（云上项目名、订阅、状态、标签、所属域、本地项目、优先级）
  - 前端表格：云上项目、订阅、状态、标签、本地项目、优先级、操作
- [ ] **Task 9: 定时任务子页面**
  - 利用现有 ScheduledTask 模型，扩展关联云账户
  - 前端表格：名称、状态、启用状态、操作动作、策略详情、操作
  - 操作：编辑、启用/禁用、删除、立即执行
- [ ] **Task 10: 操作日志子页面**
  - 利用现有 OperationLog 模型，扩展云账户关联字段
  - 前端表格：#ID、操作时间、资源名称、资源类型、操作类型、服务类型、风险级别、事件类型、结果、发起人、所属项目
  - 操作：查看详情弹窗
- [ ] **Task 11: 数据库迁移与编译验证**
  - 新增模型：CloudSubscription、CloudUser、CloudUserGroup、CloudProject
  - 扩展模型：OperationLog 添加 cloud_account_id 字段
  - 后端编译验证
  - 前端编译验证
- **Status:** in_progress

### Phase 15: 费用中心功能完善
- **目标**: 完善费用中心 9 个子页面的前端功能，实现阿里云 BSS API 数据同步
- [x] **Task 1: 页面现状分析** - 确认已完成骨架代码
- [x] **Task 2: 预算管理页面完善**
  - 前端表格：显示预算列表（预算名称、类型、金额、阈值、使用量、状态）✅
  - 前端对话框：新建/编辑预算表单 ✅
  - 功能按钮：新建、编辑、删除 ✅
- [x] **Task 3: 我的订单页面完善**（已有基础实现）
- [x] **Task 4: 续费管理页面完善**
  - 前端表格：待续费资源列表 ✅
  - 筛选：天数阈值筛选 ✅
  - 统计卡片：待续费数量、预计费用 ✅
- [x] **Task 5: 成本分析页面完善**
  - 图表展示：成本趋势条形图 ✅
  - 费用分布：按产品类型分布 ✅
- [x] **Task 6: 成本报告页面完善**
  - 报告列表表格 ✅
  - 生成报告对话框 ✅
- [x] **Task 7: 异常监测页面完善**
  - 异常列表表格 ✅
  - 处理异常对话框 ✅
  - 筛选：严重程度、状态 ✅
- [x] **Task 8: 账单导出中心完善**
  - 导出历史列表 ✅
  - 导出任务创建 ✅
- [x] **Task 9: 编译验证**
  - 前端编译成功 ✅
  - 后端编译成功 ✅
- [ ] **Task 10: 阿里云 BSS API 集成**（后续迭代）
- **Status:** in_progress（页面完成，API集成待后续）

### Phase 14: UI/UX 设计优化实施
- **目标**: 根据 ui-ux-pro-max 分析结果优化三个多云管理页面的 UI/UX
- [x] **Task 1: 创建可复用空状态组件** - EmptyState.vue
- [x] **Task 2: 添加空状态到三个页面**
  - cloud-accounts/index.vue ✅
  - sync-policies/index.vue ✅
  - scheduled-tasks.vue ✅
- [x] **Task 3: 优化操作按钮布局**
  - sync-policies: 查看 + 下拉菜单 ✅
  - scheduled-tasks: 执行 + 下拉菜单 ✅
- [x] **Task 4: 添加云账户编辑对话框**
  - 实现 handleEdit/handleEditSubmit ✅
  - 添加编辑对话框模板 ✅
- [x] **Task 5: 优化表格列结构**
  - 合并状态列（status/enabled/health_status）✅
  - 删除次要列 ✅
  - 操作列宽度优化 ✅
- [x] **Task 6: 编译验证**
  - 后端编译成功 ✅
  - 前端编译成功 ✅
- **Status:** complete

### Phase 13: 多云管理三页面优化完善
- **目标**: 完善云账户管理、同步策略、定时任务三个页面，使其达到可用状态
- [x] **Task 1: 定时任务页面修复**
  - 前端API: 添加 executeScheduledTask 函数 ✅
  - 前端表单: 添加 cloud_account_id 选择字段 ✅
  - 前端表格: 添加"执行"按钮 ✅
  - handleSubmit: 修复编辑模式调用 updateScheduledTask API ✅
- [x] **Task 2: 云账户管理页面完善**
  - 后端: 实现 PATCH /cloud-accounts/:id/status ✅
  - 后端: 实现 PATCH /cloud-accounts/:id/attributes ✅
  - 前端: 添加编辑云账号对话框 (待实现 - 使用 handleEdit)
- [x] **Task 3: 同步策略页面完善**
  - 前端: loadDomains() 加载域列表 ✅ (已有)
  - 前端: loadProjects() 加载项目列表 ✅ (已有)
- [x] **Task 4: 编译验证**
  - 后端编译成功 ✅
  - 前端编译成功 ✅
- **Status:** complete (Commit: 2ace556)

### Phase 8: 后端 API 完善与云厂商适配器补全
- [x] **Azure 适配器实现 (已完成)**
  - Azure SDK 安装 (azure-sdk-for-go: azidentity, armcompute/v5, armnetwork/v4)
  - Azure Compute (VM 创建/删除/启停/列表) - 真实 SDK 调用
  - Azure Network (VNet/Subnet/SecurityGroup/PublicIP) - 真实 SDK 调用
  - 编译验证成功
- [ ] **阿里云 Database/Middleware 适配器**
  - RDS MySQL/PostgreSQL 支持
  - Redis 实例支持
- [ ] **腾讯云 Database 适配器**
  - TencentDB MySQL 支持
- [ ] **AWS Database 适配器**
  - RDS MySQL/PostgreSQL 支持
- **Status:** in_progress (Azure 完成，Database 待实现)

### Phase 9: 云账号完整流程实现
- [x] **云账号同步功能实现**
  - SyncResources 方法：同步所有资源类型
  - /sync API 端点已添加
  - 前端已有 sync 对话框，已正确对接
- [x] **凭证验证增强**
  - VerifyCredentials 方法：通过实际 API 调用验证
  - /verify-credentials API 端点已添加
  - 阿里云：验证区域列表；其他云：验证镜像列表
- [x] **定时任务执行**
  - Execute 端点：手动触发同步任务
  - /scheduled-tasks/:id/execute 路由
  - 支持 sync_cloud_account 类型任务
- [ ] **定时调度器（Phase 12）**
  - 后台 cron scheduler 自动执行定时任务
  - 任务执行状态跟踪
- **Status:** mostly_complete (核心功能完成，后台调度器待 Phase 12)

### Phase 10: 资源管理业务流程实现
- [ ] **虚拟机全生命周期**
  - 创建：前端参数 → 后端验证 → 云厂商 SDK 调用 → 状态跟踪
  - 操作：启动/停止/重启/删除 真实调用
  - 详情：实时获取云厂商数据 + 本地缓存
- [ ] **网络资源管理**
  - VPC 创建/删除 流程验证
  - Subnet 创建/删除 流程验证
  - SecurityGroup 规则配置
  - EIP 申请/释放/绑定
- [ ] **存储资源管理**
  - Block Storage 创建/挂载/卸载/删除
  - Object Storage Bucket 创建/管理
- [ ] **数据库资源管理**
  - RDS 实例创建/管理
  - Redis 实例创建/管理
- **Status:** pending

### Phase 11: 前后端联调与功能验证
- [ ] **核心功能联调**
  - 云账号添加 → 验证 → 同步 完整流程
  - VM 创建 → 状态跟踪 → 操作 完整流程
  - VPC/Subnet 创建 → 资源关联 完整流程
- [ ] **异常处理完善**
  - 云厂商 API 错误友好提示
  - 网络异常重试机制
  - 权限不足处理
- [ ] **UI 交互优化**
  - 加载状态反馈
  - 操作进度显示
  - 成功/失败消息提示
- **Status:** pending

### Phase 12: 项目交付准备
- [ ] **部署配置**
  - Docker 容器化
  - 环境变量配置
  - 数据库初始化脚本
- [ ] **文档完善**
  - API 文档更新
  - 部署指南
  - 用户手册
- [ ] **最终测试**
  - 功能完整性测试
  - 多云场景测试
  - 性能测试
- **Status:** pending

## Key Architecture Points

### 真实数据流
```
前端页面 → API 请求 → Handler → Service → CloudProvider Adapter → 云厂商 SDK → 真实云资源
                                ↓
                           数据库存储（账号/资源元数据/同步状态）
```

### 云厂商适配器状态
| 云厂商 | Compute | Network | Storage | Database | Middleware |
|--------|---------|---------|---------|----------|------------|
| 阿里云 | ✅ SDK  | ✅ SDK  | ✅ SDK  | ⚠️ 待实现 | ⚠️ 待实现 |
| 腾讯云 | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |
| AWS    | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |
| Azure  | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| 优先实现阿里云完整功能 | 国内主要用户群，SDK成熟，作为标杆实现 |
| Azure 作为后续扩展 | 国际需求相对较低，SDK较复杂 |
| 采用同步+缓存混合模式 | 云厂商 API 有延迟，本地缓存提升响应速度 |

## Notes
- 当前后端 service 层已经正确调用 cloudprovider.GetProvider() 获取适配器
- 前端 API 调用已准备好，只需确保后端返回真实数据
- 关键是完成云厂商适配器的真实 SDK 调用实现