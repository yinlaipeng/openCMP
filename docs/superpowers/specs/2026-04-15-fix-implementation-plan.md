# Phase 18 全量修复实施计划

## 修复日期范围：2026-04-15 至 2026-04-30

---

## Phase F: 后台调度器与同步核心流程实现

### F1: 创建后台Cron调度器
**优先级：P0**
**预计耗时：4小时**

**实现内容：**
1. 创建 `backend/pkg/scheduler/scheduler.go`
   - 使用 robfig/cron 或类似库
   - 支持标准cron表达式
   - 支持动态添加/删除任务
   - 任务执行状态跟踪

2. 创建 `backend/pkg/scheduler/task_runner.go`
   - 执行定时任务类型判断
   - 调用对应的业务逻辑（sync_cloud_account等）

3. 在 `main.go` 中启动调度器
   - 服务启动时加载所有active状态的定时任务
   - 注册任务到调度器

**涉及文件：**
- backend/pkg/scheduler/scheduler.go（新建）
- backend/pkg/scheduler/task_runner.go（新建）
- backend/cmd/server/main.go（修改）

---

### F2: 实现通用权限检查中间件
**优先级：P0**
**预计耗时：4小时**

**实现内容：**
1. 创建 `backend/internal/middleware/permission.go`
   - 从context获取user_id
   - 查询用户角色列表
   - 查询角色权限列表
   - 检查权限格式：`<module>:<resource>:<action>`
   - 与请求路径匹配验证

2. 实现PermissionCache（可选）
   - 缓存用户权限减少数据库查询

**涉及文件：**
- backend/internal/middleware/permission.go（新建）
- backend/cmd/server/main.go（修改-注册中间件）

---

### F3: 实现项目隔离验证
**优先级：P0**
**预计耗时：4小时**

**实现内容：**
1. 创建 `backend/internal/middleware/project_isolation.go`
   - 获取用户所属项目列表
   - 判断用户类型（系统管理员跳过隔离）
   - 注入project_id过滤条件

2. 在Service层实现过滤注入
   - 资源查询方法添加project_id条件
   - 支持全局查看权限角色

**涉及文件：**
- backend/internal/middleware/project_isolation.go（新建）
- backend/internal/service/*.go（修改查询方法）

---

### F4: 修复BSS API返回真实数据
**优先级：P0**
**预计耗时：4小时**

**实现内容：**
1. 修改 `backend/pkg/cloudprovider/adapters/alibaba/billing.go`
   - 解析真实JSON响应
   - 使用正确的响应结构体
   - 返回真实账单/订单/成本数据

2. 添加错误处理
   - 处理API限流
   - 处理无权限场景

**涉及文件：**
- backend/pkg/cloudprovider/adapters/alibaba/billing.go（修改）

---

## Phase G: 资源同步核心流程完善

### G1: 实现资源标签解析与项目归属映射
**优先级：P1**
**预计耗时：6小时**

**实现内容：**
1. 创建 `backend/internal/service/resource_mapping.go`
   - 解析资源Tags
   - 查询同步策略规则
   - 按优先级匹配规则
   - 确定project_id归属

2. 修改 `SyncResources` 方法
   - 获取账号绑定的同步策略
   - 对每个资源应用映射规则

**涉及文件：**
- backend/internal/service/resource_mapping.go（新建）
- backend/internal/service/cloud_account.go（修改）

---

### G2: 实现增量/全量同步差异逻辑
**优先级：P1**
**预计耗时：4小时**

**实现内容：**
1. 修改 `SyncResources` 方法支持mode参数
   - 增量同步：只INSERT新资源
   - 全量同步：INSERT新资源 + UPDATE已有资源 + 标记已删除资源

2. 实现已删除资源标记
   - 全量同步对比云厂商返回列表
   - 本地存在但云厂商不存在 -> state=terminated

**涉及文件：**
- backend/internal/service/cloud_account.go（修改）
- backend/internal/model/*.go（添加状态字段）

---

### G3: 实现同步日志记录
**优先级：P1**
**预计耗时：3小时**

**实现内容：**
1. 创建 `backend/internal/model/sync_log.go`
   ```go
   type SyncLog struct {
       ID              uint
       CloudAccountID  uint
       TaskID          *uint
       SyncType        string  // incremental/full
       StartTime       time.Time
       EndTime         *time.Time
       Status          string  // success/failed/partial
       ResourcesSynced int
       ResourcesNew    int
       ResourcesUpdated int
       ResourcesDeleted int
       ErrorMessage    string
   }
   ```

2. 创建 `backend/internal/service/sync_log.go`
   - 记录同步开始/结束
   - 更新统计信息

**涉及文件：**
- backend/internal/model/sync_log.go（新建）
- backend/internal/service/sync_log.go（新建）
- backend/cmd/server/main.go（添加迁移）

---

### G4: 扩展同步资源类型
**优先级：P1**
**预计耗时：4小时**

**实现内容：**
1. 添加Disk/Snapshot同步
   - 调用ListDisks/ListSnapshots
   - 存储到本地数据库

2. 添加RDS/Redis同步（待适配器实现后）
   - 调用ListRDSInstances/ListCacheInstances

**涉及文件：**
- backend/internal/service/cloud_account.go（修改）

---

## Phase H: 模型扩展与完善

### H1: 完善UpdateAttributes持久化
**优先级：P2**
**预计耗时：2小时**

**实现内容：**
1. 扩展CloudAccount模型
   ```go
   AutoSync       bool
   SyncPolicyID   *uint
   SyncInterval   int
   LastSyncTime   *time.Time
   ```

2. 修改UpdateAttributes handler/service
   - 持久化到数据库

**涉及文件：**
- backend/internal/model/cloud_account.go（修改）
- backend/internal/handler/cloud_account.go（修改）

---

### H2: 实现阿里云Database适配器
**优先级：P2**
**预计耗时：6小时**

**实现内容：**
1. 创建 `backend/pkg/cloudprovider/adapters/alibaba/rds.go`
   - ListRDSInstances
   - GetRDSInstance
   - CreateRDSBackup

2. 创建 `backend/pkg/cloudprovider/adapters/alibaba/redis.go`
   - ListCacheInstances

**涉及文件：**
- backend/pkg/cloudprovider/adapters/alibaba/rds.go（新建）
- backend/pkg/cloudprovider/adapters/alibaba/redis.go（新建）

---

### H3: 完善其他云厂商适配器
**优先级：P2**
**预计耗时：10小时**

**实现内容：**
1. 腾讯云Database适配器
2. AWS Database适配器
3. Azure Database适配器
4. 各厂商Storage适配器

**涉及文件：**
- backend/pkg/cloudprovider/adapters/tencent/*.go
- backend/pkg/cloudprovider/adapters/aws/*.go
- backend/pkg/cloudprovider/adapters/azure/*.go

---

## 实施顺序

```
Week 1 (2026-04-15 ~ 2026-04-19):
  F1: Cron调度器
  F2: 权限检查中间件
  F3: 项目隔离验证
  F4: BSS API真实数据

Week 2 (2026-04-20 ~ 2026-04-24):
  G1: 标签解析与项目归属
  G2: 增量/全量同步差异
  G3: 同步日志记录
  G4: 扩展同步资源类型

Week 3 (2026-04-25 ~ 2026-04-30):
  H1: UpdateAttributes持久化
  H2: 阿里云Database适配器
  H3: 其他云厂商适配器（可选）
```

---

## 验证测试清单

### F阶段验证
- [ ] 定时任务自动执行（检查数据库同步状态更新）
- [ ] 权限检查阻止未授权访问
- [ ] 项目隔离限制资源可见范围
- [ ] BSS API返回真实账单数据

### G阶段验证
- [ ] 资源Tags正确解析
- [ ] 项目归属按规则正确映射
- [ ] 增量同步只新增资源
- [ ] 全量同步标记已删除资源
- [ ] 同步日志正确记录

### H阶段验证
- [ ] 云账号属性持久化
- [ ] RDS/Redis列表同步
- [ ] 其他云厂商Database接口工作