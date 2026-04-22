import request from '@/utils/request'

// 应用程序服务实例接口
export interface WebappInstance {
  id: number
  name: string
  status: string
  stack: string
  os_type: string
  ip_addr: string
  domain: string
  server_farm: string
  platform: string
  cloud_account_id: number
  cloud_account?: {
    id: number
    name: string
  }
  region_id: string
  project_id: number
  external_id: string
  tags: Record<string, string> | null
  description: string
  enabled: boolean
  sync_time: string | null
  created_at: string
  updated_at: string
}

// 应用程序服务列表查询参数
export interface WebappListParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  platform?: string
  cloud_account_id?: number | string
  project_id?: number | string
  stack?: string
}

// 创建应用程序服务参数
export interface CreateWebappParams {
  name: string
  stack?: string
  os_type?: string
  ip_addr?: string
  domain?: string
  server_farm?: string
  platform?: string
  cloud_account_id?: number
  region_id?: string
  project_id?: number
  tags?: Record<string, string>
  description?: string
}

// 更新应用程序服务参数
export interface UpdateWebappParams {
  name?: string
  stack?: string
  os_type?: string
  ip_addr?: string
  domain?: string
  server_farm?: string
  description?: string
  tags?: Record<string, string>
  enabled?: boolean
}

// 批量删除参数
export interface BatchDeleteParams {
  ids: number[]
}

// 获取应用程序服务列表
export function getWebappList(params: WebappListParams) {
  return request.get<{
    items: WebappInstance[]
    total: number
    page: number
    page_size: number
  }>('/webapp', { params })
}

// 获取单个应用程序服务详情
export function getWebappDetail(id: number) {
  return request.get<WebappInstance>(`/webapp/${id}`)
}

// 创建应用程序服务
export function createWebapp(data: CreateWebappParams) {
  return request.post<WebappInstance>('/webapp', data)
}

// 更新应用程序服务
export function updateWebapp(id: number, data: UpdateWebappParams) {
  return request.put<WebappInstance>(`/webapp/${id}`, data)
}

// 删除应用程序服务
export function deleteWebapp(id: number) {
  return request.delete(`/webapp/${id}`)
}

// 批量删除应用程序服务
export function batchDeleteWebapp(data: BatchDeleteParams) {
  return request.post<{ message: string; count: number }>('/webapp/batch-delete', data)
}

// 同步应用程序服务状态
export function syncWebappStatus(id: number) {
  return request.post<{ message: string }>(`/webapp/${id}/sync`)
}