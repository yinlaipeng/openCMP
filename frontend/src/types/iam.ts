// IAM 相关类型定义

export * from './permission'

export interface Domain {
  id: number | string
  name: string
  description: string
  enabled: boolean
  user_count?: number
  group_count?: number
  project_count?: number
  role_count?: number
  policy_count?: number
  idp_count?: number
  is_sso?: boolean
  ext_resource?: {
    cloudaccounts?: number
    cloudroles?: number
    cloudusers?: number
    hosts?: number
    storages?: number
    vpcs?: number
    route_tables?: number
    wires?: number
    proxysettings?: number
    project_mappings?: number
  }
  can_delete?: boolean
  can_update?: boolean
  delete_fail_reason?: {
    class?: string
    code?: number
    details?: string
  }
  update_fail_reason?: {
    class?: string
    code?: number
    details?: string
  }
  created_at: string
  updated_at: string
}

export interface Project {
  id: number | string
  name: string
  description: string
  domain_id: number | string
  manager_id?: number
  admin?: string
  admin_id?: string
  enabled: boolean
  user_count?: number
  group_count?: number
  is_system?: boolean
  can_delete?: boolean
  can_update?: boolean
  ext_resource?: {
    servers?: number
    disks?: number
    networks?: number
    secgroups?: number
    scheduledtasks?: number
    loadbalancercertificates?: number
    dns_zones?: number
    externalprojects?: number
  }
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
  user_count?: number
  project_count?: number
  is_sso?: boolean
  can_delete?: boolean
  can_update?: boolean
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
    user_enabled_status?: 'enabled' | 'disabled'
    protocol?: string
    auth_type?: 'ad_single' | 'ad_multiple' | 'openldap'
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