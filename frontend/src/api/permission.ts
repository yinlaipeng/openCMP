import request from '@/utils/request'

// 获取权限列表
export function getPermissions(params?: {
  domain_id?: number
  limit?: number
  offset?: number
  keyword?: string
  type?: string
  enabled?: boolean
  scope?: string
}) {
  return request({
    url: '/permissions',
    method: 'get',
    params
  })
}

// 获取权限详情
export function getPermission(id: number) {
  return request({
    url: `/permissions/${id}`,
    method: 'get'
  })
}

// 创建权限
export function createPermission(data: {
  name: string
  display_name?: string
  description?: string
  type?: string
  scope?: string
  domain_id?: number
  project_id?: number
  resource: string
  action: string
  conditions?: Record<string, any>
  enabled?: boolean
  is_public?: boolean
}) {
  return request({
    url: '/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(id: number, data: {
  name?: string
  display_name?: string
  description?: string
  type?: string
  scope?: string
  domain_id?: number
  project_id?: number
  resource?: string
  action?: string
  conditions?: Record<string, any>
  enabled?: boolean
  is_public?: boolean
}) {
  return request({
    url: `/permissions/${id}`,
    method: 'put',
    data
  })
}

// 删除权限
export function deletePermission(id: number) {
  return request({
    url: `/permissions/${id}`,
    method: 'delete'
  })
}

// 启用权限
export function enablePermission(id: number) {
  return request({
    url: `/permissions/${id}/enable`,
    method: 'post'
  })
}

// 禁用权限
export function disablePermission(id: number) {
  return request({
    url: `/permissions/${id}/disable`,
    method: 'post'
  })
}

// 克隆权限
export function clonePermission(id: number, data?: { name?: string; display_name?: string }) {
  return request({
    url: `/permissions/${id}/clone`,
    method: 'post',
    data
  })
}

// 设置权限公开
export function makePermissionPublic(id: number) {
  return request({
    url: `/permissions/${id}/public`,
    method: 'post'
  })
}