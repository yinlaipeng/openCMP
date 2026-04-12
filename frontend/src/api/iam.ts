import request from '@/utils/request'

// Import permission APIs
import {
  getPermissions,
  createPermission,
  updatePermission,
  deletePermission,
  enablePermission,
  disablePermission,
  clonePermission,
  makePermissionPublic,
  getPermission
} from '@/api/permission'


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
  return request<{ items: any[], total: number }>({
    url: `/roles/${roleId}/users`,
    method: 'get',
    params
  })
}

// 获取角色的组列表
export function getRoleGroups(roleId: number, params?: { limit?: number; offset?: number }) {
  return request<{ items: any[], total: number }>({
    url: `/roles/${roleId}/groups`,
    method: 'get',
    params
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

// 测试 LDAP 用户查询
export function testLdapUsers(config: any) {
  return request({
    url: '/auth-sources/test-ldap-users',
    method: 'post',
    data: config
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

// 获取域的云账号列表
export function getDomainCloudAccounts(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/cloud-accounts`,
    method: 'get',
    params
  })
}

// 获取域的操作日志列表
export function getDomainOperationLogs(domainId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/domains/${domainId}/operation-logs`,
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

// 设置项目管理员
export function setProjectManager(projectId: number, userId: number) {
  return request({
    url: `/projects/${projectId}/manager`,
    method: 'put',
    data: { user_id: userId }
  })
}

// ============= 用户 API =============

// 获取用户列表
export function getUsers(params?: { domain_id?: number; keyword?: string; email?: string; enabled?: boolean; limit?: number; offset?: number }) {
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
  remark?: string
  email?: string
  phone?: string
  password: string
  domain_id: number
  console_login?: boolean
  mfa_enabled?: boolean
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
  remark?: string
  email?: string
  phone?: string
  console_login?: boolean
  mfa_enabled?: boolean
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

// 将用户分配到项目
export function assignUserToProject(userId: number, projectId: number, roleId: number) {
  return request({
    url: `/users/${userId}/projects`,
    method: 'post',
    data: { project_id: projectId, role_id: roleId }
  })
}

// 获取用户所属的项目
export function getUserProjects(userId: number) {
  return request({
    url: `/users/${userId}/projects`,
    method: 'get'
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

// 添加项目到用户组
export function addGroupToProject(groupId: number, projectId: number) {
  return request({
    url: `/groups/${groupId}/projects`,
    method: 'post',
    data: { project_id: projectId }
  })
}

// 从项目移除用户组
export function removeGroupFromProject(groupId: number, projectId: number) {
  return request({
    url: `/groups/${groupId}/projects?project_id=${projectId}`,
    method: 'delete'
  })
}

// ============= 权限 API =============

export {
  getPermissions,
  getPermission,
  createPermission,
  updatePermission,
  deletePermission,
  enablePermission,
  disablePermission,
  clonePermission,
  makePermissionPublic
}

// 获取用户组的项目列表
export function getGroupProjects(groupId: number, params?: { limit?: number; offset?: number }) {
  return request({
    url: `/groups/${groupId}/projects`,
    method: 'get',
    params
  })
}
// Get project-specific alerts
export function getProjectSecurityAlerts(projectId: number, params?: any) {
  return request({
    url: '/alerts',
    method: 'get',
    params: { ...params, project_id: projectId }
  })
}

// Get project-specific messages
export function getProjectMessages(projectId: number, params?: any) {
  return request({
    url: '/messages',
    method: 'get',
    params: { ...params, project_id: projectId }
  })
}

// Get project-specific robots
export function getProjectRobots(projectId: number, params?: any) {
  return request({
    url: '/robots',
    method: 'get',
    params: { ...params, project_id: projectId }
  })
}

// Robot management functions
export function createRobot(data: any) {
  return request({
    url: '/robots',
    method: 'post',
    data
  })
}

export function updateRobot(id: number, data: any) {
  return request({
    url: `/robots/${id}`,
    method: 'put',
    data
  })
}

export function deleteRobot(id: number) {
  return request({
    url: `/robots/${id}`,
    method: 'delete'
  })
}

export function toggleRobotStatus(id: number, enabled: boolean) {
  return request({
    url: `/robots/${id}/status`,
    method: 'post',
    data: { enabled }
  })
}

// Get project-specific receivers
export function getProjectReceivers(projectId: number, params?: any) {
  return request({
    url: '/receivers',
    method: 'get',
    params: { ...params, project_id: projectId }
  })
}

// Get project-specific subscriptions
export function getProjectSubscriptions(projectId: number, params?: any) {
  return request({
    url: '/subscriptions',
    method: 'get',
    params: { ...params, project_id: projectId }
  })
}

// Get operation logs
export function getOperationLogs(params?: any) {
  return request({
    url: '/operation-logs',
    method: 'get',
    params
  })
}

// Get operation logs for specific resource
export function getResourceOperationLogs(resourceType: string, resourceId: number, params?: any) {
  return request({
    url: `/operation-logs/${resourceType}/${resourceId}`,
    method: 'get',
    params
  })
}
