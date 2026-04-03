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
export function getRoles(params?: {
  domain_id?: number
  limit?: number
  offset?: number
  keyword?: string
  type?: string
  enabled?: boolean
}) {
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

// 启用角色
export function enableRole(id: number) {
  return request({
    url: `/roles/${id}/enable`,
    method: 'post'
  })
}

// 禁用角色
export function disableRole(id: number) {
  return request({
    url: `/roles/${id}/disable`,
    method: 'post'
  })
}

// 公开角色
export function makeRolePublic(id: number) {
  return request({
    url: `/roles/${id}/public`,
    method: 'post'
  })
}

// 获取角色的用户列表
export function getRoleUsers(roleId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/roles/${roleId}/users`,
    method: 'get',
    params
  })
}

// 获取角色的用户组列表
export function getRoleGroups(roleId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/roles/${roleId}/groups`,
    method: 'get',
    params
  })
}

// ============= 策略管理 API =============

// 获取策略列表
export function getPolicies(params?: {
  scope?: string
  domain_id?: string
  limit?: number
  offset?: number
}) {
  return request({
    url: '/policies',
    method: 'get',
    params
  })
}

// 获取策略详情
export function getPolicy(id: string) {
  return request({
    url: `/policies/${id}`,
    method: 'get'
  })
}

// 创建策略
export function createPolicy(data: {
  name: string
  scope: string
  description?: string
  domain_id?: string
  policy: Record<string, any>
  is_system?: boolean
  is_public?: boolean
}) {
  return request({
    url: '/policies',
    method: 'post',
    data
  })
}

// 更新策略
export function updatePolicy(id: string, data: Record<string, any>) {
  return request({
    url: `/policies/${id}`,
    method: 'put',
    data
  })
}

// 删除策略
export function deletePolicy(id: string) {
  return request({
    url: `/policies/${id}`,
    method: 'delete'
  })
}

// 获取角色的策略列表
export function getRolePolicies(roleId: number) {
  return request({
    url: `/roles/${roleId}/policies`,
    method: 'get'
  })
}

// 分配策略给角色
export function assignPolicyToRole(roleId: number, policyId: string) {
  return request({
    url: `/roles/${roleId}/policies`,
    method: 'post',
    data: { policy_id: policyId }
  })
}

// 从角色撤销策略
export function revokePolicyFromRole(roleId: number, policyId: string) {
  return request({
    url: `/roles/${roleId}/policies?policy_id=${policyId}`,
    method: 'delete'
  })
}

// ============= 认证源 API =============

// 获取认证源列表
export function getAuthSources(params?: {
  keyword?: string
  type?: string
  enabled?: boolean
  limit?: number
  offset?: number
}) {
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
  auto_create?: boolean
}) {
  return request({
    url: '/auth-sources',
    method: 'post',
    data
  })
}

// 更新认证源
export function updateAuthSource(id: number, data: {
  name?: string
  description?: string
  config?: Record<string, any>
  enabled?: boolean
  auto_create?: boolean
}) {
  return request({
    url: `/auth-sources/${id}`,
    method: 'put',
    data
  })
}

// 删除认证源
export function deleteAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}`,
    method: 'delete'
  })
}

// 测试认证源
export function testAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}/test`,
    method: 'post'
  })
}

// 启用认证源
export function enableAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}/enable`,
    method: 'post'
  })
}

// 禁用认证源
export function disableAuthSource(id: number) {
  return request({
    url: `/auth-sources/${id}/disable`,
    method: 'post'
  })
}

// ============= 域管理 API =============

// 获取域列表
export function getDomains(params?: {
  keyword?: string
  enabled?: boolean
  limit?: number
  offset?: number
}) {
  return request({
    url: '/domains',
    method: 'get',
    params
  })
}

// 获取域详情
export function getDomain(id: number) {
  return request({
    url: `/domains/${id}`,
    method: 'get'
  })
}

// 创建域
export function createDomain(data: {
  name: string
  description?: string
  enabled?: boolean
}) {
  return request({
    url: '/domains',
    method: 'post',
    data
  })
}

// 更新域
export function updateDomain(id: number, data: {
  name?: string
  description?: string
  enabled?: boolean
}) {
  return request({
    url: `/domains/${id}`,
    method: 'put',
    data
  })
}

// 删除域
export function deleteDomain(id: number) {
  return request({
    url: `/domains/${id}`,
    method: 'delete'
  })
}

// 启用域
export function enableDomain(id: number) {
  return request({
    url: `/domains/${id}/enable`,
    method: 'post'
  })
}

// 禁用域
export function disableDomain(id: number) {
  return request({
    url: `/domains/${id}/disable`,
    method: 'post'
  })
}

// 获取域的用户列表
export function getDomainUsers(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/users`,
    method: 'get',
    params
  })
}

// 获取域的用户组列表
export function getDomainGroups(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/groups`,
    method: 'get',
    params
  })
}

// 获取域的项目列表
export function getDomainProjects(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/projects`,
    method: 'get',
    params
  })
}

// 获取域的角色列表
export function getDomainRoles(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/roles`,
    method: 'get',
    params
  })
}

// ============= 项目管理 API =============

// 获取项目列表
export function getProjects(params?: {
  domain_id?: number
  keyword?: string
  enabled?: boolean
  limit?: number
  offset?: number
}) {
  return request({
    url: '/projects',
    method: 'get',
    params
  })
}

// 获取项目详情
export function getProject(id: number) {
  return request({
    url: `/projects/${id}`,
    method: 'get'
  })
}

// 创建项目
export function createProject(data: {
  name: string
  description?: string
  domain_id: number
  parent_id?: number
}) {
  return request({
    url: '/projects',
    method: 'post',
    data
  })
}

// 更新项目
export function updateProject(id: number, data: {
  name?: string
  description?: string
  domain_id?: number
  parent_id?: number
}) {
  return request({
    url: `/projects/${id}`,
    method: 'put',
    data
  })
}

// 删除项目
export function deleteProject(id: number) {
  return request({
    url: `/projects/${id}`,
    method: 'delete'
  })
}

// 启用项目
export function enableProject(id: number) {
  return request({
    url: `/projects/${id}/enable`,
    method: 'post'
  })
}

// 禁用项目
export function disableProject(id: number) {
  return request({
    url: `/projects/${id}/disable`,
    method: 'post'
  })
}

// 加入项目（分配用户和角色）
export function joinProject(projectId: number, users: number[], roles: number[]) {
  return request({
    url: `/projects/${projectId}/join`,
    method: 'post',
    data: { users, roles }
  })
}

// 获取项目的用户列表
export function getProjectUsers(projectId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/projects/${projectId}/users`,
    method: 'get',
    params
  })
}

// 获取项目的角色列表
export function getProjectRoles(projectId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/projects/${projectId}/roles`,
    method: 'get',
    params
  })
}

// 从项目移除用户
export function removeUserFromProject(projectId: number, userId: number, roleId?: number) {
  const params: any = { user_id: userId }
  if (roleId) params.role_id = roleId
  return request({
    url: `/projects/${projectId}/users`,
    method: 'delete',
    params
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

// 更新用户
export function updateUser(id: number, data: {
  display_name?: string
  email?: string
  phone?: string
}) {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

// 重置用户密码
export function resetUserPassword(id: number, password: string) {
  return request({
    url: `/users/${id}/reset-password`,
    method: 'post',
    data: { password }
  })
}

// 获取用户角色
export function getUserRoles(userId: number, params?: { domain_id?: number }) {
  return request({
    url: `/users/${userId}/roles`,
    method: 'get',
    params
  })
}

// 分配角色给用户
export function assignRoleToUser(userId: number, roleId: number, domainId: number) {
  return request({
    url: `/users/${userId}/roles`,
    method: 'post',
    data: { role_id: roleId, domain_id: domainId }
  })
}

// 撤销用户角色
export function revokeRoleFromUser(userId: number, roleId: number, domainId: number) {
  return request({
    url: `/users/${userId}/roles?role_id=${roleId}&domain_id=${domainId}`,
    method: 'delete'
  })
}

// 获取用户组
export function getUserGroups(userId: number) {
  return request({
    url: `/users/${userId}/groups`,
    method: 'get'
  })
}

// 加入用户组
export function joinGroup(userId: number, groupId: number) {
  return request({
    url: `/users/${userId}/groups`,
    method: 'post',
    data: { group_id: groupId }
  })
}

// 离开用户组
export function leaveGroup(userId: number, groupId: number) {
  return request({
    url: `/users/${userId}/groups?group_id=${groupId}`,
    method: 'delete'
  })
}

// 获取用户组列表
export function getGroups(params?: { limit?: number; offset?: number }) {
  return request({
    url: '/groups',
    method: 'get',
    params
  })
}

// 获取用户组详情
export function getGroup(id: number) {
  return request({
    url: `/groups/${id}`,
    method: 'get'
  })
}

// 创建用户组
export function createGroup(data: {
  name: string
  description?: string
  domain_id: number
}) {
  return request({
    url: '/groups',
    method: 'post',
    data
  })
}

// 更新用户组
export function updateGroup(id: number, data: {
  description?: string
}) {
  return request({
    url: `/groups/${id}`,
    method: 'put',
    data
  })
}

// 删除用户组
export function deleteGroup(id: number) {
  return request({
    url: `/groups/${id}`,
    method: 'delete'
  })
}

// 获取用户组的用户列表
export function getGroupUsers(groupId: number) {
  return request({
    url: `/groups/${groupId}/users`,
    method: 'get'
  })
}

// 添加用户到用户组
export function addUserToGroup(groupId: number, userId: number) {
  return request({
    url: `/groups/${groupId}/users`,
    method: 'post',
    data: { user_id: userId }
  })
}

// 从用户组移除用户
export function removeUserFromGroup(groupId: number, userId: number) {
  return request({
    url: `/groups/${groupId}/users?user_id=${userId}`,
    method: 'delete'
  })
}