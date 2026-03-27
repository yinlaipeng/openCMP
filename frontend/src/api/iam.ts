import request from '@/utils/request'

// 认证源 API
export function getAuthSources(params?: { limit?: number; offset?: number }) {
  return request({
    url: '/auth-sources',
    method: 'get',
    params
  })
}

export function getAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}`,
    method: 'get'
  })
}

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

export function testAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}/test`,
    method: 'post'
  })
}

// 用户 API
export function getUsers(params?: { domain_id?: number; limit?: number; offset?: number }) {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

export function getUser(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

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

export function deleteUser(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

export function enableUser(id: number) {
  return request({
    url: `/users/${id}/enable`,
    method: 'post'
  })
}

export function disableUser(id: number) {
  return request({
    url: `/users/${id}/disable`,
    method: 'post'
  })
}

// 角色 API
export function getRoles(params?: { domain_id?: number; limit?: number; offset?: number }) {
  return request({
    url: '/roles',
    method: 'get',
    params
  })
}

export function getRole(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'get'
  })
}

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

export function deleteRole(id: number) {
  return request({
    url: `/roles/${id}`,
    method: 'delete'
  })
}

export function getPermissions(params?: { limit?: number; offset?: number }) {
  return request({
    url: '/roles/permissions',
    method: 'get',
    params
  })
}

export function assignPermission(roleId: number, permissionId: number) {
  return request({
    url: `/roles/${roleId}/permissions`,
    method: 'post',
    data: { permission_id: permissionId }
  })
}

export function getRolePermissions(roleId: number) {
  return request({
    url: `/roles/${roleId}/permissions`,
    method: 'get'
  })
}
