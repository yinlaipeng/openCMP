# Phase 18 全量前后端API交互审查报告

## 审查日期：2026-04-15

## 审查范围
- 云账户管理模块
- 同步策略模块
- 定时任务模块
- 资源同步核心流程
- API权限验证与IAM
- 费用中心模块
- 云厂商适配器

---

## 1. 云账户管理模块

### API端点清单与实现状态

| API端点 | 前端API | 后端Handler | 后端Service | 状态 |
|---------|---------|-------------|-------------|------|
| GET /cloud-accounts | ✅ getCloudAccounts | ✅ List | ✅ ListCloudAccounts | ✅ 完成 |
| POST /cloud-accounts | ✅ createCloudAccount | ✅ Create | ✅ CreateCloudAccount | ✅ 完成 |
| GET /cloud-accounts/:id | ✅ getCloudAccount | ✅ Get | ✅ GetCloudAccount | ✅ 完成 |
| PUT /cloud-accounts/:id | ✅ updateCloudAccount | ✅ Update | ✅ UpdateCloudAccount | ✅ 完成 |
| DELETE /cloud-accounts/:id | ✅ deleteCloudAccount | ✅ Delete | ✅ DeleteCloudAccount | ✅ 完成 |
| POST /cloud-accounts/:id/verify | ✅ verifyCloudAccount | ✅ Verify | ✅ VerifyCredentials | ✅ 完成 |
| POST /cloud-accounts/:id/sync | ✅ syncCloudAccount | ✅ Sync | ✅ SyncResources | ✅ 完成 |
| POST /cloud-accounts/:id/test-connection | ✅ testConnection | ✅ TestConnection | ✅ TestConnection | ✅ 完成 |
| GET /cloud-accounts/:id/verify-credentials | ✅ verifyCredentials | ✅ VerifyCredentials | ✅ VerifyCredentials | ✅ 完成 |
| PATCH /cloud-accounts/:id/status | ✅ updateCloudAccountStatus | ✅ UpdateStatus | ✅ UpdateCloudAccount | ✅ 完成 |
| PATCH /cloud-accounts/:id/attributes | ✅ updateCloudAccountAttributes | ✅ UpdateAttributes | ⚠️ 未持久化 | ⚠️ 部分 |
| POST /cloud-accounts/:id/test-connection-with-credentials | ✅ testConnectionWithCredentials | ✅ TestConnectionWithCredentials | ✅ TestConnectionWithCredentials | ✅ 完成 |

### 同步功能分析

**SyncResources 方法实现状态：**
- ✅ 同步虚拟机 (ListVMs)
- ✅ 同步 VPC (ListVPCs)
- ✅ 同步子网 (ListSubnets)
- ✅ 同步安全组 (ListSecurityGroups)
- ✅ 同步 EIP (ListEIPs)
- ✅ 同步镜像 (ListImages)
- ❌ 同步磁盘 (未调用)
- ❌ 同步快照 (未调用)
- ❌ 同步 RDS (未实现)
- ❌ 同步 Redis (未实现)
- ❌ 同步 Bucket (未实现)

### 发现的问题

| 问题 | 严重程度 | 描述 |
|------|----------|------|
| UpdateAttributes未持久化 | 中 | 属性字段（auto_sync, sync_policy, sync_interval）未保存到数据库，model缺少对应字段 |
| 同步资源类型不完整 | 中 | 未同步Disk/Snapshot/RDS/Redis/Bucket等资源 |

---

## 2. 同步策略模块

### API端点清单与实现状态

| API端点 | 前端API | 后端Handler | 后端Service | 状态 |
|---------|---------|-------------|-------------|------|
| GET /sync-policies | ✅ getSyncPolicies | ✅ List | ✅ ListSyncPolicies | ✅ 完成 |
| POST /sync-policies | ✅ createSyncPolicy | ✅ Create | ✅ CreateSyncPolicy | ✅ 完成 |
| GET /sync-policies/:id | ✅ getSyncPolicy | ✅ Get | ✅ GetSyncPolicy | ✅ 完成 |
| PUT /sync-policies/:id | ✅ updateSyncPolicy | ✅ Update | ✅ UpdateSyncPolicy | ✅ 完成 |
| DELETE /sync-policies/:id | ✅ deleteSyncPolicy | ✅ Delete | ✅ DeleteSyncPolicy | ✅ 完成 |
| POST /sync-policies/:id/status | ✅ updateSyncPolicyStatus | ✅ UpdateStatus | ✅ ToggleSyncPolicyStatus | ✅ 完成 |

### 前端表单字段分析

**创建同步策略弹窗已有字段：**
- ✅ 策略名称 (name)
- ✅ 应用范围 (scope)
- ✅ 备注 (remarks)
- ✅ **所属域 (domain_id)** - 已实现，有域下拉选择
- ✅ 规则配置 (rules)
- ✅ 启用状态 (enabled)

### 发现的问题

| 问题 | 严重程度 | 描述 |
|------|----------|------|
| 同步策略未应用到同步流程 | 高 | SyncResources方法未读取和应用同步策略的标签映射规则 |
| 规则未与云账号关联 | 中 | 用户需求：策略需关联域，让特定域下的云账号使用 |

---

## 3. 定时任务模块

### API端点清单与实现状态

| API端点 | 前端API | 后端Handler | 后端Service | 状态 |
|---------|---------|-------------|-------------|------|
| GET /scheduled-tasks | ✅ getScheduledTasks | ✅ List | ✅ ListScheduledTasks | ✅ 完成 |
| POST /scheduled-tasks | ✅ createScheduledTask | ✅ Create | ✅ CreateScheduledTask | ✅ 完成 |
| GET /scheduled-tasks/:id | ✅ getScheduledTask | ✅ Get | ✅ GetScheduledTask | ✅ 完成 |
| PUT /scheduled-tasks/:id | ✅ updateScheduledTask | ✅ Update | ✅ UpdateScheduledTask | ✅ 完成 |
| DELETE /scheduled-tasks/:id | ✅ deleteScheduledTask | ✅ Delete | ✅ DeleteScheduledTask | ✅ 完成 |
| POST /scheduled-tasks/:id/status | ✅ updateScheduledTaskStatus | ✅ UpdateStatus | ✅ UpdateScheduledTaskStatus | ✅ 完成 |
| POST /scheduled-tasks/:id/execute | ✅ executeScheduledTask | ✅ Execute | ✅ SyncResources | ✅ 完成 |

### 发现的问题

| 问题 | 严重程度 | 描述 |
|------|----------|------|
| **缺少后台Cron调度器** | **严重** | 没有自动执行定时任务的调度器，只能手动触发Execute |
| 任务执行无日志记录 | 中 | Execute方法执行后未写入同步日志 |
| 任务执行未读取同步策略 | 高 | Execute直接调用SyncResources，未读取账号绑定的同步策略 |

---

## 4. 资源同步核心流程

### 业务逻辑要求

```
[定时任务调度器触发] 
  -> [读取账号绑定的定时任务配置] 
  -> [读取账号绑定的同步策略] 
  -> [解析资源标签(Tags)] 
  -> [应用同步策略的标签映射规则 → 确定项目归属] 
  -> [状态对比 + 数据处理]
  -> [写入同步日志] 
  -> [更新云账号同步状态]
```

### 实现状态

| 步骤 | 状态 | 描述 |
|------|------|------|
| 定时任务调度器触发 | ❌ 未实现 | 没有后台Cron Scheduler |
| 读取账号定时任务配置 | ✅ Execute实现 | 手动触发时读取 |
| 读取账号同步策略 | ❌ 未实现 | SyncResources不读取同步策略 |
| 解析资源标签 | ❌ 未实现 | 资源Tags未解析 |
| 应用标签映射规则 | ❌ 未实现 | 无标签匹配和项目归属逻辑 |
| 增量/全量同步差异 | ⚠️ 部分 | 前端传递mode参数，后端未区分处理 |
| 已删除资源标记 | ❌ 未实现 | 全量同步时未标记state=terminated |
| 写入同步日志 | ❌ 未实现 | 无同步日志写入 |
| 更新云账号同步状态 | ⚠️ 部分 | 只更新Status字段，无同步时间等 |

---

## 5. API权限验证与IAM

### JWT认证中间件实现状态

| 步骤 | 状态 | 描述 |
|------|------|------|
| 提取Token | ✅ 完成 | Authorization header解析 |
| 解析user_id | ✅ 完成 | 使用ParseJWTToken |
| 验证Token有效性 | ⚠️ 部分 | JWT签名验证完成，无Redis验证 |
| 注入domain_id | ✅ 完成 | claims.DomainID注入context |
| 注入auth_source_id | ✅ 完成 | claims.AuthSourceID注入context |

### 权限检查中间件实现状态

| 步骤 | 状态 | 描述 |
|------|------|------|
| 查询用户角色 | ⚠️ AdminOnlyMiddleware | 只检查admin角色，无通用权限检查 |
| 查询角色权限 | ❌ 未实现 | 无权限检查中间件 |
| 项目隔离验证 | ❌ 未实现 | 无WHERE project_id IN (user_projects)注入 |
| RBAC权限矩阵 | ⚠️ 部分 | Role/Permission模型存在，无检查逻辑 |

### 发现的问题

| 问题 | 严重程度 | 描述 |
|------|----------|------|
| **缺少通用权限检查中间件** | **严重** | 只有AdminOnlyMiddleware，无基于权限的访问控制 |
| **缺少项目隔离验证** | **严重** | 用户只能看到所属项目的资源，此逻辑未实现 |
| Token无Redis验证 | 中 | JWT只验证签名，无Redis存储验证（可选增强） |

---

## 6. 费用中心模块

### BSS API实现状态

| API | 后端实现 | SDK调用 | 数据来源 |
|-----|----------|---------|---------|
| ListBills | ✅ | ⚠️ 发送请求 | **返回Mock数据** |
| ListOrders | ✅ | ⚠️ 发送请求 | **返回Mock数据** |
| ListRenewalResources | ✅ | ⚠️ 发送请求 | **返回Mock数据** |
| GetCostAnalysis | ✅ | ⚠️ 发送请求 | **返回Mock数据** |
| GetAccountBalance | ✅ | ⚠️ 发送请求 | **返回Mock数据** |

### 发现的问题

| 问题 | 严重程度 | 描述 |
|------|----------|------|
| **BSS API返回Mock数据** | **严重** | parseXxxResponse函数返回硬编码数据，不是真实API解析 |

---

## 7. 云厂商适配器完整性

### 阿里云适配器

| 资源类型 | 接口 | SDK调用 | 状态 |
|----------|------|---------|------|
| VM | CreateVM/DeleteVM/StartVM/StopVM/RebootVM/ListVMs/GetVM | ✅ ECS SDK | ✅ 真实 |
| VPC | CreateVPC/DeleteVPC/ListVPCs | ✅ VPC SDK | ✅ 真实 |
| Subnet | CreateSubnet/DeleteSubnet/ListSubnets | ✅ VPC SDK | ✅ 真实 |
| SecurityGroup | Create/Delete/List/Authorize/Revoke | ✅ ECS SDK | ✅ 真实 |
| EIP | Allocate/Release/Associate/Dissociate/List | ✅ VPC SDK | ✅ 真实 |
| Image | ListImages/GetImage | ✅ ECS SDK | ✅ 真实 |
| Disk | ✅ disk.go存在 | ⚠️ 待验证 | ⚠️ 需检查 |
| Bucket/Object | ❌ | ErrUnsupportedOperation | ❌ 未实现 |
| RDS | ❌ | ErrUnsupportedOperation | ❌ 未实现 |
| Redis | ❌ | ErrUnsupportedOperation | ❌ 未实现 |
| LoadBalancer | ❌ | ErrUnsupportedOperation | ❌ 未实现 |
| DNS | ❌ | ErrUnsupportedOperation | ❌ 未实现 |

### 其他云厂商适配器（腾讯云/AWS/Azure）

| 云厂商 | Compute | Network | Storage | Database | Billing |
|--------|---------|---------|---------|----------|---------|
| 腾讯云 | ⚠️ SDK调用 | ⚠️ SDK调用 | ❌ Stub | ❌ Stub | ❌ 未实现 |
| AWS | ⚠️ SDK调用 | ⚠️ SDK调用 | ❌ Stub | ❌ Stub | ❌ 未实现 |
| Azure | ⚠️ SDK调用 | ⚠️ SDK调用 | ❌ Stub | ❌ Stub | ❌ 未实现 |

---

## 8. 修复计划优先级

### P0 - 严重问题（立即修复）

1. **创建后台Cron调度器**
   - 文件：`backend/pkg/scheduler/scheduler.go`
   - 功能：自动执行定时任务，支持cron表达式

2. **实现通用权限检查中间件**
   - 文件：`backend/internal/middleware/permission.go`
   - 功能：基于permission检查API访问权限

3. **实现项目隔离验证**
   - 文件：`backend/internal/middleware/project_isolation.go`
   - 功能：注入WHERE project_id IN (user_projects)

4. **修复BSS API返回真实数据**
   - 文件：`backend/pkg/cloudprovider/adapters/alibaba/billing.go`
   - 功能：解析真实API响应，返回真实数据

### P1 - 高优先级（本周修复）

5. **实现资源标签解析与项目归属映射**
   - 文件：`backend/internal/service/resource_sync.go`
   - 功能：解析Tags，匹配同步策略规则，确定project_id

6. **实现增量/全量同步差异逻辑**
   - 文件：`backend/internal/service/cloud_account.go`
   - 功能：区分mode参数，全量同步标记已删除资源

7. **实现同步日志记录**
   - 文件：`backend/internal/model/sync_log.go`
   - 功能：记录每次同步的详细日志

8. **扩展同步资源类型**
   - 文件：`backend/internal/service/cloud_account.go`
   - 功能：同步Disk/Snapshot/RDS/Redis

### P2 - 中优先级（后续迭代）

9. **完善UpdateAttributes持久化**
   - 文件：`backend/internal/model/cloud_account.go`
   - 功能：添加auto_sync/sync_policy_id/sync_interval字段

10. **实现阿里云Database/Middleware适配器**
    - 文件：`backend/pkg/cloudprovider/adapters/alibaba/rds.go`
    - 功能：实现RDS/Redis接口

11. **完善其他云厂商适配器**
    - 文件：`backend/pkg/cloudprovider/adapters/tencent/aws/azure/`
    - 功能：补全Storage/Database接口

---

## 9. 总结

### 已实现功能（真实API调用）
- 云账户管理完整CRUD
- 阿里云VM/VPC/Subnet/SG/EIP/Image管理
- 同步策略CRUD（含域选择）
- 定时任务CRUD及手动执行
- JWT认证
- 前端完整页面功能

### 未实现功能（需修复）
- 后台Cron调度器
- 通用权限检查中间件
- 项目隔离验证
- 资源标签解析与项目归属
- 增量/全量同步差异处理
- 同步日志记录
- BSS API真实数据解析
- Database/Middleware适配器
- Storage适配器（Bucket/Object）