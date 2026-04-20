import request from '@/utils/request'
import type { SyncPolicy, CreateSyncPolicyRequest } from '@/types/sync-policy'

export function getSyncPolicies(params?: { limit?: number; offset?: number; name?: string; condition_type?: string; enabled?: boolean }) {
  return request<{ items: SyncPolicy[], total: number }>({ url: '/sync-policies', method: 'get', params })
}

export function getSyncPolicy(id: number) {
  return request<SyncPolicy>({ url: `/sync-policies/${id}`, method: 'get' })
}

export function createSyncPolicy(data: CreateSyncPolicyRequest) {
  return request<SyncPolicy>({ url: '/sync-policies', method: 'post', data })
}

export function updateSyncPolicy(id: number, data: Partial<CreateSyncPolicyRequest>) {
  return request<SyncPolicy>({ url: `/sync-policies/${id}`, method: 'put', data })
}

export function deleteSyncPolicy(id: number) {
  return request({ url: `/sync-policies/${id}`, method: 'delete' })
}

export function updateSyncPolicyStatus(id: number, enabled: boolean) {
  return request({ url: `/sync-policies/${id}/status`, method: 'post', data: { enabled } })
}

// 执行同步策略
export interface ExecuteResult {
  id: number
  sync_policy_id: number
  execution_time: string
  trigger_type: string
  resource_count: number
  matched_count: number
  mapped_count: number
  result: string
  duration: number
}

export function executeSyncPolicy(id: number, data?: { cloud_account_id?: number; operator?: string }) {
  return request<ExecuteResult>({ url: `/sync-policies/${id}/execute`, method: 'post', data })
}

// 获取执行日志
export interface ExecutionLog {
  id: number
  sync_policy_id: number
  execution_time: string
  trigger_type: string
  resource_count: number
  matched_count: number
  mapped_count: number
  result: string
  duration: number
  error_message?: string
  operator: string
}

export function getExecutionLogs(id: number, params?: { limit?: number; offset?: number; result?: string }) {
  return request<{ items: ExecutionLog[], total: number }>({ url: `/sync-policies/${id}/execution-logs`, method: 'get', params })
}

// 获取映射结果
export interface MappingResult {
  id: number
  sync_policy_id: number
  execution_log_id: number
  resource_name: string
  resource_type: string
  cloud_account_id: number
  matched_rule_id: number
  matched_tags: string
  target_project_id: number
  target_project_name: string
  mapped_at: string
}

export function getMappingResults(id: number, params?: { limit?: number; offset?: number; project_id?: string }) {
  return request<{ items: MappingResult[], total: number }>({ url: `/sync-policies/${id}/mapping-results`, method: 'get', params })
}

// 批量操作
export function batchEnableSyncPolicies(ids: number[]) {
  return request<{ message: string; count: number }>({ url: '/sync-policies/batch-enable', method: 'post', data: { ids } })
}

export function batchDisableSyncPolicies(ids: number[]) {
  return request<{ message: string; count: number }>({ url: '/sync-policies/batch-disable', method: 'post', data: { ids } })
}

export function batchDeleteSyncPolicies(ids: number[]) {
  return request<{ message: string; count: number }>({ url: '/sync-policies/batch-delete', method: 'post', data: { ids } })
}

// 导出
export function exportSyncPolicies() {
  return request<{ items: SyncPolicy[], total: number }>({ url: '/sync-policies/export', method: 'get' })
}