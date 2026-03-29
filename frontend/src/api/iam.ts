import request from '@/utils/request'

// ============= 权限管理 API =============

// 获取权限列表
export function getPermissions(params?: {
  resource?: string
  action?: string
  type?: string
  keyword?: string
  limit?: number
  offset?: number
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
  display_name: string
  resource: string
  action: string
  type?: string
  description?: string
}) {
  return request({
    url: '/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(id: number, data: {
  display_name?: string
  description?: string
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

// 获取资源类型列表
export function getResources() {
  return request({
    url: '/permissions/resources',
    method: 'get'
  })
}

// 获取操作类型列表
export function getActions() {
  return request({
    url: '/permissions/actions',
    method: 'get'
  })
}

// ============= 角色管理 API =============

// 获取角色列表
export function getRoles(params?: { domain_id?: number; limit?: number; offset?: number }) {
  return request({
    url: '/roles',
    method: 'get',
    params
  })
}

// 获取角色详情
export function getRole(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'get'
  })
}

// 创建角色
export function createRole(data: {
  name: string
  display_name?: string
  description?: string
  domain_id?: number
  type?: string
}) {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id: number, data: {
  name: string
  display_name?: string
  description?: string
  domain_id?: number
  type?: string
}) {
  return request({
    url: `/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'delete'
  })
}

// 获取角色权限
export function getRolePermissions(roleId: number) {
  return request({
    url: `/roles/${roleId}/permissions`,
    method: 'get'
  })
}

// 分配权限给角色
export function assignPermission(roleId: number, permissionId: number) {
  return request({
    url: `/roles/${roleId}/permissions`,
    method: 'post',
    data: { permission_id: permissionId }
  })
}

// 从角色撤销权限
export function revokePermission(roleId: number, permissionId: number) {
  return request({
    url: `/roles/${roleId}/permissions?permission_id=${permissionId}`,
    method: 'delete'
  })
}

// ============= 认证源 API =============

// 获取认证源列表
export function getAuthSources(params?: { limit?: number; offset?: number }) {
  return request({
    url: '/auth-sources',
    method: 'get',
    params
  })
}

// 获取认证源详情
export function getAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}`,
    method: 'get'
  })
}

// 创建认证源
export function createAuthSource(data: {
  name: string
  type: string
  description?: string
  config?: Record<string, any>
  enabled?: boolean
}) {
  return request({
    url: '/auth-sources',
    method: 'post',
    data
  })
}

// 测试认证源
export function testAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}/test`,
    method: 'post'
  })
}

// ============= 用户 API =============

// 获取用户列表
export function getUsers(params?: { domain_id?: number; limit?: number; offset?: number }) {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUser(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser(data: {
  name: string
  display_name?: string
  email?: string
  phone?: string
  password: string
  domain_id: number
}) {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

// 删除用户
export function deleteUser(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

// 启用用户
export function enableUser(id: number) {
  return request({
    url: `/users/${id}/enable`,
    method: 'post'
  })
}

// 禁用用户
export function disableUser(id: number) {
  return request({
    url: `/users/${id}/disable`,
    method: 'post'
  })
}
