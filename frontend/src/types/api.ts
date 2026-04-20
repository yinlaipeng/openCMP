/**
 * 共享API类型定义
 * 统一API响应和请求参数类型
 */

/**
 * 分页响应基础类型
 */
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

/**
 * 单一资源响应类型
 */
export interface SingleResponse<T> {
  data: T
  message?: string
}

/**
 * 操作结果响应类型
 */
export interface OperationResponse {
  success: boolean
  message: string
}

/**
 * 分页请求参数
 */
export interface PaginationParams {
  page?: number
  page_size?: number
}

/**
 * 资源同步选项
 */
export interface SyncOptions {
  mode: 'full' | 'incremental'
  resource_types: string[]
}

/**
 * 批量操作结果
 */
export interface BatchOperationResult<T> {
  success: number
  failed: number
  total: number
  results: Array<{
    id: number
    success: boolean
    message: string
    data?: T
  }>
}

/**
 * 区域信息
 */
export interface RegionInfo {
  id: string
  name: string
  status: string
}

/**
 * 可用区信息
 */
export interface ZoneInfo {
  id: string
  name: string
  region_id: string
  status: string
}