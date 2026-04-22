import request from '@/utils/request'

// 代理类型定义
export interface Proxy {
  id: number
  name: string
  https_proxy?: string
  http_proxy?: string
  no_proxy?: string
  owner_domain?: string
  domain_id?: number
  shared_scope?: string
  created_at?: string
  updated_at?: string
}

export interface CreateProxyRequest {
  name: string
  domain_id: number
  https_proxy?: string
  http_proxy?: string
  no_proxy?: string
  shared_scope?: string
}

export interface UpdateProxyRequest {
  name?: string
  https_proxy?: string
  http_proxy?: string
  no_proxy?: string
  shared_scope?: string
}

export interface ProxySearchParams {
  page?: number
  page_size?: number
  name?: string
  domain_id?: number
}

// 获取代理列表
export function getProxies(params?: ProxySearchParams) {
  return request({
    url: '/proxies',
    method: 'get',
    params
  })
}

// 获取单个代理详情
export function getProxy(id: number) {
  return request({
    url: `/proxies/${id}`,
    method: 'get'
  })
}

// 创建代理
export function createProxy(data: CreateProxyRequest) {
  return request({
    url: '/proxies',
    method: 'post',
    data
  })
}

// 更新代理
export function updateProxy(id: number, data: UpdateProxyRequest) {
  return request({
    url: `/proxies/${id}`,
    method: 'put',
    data
  })
}

// 删除代理
export function deleteProxy(id: number) {
  return request({
    url: `/proxies/${id}`,
    method: 'delete'
  })
}

// 批量删除代理
export function batchDeleteProxies(ids: number[]) {
  return request({
    url: '/proxies/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// 设置共享
export function setProxySharing(id: number, sharedScope: string) {
  return request({
    url: `/proxies/${id}/sharing`,
    method: 'put',
    data: { shared_scope: sharedScope }
  })
}