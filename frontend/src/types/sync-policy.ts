export interface Rule {
  id?: number
  condition_type: 'all_match' | 'any_match' | 'key_match' // 全部匹配标签、至少一个标签、根据标签key匹配
  tags: Array<{
    key: string
    value: string
  }> // 标签键值对
  resource_mapping: 'specify_project' | 'specify_name' // 资源映射方式
  target_project_id?: number // 指定项目ID
  target_project_name?: string // 指定项目名称
}

export interface SyncPolicy {
  id: number
  name: string
  remarks?: string
  status: string
  enabled: boolean
  rules: Rule[] // 规则数组
  scope: string
  domain_id: number
  created_at: string
  updated_at: string
}

export interface CreateSyncPolicyRequest {
  name: string
  remarks?: string
  status?: string
  enabled?: boolean
  rules: Rule[] // 规则数组
  scope: string
  domain_id: number
}