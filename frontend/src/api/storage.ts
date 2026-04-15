// frontend/src/api/storage.ts
import request from '@/utils/request'

// ========== 云硬盘 API ==========

export interface CloudDisk {
  id: number
  cloud_account_id: number
  disk_id: string
  name: string
  size: number
  type: string
  status: string
  vm_id: string
  zone_id: string
  provider_type: string
  created_at: string
  updated_at: string
}

export function getCloudDisks(params?: {
  cloud_account_id?: number
  status?: string
  page?: number
  page_size?: number
}) {
  return request<{ items: CloudDisk[]; total: number; page: number; page_size: number }>({
    url: '/storage/cloud-disks',
    method: 'get',
    params
  })
}

export function createCloudDisk(params: {
  cloud_account_id: number
  name: string
  size: number
  type?: string
  zone_id: string
}) {
  return request<CloudDisk>({
    url: '/storage/cloud-disks',
    method: 'post',
    params: { cloud_account_id: params.cloud_account_id },
    data: {
      name: params.name,
      size: params.size,
      type: params.type,
      zone_id: params.zone_id
    }
  })
}

export function deleteCloudDisk(id: number) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}`,
    method: 'delete'
  })
}

export function attachCloudDisk(id: number, vmId: string) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}/attach`,
    method: 'post',
    data: { vm_id: vmId }
  })
}

export function detachCloudDisk(id: number) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}/detach`,
    method: 'post'
  })
}

export function resizeCloudDisk(id: number, newSize: number) {
  return request<{ message: string; new_size: number }>({
    url: `/storage/cloud-disks/${id}/resize`,
    method: 'post',
    data: { new_size: newSize }
  })
}

export function syncCloudDisks(cloudAccountId: number) {
  return request<{ message: string; count: number; total: number }>({
    url: '/storage/cloud-disks/sync',
    method: 'post',
    params: { cloud_account_id: cloudAccountId }
  })
}

// ========== 云快照 API ==========

export interface CloudSnapshot {
  id: number
  cloud_account_id: number
  snapshot_id: string
  name: string
  disk_id: number
  size: number
  status: string
  provider_type: string
  created_at: string
}

export function getCloudSnapshots(params?: {
  cloud_account_id?: number
  page?: number
  page_size?: number
}) {
  return request<{ items: CloudSnapshot[]; total: number; page: number; page_size: number }>({
    url: '/storage/cloud-snapshots',
    method: 'get',
    params
  })
}

export function createCloudSnapshot(diskId: number, name: string) {
  return request<CloudSnapshot>({
    url: '/storage/cloud-snapshots',
    method: 'post',
    params: { disk_id: diskId },
    data: { name }
  })
}

export function deleteCloudSnapshot(id: number) {
  return request<{ message: string }>({
    url: `/storage/cloud-snapshots/${id}`,
    method: 'delete'
  })
}