# openCMP 项目全面完善计划

> 创建日期: 2026-04-14
> 版本: v1.0

## 一、项目现状概述

### 已完成模块
- IAM 栨块：Domain, Project, Group, User, Role, Permission, AuthSource, Policy ✅
- 消息中心：Message, NotificationChannel, Robot, Receiver, Subscription ✅
- 多云管理：CloudAccount, SyncPolicy, ScheduledTask ✅
- 云账户详情弹窗：8个子页面框架 ✅
- 云厂商适配器：阿里云 Compute/Network/Storage ✅，腾讯云/AWS/Azure Compute/Network 基础 ✅

### 待完善模块
- 云厂商适配器：Database/Middleware 全部待实现
- 费用中心：9个子页面骨架完成，真实API调用待实现
- 云账户详情：子页面功能待完善（编辑/删除等）
- 网络/存储/数据库页面：大量功能开发中

---

## 二、功能不全问题清单

### 2.1 云账户详情子页面（高优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| SubscriptionTab 编辑/删除功能缺失 | 高 | `tabs/SubscriptionTab.vue:156-159` | 实现 handleChangeProject/handleSyncPolicy/handleToggle/handleDelete |
| SubscriptionTab 更改项目功能缺失 | 高 | 同上 | 添加项目选择对话框 |
| SubscriptionTab 同步策略设置缺失 | 高 | 同上 | 添加同步策略配置表单 |
| CloudUserTab 编辑/删除功能缺失 | 中 | `tabs/CloudUserTab.vue:52-53` | 后端添加 CRUD API，前端实现 |
| CloudUserGroupTab 编辑/删除功能缺失 | 中 | `tabs/CloudUserGroupTab.vue:46-47` | 后端添加 CRUD API，前端实现 |
| CloudProjectTab 编辑/映射功能缺失 | 中 | `tabs/CloudProjectTab.vue:68-69` | 实现编辑表单和本地项目映射选择 |
| ScheduledTaskTab 执行/编辑/删除功能缺失 | 高 | `tabs/ScheduledTaskTab.vue:66-69` | 调用真实 scheduled-task API |
| DetailTab 权限刷新获取真实数据 | 中 | `tabs/DetailTab.vue` | 后端实现 GetPermissions API 调用云厂商 |

### 2.2 费用中心模块（中优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| 续费管理同步功能缺失 | 中 | `finance/orders/renewals/index.vue:154` | 调用云厂商 API 获取续费数据 |
| 账单导出下载功能缺失 | 中 | `finance/bills/export/index.vue:175` | 实现导出文件下载 API |
| 成本报告下载功能缺失 | 中 | `finance/cost/reports/index.vue:196` | 实现报告生成和下载 |
| 预算管理真实API对接 | 中 | `finance/cost/budgets/index.vue` | 对接后端 budget API |
| 阿里云 BSS API 集成 | 高 | `backend/service/finance.go` | 实现阿里云账单/订单/成本 API |

### 2.3 网络模块（中优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| Subnet 修改属性功能缺失 | 中 | `network/subnets/index.vue:340` | 添加属性修改对话框 |
| Subnet 更改项目功能缺失 | 中 | `network/subnets/index.vue:344` | 添加项目选择对话框 |
| Subnet 分割IP子网功能缺失 | 低 | `network/subnets/index.vue:348` | 后端实现子网分割逻辑 |
| Subnet 预留IP功能缺失 | 低 | `network/subnets/index.vue:352` | 后端实现IP预留逻辑 |
| DNS 更改域/设置共享功能缺失 | 低 | `network/services/dns/index.vue:270-274` | 实现 DNS 域管理 API |
| NAT Gateway 添加规则功能缺失 | 低 | `network/services/nat-gateways/index.vue:226` | 实现 NAT 规则管理 |
| LoadBalancer 添加监听功能缺失 | 中 | `network/loadbalancer/instances/index.vue:282` | 实现 LB 监听器配置 |

### 2.4 监控模块（低优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| 监控查询导出数据功能缺失 | 低 | `monitoring/monitor/query/index.vue:165` | 实现监控数据导出 |
| VM监控新增策略功能缺失 | 低 | `monitoring/resources/vms/index.vue:249` | 实现告警策略配置 |

### 2.5 计算模块（中优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| 主机模板部署虚拟机功能缺失 | 中 | `compute/host-templates/index.vue:425` | 实现基于模板创建 VM |

### 2.6 Layout/用户中心（低优先级）

| 问题 | 严重程度 | 文件位置 | 建议方案 |
|------|---------|---------|---------|
| 修改密码功能缺失 | 低 | `layout/index.vue:550` | 实现密码修改 API 和表单 |
| 个人信息功能缺失 | 低 | `layout/index.vue:553` | 实现用户信息编辑 |

---

## 三、API不合理问题清单

### 3.1 缺失的后端 API（高优先级）

| 模块 | 缺失 API | 建议方案 |
|------|---------|---------|
| 云账户详情 | CloudSubscription CRUD (PUT/DELETE) | 添加 updateSubscription/deleteSubscription |
| 云账户详情 | CloudUser CRUD | 添加 createCloudUser/updateCloudUser/deleteCloudUser |
| 云账户详情 | CloudUserGroup CRUD | 添加 createCloudUserGroup/updateCloudUserGroup/deleteCloudUserGroup |
| 云账户详情 | CloudProject CRUD + 映射 | 添加 updateCloudProject/mapToLocalProject |
| 云账户详情 | 定时任务筛选 | 支持按 cloud_account_id 筛选 |
| 费用中心 | 阿里云 BSS 集成 | 实现账单/订单/成本查询 |
| 网络 | Subnet 属性修改/分割/预留 | 扩展 Subnet Service |
| 网络 | DNS Zone/Record 管理 | 添加 DNS Handler |
| 网络 | LoadBalancer 监听器 | 添加 LB Listener Handler |

### 3.2 API 响应格式不统一（中优先级）

| 问题 | 文件位置 | 建议方案 |
|------|---------|---------|
| 分页参数命名不一致 | handler/*.go | 统一使用 page/page_size，废弃 limit/offset |
| 响应结构不一致 | 部分 handler | 统一使用 `{ items: [], total: N, page: N, page_size: N }` |
| 错误响应格式不一致 | handler/*.go | 统一使用 `{ error: "message" }` |

### 3.3 缺失的路由端点（高优先级）

| 端点 | 状态 | 建议 |
|------|------|------|
| PUT /cloud-accounts/:id/subscriptions/:sid | 缺失 | 添加 |
| DELETE /cloud-accounts/:id/subscriptions/:sid | 缺失 | 添加 |
| POST /cloud-accounts/:id/cloud-users | 缺失 | 添加 |
| PUT /cloud-accounts/:id/cloud-users/:uid | 缺失 | 添加 |
| DELETE /cloud-accounts/:id/cloud-users/:uid | 缺失 | 添加 |
| POST /cloud-accounts/:id/cloud-user-groups | 缺失 | 添加 |
| PUT /cloud-accounts/:id/cloud-user-groups/:gid | 缺失 | 添加 |
| DELETE /cloud-accounts/:id/cloud-user-groups/:gid | 缺失 | 添加 |
| PUT /cloud-accounts/:id/cloud-projects/:pid | 缺失 | 添加 |
| DELETE /cloud-accounts/:id/cloud-projects/:pid | 缺失 | 添加 |

---

## 四、云厂商适配器问题清单

### 4.1 未实现的适配器方法（高优先级）

| 云厂商 | 缺失方法数量 | 主要缺失 |
|--------|-------------|---------|
| AWS | ~40个 | Storage/Database/Middleware/LoadBalancer/DNS 全部未实现 |
| 腾讯云 | ~40个 | Storage/Database/Middleware 全部未实现 |
| Azure | ~40个 | Storage/Database/Middleware 全部未实现 |
| 阿里云 | ~20个 | Database/Middleware 未实现 |

### 4.2 具体缺失方法

**Storage 接口（所有云厂商）**
- CreateDisk/DeleteDisk/AttachDisk/DetachDisk/ResizeDisk
- CreateSnapshot/DeleteSnapshot/ListDisks/ListSnapshots
- CreateBucket/DeleteBucket/ListBuckets/PutObject/GetObject/DeleteObject/ListObjects
- CreateFileSystem/DeleteFileSystem/MountFileSystem/UnmountFileSystem/ListFileSystems

**Database 接口（所有云厂商）**
- CreateRDSInstance/DeleteRDSInstance/StartRDSInstance/StopRDSInstance
- RebootRDSInstance/ResizeRDSInstance/CreateRDSBackup/RestoreRDSFromBackup
- ListRDSInstances/ListRDSBackups
- CreateCacheInstance/DeleteCacheInstance/RebootCacheInstance/ResizeCacheInstance
- CreateCacheBackup/ListCacheInstances

**Middleware 接口（所有云厂商）**
- 未定义具体接口，需要设计

**高级网络接口（所有云厂商）**
- CreateLoadBalancer/DeleteLoadBalancer/CreateListener/DeleteListener/ListLoadBalancers
- CreateDNSZone/DeleteDNSZone/CreateDNSRecord/DeleteDNSRecord/ListDNSZones/ListDNSRecords
- CreateVPCInterconnect/DeleteVPCInterconnect/ListVPCInterconnects

---

## 五、数据模型问题清单

### 5.1 模型字段缺失（中优先级）

| 模型 | 缺失字段 | 建议 |
|------|---------|------|
| CloudAccount | auto_sync, sync_interval, sync_policy_id | 添加自动同步配置字段 |
| CloudAccount | health_check_time | 添加健康检查时间戳 |
| CloudSubscription | sync_policy_config | 添加同步策略详细配置 |
| ScheduledTask | enabled | 添加启用状态字段（目前只有 status） |
| OperationLog | request_detail, response_detail | 添加请求/响应详情（已规划） |

### 5.2 模型关联缺失（中优先级）

| 模型 | 缺失关联 | 建议 |
|------|---------|------|
| CloudSubscription -> SyncPolicy | 外键关联 | 添加 sync_policy_id 外键 |
| CloudUser -> User | 本地用户关联 | 添加 local_user_id 外键 |
| CloudProject -> Project | 本地项目关联 | 添加 local_project_id 外键（已有） |
| OperationLog -> CloudAccount | 云账户关联 | 添加 cloud_account_id 外键（已添加） |

---

## 六、前端页面结构问题

### 6.1 页面骨架但功能不全（统计）

| 类别 | 数量 | 说明 |
|------|------|------|
| 功能开发中占位 | 30+ | 使用 ElMessage.info('功能开发中') |
| 使用模拟数据 | 15+ | 硬编码数据，未调用真实 API |
| 缺少创建对话框 | 5 | 只有列表，无创建入口 |
| 缺少编辑对话框 | 10 | 编辑按钮显示"功能开发中" |

### 6.2 页面分组统计

| 模块 | 完整页面数 | 骨架页面数 | 功能不全页面数 |
|------|-----------|-----------|---------------|
| IAM | 8 | 0 | 0 |
| 消息中心 | 6 | 0 | 0 |
| 多云管理 | 4 | 0 | 5（子页面） |
| Compute | 8 | 4 | 3 |
| Network | 15 | 5 | 8 |
| Storage | 5 | 3 | 3 |
| Database | 3 | 3 | 3 |
| Middleware | 2 | 2 | 2 |
| Monitoring | 8 | 2 | 3 |
| Finance | 9 | 0 | 4 |

---

## 七、实施优先级排序

### Phase A: 云账户详情完善（最高优先级）
1. **A1**: SubscriptionTab 完整 CRUD
2. **A2**: ScheduledTaskTab 调用真实 API
3. **A3**: CloudUser/CloudUserGroup/CloudProject CRUD API + 前端
4. **A4**: DetailTab 权限真实获取

### Phase B: 费用中心真实API对接（高优先级）
1. **B1**: 阿里云 BSS SDK 集成
2. **B2**: 账单/订单/成本 API 实现
3. **B3**: 前端对接真实数据

### Phase C: 云厂商适配器完善（中优先级）
1. **C1**: 阿里云 Database 适配器
2. **C2**: AWS Storage 适配器
3. **C3**: 腾讯云/Azure Storage 适配器
4. **C4**: Database 适配器（所有云厂商）

### Phase D: 网络模块功能完善（中优先级）
1. **D1**: Subnet 属性修改/项目更改
2. **D2**: LoadBalancer 监听器配置
3. **D3**: DNS Zone/Record 管理
4. **D4**: NAT Gateway 规则管理

### Phase E: 监控/用户中心完善（低优先级）
1. **E1**: 监控数据导出
2. **E2**: 告警策略配置
3. **E3**: 修改密码功能
4. **E4**: 个人信息编辑

### Phase F: 代码质量优化（持续）
1. **F1**: API 响应格式统一
2. **F2**: 错误处理统一
3. **F3**: 代码复用优化
4. **F4**: 单元测试补充

---

## 八、预估工作量

| Phase | 任务数 | 预估时间 | 依赖 |
|-------|-------|---------|------|
| A | 15 | 2-3天 | 无 |
| B | 10 | 3-4天 | 阿里云 SDK |
| C | 20 | 5-7天 | 云厂商 SDK |
| D | 12 | 2-3天 | 无 |
| E | 8 | 1-2天 | 无 |
| F | 10 | 持续 | 无 |

**总计**: 约 15-20 天工作量（可并行分阶段实施）

---

## 九、风险与注意事项

1. **云厂商 API 限制**: 部分 API 可能需要特殊权限或付费才能调用
2. **SDK 版本兼容**: 云厂商 SDK 可能存在版本兼容问题
3. **数据一致性**: 多表关联需要确保外键约束正确
4. **性能问题**: 大量数据同步可能影响系统性能
5. **安全考虑**: 凭证存储需要加密，敏感信息不返回前端

---

## 十、建议实施策略

1. **优先完成 Phase A**: 云账户详情是用户高频操作，优先完善
2. **并行实施 Phase B 和 D**: 费用中心和网络模块可并行开发
3. **分批实施 Phase C**: 云厂商适配器按需实现，先完成阿里云
4. **持续优化 Phase F**: 代码质量改进穿插在各阶段中

---

## 十一、确认清单

- [ ] 是否同意优先级排序？
- [ ] 是否有特定模块需要优先完善？
- [ ] 是否需要调整预估工作量？
- [ ] 是否有其他未列出的问题需要补充？