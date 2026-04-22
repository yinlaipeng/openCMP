// frontend/src/api/storage.ts
import request from '@/utils/request'

// ========== 云硬盘 API ==========

export interface CloudDisk {
  id: number | string
  cloud_account_id: number
  disk_id: string
  name: string
  size: number // GB
  max_iops: number
  disk_format: string
  type: string // system/data
  storage_type: string // SSD/HDD
  medium_type: string
  status: string // available/in_use/creating/error
  vm_id: string
  vm_name: string
  device_name: string
  primary_storage: string
  zone_id: string
  region_id: string
  billing_type: string
  shutdown_reset: boolean
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface DiskListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  type?: string
  region?: string
}

export function getCloudDisks(params?: DiskListParams) {
  return request<{ items: CloudDisk[]; total: number; page: number; page_size: number }>({
    url: '/storage/cloud-disks',
    method: 'get',
    params
  })
}

export function getCloudDisk(id: number | string) {
  return request<CloudDisk>({
    url: `/storage/cloud-disks/${id}`,
    method: 'get'
  })
}

export interface CreateDiskParams {
  cloud_account_id: number
  name: string
  size: number
  type?: string
  zone_id: string
  project_id?: number
}

export function createCloudDisk(params: CreateDiskParams) {
  return request<CloudDisk>({
    url: '/storage/cloud-disks',
    method: 'post',
    params: { cloud_account_id: params.cloud_account_id },
    data: {
      name: params.name,
      size: params.size,
      type: params.type,
      zone_id: params.zone_id,
      project_id: params.project_id
    }
  })
}

export function updateCloudDisk(id: number | string, data: { name?: string; description?: string }) {
  return request<CloudDisk>({
    url: `/storage/cloud-disks/${id}`,
    method: 'put',
    data
  })
}

export function deleteCloudDisk(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}`,
    method: 'delete'
  })
}

export function batchDeleteCloudDisks(ids: (number | string)[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/storage/cloud-disks/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export function attachCloudDisk(id: number | string, vmId: string) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}/attach`,
    method: 'post',
    data: { vm_id: vmId }
  })
}

export function detachCloudDisk(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/cloud-disks/${id}/detach`,
    method: 'post'
  })
}

export function resizeCloudDisk(id: number | string, newSize: number) {
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

// ========== 云硬盘快照 API ==========

export interface CloudDiskSnapshot {
  id: number | string
  cloud_account_id: number
  snapshot_id: string
  name: string
  disk_id: string
  disk_name: string
  disk_type: string // system/data
  size: number // GB
  status: string // available/creating/rollbacking/error
  progress: number // 0-100
  vm_id: string
  vm_name: string
  provider_type: string
  account_name: string
  project_name: string
  region_id: string
  tags: { key: string; value: string }[]
  created_at: string
  description?: string
}

export interface SnapshotListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  disk_type?: string
}

export function getCloudDiskSnapshots(params?: SnapshotListParams) {
  return request<{ items: CloudDiskSnapshot[]; total: number; page: number; page_size: number }>({
    url: '/storage/cloud-snapshots',
    method: 'get',
    params
  })
}

export function getCloudDiskSnapshot(id: number | string) {
  return request<CloudDiskSnapshot>({
    url: `/storage/cloud-snapshots/${id}`,
    method: 'get'
  })
}

export interface CreateSnapshotParams {
  disk_id: number | string
  name: string
  description?: string
}

export function createCloudDiskSnapshot(params: CreateSnapshotParams) {
  return request<CloudDiskSnapshot>({
    url: '/storage/cloud-snapshots',
    method: 'post',
    params: { disk_id: params.disk_id },
    data: { name: params.name, description: params.description }
  })
}

export function deleteCloudDiskSnapshot(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/cloud-snapshots/${id}`,
    method: 'delete'
  })
}

export function batchDeleteCloudDiskSnapshots(ids: (number | string)[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/storage/cloud-snapshots/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export function rollbackCloudDiskSnapshot(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/cloud-snapshots/${id}/rollback`,
    method: 'post'
  })
}

export function createDiskFromSnapshot(snapshotId: number | string, data: { name: string; disk_type?: string }) {
  return request<CloudDisk>({
    url: `/storage/cloud-snapshots/${snapshotId}/create-disk`,
    method: 'post',
    data
  })
}

// ========== 主机快照 API ==========

export interface InstanceSnapshot {
  id: number | string
  cloud_account_id: number
  snapshot_id: string
  name: string
  instance_id: string
  instance_name: string
  disk_snapshots: number
  memory_snapshot: boolean
  cpu_arch: string
  size: number // GB
  status: string // available/creating/rollbacking/error
  provider_type: string
  account_name: string
  project_name: string
  region_id: string
  tags: { key: string; value: string }[]
  created_at: string
  description?: string
}

export interface InstanceSnapshotListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
}

export function getInstanceSnapshots(params?: InstanceSnapshotListParams) {
  return request<{ items: InstanceSnapshot[]; total: number; page: number; page_size: number }>({
    url: '/storage/instance-snapshots',
    method: 'get',
    params
  })
}

export function getInstanceSnapshot(id: number | string) {
  return request<InstanceSnapshot>({
    url: `/storage/instance-snapshots/${id}`,
    method: 'get'
  })
}

export interface CreateInstanceSnapshotParams {
  instance_id: string
  name: string
  memory_snapshot?: boolean
  description?: string
}

export function createInstanceSnapshot(params: CreateInstanceSnapshotParams) {
  return request<InstanceSnapshot>({
    url: '/storage/instance-snapshots',
    method: 'post',
    data: params
  })
}

export function deleteInstanceSnapshot(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/instance-snapshots/${id}`,
    method: 'delete'
  })
}

export function batchDeleteInstanceSnapshots(ids: (number | string)[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/storage/instance-snapshots/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export function rollbackInstanceSnapshot(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/instance-snapshots/${id}/rollback`,
    method: 'post'
  })
}

export function createVMFromSnapshot(snapshotId: number | string, data: { name: string; zone_id?: string }) {
  return request<{ message: string; vm_id: string }>({
    url: `/storage/instance-snapshots/${snapshotId}/create-vm`,
    method: 'post',
    data
  })
}

// ========== 快照策略 API ==========

export interface SnapshotPolicy {
  id: number | string
  cloud_account_id: number
  policy_id: string
  name: string
  status: string // active/inactive
  resource_type: string // disk/instance
  schedule_type: string // daily/weekly/monthly
  execute_time: string // HH:mm
  week_day: string // 1-7 (周一到周日)
  month_day: string // 1-28
  retention_days: number
  associated_count: number
  provider_type: string
  account_name: string
  region_id: string
  created_at: string
}

export interface SnapshotPolicyListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  resource_type?: string
}

export function getSnapshotPolicies(params?: SnapshotPolicyListParams) {
  return request<{ items: SnapshotPolicy[]; total: number; page: number; page_size: number }>({
    url: '/storage/snapshot-policies',
    method: 'get',
    params
  })
}

export function getSnapshotPolicy(id: number | string) {
  return request<SnapshotPolicy>({
    url: `/storage/snapshot-policies/${id}`,
    method: 'get'
  })
}

export interface CreateSnapshotPolicyParams {
  cloud_account_id: number
  name: string
  resource_type: string // disk/instance
  schedule_type: string // daily/weekly/monthly
  execute_time: string // HH:mm
  week_day?: string
  month_day?: string
  retention_days: number
  region_id?: string
}

export function createSnapshotPolicy(params: CreateSnapshotPolicyParams) {
  return request<SnapshotPolicy>({
    url: '/storage/snapshot-policies',
    method: 'post',
    data: params
  })
}

export function updateSnapshotPolicy(id: number | string, data: Partial<CreateSnapshotPolicyParams>) {
  return request<SnapshotPolicy>({
    url: `/storage/snapshot-policies/${id}`,
    method: 'put',
    data
  })
}

export function deleteSnapshotPolicy(id: number | string) {
  return request<{ message: string }>({
    url: `/storage/snapshot-policies/${id}`,
    method: 'delete'
  })
}

export function batchDeleteSnapshotPolicies(ids: (number | string)[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/storage/snapshot-policies/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export function toggleSnapshotPolicy(id: number | string, enabled: boolean) {
  return request<{ message: string }>({
    url: `/storage/snapshot-policies/${id}/toggle`,
    method: 'post',
    data: { enabled }
  })
}

export function associateResourcesToPolicy(policyId: number | string, resourceIds: string[]) {
  return request<{ message: string; associated_count: number }>({
    url: `/storage/snapshot-policies/${policyId}/associate`,
    method: 'post',
    data: { resource_ids: resourceIds }
  })
}

export function disassociateResourceFromPolicy(policyId: number | string, resourceId: string) {
  return request<{ message: string }>({
    url: `/storage/snapshot-policies/${policyId}/disassociate`,
    method: 'post',
    data: { resource_id: resourceId }
  })
}