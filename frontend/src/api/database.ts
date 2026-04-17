import request from '@/utils/request'

// RDS Types
export interface RDSInstance {
  id: string
  name: string
  engine: string
  engine_version: string
  instance_type: string
  storage_size: number
  storage_type: string
  status: string
  vpc_id: string
  subnet_id: string
  endpoint: string
  port: number
  master_username: string
  tags: Record<string, string>
  created_at: string
  zone_id: string
}

export interface RDSConfig {
  account_id: number
  name: string
  engine: string
  engine_version: string
  instance_type: string
  storage_size: number
  storage_type: string
  vpc_id: string
  subnet_id: string
  zone_id: string
  master_username?: string
  master_password?: string
  tags?: Record<string, string>
}

export interface RDSBackup {
  id: string
  name: string
  instance_id: string
  status: string
  size: number
  start_time: string
  end_time: string
}

export interface RDSFilter {
  account_id: number
  instance_id?: string
  engine?: string
  status?: string
  vpc_id?: string
}

// Redis/Cache Types
export interface CacheInstance {
  id: string
  name: string
  engine: string
  engine_version: string
  instance_type: string
  status: string
  vpc_id: string
  subnet_id: string
  endpoint: string
  port: number
  tags: Record<string, string>
  created_at: string
  zone_id: string
}

export interface CacheConfig {
  account_id: number
  name: string
  engine: string
  engine_version: string
  instance_type: string
  vpc_id: string
  subnet_id: string
  zone_id: string
  tags?: Record<string, string>
}

export interface CacheBackup {
  id: string
  name: string
  instance_id: string
  status: string
  start_time: string
  end_time: string
}

export interface CacheFilter {
  account_id: number
  instance_id?: string
  engine?: string
  status?: string
}

// RDS API
export function listRDS(filter: RDSFilter) {
  return request.get<RDSInstance[]>('/database/rds', { params: filter })
}

export function getRDS(accountId: number, instanceId: string) {
  return request.get<RDSInstance>(`/database/rds/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function createRDS(config: RDSConfig) {
  return request.post<RDSInstance>('/database/rds', config)
}

export function deleteRDS(accountId: number, instanceId: string) {
  return request.delete(`/database/rds/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function rdsAction(accountId: number, instanceId: string, action: string) {
  return request.post(`/database/rds/${instanceId}/action`, { action }, {
    params: { account_id: accountId }
  })
}

export function resizeRDS(accountId: number, instanceId: string, instanceType: string, storageSize?: number) {
  return request.post(`/database/rds/${instanceId}/resize`, {
    instance_type: instanceType,
    storage_size: storageSize
  }, {
    params: { account_id: accountId }
  })
}

export function createRDSBackup(accountId: number, instanceId: string, name?: string) {
  return request.post<RDSBackup>(`/database/rds/${instanceId}/backups`, null, {
    params: { account_id: accountId, name }
  })
}

export function listRDSBackups(accountId: number, instanceId: string) {
  return request.get<RDSBackup[]>(`/database/rds/${instanceId}/backups`, {
    params: { account_id: accountId }
  })
}

// Cache/Redis API
export function listCache(filter: CacheFilter) {
  return request.get<CacheInstance[]>('/database/cache', { params: filter })
}

export function createCache(config: CacheConfig) {
  return request.post<CacheInstance>('/database/cache', config)
}

export function deleteCache(accountId: number, instanceId: string) {
  return request.delete(`/database/cache/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function cacheAction(accountId: number, instanceId: string, action: string) {
  return request.post(`/database/cache/${instanceId}/action`, { action }, {
    params: { account_id: accountId }
  })
}

export function resizeCache(accountId: number, instanceId: string, instanceType: string) {
  return request.post(`/database/cache/${instanceId}/resize`, {
    instance_type: instanceType
  }, {
    params: { account_id: accountId }
  })
}

export function createCacheBackup(accountId: number, instanceId: string) {
  return request.post<CacheBackup>(`/database/cache/${instanceId}/backups`, null, {
    params: { account_id: accountId }
  })
}