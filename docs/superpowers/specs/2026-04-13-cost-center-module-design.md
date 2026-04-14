# 费用中心模块设计文档

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现费用中心一级菜单及 9 个子页面的骨架结构，为后续费用数据同步功能奠定基础。

**Architecture:** 三级菜单嵌套结构（费用中心 -> 二级菜单 -> 子页面），数据关联现有云账号，通过云厂商 API 同步获取。

**Tech Stack:** Vue 3 + Element Plus + TypeScript (前端), Go + Gin + Gorm (后端)

---

## 菜单结构设计

```
费用中心 (Cost Center) [icon: Wallet]
├── 订单管理 (Orders)
│   ├── 我的订单 (My Orders) [/finance/orders/my-orders]
│   └── 续费管理 (Renewals) [/finance/orders/renewals]
├── 费用账单 (Bills)
│   ├── 账单查看 (Bill View) [/finance/bills/view]
│   └── 账单导出中心 (Export Center) [/finance/bills/export]
└── 成本管理 (Cost Management)
    ├── 成本分析 (Analysis) [/finance/cost/analysis]
    ├── 成本报告 (Reports) [/finance/cost/reports]
    ├── 预算管理 (Budgets) [/finance/cost/budgets]
    └── 异常监测 (Anomaly Detection) [/finance/cost/anomaly]
```

---

## 前端实现

### 路由配置 (router.ts)

添加费用中心路由：

```typescript
{
  path: '/finance',
  name: 'Finance',
  meta: { title: '费用中心', icon: 'Wallet' },
  children: [
    {
      path: 'orders',
      name: 'FinanceOrders',
      meta: { title: '订单管理' },
      children: [
        {
          path: 'my-orders',
          name: 'FinanceMyOrders',
          component: () => import('@/views/finance/orders/my-orders/index.vue'),
          meta: { title: '我的订单' }
        },
        {
          path: 'renewals',
          name: 'FinanceRenewals',
          component: () => import('@/views/finance/orders/renewals/index.vue'),
          meta: { title: '续费管理' }
        }
      ]
    },
    {
      path: 'bills',
      name: 'FinanceBills',
      meta: { title: '费用账单' },
      children: [
        {
          path: 'view',
          name: 'FinanceBillView',
          component: () => import('@/views/finance/bills/view/index.vue'),
          meta: { title: '账单查看' }
        },
        {
          path: 'export',
          name: 'FinanceBillExport',
          component: () => import('@/views/finance/bills/export/index.vue'),
          meta: { title: '账单导出中心' }
        }
      ]
    },
    {
      path: 'cost',
      name: 'FinanceCost',
      meta: { title: '成本管理' },
      children: [
        {
          path: 'analysis',
          name: 'FinanceCostAnalysis',
          component: () => import('@/views/finance/cost/analysis/index.vue'),
          meta: { title: '成本分析' }
        },
        {
          path: 'reports',
          name: 'FinanceCostReports',
          component: () => import('@/views/finance/cost/reports/index.vue'),
          meta: { title: '成本报告' }
        },
        {
          path: 'budgets',
          name: 'FinanceCostBudgets',
          component: () => import('@/views/finance/cost/budgets/index.vue'),
          meta: { title: '预算管理' }
        },
        {
          path: 'anomaly',
          name: 'FinanceCostAnomaly',
          component: () => import('@/views/finance/cost/anomaly/index.vue'),
          meta: { title: '异常监测' }
        }
      ]
    }
  ]
}
```

### 页面目录结构

```
frontend/src/views/finance/
├── orders/
│   ├── my-orders/
│   │   └── index.vue    # 我的订单骨架页面
│   └── renewals/
│       └── index.vue    # 续费管理骨架页面
├── bills/
│   ├── view/
│   │   └── index.vue    # 账单查看骨架页面
│   └── export/
│       └── index.vue    # 账单导出中心骨架页面
└── cost/
    ├── analysis/
    │   └── index.vue    # 成本分析骨架页面
    ├── reports/
    │   └── index.vue    # 成本报告骨架页面
    ├── budgets/
    │   └── index.vue    # 预算管理骨架页面
    └── anomaly/
        └── index.vue    # 异常监测骨架页面
```

### 骨架页面组件模板

每个骨架页面包含：
- 页面标题和描述
- 空状态组件 (EmptyState.vue)
- 云账号筛选下拉框
- 功能占位提示

示例模板：

```vue
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">{{ pageTitle }}</span>
          <el-button type="primary" @click="handleSync">同步数据</el-button>
        </div>
      </template>
      
      <!-- 云账号筛选 -->
      <div class="filter-bar">
        <CloudAccountSelector v-model="selectedAccountId" />
      </div>
      
      <!-- 空状态 -->
      <EmptyState
        v-if="!data.length"
        :icon="pageIcon"
        :title="pageTitle"
        :description="pageDescription"
        createButtonText="同步数据"
        @create="handleSync"
      />
      
      <!-- 数据表格（骨架） -->
      <el-table v-else :data="data" v-loading="loading">
        <!-- 根据页面类型定义列 -->
      </el-table>
    </el-card>
  </div>
</template>
```

---

## 后端实现

### 数据模型设计

```go
// model/finance.go

// Bill 账单记录
type Bill struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint      `gorm:"index;not null" json:"cloud_account_id"`
    BillingCycle    string    `gorm:"type:varchar(20);not null" json:"billing_cycle"` // 2026-04
    ProductType     string    `gorm:"type:varchar(50)" json:"product_type"`           // ECS, RDS, OSS
    ProductName     string    `gorm:"size:200" json:"product_name"`
    InstanceID      string    `gorm:"size:100" json:"instance_id"`
    UsageAmount     float64   `json:"usage_amount"`
    UnitPrice       float64   `json:"unit_price"`
    TotalCost       float64   `json:"total_cost"`
    Currency        string    `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
    BillingMethod   string    `gorm:"type:varchar(20)" json:"billing_method"`         // 按量/包年包月
    Status          string    `gorm:"type:varchar(20)" json:"status"`                 // 已支付/待支付
    ProviderType    string    `gorm:"type:varchar(20)" json:"provider_type"`          // alibaba/tencent/aws/azure
    CreatedAt       time.Time `json:"created_at"`
}

// Order 订单记录
type Order struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint      `gorm:"index;not null" json:"cloud_account_id"`
    OrderID         string    `gorm:"uniqueIndex;size:100" json:"order_id"`           // 云厂商订单号
    OrderType       string    `gorm:"type:varchar(20)" json:"order_type"`             // 新购/续费/升降配
    ProductType     string    `gorm:"type:varchar(50)" json:"product_type"`
    ProductName     string    `gorm:"size:200" json:"product_name"`
    InstanceID      string    `gorm:"size:100" json:"instance_id"`
    Amount          float64   `json:"amount"`
    Currency        string    `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
    Status          string    `gorm:"type:varchar(20)" json:"status"`                 // 待支付/已支付/已取消
    PaymentTime     *time.Time `json:"payment_time"`
    EffectiveTime   time.Time `json:"effective_time"`
    ExpireTime      *time.Time `json:"expire_time"`
    ProviderType    string    `gorm:"type:varchar(20)" json:"provider_type"`
    CreatedAt       time.Time `json:"created_at"`
}

// Budget 预算配置
type Budget struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint      `gorm:"index" json:"cloud_account_id"`                  // 0 表示全局预算
    Name            string    `gorm:"size:100;not null" json:"name"`
    Type            string    `gorm:"type:varchar(20)" json:"type"`                   // monthly/quarterly/yearly
    Amount          float64   `json:"amount"`
    AlertThreshold  float64   `json:"alert_threshold"`                                // 预警阈值百分比 (80%)
    CurrentUsage    float64   `json:"current_usage"`
    Status          string    `gorm:"type:varchar(20);default:'active'" json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// CostAnomaly 成本异常记录
type CostAnomaly struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    CloudAccountID  uint      `gorm:"index;not null" json:"cloud_account_id"`
    AnomalyType     string    `gorm:"type:varchar(50)" json:"anomaly_type"`           // spike/drop/unusual_pattern
    DetectedAt      time.Time `json:"detected_at"`
    Period          string    `gorm:"size:20" json:"period"`                          // 2026-04-13
    ExpectedCost    float64   `json:"expected_cost"`
    ActualCost      float64   `json:"actual_cost"`
    DeviationRate   float64   `json:"deviation_rate"`                                 // 偏差率
    Severity        string    `gorm:"type:varchar(20)" json:"severity"`               // low/medium/high
    Status          string    `gorm:"type:varchar(20)" json:"status"`                 // new/confirmed/resolved
    Resolution      string    `gorm:"size:500" json:"resolution"`
    CreatedAt       time.Time `json:"created_at"`
}
```

### API 端点设计

| 功能 | 端点 | 方法 | 说明 |
|------|------|------|------|
| **账单查看** | `/finance/bills` | GET | 获取账单列表 |
| | `/finance/bills/sync` | POST | 同步云厂商账单 |
| **账单导出** | `/finance/bills/export` | POST | 导出账单 CSV/Excel |
| **我的订单** | `/finance/orders` | GET | 获取订单列表 |
| | `/finance/orders/sync` | POST | 同步云厂商订单 |
| **续费管理** | `/finance/renewals` | GET | 获取待续费资源 |
| **成本分析** | `/finance/cost/analysis` | GET | 成本趋势分析 |
| **成本报告** | `/finance/cost/reports` | GET | 成本报告列表 |
| | `/finance/cost/reports/generate` | POST | 生成成本报告 |
| **预算管理** | `/finance/budgets` | GET/POST | 预算 CRUD |
| | `/finance/budgets/:id` | PUT/DELETE | 预算更新/删除 |
| **异常监测** | `/finance/anomalies` | GET | 异常列表 |
| | `/finance/anomalies/:id/resolve` | POST | 处理异常 |

---

## 云厂商 API 对接

### 阿里云账单 API

使用阿里云 BSS (Business Support System) API：
- `QueryAccountBill`: 查询账户账单
- `QueryInstanceBill`: 查询实例账单
- `QueryBill`: 查询账单详情
- `QueryOrders`: 查询订单列表

### 费用同步流程

```
用户点击"同步数据" -> 前端调用 /finance/bills/sync
    -> 后端获取云账号配置 -> 初始化云厂商 SDK
    -> 调用云厂商账单 API -> 解析返回数据
    -> 存储到数据库 -> 返回同步结果
```

---

## 实现任务清单

- [ ] **Task 1: 添加费用中心路由** - router.ts 添加三级菜单路由
- [ ] **Task 2: 创建 9 个骨架页面目录和组件**
- [ ] **Task 3: 添加费用中心数据模型** - model/finance.go
- [ ] **Task 4: 创建费用中心 API 函数** - api/finance.ts
- [ ] **Task 5: 创建费用中心 Handler** - handler/finance.go
- [ ] **Task 6: 创建费用中心 Service** - service/finance.go
- [ ] **Task 7: 注册路由** - main.go
- [ ] **Task 8: 编译验证** - 前端和后端