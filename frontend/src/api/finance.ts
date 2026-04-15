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

// 同步待续费资源
export function syncRenewals(cloudAccountId: number, daysThreshold?: number) {
  return request<{ message: string; count: number }>({
    url: '/finance/renewals/sync',
    method: 'post',
    data: { cloud_account_id: cloudAccountId, days_threshold: daysThreshold || 30 }
  })
}

// ========== 账户余额 ==========

export interface AccountBalance {
  balance: number
  currency: string
  credit_limit: number
  status: string
}

export function getAccountBalance(cloudAccountId: number) {
  return request<AccountBalance>({
    url: '/finance/account-balance',
    method: 'get',
    params: { cloud_account_id: cloudAccountId }
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

// 生成成本报告
export function generateCostReport(params: {
  cloud_account_id?: number
  start_date?: string
  end_date?: string
  report_type?: string
}) {
  return request<{ report_id: string; message: string }>({
    url: '/finance/cost/reports/generate',
    method: 'post',
    data: params
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