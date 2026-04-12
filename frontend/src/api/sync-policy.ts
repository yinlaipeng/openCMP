import request from '@/utils/request'
import type { SyncPolicy, CreateSyncPolicyRequest } from '@/types/sync-policy'

export function getSyncPolicies(params?: { limit?: number; offset?: number }) {
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