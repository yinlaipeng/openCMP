import request from '@/utils/request'

// 同步日志 API

export interface SyncLog {
  id: number
  cloud_account_id: number
  cloud_account_name: string
  sync_type: string
  sync_mode: string
  resource_type: string
  sync_start_time: string
  sync_end_time: string | null
  sync_duration: number
  status: string
  total_count: number
  new_count: number
  updated_count: number
  deleted_count: number
  skipped_count: number
  error_count: number
  error_message: string
  triggered_by: string
  scheduled_task_id: number | null
}

export interface SyncStatistics {
  cloud_account_id: number
  period_days: number
  total_sync_count: number
  success_count: number
  failure_count: number
  avg_duration: number
  total_new: number
  total_updated: number
  total_deleted: number
}

// 获取同步日志列表
export function getSyncLogs(cloudAccountId: number, limit?: number) {
  return request<{ items: SyncLog[], total: number }>({
    url: '/sync-logs',
    method: 'get',
    params: {
      cloud_account_id: cloudAccountId,
      limit: limit || 20
    }
  })
}

// 获取同步日志详情
export function getSyncLog(id: number) {
  return request<SyncLog>({
    url: `/sync-logs/${id}`,
    method: 'get'
  })
}

// 获取同步统计信息
export function getSyncStatistics(cloudAccountId: number, days?: number) {
  return request<SyncStatistics>({
    url: '/sync-logs/statistics',
    method: 'get',
    params: {
      cloud_account_id: cloudAccountId,
      days: days || 30
    }
  })
}

// 获取最新同步日志
export function getLatestSyncLog(cloudAccountId: number) {
  return request<SyncLog>({
    url: '/sync-logs/latest',
    method: 'get',
    params: {
      cloud_account_id: cloudAccountId
    }
  })
}