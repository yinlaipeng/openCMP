# 费用中心模块骨架实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现费用中心一级菜单及 9 个子页面的骨架结构，包含路由、前端页面组件、后端数据模型和基础 API。

**Architecture:** 三级菜单嵌套结构（费用中心 -> 二级菜单 -> 子页面），数据关联现有云账号，使用 Element Plus 组件构建页面骨架。

**Tech Stack:** Vue 3 + Element Plus + TypeScript (前端), Go + Gin + Gorm (后端)

---

## 文件结构

**前端新增文件：**
```
frontend/src/views/finance/
├── orders/
│   ├── my-orders/index.vue
│   └── renewals/index.vue
├── bills/
│   ├── view/index.vue
│   └── export/index.vue
└── cost/
    ├── analysis/index.vue
    ├── reports/index.vue
    ├── budgets/index.vue
    └── anomaly/index.vue

frontend/src/api/finance.ts
frontend/src/types/finance.ts
```

**后端新增/修改文件：**
```
backend/internal/model/finance.go
backend/internal/handler/finance.go
backend/internal/service/finance.go
backend/cmd/server/main.go (修改 - 注册路由)
```

---

## Task 1: 添加费用中心路由配置

**Files:**
- Modify: `frontend/src/router.ts`

- [ ] **Step 1: 添加费用中心路由**

在 `router.ts` 的 `routes` 数组中，在 `message-center` 路由之后添加费用中心路由：

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

- [ ] **Step 2: 添加 Wallet 图标导入**

在 router.ts 开头添加图标导入（如果未导入）：

```typescript
import { Wallet } from '@element-plus/icons-vue'
```

---

## Task 2: 创建费用中心 TypeScript 类型定义

**Files:**
- Create: `frontend/src/types/finance.ts`

- [ ] **Step 1: 创建类型定义文件**

```typescript
// frontend/src/types/finance.ts

// 账单记录
export interface Bill {
  id: number
  cloud_account_id: number
  billing_cycle: string
  product_type: string
  product_name: string
  instance_id: string
  usage_amount: number
  unit_price: number
  total_cost: number
  currency: string
  billing_method: string
  status: string
  provider_type: string
  created_at: string
}

// 订单记录
export interface Order {
  id: number
  cloud_account_id: number
  order_id: string
  order_type: string
  product_type: string
  product_name: string
  instance_id: string
  amount: number
  currency: string
  status: string
  payment_time?: string
  effective_time: string
  expire_time?: string
  provider_type: string
  created_at: string
}

// 预算配置
export interface Budget {
  id: number
  cloud_account_id: number
  name: string
  type: string
  amount: number
  alert_threshold: number
  current_usage: number
  status: string
  created_at: string
  updated_at: string
}

// 成本异常
export interface CostAnomaly {
  id: number
  cloud_account_id: number
  anomaly_type: string
  detected_at: string
  period: string
  expected_cost: number
  actual_cost: number
  deviation_rate: number
  severity: string
  status: string
  resolution?: string
  created_at: string
}

// 续费资源
export interface RenewalResource {
  id: number
  cloud_account_id: number
  instance_id: string
  instance_name: string
  product_type: string
  expire_time: string
  days_remaining: number
  renewal_price: number
  status: string
}

// 成本分析数据
export interface CostAnalysisData {
  period: string
  total_cost: number
  product_costs: Record<string, number>
  trend: number
}
```

---

## Task 3: 创建费用中心 API 函数

**Files:**
- Create: `frontend/src/api/finance.ts`

- [ ] **Step 1: 创建 API 文件**

```typescript
// frontend/src/api/finance.ts
import request from '@/utils/request'
import type { Bill, Order, Budget, CostAnomaly, RenewalResource, CostAnalysisData } from '@/types/finance'

// ========== 账单相关 ==========

// 获取账单列表
export function getBills(params: {
  cloud_account_id?: number
  billing_cycle?: string
  page?: number
  page_size?: number
}) {
  return request<{ items: Bill[]; total: number }>({
    url: '/finance/bills',
    method: 'get',
    params
  })
}

// 同步账单数据
export function syncBills(cloudAccountId: number) {
  return request<{ message: string; count: number }>({
    url: '/finance/bills/sync',
    method: 'post',
    data: { cloud_account_id: cloudAccountId }
  })
}

// 导出账单
export function exportBills(params: {
  cloud_account_id?: number
  billing_cycle?: string
  format?: 'csv' | 'excel'
}) {
  return request<{ download_url: string }>({
    url: '/finance/bills/export',
    method: 'post',
    data: params
  })
}

// ========== 订单相关 ==========

// 获取订单列表
export function getOrders(params: {
  cloud_account_id?: number
  status?: string
  page?: number
  page_size?: number
}) {
  return request<{ items: Order[]; total: number }>({
    url: '/finance/orders',
    method: 'get',
    params
  })
}

// 同步订单数据
export function syncOrders(cloudAccountId: number) {
  return request<{ message: string; count: number }>({
    url: '/finance/orders/sync',
    method: 'post',
    data: { cloud_account_id: cloudAccountId }
  })
}

// ========== 续费管理 ==========

// 获取待续费资源列表
export function getRenewalResources(params: {
  cloud_account_id?: number
  days_threshold?: number
  page?: number
  page_size?: number
}) {
  return request<{ items: RenewalResource[]; total: number }>({
    url: '/finance/renewals',
    method: 'get',
    params
  })
}

// ========== 成本分析 ==========

// 获取成本分析数据
export function getCostAnalysis(params: {
  cloud_account_id?: number
  start_date?: string
  end_date?: string
}) {
  return request<CostAnalysisData[]>({
    url: '/finance/cost/analysis',
    method: 'get',
    params
  })
}

// ========== 成本报告 ==========

// 获取成本报告列表
export function getCostReports(params: {
  page?: number
  page_size?: number
}) {
  return request<{ items: any[]; total: number }>({
    url: '/finance/cost/reports',
    method: 'get',
    params
  })
}

// ========== 预算管理 ==========

// 获取预算列表
export function getBudgets(params?: { cloud_account_id?: number }) {
  return request<Budget[]>({
    url: '/finance/budgets',
    method: 'get',
    params
  })
}

// 创建预算
export function createBudget(data: Partial<Budget>) {
  return request<Budget>({
    url: '/finance/budgets',
    method: 'post',
    data
  })
}

// 更新预算
export function updateBudget(id: number, data: Partial<Budget>) {
  return request<Budget>({
    url: `/finance/budgets/${id}`,
    method: 'put',
    data
  })
}

// 删除预算
export function deleteBudget(id: number) {
  return request<{ message: string }>({
    url: `/finance/budgets/${id}`,
    method: 'delete'
  })
}

// ========== 异常监测 ==========

// 获取异常列表
export function getAnomalies(params: {
  cloud_account_id?: number
  status?: string
  severity?: string
  page?: number
  page_size?: number
}) {
  return request<{ items: CostAnomaly[]; total: number }>({
    url: '/finance/anomalies',
    method: 'get',
    params
  })
}

// 处理异常
export function resolveAnomaly(id: number, resolution: string) {
  return request<{ message: string }>({
    url: `/finance/anomalies/${id}/resolve`,
    method: 'post',
    data: { resolution }
  })
}
```

---

## Task 4: 创建后端数据模型

**Files:**
- Create: `backend/internal/model/finance.go`

- [ ] **Step 1: 创建模型文件**

```go
// backend/internal/model/finance.go
package model

import (
	"time"
)

// Bill 账单记录
type Bill struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	BillingCycle   string    `gorm:"type:varchar(20);not null" json:"billing_cycle"` // 2026-04
	ProductType    string    `gorm:"type:varchar(50)" json:"product_type"`
	ProductName    string    `gorm:"size:200" json:"product_name"`
	InstanceID     string    `gorm:"size:100" json:"instance_id"`
	UsageAmount    float64   `json:"usage_amount"`
	UnitPrice      float64   `json:"unit_price"`
	TotalCost      float64   `json:"total_cost"`
	Currency       string    `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	BillingMethod  string    `gorm:"type:varchar(20)" json:"billing_method"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
	ProviderType   string    `gorm:"type:varchar(20)" json:"provider_type"`
	CreatedAt      time.Time `json:"created_at"`
}

// Order 订单记录
type Order struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	OrderID       string     `gorm:"uniqueIndex;size:100" json:"order_id"`
	OrderType     string     `gorm:"type:varchar(20)" json:"order_type"`
	ProductType   string     `gorm:"type:varchar(50)" json:"product_type"`
	ProductName   string     `gorm:"size:200" json:"product_name"`
	InstanceID    string     `gorm:"size:100" json:"instance_id"`
	Amount        float64    `json:"amount"`
	Currency      string     `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	Status        string     `gorm:"type:varchar(20)" json:"status"`
	PaymentTime   *time.Time `json:"payment_time"`
	EffectiveTime time.Time  `json:"effective_time"`
	ExpireTime    *time.Time `json:"expire_time"`
	ProviderType  string     `gorm:"type:varchar(20)" json:"provider_type"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Budget 预算配置
type Budget struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index" json:"cloud_account_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Type           string    `gorm:"type:varchar(20)" json:"type"`
	Amount         float64   `json:"amount"`
	AlertThreshold float64   `json:"alert_threshold"`
	CurrentUsage   float64   `json:"current_usage"`
	Status         string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CostAnomaly 成本异常记录
type CostAnomaly struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	AnomalyType    string    `gorm:"type:varchar(50)" json:"anomaly_type"`
	DetectedAt     time.Time `json:"detected_at"`
	Period         string    `gorm:"size:20" json:"period"`
	ExpectedCost   float64   `json:"expected_cost"`
	ActualCost     float64   `json:"actual_cost"`
	DeviationRate  float64   `json:"deviation_rate"`
	Severity       string    `gorm:"type:varchar(20)" json:"severity"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
	Resolution     string    `gorm:"size:500" json:"resolution"`
	CreatedAt      time.Time `json:"created_at"`
}

// RenewalResource 待续费资源
type RenewalResource struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	InstanceID    string     `gorm:"size:100" json:"instance_id"`
	InstanceName  string     `gorm:"size:200" json:"instance_name"`
	ProductType   string     `gorm:"type:varchar(50)" json:"product_type"`
	ExpireTime    time.Time  `json:"expire_time"`
	DaysRemaining int        `json:"days_remaining"`
	RenewalPrice  float64    `json:"renewal_price"`
	Status        string     `gorm:"type:varchar(20)" json:"status"`
}

// TableName 方法
func (Bill) TableName() string { return "finance_bills" }
func (Order) TableName() string { return "finance_orders" }
func (Budget) TableName() string { return "finance_budgets" }
func (CostAnomaly) TableName() string { return "finance_cost_anomalies" }
func (RenewalResource) TableName() string { return "finance_renewal_resources" }
```

---

## Task 5: 创建后端 Handler

**Files:**
- Create: `backend/internal/handler/finance.go`

- [ ] **Step 1: 创建 Handler 文件**

```go
// backend/internal/handler/finance.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

type FinanceHandler struct {
	service *service.FinanceService
	logger  *zap.Logger
}

func NewFinanceHandler(db *gorm.DB, logger *zap.Logger) *FinanceHandler {
	return &FinanceHandler{
		service: service.NewFinanceService(db),
		logger:  logger,
	}
}

// ========== 账单相关 ==========

// GetBills 获取账单列表
func (h *FinanceHandler) GetBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	billingCycle := c.Query("billing_cycle")

	bills, total, err := h.service.GetBills(c.Request.Context(), uint(cloudAccountID), billingCycle, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": bills,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// SyncBills 同步账单数据
func (h *FinanceHandler) SyncBills(c *gin.Context) {
	var req struct {
		CloudAccountID uint `json:"cloud_account_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.SyncBills(c.Request.Context(), req.CloudAccountID)
	if err != nil {
		h.logger.Error("failed to sync bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

// ExportBills 导出账单
func (h *FinanceHandler) ExportBills(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id"`
		BillingCycle   string `json:"billing_cycle"`
		Format         string `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.ExportBills(c.Request.Context(), req.CloudAccountID, req.BillingCycle, req.Format)
	if err != nil {
		h.logger.Error("failed to export bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// ========== 订单相关 ==========

// GetOrders 获取订单列表
func (h *FinanceHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")

	orders, total, err := h.service.GetOrders(c.Request.Context(), uint(cloudAccountID), status, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": orders,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// SyncOrders 同步订单数据
func (h *FinanceHandler) SyncOrders(c *gin.Context) {
	var req struct {
		CloudAccountID uint `json:"cloud_account_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.SyncOrders(c.Request.Context(), req.CloudAccountID)
	if err != nil {
		h.logger.Error("failed to sync orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

// ========== 续费管理 ==========

// GetRenewals 获取待续费资源列表
func (h *FinanceHandler) GetRenewals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	daysThreshold, _ := strconv.Atoi(c.DefaultQuery("days_threshold", "30"))

	renewals, total, err := h.service.GetRenewals(c.Request.Context(), uint(cloudAccountID), daysThreshold, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get renewals", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": renewals,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// ========== 成本分析 ==========

// GetCostAnalysis 获取成本分析数据
func (h *FinanceHandler) GetCostAnalysis(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.service.GetCostAnalysis(c.Request.Context(), uint(cloudAccountID), startDate, endDate)
	if err != nil {
		h.logger.Error("failed to get cost analysis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ========== 成本报告 ==========

// GetCostReports 获取成本报告列表
func (h *FinanceHandler) GetCostReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reports, total, err := h.service.GetCostReports(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("failed to get cost reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": reports,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// ========== 预算管理 ==========

// GetBudgets 获取预算列表
func (h *FinanceHandler) GetBudgets(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	budgets, err := h.service.GetBudgets(c.Request.Context(), uint(cloudAccountID))
	if err != nil {
		h.logger.Error("failed to get budgets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// CreateBudget 创建预算
func (h *FinanceHandler) CreateBudget(c *gin.Context) {
	var budget model.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBudget(c.Request.Context(), &budget); err != nil {
		h.logger.Error("failed to create budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// UpdateBudget 更新预算
func (h *FinanceHandler) UpdateBudget(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var budget model.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget.ID = uint(id)

	if err := h.service.UpdateBudget(c.Request.Context(), &budget); err != nil {
		h.logger.Error("failed to update budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// DeleteBudget 删除预算
func (h *FinanceHandler) DeleteBudget(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.service.DeleteBudget(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ========== 异常监测 ==========

// GetAnomalies 获取异常列表
func (h *FinanceHandler) GetAnomalies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	severity := c.Query("severity")

	anomalies, total, err := h.service.GetAnomalies(c.Request.Context(), uint(cloudAccountID), status, severity, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get anomalies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": anomalies,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// ResolveAnomaly 处理异常
func (h *FinanceHandler) ResolveAnomaly(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Resolution string `json:"resolution" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResolveAnomaly(c.Request.Context(), uint(id), req.Resolution); err != nil {
		h.logger.Error("failed to resolve anomaly", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "resolved"})
}
```

---

## Task 6: 创建后端 Service

**Files:**
- Create: `backend/internal/service/finance.go`

- [ ] **Step 1: 创建 Service 文件**

```go
// backend/internal/service/finance.go
package service

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

type FinanceService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{
		db:     db,
		logger: zap.L(),
	}
}

// ========== 账单相关 ==========

func (s *FinanceService) GetBills(ctx context.Context, cloudAccountID uint, billingCycle string, page, pageSize int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	query := s.db.Model(&model.Bill{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if billingCycle != "" {
		query = query.Where("billing_cycle = ?", billingCycle)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (s *FinanceService) SyncBills(ctx context.Context, cloudAccountID uint) (int, error) {
	// TODO: 实现云厂商 API 调用同步账单
	// 当前返回模拟数据
	return 0, nil
}

func (s *FinanceService) ExportBills(ctx context.Context, cloudAccountID uint, billingCycle, format string) (string, error) {
	// TODO: 实现账单导出功能
	return "", nil
}

// ========== 订单相关 ==========

func (s *FinanceService) GetOrders(ctx context.Context, cloudAccountID uint, status string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := s.db.Model(&model.Order{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *FinanceService) SyncOrders(ctx context.Context, cloudAccountID uint) (int, error) {
	// TODO: 实现云厂商 API 调用同步订单
	return 0, nil
}

// ========== 续费管理 ==========

func (s *FinanceService) GetRenewals(ctx context.Context, cloudAccountID uint, daysThreshold, page, pageSize int) ([]model.RenewalResource, int64, error) {
	var renewals []model.RenewalResource
	var total int64

	query := s.db.Model(&model.RenewalResource{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	query = query.Where("days_remaining <= ?", daysThreshold)

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("expire_time asc").Find(&renewals).Error; err != nil {
		return nil, 0, err
	}

	return renewals, total, nil
}

// ========== 成本分析 ==========

func (s *FinanceService) GetCostAnalysis(ctx context.Context, cloudAccountID uint, startDate, endDate string) ([]map[string]interface{}, error) {
	// TODO: 实现成本分析聚合查询
	return []map[string]interface{}{}, nil
}

// ========== 成本报告 ==========

func (s *FinanceService) GetCostReports(ctx context.Context, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// TODO: 实现成本报告列表
	return []map[string]interface{}{}, 0, nil
}

// ========== 预算管理 ==========

func (s *FinanceService) GetBudgets(ctx context.Context, cloudAccountID uint) ([]model.Budget, error) {
	var budgets []model.Budget
	query := s.db.Model(&model.Budget{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if err := query.Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (s *FinanceService) CreateBudget(ctx context.Context, budget *model.Budget) error {
	return s.db.Create(budget).Error
}

func (s *FinanceService) UpdateBudget(ctx context.Context, budget *model.Budget) error {
	return s.db.Save(budget).Error
}

func (s *FinanceService) DeleteBudget(ctx context.Context, id uint) error {
	return s.db.Delete(&model.Budget{}, id).Error
}

// ========== 异常监测 ==========

func (s *FinanceService) GetAnomalies(ctx context.Context, cloudAccountID uint, status, severity string, page, pageSize int) ([]model.CostAnomaly, int64, error) {
	var anomalies []model.CostAnomaly
	var total int64

	query := s.db.Model(&model.CostAnomaly{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("detected_at desc").Find(&anomalies).Error; err != nil {
		return nil, 0, err
	}

	return anomalies, total, nil
}

func (s *FinanceService) ResolveAnomaly(ctx context.Context, id uint, resolution string) error {
	return s.db.Model(&model.CostAnomaly{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "resolved",
		"resolution": resolution,
	}).Error
}
```

---

## Task 7: 注册路由

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 导入 Finance Handler**

在 main.go 的 handler 初始化部分添加：

```go
financeHandler := handler.NewFinanceHandler(db, logger)
```

- [ ] **Step 2: 注册费用中心路由组**

在路由注册部分添加：

```go
// 费用中心路由
finance := api.Group("/finance")
{
	// 账单
	finance.GET("/bills", financeHandler.GetBills)
	finance.POST("/bills/sync", financeHandler.SyncBills)
	finance.POST("/bills/export", financeHandler.ExportBills)
	
	// 订单
	finance.GET("/orders", financeHandler.GetOrders)
	finance.POST("/orders/sync", financeHandler.SyncOrders)
	
	// 续费
	finance.GET("/renewals", financeHandler.GetRenewals)
	
	// 成本分析
	finance.GET("/cost/analysis", financeHandler.GetCostAnalysis)
	finance.GET("/cost/reports", financeHandler.GetCostReports)
	
	// 预算
	finance.GET("/budgets", financeHandler.GetBudgets)
	finance.POST("/budgets", financeHandler.CreateBudget)
	finance.PUT("/budgets/:id", financeHandler.UpdateBudget)
	finance.DELETE("/budgets/:id", financeHandler.DeleteBudget)
	
	// 异常
	finance.GET("/anomalies", financeHandler.GetAnomalies)
	finance.POST("/anomalies/:id/resolve", financeHandler.ResolveAnomaly)
}
```

---

## Task 8: 创建前端骨架页面 - 我的订单

**Files:**
- Create: `frontend/src/views/finance/orders/my-orders/index.vue`

- [ ] **Step 1: 创建页面组件**

```vue
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">我的订单</span>
          <el-button type="primary" @click="handleSync">同步数据</el-button>
        </div>
      </template>
      
      <!-- 云账号筛选 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px;">
          <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
        </el-select>
        <el-select v-model="selectedStatus" placeholder="订单状态" clearable style="width: 150px; margin-left: 8px;">
          <el-option label="待支付" value="pending" />
          <el-option label="已支付" value="paid" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </div>
      
      <!-- 数据表格 -->
      <el-table :data="orders" v-loading="loading" style="width: 100%">
        <el-table-column prop="order_id" label="订单号" width="180" />
        <el-table-column prop="order_type" label="订单类型" width="100" />
        <el-table-column prop="product_name" label="产品名称" min-width="200" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="effective_time" label="生效时间" width="160" />
        <el-table-column prop="expire_time" label="到期时间" width="160" />
        <el-table-column prop="provider_type" label="云平台" width="100" />
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getOrders, syncOrders } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { Order } from '@/types/finance'

const orders = ref<Order[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])
const selectedAccountId = ref<number | undefined>()
const selectedStatus = ref<string>('')

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消'
  }
  return map[status] || status
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders({
      cloud_account_id: selectedAccountId.value,
      status: selectedStatus.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    orders.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const handleSync = async () => {
  if (!selectedAccountId.value) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    await syncOrders(selectedAccountId.value)
    ElMessage.success('订单数据同步成功')
    loadOrders()
  } catch (e) {
    ElMessage.error('同步失败')
  }
}

watch([selectedAccountId, selectedStatus, pagination.currentPage, pagination.pageSize], loadOrders)

onMounted(() => {
  loadCloudAccounts()
  loadOrders()
})
</script>

<style scoped>
.finance-page {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.filter-bar {
  display: flex;
  align-items: center;
}
</style>
```

---

## Task 9: 创建前端骨架页面 - 账单查看

**Files:**
- Create: `frontend/src/views/finance/bills/view/index.vue`

- [ ] **Step 1: 创建页面组件**

```vue
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">账单查看</span>
          <el-button type="primary" @click="handleSync">同步账单</el-button>
        </div>
      </template>
      
      <!-- 筛选 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px;">
          <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
        </el-select>
        <el-date-picker
          v-model="selectedBillingCycle"
          type="month"
          placeholder="选择账单周期"
          format="YYYY-MM"
          value-format="YYYY-MM"
          style="width: 150px; margin-left: 8px;"
        />
      </div>
      
      <!-- 统计卡片 -->
      <el-row :gutter="16" style="margin-bottom: 16px;">
        <el-col :span="6">
          <el-statistic title="本月总费用" :value="totalCost" suffix="元" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="账单数量" :value="pagination.total" />
        </el-col>
      </el-row>
      
      <!-- 数据表格 -->
      <el-table :data="bills" v-loading="loading" style="width: 100%">
        <el-table-column prop="billing_cycle" label="账期" width="100" />
        <el-table-column prop="product_type" label="产品类型" width="120" />
        <el-table-column prop="product_name" label="产品名称" min-width="200" />
        <el-table-column prop="instance_id" label="实例ID" width="150" />
        <el-table-column prop="usage_amount" label="用量" width="100" />
        <el-table-column prop="total_cost" label="费用" width="120">
          <template #default="{ row }">
            ¥{{ row.total_cost?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="billing_method" label="计费方式" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'paid' ? 'success' : 'warning'">
              {{ row.status === 'paid' ? '已支付' : '待支付' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getBills, syncBills } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { Bill } from '@/types/finance'

const bills = ref<Bill[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])
const selectedAccountId = ref<number | undefined>()
const selectedBillingCycle = ref<string>('')

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const totalCost = computed(() => {
  return bills.value.reduce((sum, bill) => sum + (bill.total_cost || 0), 0)
})

const loadBills = async () => {
  loading.value = true
  try {
    const res = await getBills({
      cloud_account_id: selectedAccountId.value,
      billing_cycle: selectedBillingCycle.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    bills.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const handleSync = async () => {
  if (!selectedAccountId.value) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    await syncBills(selectedAccountId.value)
    ElMessage.success('账单数据同步成功')
    loadBills()
  } catch (e) {
    ElMessage.error('同步失败')
  }
}

watch([selectedAccountId, selectedBillingCycle, pagination.currentPage, pagination.pageSize], loadBills)

onMounted(() => {
  loadCloudAccounts()
  loadBills()
})
</script>
```

---

## Task 10: 创建剩余 7 个骨架页面

**Files:**
- Create: `frontend/src/views/finance/orders/renewals/index.vue`
- Create: `frontend/src/views/finance/bills/export/index.vue`
- Create: `frontend/src/views/finance/cost/analysis/index.vue`
- Create: `frontend/src/views/finance/cost/reports/index.vue`
- Create: `frontend/src/views/finance/cost/budgets/index.vue`
- Create: `frontend/src/views/finance/cost/anomaly/index.vue`

- [ ] **Step 1: 创建续费管理页面**

```vue
<!-- frontend/src/views/finance/orders/renewals/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">续费管理</span>
        </div>
      </template>
      
      <el-table :data="renewals" v-loading="loading">
        <el-table-column prop="instance_name" label="实例名称" min-width="200" />
        <el-table-column prop="product_type" label="产品类型" width="120" />
        <el-table-column prop="expire_time" label="到期时间" width="160" />
        <el-table-column prop="days_remaining" label="剩余天数" width="100">
          <template #default="{ row }">
            <el-tag :type="row.days_remaining <= 7 ? 'danger' : row.days_remaining <= 30 ? 'warning' : 'success'">
              {{ row.days_remaining }}天
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="renewal_price" label="续费价格" width="120">
          <template #default="{ row }">¥{{ row.renewal_price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getRenewalResources } from '@/api/finance'
import type { RenewalResource } from '@/types/finance'

const renewals = ref<RenewalResource[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getRenewalResources({ days_threshold: 30 })
    renewals.value = res.items || []
  } finally {
    loading.value = false
  }
})
</script>
```

- [ ] **Step 2: 创建账单导出中心页面**

```vue
<!-- frontend/src/views/finance/bills/export/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <span class="title">账单导出中心</span>
      </template>
      
      <el-form :model="exportForm" label-width="100px">
        <el-form-item label="云账号">
          <el-select v-model="exportForm.cloud_account_id" placeholder="选择云账号" clearable>
            <el-option v-for="a in cloudAccounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单周期">
          <el-date-picker v-model="exportForm.billing_cycle" type="month" format="YYYY-MM" value-format="YYYY-MM" />
        </el-form-item>
        <el-form-item label="导出格式">
          <el-radio-group v-model="exportForm.format">
            <el-radio label="csv">CSV</el-radio>
            <el-radio label="excel">Excel</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleExport">导出账单</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { exportBills } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'

const cloudAccounts = ref<any[]>([])
const exportForm = ref({
  cloud_account_id: undefined as number | undefined,
  billing_cycle: '',
  format: 'excel'
})

const handleExport = async () => {
  try {
    const res = await exportBills(exportForm.value)
    ElMessage.success('导出成功：' + res.download_url)
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

onMounted(async () => {
  const res = await getCloudAccounts({ page: 1, page_size: 100 })
  cloudAccounts.value = res.items || []
})
</script>
```

- [ ] **Step 3: 创建成本分析页面**

```vue
<!-- frontend/src/views/finance/cost/analysis/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <span class="title">成本分析</span>
      </template>
      
      <el-empty description="成本分析图表开发中..." />
    </el-card>
  </div>
</template>
```

- [ ] **Step 4: 创建成本报告页面**

```vue
<!-- frontend/src/views/finance/cost/reports/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <span class="title">成本报告</span>
      </template>
      <el-empty description="成本报告功能开发中..." />
    </el-card>
  </div>
</template>
```

- [ ] **Step 5: 创建预算管理页面**

```vue
<!-- frontend/src/views/finance/cost/budgets/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">预算管理</span>
          <el-button type="primary">新建预算</el-button>
        </div>
      </template>
      <el-empty description="预算管理功能开发中..." />
    </el-card>
  </div>
</template>
```

- [ ] **Step 6: 创建异常监测页面**

```vue
<!-- frontend/src/views/finance/cost/anomaly/index.vue -->
<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <span class="title">异常监测</span>
      </template>
      <el-empty description="异常监测功能开发中..." />
    </el-card>
  </div>
</template>
```

---

## Task 11: 编译验证

- [ ] **Step 1: 后端编译验证**

Run: `cd /Users/aurora/Desktop/xtwork/git/openCMP/backend && go build ./cmd/server`
Expected: 编译成功，无错误

- [ ] **Step 2: 前端编译验证**

Run: `cd /Users/aurora/Desktop/xtwork/git/openCMP/frontend && node_modules/.bin/vite build`
Expected: 编译成功，无错误

---

## Task 12: 提交代码

- [ ] **Step 1: Git 提交**

```bash
git add frontend/src/router.ts frontend/src/types/finance.ts frontend/src/api/finance.ts frontend/src/views/finance/
git add backend/internal/model/finance.go backend/internal/handler/finance.go backend/internal/service/finance.go backend/cmd/server/main.go
git commit -m "feat: add cost center module skeleton with menu and basic pages

- Add finance center menu with 3 sub-menus: Orders, Bills, Cost Management
- Create 9 skeleton pages for finance module
- Add backend models: Bill, Order, Budget, CostAnomaly, RenewalResource
- Add backend handler and service for finance APIs
- Register finance routes in main.go

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```