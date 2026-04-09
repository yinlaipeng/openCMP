// IAM 相关类型定义

export * from './permission'

export interface Domain {
  id: number
  name: string
  description: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface Project {
  id: number
  name: string
  description: string
  domain_id: number
  manager_id?: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  name: string
  display_name: string
  remark: string
  email: string
  phone: string
  enabled: boolean
  console_login: boolean
  mfa_enabled: boolean
  domain_id: number
  created_at: string
  updated_at: string
}

export interface Group {
  id: number
  name: string
  description: string
  domain_id: number
  created_at: string
  updated_at: string
}

export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  type: 'system' | 'custom'
  enabled: boolean
  is_public: boolean
  domain_id: number
  created_at: string
  updated_at: string
}

export interface AuthSource {
  id: number
  name: string
  type: 'ldap' | 'local' | 'sql'
  scope: 'system' | 'domain'
  description: string
  config: Record<string, any>
  enabled: boolean
  running: boolean
  auto_create: boolean
  domain_id: number
  created_at: string
  updated_at: string
  protocol: string
  sync_status: string
  is_default: boolean
}