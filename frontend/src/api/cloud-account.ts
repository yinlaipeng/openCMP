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
