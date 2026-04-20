import request from '@/utils/request'
import type { CloudAccount, CreateCloudAccountRequest } from '@/types'

// 云账号搜索参数
export interface CloudAccountSearchParams {
  page?: number
  page_size?: number
  id?: string           // 支持多ID用|分隔
  name?: string         // 模糊搜索，自动匹配IP或ID
  remarks?: string      // 备注模糊搜索
  provider_type?: string // 平台精确匹配
  status?: string       // 连接状态精确匹配
  enabled?: boolean     // 启用状态精确匹配
  health_status?: string // 健康状态精确匹配
  account_number?: string // 账号模糊搜索
  domain_id?: number    // 域ID精确匹配
}

export function getCloudAccounts(params?: CloudAccountSearchParams) {
  return request<{ items: CloudAccount[], total: number, page: number, page_size: number }>({
    url: '/cloud-accounts',
    method: 'get',
    params
  })
}

export function getCloudAccount(id: number) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}`,
    method: 'get'
  })
}

export function createCloudAccount(data: CreateCloudAccountRequest) {
  return request<CloudAccount>({
    url: '/cloud-accounts',
    method: 'post',
    data
  })
}

export function updateCloudAccount(id: number, data: Partial<CreateCloudAccountRequest>) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}`,
    method: 'put',
    data
  })
}

export function deleteCloudAccount(id: number) {
  return request({
    url: `/cloud-accounts/${id}`,
    method: 'delete'
  })
}

export function verifyCloudAccount(id: number) {
  return request({
    url: `/cloud-accounts/${id}/verify`,
    method: 'post'
  })
}

// 同步云账号
export interface SyncCloudAccountOptions {
  mode: 'full' | 'incremental'
  resource_types: string[]
}

export function syncCloudAccount(id: number, options: SyncCloudAccountOptions) {
  return request({
    url: `/cloud-accounts/${id}/sync`,
    method: 'post',
    data: options
  })
}

// 测试连接
export function testConnection(id: number) {
  return request<{ connected: boolean; message: string }>({
    url: `/cloud-accounts/${id}/test-connection`,
    method: 'post'
  })
}

// 验证凭证（实际调用云厂商 API）
export function verifyCredentials(id: number) {
  return request<{ valid: boolean; message: string }>({
    url: `/cloud-accounts/${id}/verify-credentials`,
    method: 'get'
  })
}

// 更新云账号状态
export function updateCloudAccountStatus(id: number, enabled: boolean) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}/status`,
    method: 'patch',
    data: { enabled }
  })
}

// 更新云账号属性
export interface UpdateCloudAccountAttributes {
  auto_sync?: boolean
  sync_policy_id?: number | null
  sync_interval?: number
  sync_resource_types?: string[]
}

export function updateCloudAccountAttributes(id: number, attrs: UpdateCloudAccountAttributes) {
  return request<CloudAccount>({
    url: `/cloud-accounts/${id}/attributes`,
    method: 'patch',
    data: attrs
  })
}

// 使用新凭证测试连接
export interface TestConnectionWithCredentialsRequest {
  access_key_id: string
  access_key_secret: string
}

export interface TestConnectionWithCredentialsResponse {
  connected: boolean
  message: string
  regions: string[]
}

export function testConnectionWithCredentials(id: number, credentials: TestConnectionWithCredentialsRequest) {
  return request<TestConnectionWithCredentialsResponse>({
    url: `/cloud-accounts/${id}/test-connection-with-credentials`,
    method: 'post',
    data: credentials
  })
}

// 获取资源统计
export interface ResourceStatsResponse {
  resources: {
    vms: number
    rds: number
    redis: number
    buckets: number
    eips: number
    public_ips: number
    snapshots: number
    vpcs: number
    subnets: number
    total_ips: number
    vms_running: number
    eips_bound: number
    ips_used: number
    disks: number
    disks_mounted: number
  }
  usage_rates: {
    vm_running_rate: number
    disk_mounted_rate: number
    eip_bound_rate: number
    ip_used_rate: number
  }
}

export function getResourceStats(id: number) {
  return request<ResourceStatsResponse>({
    url: `/cloud-accounts/${id}/resource-stats`,
    method: 'get'
  })
}

// 获取支持的资源类型列表
export interface ResourceType {
  id: string
  name: string
}

export function getSupportedResourceTypes() {
  return request<{ items: ResourceType[], total: number }>({
    url: '/cloud-accounts/resource-types',
    method: 'get'
  })
}

// 获取权限列表
export function getPermissions(id: number) {
  return request<{ permissions: Array<{ name: string; description: string; granted: boolean }> }>({
    url: `/cloud-accounts/${id}/permissions`,
    method: 'get'
  })
}

// 订阅 API
export interface CloudSubscription {
  id: number
  cloud_account_id: number
  name: string
  subscription_id: string
  enabled: boolean
  status: string
  sync_time: string | null
  sync_duration: number
  sync_status: string
  domain_id: number
  default_project_id: number | null
}

export function getSubscriptions(accountId: number) {
  return request<{ items: CloudSubscription[]; total: number }>({
    url: `/cloud-accounts/${accountId}/subscriptions`,
    method: 'get'
  })
}

export function createSubscription(accountId: number, data: Partial<CloudSubscription>) {
  return request<CloudSubscription>({
    url: `/cloud-accounts/${accountId}/subscriptions`,
    method: 'post',
    data
  })
}

export function updateSubscription(accountId: number, subscriptionId: number, data: Partial<CloudSubscription>) {
  return request<CloudSubscription>({
    url: `/cloud-accounts/${accountId}/subscriptions/${subscriptionId}`,
    method: 'put',
    data
  })
}

export function deleteSubscription(accountId: number, subscriptionId: number) {
  return request<{ message: string }>({
    url: `/cloud-accounts/${accountId}/subscriptions/${subscriptionId}`,
    method: 'delete'
  })
}

export function toggleSubscription(accountId: number, subscriptionId: number, enabled: boolean) {
  return request<{ message: string; enabled: boolean }>({
    url: `/cloud-accounts/${accountId}/subscriptions/${subscriptionId}/toggle`,
    method: 'post',
    data: { enabled }
  })
}

export function syncSubscription(accountId: number, subscriptionId: number) {
  return request<{ message: string; subscription: CloudSubscription }>({
    url: `/cloud-accounts/${accountId}/subscriptions/${subscriptionId}/sync`,
    method: 'post'
  })
}

export function updateSubscriptionProject(accountId: number, subscriptionId: number, projectId: number) {
  return request<{ message: string; project_id: number }>({
    url: `/cloud-accounts/${accountId}/subscriptions/${subscriptionId}/project`,
    method: 'put',
    data: { project_id: projectId }
  })
}

// 云用户 API
export interface CloudUser {
  id: number
  cloud_account_id: number
  username: string
  console_login: boolean
  status: string
  password: string
  login_url: string
  local_user_id: number | null
  platform: string
}

export function getCloudUsers(accountId: number) {
  return request<{ items: CloudUser[]; total: number }>({
    url: `/cloud-accounts/${accountId}/cloud-users`,
    method: 'get'
  })
}

export function createCloudUser(accountId: number, data: Partial<CloudUser>) {
  return request<CloudUser>({
    url: `/cloud-accounts/${accountId}/cloud-users`,
    method: 'post',
    data
  })
}

export function updateCloudUser(accountId: number, userId: number, data: Partial<CloudUser>) {
  return request<CloudUser>({
    url: `/cloud-accounts/${accountId}/cloud-users/${userId}`,
    method: 'put',
    data
  })
}

export function deleteCloudUser(accountId: number, userId: number) {
  return request<{ message: string }>({
    url: `/cloud-accounts/${accountId}/cloud-users/${userId}`,
    method: 'delete'
  })
}

// 云用户组 API
export interface CloudUserGroup {
  id: number
  cloud_account_id: number
  name: string
  status: string
  permissions: string
  platform: string
  domain_id: number
}

export function getCloudUserGroups(accountId: number) {
  return request<{ items: CloudUserGroup[]; total: number }>({
    url: `/cloud-accounts/${accountId}/cloud-user-groups`,
    method: 'get'
  })
}

export function createCloudUserGroup(accountId: number, data: Partial<CloudUserGroup>) {
  return request<CloudUserGroup>({
    url: `/cloud-accounts/${accountId}/cloud-user-groups`,
    method: 'post',
    data
  })
}

export function updateCloudUserGroup(accountId: number, groupId: number, data: Partial<CloudUserGroup>) {
  return request<CloudUserGroup>({
    url: `/cloud-accounts/${accountId}/cloud-user-groups/${groupId}`,
    method: 'put',
    data
  })
}

export function deleteCloudUserGroup(accountId: number, groupId: number) {
  return request<{ message: string }>({
    url: `/cloud-accounts/${accountId}/cloud-user-groups/${groupId}`,
    method: 'delete'
  })
}

// 云上项目 API
export interface CloudProject {
  id: number
  cloud_account_id: number
  name: string
  subscription_id: number | null
  status: string
  tags: string
  domain_id: number
  local_project_id: number | null
  priority: number
  sync_time: string | null
}

export function getCloudProjects(accountId: number) {
  return request<{ items: CloudProject[]; total: number }>({
    url: `/cloud-accounts/${accountId}/cloud-projects`,
    method: 'get'
  })
}

export function createCloudProject(accountId: number, data: Partial<CloudProject>) {
  return request<CloudProject>({
    url: `/cloud-accounts/${accountId}/cloud-projects`,
    method: 'post',
    data
  })
}

export function updateCloudProject(accountId: number, projectId: number, data: Partial<CloudProject>) {
  return request<CloudProject>({
    url: `/cloud-accounts/${accountId}/cloud-projects/${projectId}`,
    method: 'put',
    data
  })
}

export function deleteCloudProject(accountId: number, projectId: number) {
  return request<{ message: string }>({
    url: `/cloud-accounts/${accountId}/cloud-projects/${projectId}`,
    method: 'delete'
  })
}

export function mapCloudProjectToLocal(accountId: number, projectId: number, localProjectId: number) {
  return request<{ message: string; local_project_id: number }>({
    url: `/cloud-accounts/${accountId}/cloud-projects/${projectId}/map`,
    method: 'post',
    data: { local_project_id: localProjectId }
  })
}

// 操作日志 API（云账户相关）
export function getOperationLogsByAccount(accountId: number, params?: { page?: number; page_size?: number }) {
  return request<{ items: any[]; total: number; page: number; page_size: number }>({
    url: `/cloud-accounts/${accountId}/operation-logs`,
    method: 'get',
    params
  })
}

// 区域信息类型
export interface RegionInfo {
  id: string
  name: string
  status: string
}

// 获取可同步区域列表
export function getAvailableRegions(accountId: number) {
  return request<{ items: RegionInfo[]; total: number }>({
    url: `/cloud-accounts/${accountId}/regions`,
    method: 'get'
  })
}

// 批量同步请求
export interface BatchSyncRequest {
  account_ids: number[]
  mode: 'full' | 'incremental'
  resource_types: string[]
}

export interface BatchSyncResult {
  account_id: number
  success: boolean
  message: string
  statistics?: any
}

// 批量同步云账号
export function batchSyncCloudAccounts(data: BatchSyncRequest) {
  return request<{ message: string; total: number; success: number; results: BatchSyncResult[] }>({
    url: '/cloud-accounts/batch-sync',
    method: 'post',
    data
  })
}

// 导出云账号列表
export function exportCloudAccounts() {
  return request<{ items: CloudAccount[]; total: number }>({
    url: '/cloud-accounts/export',
    method: 'get'
  })
}
