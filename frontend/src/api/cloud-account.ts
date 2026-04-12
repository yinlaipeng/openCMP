import request from '@/utils/request'
import type { CloudAccount, CreateCloudAccountRequest } from '@/types'

export function getCloudAccounts(params?: { page?: number; page_size?: number }) {
  return request<{ items: CloudAccount[], total: number }>({
    url: '/cloud-accounts',
    method: 'get',
    params
  })
}

export function getCloudAccount(id: number) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}`,
    method: 'get'
  })
}

export function createCloudAccount(data: CreateCloudAccountRequest) {
  return request<CloudAccount>({
    url: '/cloud-accounts',
    method: 'post',
    data
  })
}

export function updateCloudAccount(id: number, data: Partial<CreateCloudAccountRequest>) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}`,
    method: 'put',
    data
  })
}

export function deleteCloudAccount(id: number) {
  return request({
    url: `/cloud-accounts/${id}`,
    method: 'delete'
  })
}

export function verifyCloudAccount(id: number) {
  return request({
    url: `/cloud-accounts/${id}/verify`,
    method: 'post'
  })
}

// 同步云账号
export interface SyncCloudAccountOptions {
  mode: 'full' | 'incremental'
  resource_types: string[]
}

export function syncCloudAccount(id: number, options: SyncCloudAccountOptions) {
  return request({
    url: `/cloud-accounts/${id}/sync`,
    method: 'post',
    data: options
  })
}

// 测试连接
export function testConnection(id: number) {
  return request<{ connected: boolean; message: string }>({
    url: `/cloud-accounts/${id}/test-connection`,
    method: 'post'
  })
}

// 验证凭证（实际调用云厂商 API）
export function verifyCredentials(id: number) {
  return request<{ valid: boolean; message: string }>({
    url: `/cloud-accounts/${id}/verify-credentials`,
    method: 'get'
  })
}

// 更新云账号状态
export function updateCloudAccountStatus(id: number, enabled: boolean) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}/status`,
    method: 'patch',
    data: { enabled }
  })
}

// 更新云账号属性
export interface UpdateCloudAccountAttributes {
  auto_sync?: boolean
  sync_policy?: string
  sync_interval?: number
  sync_resource_types?: string[]
}

export function updateCloudAccountAttributes(id: number, attrs: UpdateCloudAccountAttributes) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}/attributes`,
    method: 'patch',
    data: attrs
  })
}
