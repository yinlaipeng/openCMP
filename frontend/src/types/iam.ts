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
  config: {
    url?: string
    base_dn?: string
    bind_dn?: string
    bind_password?: string
    user_filter?: string
    user_id_attr?: string
    user_name_attr?: string
    user_search_base?: string
    group_search_base?: string
    user_enabled_attribute?: string
    protocol?: string
    auth_type?: string
    target_domain?: string
  }
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

export interface OperationLog {
  id: number
  operation_time: string
  resource_name: string
  resource_type: string
  operation_type: string
  service_type: string
  risk_level: string
  time_type: string
  result: string
  operator: string
  project_id: number
  created_at: string
  updated_at: string
}