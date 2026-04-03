// IAM 相关类型定义

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
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  name: string
  display_name: string
  email: string
  phone: string
  enabled: boolean
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

export interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
  type: 'system' | 'custom'
  description: string
  domain_id: number
  created_at: string
  updated_at: string
}

export interface Policy {
  id: string
  name: string
  description: string
  scope: 'system' | 'domain' | 'project'
  type: 'system' | 'custom'
  enabled: boolean
  is_system: boolean
  is_public: boolean
  domain_id: string
  policy: Record<string, any>
  created_at: string
  updated_at: string
  can_update: boolean
  can_delete: boolean
  delete_fail_reason?: {
    details: string
  }
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

export interface PolicyStatement {
  effect: 'allow' | 'deny'
  resource: string
  action: string[]
  service: string
}