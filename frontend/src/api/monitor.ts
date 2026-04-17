import request from '@/utils/request'

// 告警策略 Types
export interface AlertPolicy {
  id: number
  name: string
  status: string
  enabled: boolean
  resource_type: string
  metric: string
  threshold: number
  duration: number
  level: string
  owner: string
  domain_id: number
  project_id: number
  description: string
  notify_channel: string
  created_at: string
  updated_at: string
}

export interface AlertPolicyRequest {
  name: string
  resource_type: string
  metric: string
  threshold: number
  duration?: number
  level?: string
  enabled?: boolean
  owner?: string
  domain_id?: number
  project_id?: number
  description?: string
  notify_channel?: string
}

export interface AlertPolicyFilter {
  resource_type?: string
  owner?: string
  level?: string
  enabled?: string
}

// 告警历史 Types
export interface AlertHistory {
  id: number
  policy_id: number
  policy_name: string
  resource_id: string
  resource_name: string
  resource_type: string
  level: string
  status: string
  metric_value: number
  message: string
  triggered_at: string
  resolved_at?: string
  domain_id: number
  project_id: number
}

export interface AlertHistoryFilter {
  status?: string
  level?: string
  resource_type?: string
  start_time?: string
  end_time?: string
  page?: number
  page_size?: number
}

// 监控资源 Types
export interface MonitorResource {
  id: number
  resource_id: string
  resource_name: string
  resource_type: string
  monitor_status: string
  account_id: number
  platform: string
  region: string
  project_id: number
  last_sync_at: string
  metrics: string
}

export interface ResourceMetrics {
  resource_id: string
  timestamp: number
  cpu_usage: { value: number; unit: string; timestamp: number }
  memory_usage: { value: number; unit: string; timestamp: number }
  disk_usage: { value: number; unit: string; timestamp: number }
  network_in: { value: number; unit: string; timestamp: number }
  network_out: { value: number; unit: string; timestamp: number }
}

// 告警策略 API
export function listAlertPolicies(filter?: AlertPolicyFilter) {
  return request.get<AlertPolicy[]>('/monitor/alert-policies', { params: filter })
}

export function createAlertPolicy(data: AlertPolicyRequest) {
  return request.post<AlertPolicy>('/monitor/alert-policies', data)
}

export function getAlertPolicy(id: number) {
  return request.get<AlertPolicy>(`/monitor/alert-policies/${id}`)
}

export function updateAlertPolicy(id: number, data: AlertPolicyRequest) {
  return request.put<AlertPolicy>(`/monitor/alert-policies/${id}`, data)
}

export function deleteAlertPolicy(id: number) {
  return request.delete(`/monitor/alert-policies/${id}`)
}

export function toggleAlertPolicy(id: number) {
  return request.post<AlertPolicy>(`/monitor/alert-policies/${id}/toggle`)
}

// 告警历史 API
export function listAlertHistory(filter?: AlertHistoryFilter) {
  return request.get<{ data: AlertHistory[]; total: number; page: number; page_size: number }>('/monitor/alert-history', { params: filter })
}

export function resolveAlert(id: number) {
  return request.post<AlertHistory>(`/monitor/alert-history/${id}/resolve`)
}

export function ignoreAlert(id: number) {
  return request.post<AlertHistory>(`/monitor/alert-history/${id}/ignore`)
}

// 监控资源 API
export function listMonitorResources(accountId?: number, resourceType?: string, monitorStatus?: string) {
  return request.get<MonitorResource[]>('/monitor/resources', {
    params: {
      account_id: accountId,
      resource_type: resourceType,
      monitor_status: monitorStatus
    }
  })
}

export function syncMonitorResources(accountId: number) {
  return request.post<{ message: string; count: number }>('/monitor/resources/sync', null, {
    params: { account_id: accountId }
  })
}

export function getResourceMetrics(accountId: number, resourceId: string) {
  return request.get<ResourceMetrics>(`/monitor/resources/${resourceId}/metrics`, {
    params: { account_id: accountId }
  })
}