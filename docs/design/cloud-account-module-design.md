# 云账号模块页面设计方案

## 设计原则
1. **延续 Element Plus 风格**：使用 el-* 组件系列，保持视觉一致性
2. **遵循 Enterprise Gateway Pattern**：企业级平台风格，导航清晰，信任指标突出
3. **Flat Design 样式**：无阴影，简洁线条，图标为主，响应快
4. **8dp 间距系统**：使用 4/8/12/16/20/24/32 的间距节奏

## 颜色 Token（基于现有 Element Plus 主题）

| 角色 | Element Plus 变量 | 用途 |
|------|------------------|------|
| Primary | `--el-color-primary` | 主按钮、链接 |
| Success | `--el-color-success` | 启用、健康状态 |
| Warning | `--el-color-warning` | 告警、待处理 |
| Danger | `--el-color-danger` | 错误、禁用、删除 |
| Info | `--el-color-info` | 次要信息、标签 |

---

## 一、云账号列表页

### 页面结构（延续现有 .cloud-accounts-container 结构）

```
┌──────────────────────────────────────────────────────────────────────┐
│  .page-header { display: flex; justify: space-between; margin-bottom: 20px }  │
│  ├── h2 { font-size: 18px; font-weight: 600 } "云账户管理"              │
│  └── 工具区 .tool-bar { gap: 8px }                                     │
│      ├── el-button (刷新) icon="Refresh"                               │
│      ├── el-button type="primary" (新建) icon="Plus"                   │
│      ├── el-button (批量操作) :disabled="selectedIds.length === 0"     │
│      ├── el-button (导出) icon="Download"                              │
│      └── el-button (设置) icon="Setting"                               │
├──────────────────────────────────────────────────────────────────────┤
│  el-tabs .top-tabs { margin-bottom: 16px }                             │
│  ├── el-tab-pane label="全部" name="all"                               │
│  └── el-tab-pane label="公有云" name="public"                          │
├──────────────────────────────────────────────────────────────────────┤
│  el-card .filter-card { margin-bottom: 20px }                          │
│  └── el-form :inline="true"                                            │
│      ├── el-form-item label="名称/ID/IP"                               │
│      │   └── el-input placeholder="支持名称搜索，自动匹配IP或ID..."     │
│      │   └── el-tooltip 提示："默认为名称搜索，自动匹配 IP 或 ID 搜索项...│
│      ├── el-form-item label="平台"                                     │
│      │   └── el-select (阿里云/腾讯云/AWS/Azure)                        │
│      ├── el-form-item label="状态"                                     │
│      │   └── el-select (已连接/未连接/错误)                             │
│      ├── el-form-item label="启用状态"                                 │
│      │   └── el-select (启用/禁用)                                     │
│      └── el-form-item                                                  │
│          ├── el-button type="primary" (查询)                           │
│          └── el-button (重置)                                          │
├──────────────────────────────────────────────────────────────────────┤
│  el-table :data="accounts" row-key="id"                                │
│  ├── el-table-column type="selection" width="55"                       │
│  ├── el-table-column prop="name" label="名称" min-width="150"          │
│  │   └── template: el-link @click="openDetailDrawer(row)"              │
│  ├── el-table-column label="状态" width="100"                          │
│  │   └── el-tag (status) + el-tag (enabled) + el-tag (health_status)  │
│  ├── el-table-column prop="balance" label="余额" width="100"           │
│  │   └── ¥{{ row.balance?.toFixed(2) || '0.00' }}                      │
│  ├── el-table-column prop="provider_type" label="平台" width="100"     │
│  │   └── el-tag :type="getProviderType(row.provider_type)"             │
│  ├── el-table-column prop="account_number" label="账号" width="150"    │
│  ├── el-table-column label="上次同步耗时" width="120"                   │
│  │   └── {{ row.sync_duration ? row.sync_duration + 's' : '-' }}       │
│  ├── el-table-column prop="last_sync" label="同步时间" width="160"     │
│  ├── el-table-column prop="domain_id" label="所属域" width="100"       │
│  ├── el-table-column label="资源归属方式" width="140"                  │
│  │   └── {{ getResourceAssignmentText(row.resource_assignment_method) }}│
│  └── el-table-column label="操作" width="160" fixed="right"            │
│      └── template                                                      │
│          ├── el-button size="small" type="primary" (同步)              │
│          └── el-dropdown (更多)                                        │
│              ├── el-dropdown-item (状态设置)                           │
│              ├── el-dropdown-item (属性设置)                           │
│              ├── el-dropdown-item divided (启用/禁用)                  │
│              ├── el-dropdown-item (连接测试)                           │
│              ├── el-dropdown-item (更新账号)                           │
│              ├── el-dropdown-item (设置同步归属策略)                   │
│              ├── el-dropdown-item (设置同步策略)                       │
│              ├── el-dropdown-item (设置代理)                           │
│              ├── el-dropdown-item (设置免密登录)                       │
│              ├── el-dropdown-item (只读模式)                           │
│              └── el-dropdown-item divided danger (删除)                │
├──────────────────────────────────────────────────────────────────────┤
│  el-pagination .pagination { margin-top: 20px; text-align: right }     │
│  ├── layout="total, sizes, prev, pager, next, jumper"                 │
│  └── :page-sizes="[10, 20, 50, 100]"                                   │
└──────────────────────────────────────────────────────────────────────┘
```

### 行操作下拉菜单分组设计

```
更多 ▼
├── 状态设置          → 打开 StatusSettingDialog
├── 属性设置          → 打开 AttributeSettingDialog
├── ────────────────
├── 启用 / 禁用       → 直接 toggle enabled 状态
├── 连接测试          → 调用 testConnection API
├── 更新账号          → 打开 EditAccountDialog
├── ────────────────
├── 设置同步归属策略   → 打开 SetSyncAttributionDialog
├── 设置同步策略      → 打开 SetSyncPolicyDialog
├── 设置代理          → 打开 SetProxyDialog
├── 设置免密登录      → 确认弹窗后 toggle
├── 只读模式          → 确认弹窗后 toggle
├── ────────────────
├── 删除 (红色)       → 确认删除弹窗
```

---

## 二、新建云账号4步向导

### 向导对话框结构

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="添加云账户" width="80%" :fullscreen="isMobile"      │
│  ├── el-steps :active="wizardStep" finish-status="success" align-center│
│  │   ├── el-step title="选择云平台"                                    │
│  │   ├── el-step title="配置云账号"                                    │
│  │   ├── el-step title="配置同步资源区域"                              │
│  │   └── el-step title="定时同步任务设置"                              │
│  ├── wizard-step-content { margin: 30px 0 }                            │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【第1步：选择云平台】                                              │
│  │   ├── h3 "请选择您要添加的云平台"                                   │
│  │   ├── .provider-category { margin-bottom: 24px }                    │
│  │   │   ├── h4 "公有云"                                               │
│  │   │   └── .provider-grid { grid-template-columns: repeat(4, 1fr) } │
│  │   │       ├── el-card.provider-card @click="selectProvider('alibaba')│
│  │   │       │   └── .provider-item { text-align: center }             │
│  │   │       │       ├── el-icon :size="40" (阿里云图标)               │
│  │   │       │       └── h4 "阿里云"                                   │
│  │   │       ├── ... (AWS/Azure/华为云/腾讯云/Google/天翼云/金山云...)  │
│  │   ├── .provider-category                                            │
│  │   │   ├── h4 "私有云 & 虚拟化平台"                                  │
│  │   │   └── .provider-grid                                            │
│  │   │       ├── ... (VMware/OpenStack/Cloudpods/Nutanix/Proxmox)      │
│  │   ├── .provider-category                                            │
│  │   │   ├── h4 "对象存储"                                             │
│  │   │   └── .provider-grid                                            │
│  │   │       ├── ... (S3/Ceph/XSKY)                                    │
│  │   └── .selected-indicator { margin-top: 16px }                      │
│  │       └── el-tag type="success" "已选择：阿里云"                    │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【第2步：配置云账号】                                              │
│  │   ├── el-form :model="wizardForm" :rules="wizardRules" label-width="150px"│
│  │   │   ├── el-form-item label="名称" prop="name"                     │
│  │   │   │   └── el-input v-model="wizardForm.name"                   │
│  │   │   ├── el-form-item label="备注"                                 │
│  │   │   │   └── el-input type="textarea" v-model="wizardForm.remarks" │
│  │   │   ├── el-form-item label="账号类型"                             │
│  │   │   │   └── el-select v-model="wizardForm.accountType"           │
│  │   │   │       ├── el-option label="主账号" value="primary"          │
│  │   │   │       └── el-option label="子账号" value="sub"              │
│  │   │   ├── el-form-item label="密钥 ID" prop="accessKeyId"           │
│  │   │   │   └── el-input v-model="wizardForm.accessKeyId"            │
│  │   │   ├── el-form-item label="密码 / Secret" prop="accessKeySecret" │
│  │   │   │   └── el-input type="password" show-password               │
│  │   │   ├── el-divider "资源归属设置"                                 │
│  │   │   ├── el-form-item label="资源归属方式"                         │
│  │   │   │   └── el-checkbox-group v-model="wizardForm.resourceAssignment"│
│  │   │   │       ├── el-checkbox label="根据同步策略归属"               │
│  │   │   │       ├── el-checkbox label="根据云上项目归属"               │
│  │   │   │       ├── el-checkbox label="根据云订阅归属"                │
│  │   │   │       ├── el-checkbox label="指定项目"                      │
│  │   │   ├── el-form-item label="同步策略" v-if="hasSyncPolicyAssignment"│
│  │   │   │   └── el-select v-model="wizardForm.syncPolicyId"          │
│  │   │   ├── el-form-item label="同步策略生效范围"                     │
│  │   │   │   ├── el-checkbox-group v-model="wizardForm.syncScope"     │
│  │   │   │       ├── el-checkbox label="资源标签"                      │
│  │   │   │       ├── el-checkbox label="项目标签"                      │
│  │   │   ├── el-form-item label="缺省项目"                             │
│  │   │   │   └── el-select v-model="wizardForm.defaultProjectId"      │
│  │   │   ├── el-divider "高级设置"                                     │
│  │   │   ├── el-form-item label="屏蔽同步资源"                         │
│  │   │   │   └── el-switch v-model="wizardForm.blockSyncResources"    │
│  │   │   ├── el-form-item label="代理"                                 │
│  │   │   │   ├── el-switch v-model="wizardForm.enableProxy"           │
│  │   │   │   └── el-select v-if="wizardForm.enableProxy"              │
│  │   │   ├── el-form-item label="开启免密登录"                         │
│  │   │   │   └── el-switch v-model="wizardForm.enablePasswordless"    │
│  │   │   ├── el-form-item label="只读模式"                             │
│  │   │   │   └── el-switch v-model="wizardForm.readOnlyMode"          │
│  │   │   ├── el-form-item label="连接测试"                             │
│  │   │   │   ├── el-button @click="testConnectionInWizard" :loading    │
│  │   │   │   └── el-tag :type="testResultType" v-if="testResult"      │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【第3步：配置同步资源区域】                                        │
│  │   ├── el-alert type="info" :closable="false"                        │
│  │   │   └── "请选择要同步资源的区域，该配置可在云账号导入后在云账号... │
│  │   ├── .region-actions { margin-bottom: 16px }                       │
│  │   │   ├── el-button @click="selectAllRegions" (全选)               │
│  │   │   └── el-button @click="clearAllRegions" (清空)                │
│  │   ├── el-table :data="availableRegions" @selection-change           │
│  │   │   ├── el-table-column type="selection"                         │
│  │   │   ├── el-table-column prop="name" label="区域名称"             │
│  │   │   ├── el-table-column prop="status" label="状态"               │
│  │   │   │   └── el-tag (可用/不可用)                                  │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【第4步：定时同步任务设置】                                        │
│  │   ├── h3 "定时同步任务设置（可选）"                                 │
│  │   ├── el-form :model="scheduleForm" label-width="120px"             │
│  │   │   ├── el-form-item label="名称"                                 │
│  │   │   │   └── el-input v-model="scheduleForm.name"                 │
│  │   │   ├── el-form-item label="类型"                                 │
│  │   │   │   └── el-input value="同步云账号" readonly                  │
│  │   │   ├── el-form-item label="触发频率"                             │
│  │   │   │   └── el-select v-model="scheduleForm.frequency"           │
│  │   │   │       ├── el-option label="单次" value="once"               │
│  │   │   │       ├── el-option label="每天" value="daily"              │
│  │   │   │       ├── el-option label="每周" value="weekly"             │
│  │   │   │       ├── el-option label="每月" value="monthly"            │
│  │   │   │       ├── el-option label="周期" value="custom"             │
│  │   │   ├── el-form-item label="触发时间"                             │
│  │   │   │   └── el-time-picker v-model="scheduleForm.triggerTime"    │
│  │   │   ├── el-form-item label="有效时间"                             │
│  │   │   │   ├── el-date-picker type="daterange"                      │
│  │   │   │       v-model="scheduleForm.validRange"                    │
│  │   ├── .step-actions                                                 │
│  │   │   ├── el-button type="primary" (确认创建任务)                   │
│  │   │   └── el-button (跳过，不创建定时任务)                          │
│  │   └───────────────────────────────────────────────────────────────│
│  └── template #footer                                                  │
│      └── .wizard-footer { display: flex; justify-content: flex-end }   │
│          ├── el-button @click="previousStep" :disabled="wizardStep===0"│
│          ├── el-button type="primary" @click="nextStep" (下一步)       │
│          │   └── el-button type="primary" @click="submitWizard" (提交) │
│          └── el-button @click="showWizard=false" (取消)                │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 三、云账号详情抽屉

### 抽屉结构（el-drawer 右侧滑出）

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-drawer title="" size="50%" direction="rtl"                        │
│  ├── #header                                                          │
│  │   └── .drawer-header { display: flex; align-items: center }        │
│  │       ├── .account-icon { margin-right: 16px }                     │
│  │       │   └── el-avatar :size="48"                                 │
│  │       │       └── el-icon :size="32" (云平台图标)                   │
│  │       ├── .account-info                                            │
│  │       │   ├── h3 { margin: 0 } {{ account.name }}                  │
│  │       │   ├── el-tag size="small" (平台)                           │
│  │       │   ├── el-tag size="small" :type="statusType" (状态)        │
│  │       ├── .quick-actions { margin-left: auto; gap: 8px }           │
│  │       │   ├── el-button size="small" type="primary" (同步)         │
│  │       │   ├── el-button size="small" (连接测试)                    │
│  │       │   ├── el-button size="small" (更新账号)                    │
│  │       │   ├── el-dropdown (更多)                                   │
│  │       │       ├── el-dropdown-item (状态设置)                      │
│  │       │       ├── el-dropdown-item (属性设置)                      │
│  │       │       ├── el-dropdown-item danger (删除)                   │
│  ├── el-tabs v-model="activeTab" type="border-card"                   │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【Tab 1: 详情】                                                   │
│  │   ├── el-descriptions title="基本信息" :column="2" border          │
│  │   │   ├── el-descriptions-item label="ID" {{ account.id }}         │
│  │   │   ├── el-descriptions-item label="状态" el-tag                  │
│  │   │   ├── el-descriptions-item label="名称"                        │
│  │   │   ├── el-descriptions-item label="域"                          │
│  │   │   ├── el-descriptions-item label="资源归属方式"                │
│  │   │   ├── el-descriptions-item label="屏蔽资源类型"                │
│  │   │   ├── el-descriptions-item label="平台"                        │
│  │   │   ├── el-descriptions-item label="账号"                        │
│  │   │   ├── el-descriptions-item label="账号 ID"                     │
│  │   │   ├── el-descriptions-item label="代理"                        │
│  │   │   ├── el-descriptions-item label="启用状态"                    │
│  │   │   ├── el-descriptions-item label="免密登录"                    │
│  │   │   ├── el-descriptions-item label="只读模式"                    │
│  │   │   ├── el-descriptions-item label="同步时间"                    │
│  │   │   ├── el-descriptions-item label="创建时间"                    │
│  │   │   ├── el-descriptions-item label="更新时间"                    │
│  │   │   ├── el-descriptions-item label="备注" :span="2"              │
│  │   ├── el-descriptions title="账号信息" :column="2" border          │
│  │   │   ├── el-descriptions-item label="环境"                        │
│  │   │   ├── el-descriptions-item label="健康状态" el-tag              │
│  │   │   ├── el-descriptions-item label="优惠率"                      │
│  │   │   ├── el-descriptions-item label="余额" ¥{{ account.balance }} │
│  │   │   ├── el-descriptions-item label="虚拟机数量"                  │
│  │   │   ├── el-descriptions-item label="宿主机数量"                  │
│  │   ├── el-card title="权限"                                         │
│  │   │   ├── el-table :data="permissions"                             │
│  │   │   │   ├── el-table-column prop="name" label="权限名称"         │
│  │   │   │   ├── el-table-column prop="description" label="描述"      │
│  │   │   │   ├── el-table-column prop="granted" label="是否授予"      │
│  │   │   │       └── el-tag (是/否)                                   │
│  │   │   ├── .permission-actions                                      │
│  │   │   │   ├── el-button (清空权限提示)                             │
│  │   │   │   └── el-button (导出)                                     │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【Tab 2: 资源统计】                                               │
│  │   ├── .stats-cards { display: grid; grid-template-columns: repeat(4,1fr)}│
│  │   │   ├── el-card.stat-card                                        │
│  │   │   │   ├── .stat-value {{ stats.vms }}                          │
│  │   │   │   ├── .stat-label "虚拟机"                                 │
│  │   │   ├── ... (LB实例/RDS实例/Redis实例/存储桶数量/对象存储容量/EIP...)│
│  │   ├── el-card title="使用率图表"                                   │
│  │   │   ├── .chart-grid { display: grid; grid-template-columns: repeat(2,1fr)}│
│  │   │   │   ├── .chart-item                                          │
│  │   │   │   │   ├── h4 "虚拟机开机率"                                │
│  │   │   │   │   └── el-progress :percentage="usageRates.vm_running_rate"│
│  │   │   │   ├── ... (磁盘挂载率/EIP使用率/IP使用率)                   │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【Tab 3: 订阅】                                                   │
│  │   ├── el-table :data="subscriptions"                                │
│  │   │   ├── el-table-column prop="name" label="名称"                 │
│  │   │   ├── el-table-column prop="subscription_id" label="Subscription ID"│
│  │   │   ├── el-table-column label="启用状态"                         │
│  │   │   │   └── el-switch @change="toggleSubscription"               │
│  │   │   ├── el-table-column prop="status" label="状态"               │
│  │   │   ├── el-table-column prop="sync_time" label="同步时间"        │
│  │   │   ├── el-table-column label="上次同步耗时"                     │
│  │   │   ├── el-table-column prop="sync_status" label="同步状态"      │
│  │   │   ├── el-table-column prop="domain_id" label="所属域"          │
│  │   │   ├── el-table-column label="资源默认归属项目"                 │
│  │   │   ├── el-table-column label="操作"                             │
│  │   │   │   ├── el-button size="small" (更改项目)                    │
│  │   │   │   ├── el-button size="small" (同步资源)                    │
│  │   │   │   └── el-dropdown (更多)                                   │
│  │   └───────────────────────────────────────────────────────────────│
│  │   【Tab 4-8: 云用户/云用户组/云上项目/定时任务/操作日志】            │
│  │   └ 各Tab使用 el-table 结构，延续现有 tabs/*.vue 组件设计          │
│  └───────────────────────────────────────────────────────────────────────│
│  └── #footer                                                          │
│      └── el-button @click="closeDrawer" (关闭)                        │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 四、操作弹窗设计

### 4.1 设置同步归属策略弹窗

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="设置同步归属策略" width="600px"                     │
│  ├── el-form label-width="150px"                                      │
│  │   ├── el-form-item label="资源归属方式"                            │
│  │   │   └── el-checkbox-group                                        │
│  │   │       ├── el-checkbox "根据同步策略归属"                       │
│  │   │       ├── el-checkbox "根据云上项目归属"                       │
│  │   │       ├── el-checkbox "根据云订阅归属"                         │
│  │   │       ├── el-checkbox "指定项目"                               │
│  │   ├── el-form-item label="指定项目" v-if="hasSpecifyProject"       │
│  │   │   └── el-select                                                │
│  │   ├── el-form-item label="屏蔽同步资源"                            │
│  │   │   └── el-switch                                                │
│  │   ├── el-form-item label="屏蔽资源类型" v-if="blockSync"           │
│  │   │   └── el-checkbox-group                                        │
│  │   │       ├── el-checkbox "虚拟机"                                 │
│  │   │       ├── el-checkbox "磁盘"                                   │
│  │   │       ├── el-checkbox "网络"                                   │
│  │   │       ├── el-checkbox "数据库"                                 │
│  ├── template #footer                                                 │
│  │   ├── el-button type="primary" @click="save" (保存)                │
│  │   └── el-button @click="close" (取消)                              │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.2 状态设置弹窗

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="状态设置" width="500px"                             │
│  ├── el-form label-width="120px"                                      │
│  │   ├── el-form-item label="云账号名称"                              │
│  │   │   └── el-input readonly                                        │
│  │   ├── el-form-item label="当前状态"                                │
│  │   │   └── el-tag                                                   │
│  │   ├── el-form-item label="设置状态"                                │
│  │   │   └── el-radio-group                                           │
│  │   │       ├── el-radio value="active" "已连接"                     │
│  │   │       ├── el-radio value="inactive" "未连接"                   │
│  │   │       ├── el-radio value="error" "连接错误"                    │
│  ├── template #footer                                                 │
│  │   ├── el-button type="primary" (确认)                              │
│  │   └── el-button (取消)                                             │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.3 属性设置弹窗

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="属性设置" width="700px"                             │
│  ├── el-tabs                                                          │
│  │   ├── el-tab-pane label="基本属性"                                 │
│  │   │   └── el-form                                                  │
│  │   │       ├── el-form-item label="备注"                            │
│  │   │       │   └── el-input type="textarea"                         │
│  │   │       ├── el-form-item label="所属域"                          │
│  │   │       │   └── el-select                                        │
│  │   ├── el-tab-pane label="同步属性"                                 │
│  │   │   └── el-form                                                  │
│  │   │       ├── el-form-item label="自动同步"                        │
│  │   │       │   └── el-switch                                        │
│  │   │       ├── el-form-item label="同步间隔"                        │
│  │   │       │   └── el-input-number                                  │
│  │   │       ├── el-form-item label="同步资源类型"                    │
│  │   │       │   └── el-checkbox-group                                │
│  ├── template #footer                                                 │
│  │   ├── el-button type="primary" (保存)                              │
│  │   └── el-button (取消)                                             │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.4 设置代理弹窗

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="设置代理" width="500px"                             │
│  ├── el-form label-width="120px"                                      │
│  │   ├── el-form-item label="启用代理"                                │
│  │   │   └── el-switch                                                │
│  │   ├── el-form-item label="代理地址" v-if="enableProxy"             │
│  │   │   └── el-input placeholder="http://proxy.example.com:8080"    │
│  │   ├── el-form-item label="代理类型"                                │
│  │   │   └── el-select (HTTP/SOCKS5)                                  │
│  │   ├── el-form-item label="认证"                                    │
│  │   │   └── el-checkbox (需要认证)                                   │
│  │   ├── el-form-item label="用户名" v-if="needAuth"                  │
│  │   │   └── el-input                                                 │
│  │   ├── el-form-item label="密码" v-if="needAuth"                    │
│  │   │   └── el-input type="password"                                 │
│  ├── template #footer                                                 │
│  │   ├── el-button type="primary" (保存)                              │
│  │   └── el-button (取消)                                             │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.5 免密登录/只读模式确认弹窗

```
┌──────────────────────────────────────────────────────────────────────┐
│  el-dialog title="确认操作" width="400px"                             │
│  ├── .confirm-content                                                 │
│  │   ├── el-icon :size="48" (Warning)                                 │
│  │   ├── p "确定要{{ actionText }}吗？"                               │
│  │   ├── p class="warning-text"                                       │
│  │   │   └── "{{ warningText }}" (操作影响说明)                       │
│  ├── template #footer                                                 │
│  │   ├── el-button type="primary" @click="confirm" (确认)             │
│  │   └── el-button @click="close" (取消)                              │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 五、关键CSS样式规范

```css
/* 容器与间距 */
.xxx-container { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.filter-card { margin-bottom: 20px; }
.pagination { margin-top: 20px; text-align: right; }

/* 表格 */
.el-table { width: 100%; }
.el-table { row-key: id; }

/* 工具区按钮 */
.tool-bar { display: flex; gap: 8px; }

/* 云平台卡片 */
.provider-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 20px; }
.provider-card { cursor: pointer; transition: all 0.2s; }
.provider-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
.provider-card.selected { border-color: var(--el-color-primary); }

/* 统计卡片 */
.stats-cards { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; }
.stat-card { text-align: center; padding: 20px; }
.stat-value { font-size: 24px; font-weight: 600; color: var(--el-color-primary); }
.stat-label { font-size: 14px; color: var(--el-text-color-secondary); }

/* 抽屉头部 */
.drawer-header { display: flex; align-items: center; padding: 16px 20px; border-bottom: 1px solid var(--el-border-color); }
.account-icon { margin-right: 16px; }
.account-info h3 { margin: 0; font-size: 18px; }
.quick-actions { margin-left: auto; display: flex; gap: 8px; }

/* 步骤表单 */
.wizard-step-content { margin: 30px 0; min-height: 400px; }
.wizard-footer { display: flex; justify-content: flex-end; gap: 10px; }
```

---

## 六、响应式适配

| 断点 | 适配策略 |
|------|---------|
| ≥1440px | 完整布局，表格显示所有列，统计卡片4列 |
| 1024-1439px | 表格隐藏次要列（余额、域），统计卡片3列 |
| 768-1023px | 表格隐藏更多列，向导fullscreen，统计卡片2列 |
| <768px | 移动端布局，卡片替代表格，抽屉占满屏幕 |

---

## 七、交互状态设计

### 状态Tag颜色映射

| 状态 | Element Plus Type | 显示文本 |
|------|------------------|---------|
| active/connected | success | 已连接 |
| inactive/disconnected | info | 未连接 |
| error | danger | 连接错误 |
| pending | warning | 连接中 |
| enabled=true | success | 启用 |
| enabled=false | info | 禁用 |
| healthy/normal | success | 正常 |
| unhealthy/abnormal | danger | 异常 |
| warning | warning | 警告 |

### 平台Tag颜色映射

| 平台 | Element Plus Type |
|------|------------------|
| alibaba (阿里云) | warning |
| tencent (腾讯云) | primary |
| aws | danger |
| azure | success |
| huawei | warning |
| google | primary |

---

## 八、实施优先级

| 优先级 | 功能 | 说明 |
|--------|------|------|
| P0 | 列表页工具区、顶部tab | 基础交互必需 |
| P0 | 表格字段补齐 | 信息完整性 |
| P1 | 新建向导云平台分类 | 用户引导体验 |
| P1 | 新建向导表单字段完整 | 功能完整性 |
| P2 | 详情抽屉顶部设计 | 视觉优化 |
| P3 | 操作弹窗 | 高级功能 |
| P4 | 响应式适配 | 移动端体验 |