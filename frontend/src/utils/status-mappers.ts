/**
 * 共享状态映射工具
 * 统一各组件的状态显示逻辑
 */

// 状态标签类型
type TagType = '' | 'success' | 'warning' | 'danger' | 'info'

// 资源状态标签映射
export const STATUS_LABELS: Record<string, string> = {
  Running: '运行中',
  Stopped: '已停止',
  Starting: '启动中',
  Stopping: '停止中',
  Pending: '创建中',
  Error: '错误',
  Deleted: '已删除',
  Available: '可用',
  Unavailable: '不可用',
  Creating: '创建中',
  Terminating: '终止中',
  active: '活跃',
  inactive: '不活跃',
  connected: '已连接',
  disconnected: '连接断开',
  checking: '检测中'
}

// 状态标签类型映射
export const STATUS_TAG_TYPES: Record<string, TagType> = {
  Running: 'success',
  Stopped: 'info',
  Starting: 'warning',
  Stopping: 'warning',
  Pending: 'warning',
  Creating: 'warning',
  Error: 'danger',
  Deleted: 'info',
  Available: 'success',
  Unavailable: 'danger',
  active: 'success',
  inactive: 'info',
  connected: 'success',
  disconnected: 'danger',
  checking: 'warning'
}

// 云平台标签映射
export const PROVIDER_LABELS: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure',
  huawei: '华为云',
  google: 'GCP',
  openstack: 'OpenStack',
  vmware: 'VMware'
}

// 云平台标签类型映射
export const PROVIDER_TAG_TYPES: Record<string, TagType> = {
  alibaba: 'primary',
  tencent: 'warning',
  aws: 'success',
  azure: 'info',
  huawei: 'warning',
  google: 'success',
  openstack: 'info',
  vmware: 'info'
}

// 健康状态标签映射
export const HEALTH_STATUS_LABELS: Record<string, string> = {
  healthy: '正常',
  unhealthy: '异常',
  unknown: '未知'
}

// 健康状态标签类型映射
export const HEALTH_TAG_TYPES: Record<string, TagType> = {
  healthy: 'success',
  unhealthy: 'danger',
  unknown: 'info'
}

/**
 * 获取状态显示名称
 */
export function getStatusLabel(status: string): string {
  return STATUS_LABELS[status] ?? status
}

/**
 * 获取状态标签类型
 */
export function getStatusTagType(status: string): TagType {
  return STATUS_TAG_TYPES[status] ?? ''
}

/**
 * 获取云平台显示名称
 */
export function getProviderLabel(providerType: string): string {
  return PROVIDER_LABELS[providerType] ?? providerType ?? '未知'
}

/**
 * 获取云平台标签类型
 */
export function getProviderTagType(providerType: string): TagType {
  return PROVIDER_TAG_TYPES[providerType] ?? 'info'
}

/**
 * 获取健康状态显示名称
 */
export function getHealthStatusLabel(healthStatus?: string): string {
  if (!healthStatus) return '未知'
  return HEALTH_STATUS_LABELS[healthStatus] ?? '未知'
}

/**
 * 获取健康状态标签类型
 */
export function getHealthTagType(healthStatus?: string): TagType {
  if (!healthStatus) return 'info'
  return HEALTH_TAG_TYPES[healthStatus] ?? 'info'
}

/**
 * 获取启用状态显示名称
 */
export function getEnabledLabel(enabled: boolean): string {
  return enabled ? '启用' : '禁用'
}

/**
 * 获取启用状态标签类型
 */
export function getEnabledTagType(enabled: boolean): TagType {
  return enabled ? 'success' : 'info'
}