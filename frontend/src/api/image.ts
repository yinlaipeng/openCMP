import request from '@/utils/request'

export interface Image {
  id: string
  name: string
  description?: string
  status: string
  format: string
  os_name: string
  os_version?: string
  size: number
  cpu_arch: string
  architecture?: string
  image_type?: string
  share_scope?: string
  project_id?: string
  project_name?: string
  platform?: string
  account_id?: number
  account_name?: string
  created_at?: string
  updated_at?: string
}

export interface ImageListParams {
  page?: number
  page_size?: number
  details?: boolean
  is_guest_image?: boolean
  name?: string
  os_name?: string
  format?: string
  status?: string
  cpu_arch?: string
  account_id?: number
  project_id?: string
}

export interface ImageCreateParams {
  name: string
  description?: string
  os_name: string
  os_version?: string
  format: string
  cpu_arch: string
  project_id?: string
  tags?: string[]
}

export interface ImageUpdateParams {
  name?: string
  description?: string
  os_name?: string
  share_scope?: string
}

// 获取镜像列表
export function getImages(params?: ImageListParams) {
  return request<{ items: Image[]; total: number; pagination?: { total: number } }>({
    url: '/images',
    method: 'get',
    params
  })
}

// 获取单个镜像详情
export function getImage(id: string) {
  return request<Image>({
    url: `/images/${id}`,
    method: 'get'
  })
}

// 创建镜像（上传）
export function createImage(data: ImageCreateParams) {
  return request<Image>({
    url: '/images',
    method: 'post',
    data
  })
}

// 更新镜像信息
export function updateImage(id: string, data: ImageUpdateParams) {
  return request<Image>({
    url: `/images/${id}`,
    method: 'put',
    data
  })
}

// 删除镜像
export function deleteImage(id: string) {
  return request({
    url: `/images/${id}`,
    method: 'delete'
  })
}

// 批量删除镜像
export function batchDeleteImages(ids: string[]) {
  return request<{ total: number; success: number; failed: number }>({
    url: '/images/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// 同步镜像
export function syncImages(accountId: number) {
  return request({
    url: `/images/sync`,
    method: 'post',
    data: { account_id: accountId }
  })
}

// 获取社区镜像列表
export function getCommunityImages(params?: { os_name?: string }) {
  return request<Image[]>({
    url: '/images/community',
    method: 'get',
    params
  })
}

// 导入社区镜像
export function importCommunityImage(imageId: string, projectId?: string) {
  return request<Image>({
    url: `/images/community/${imageId}/import`,
    method: 'post',
    data: { project_id: projectId }
  })
}