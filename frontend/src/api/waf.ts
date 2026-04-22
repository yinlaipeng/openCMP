import request from '@/utils/request'

// WAF实例接口
export interface WAFInstance {
  id: number
  name: string
  status: string
  type: string
  platform: string
  cloud_account_id: number
  cloud_account?: {
    id: number
    name: string
  }
  domain_id: number
  region_id: string
  external_id: string
  tags: Record<string, string> | null
  description: string
  enabled: boolean
  sync_time: string | null
  created_at: string
  updated_at: string
}

// WAF列表查询参数
export interface WAFListParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  platform?: string
  cloud_account_id?: number | string
  domain_id?: number | string
}

// 创建WAF实例参数
export interface CreateWAFParams {
  name: string
  type?: string
  platform?: string
  cloud_account_id?: number
  domain_id?: number
  region_id?: string
  tags?: Record<string, string>
  description?: string
}

// 更新WAF实例参数
export interface UpdateWAFParams {
  name?: string
  description?: string
  tags?: Record<string, string>
  enabled?: boolean
}

// 批量删除参数
export interface BatchDeleteParams {
  ids: number[]
}

// 获取WAF实例列表
export function getWAFList(params: WAFListParams) {
  return request.get<{
    items: WAFInstance[]
    total: number
    page: number
    page_size: number
  }>('/waf', { params })
}

// 获取单个WAF实例详情
export function getWAFDetail(id: number) {
  return request.get<WAFInstance>(`/waf/${id}`)
}

// 创建WAF实例
export function createWAF(data: CreateWAFParams) {
  return request.post<WAFInstance>('/waf', data)
}

// 更新WAF实例
export function updateWAF(id: number, data: UpdateWAFParams) {
  return request.put<WAFInstance>(`/waf/${id}`, data)
}

// 删除WAF实例
export function deleteWAF(id: number) {
  return request.delete(`/waf/${id}`)
}

// 批量删除WAF实例
export function batchDeleteWAF(data: BatchDeleteParams) {
  return request.post<{ message: string; count: number }>('/waf/batch-delete', data)
}

// 同步WAF状态
export function syncWAFStatus(id: number) {
  return request.post<{ message: string }>(`/waf/${id}/sync`)
}