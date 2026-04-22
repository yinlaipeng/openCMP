import request from '@/utils/request'

// 云用户组类型定义
export interface CloudUserGroup {
  id: number
  name: string
  status: string
  permissions?: string
  platform?: string
  cloud_accounts?: string
  shared_scope?: string
  owner_domain?: string
  domain_id?: number
  created_at?: string
  updated_at?: string
}

export interface CreateCloudUserGroupRequest {
  name: string
  domain_id: number
  permissions?: string
  platform?: string
  cloud_account_ids?: number[]
  shared_scope?: string
}

export interface UpdateCloudUserGroupRequest {
  name?: string
  permissions?: string
  platform?: string
  cloud_account_ids?: number[]
  shared_scope?: string
}

export interface CloudUserGroupSearchParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  platform?: string
  domain_id?: number
}

// 获取云用户组列表
export function getCloudUserGroups(params?: CloudUserGroupSearchParams) {
  return request({
    url: '/cloud-user-groups',
    method: 'get',
    params
  })
}

// 获取单个云用户组详情
export function getCloudUserGroup(id: number) {
  return request({
    url: `/cloud-user-groups/${id}`,
    method: 'get'
  })
}

// 创建云用户组
export function createCloudUserGroup(data: CreateCloudUserGroupRequest) {
  return request({
    url: '/cloud-user-groups',
    method: 'post',
    data
  })
}

// 更新云用户组
export function updateCloudUserGroup(id: number, data: UpdateCloudUserGroupRequest) {
  return request({
    url: `/cloud-user-groups/${id}`,
    method: 'put',
    data
  })
}

// 删除云用户组
export function deleteCloudUserGroup(id: number) {
  return request({
    url: `/cloud-user-groups/${id}`,
    method: 'delete'
  })
}

// 批量删除云用户组
export function batchDeleteCloudUserGroups(ids: number[]) {
  return request({
    url: '/cloud-user-groups/batch-delete',
    method: 'post',
    data: { ids }
  })
}