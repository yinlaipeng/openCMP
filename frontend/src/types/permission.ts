// 权限相关类型定义

export interface Permission {
  id: number
  name: string
  display_name: string
  description: string
  type: 'system' | 'custom' // 类型：系统/自定义
  scope: 'global' | 'domain' | 'project' // 权限范围：全局/域/项目
  domain_id?: number // 域ID，当scope为domain时有效
  project_id?: number // 项目ID，当scope为project时有效
  enabled: boolean // 是否启用
  is_public: boolean // 是否公开
  action: string // 操作权限，如 read, write, delete 等
  resource: string // 资源类型，如 vm, vpc, subnet 等
  conditions?: Record<string, any> // 条件限制
  created_at: string
  updated_at: string
}

export interface PermissionListResponse {
  items: Permission[]
  total: number
  page: number
  limit: number
}