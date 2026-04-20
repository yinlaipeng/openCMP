import { computed, type Ref } from 'vue'

/**
 * 资源归属方式配置项
 */
export interface ResourceAssignmentOption {
  key: string
  label: string
  description: string // 单独勾选时的说明
  priority: number // 优先级（数字越小优先级越高）
  fallbackDescription: string // 作为兜底时的说明
}

/**
 * 预定义的资源归属方式选项
 */
const ASSIGNMENT_OPTIONS: ResourceAssignmentOption[] = [
  {
    key: 'tag_mapping',
    label: '根据同步策略归属',
    description: '资源会根据同步策略归属，若资源不匹配同步策略，则归属到缺省项目',
    priority: 1,
    fallbackDescription: '不匹配同步策略的资源'
  },
  {
    key: 'project_mapping',
    label: '根据云上项目归属',
    description: '资源会同步到与云上项目同名的本地项目中，若资源无云上项目属性，则归属到缺省项目',
    priority: 2,
    fallbackDescription: '无云上项目属性的 resources'
  },
  {
    key: 'subscription_mapping',
    label: '根据云订阅归属',
    description: '资源会同步到与云订阅同名的本地项目中，若资源不属于云订阅，则归属到缺省项目',
    priority: 3,
    fallbackDescription: '不属于云订阅的资源'
  },
  {
    key: 'specify_project',
    label: '指定项目',
    description: '所有资源均归属到指定的目标项目',
    priority: 4, // 最后兜底
    fallbackDescription: '所有未匹配的资源'
  }
]

/**
 * 根据选中的归属方式生成组合说明文案
 */
export function generateAssignmentDescription(selectedMethods: string[]): string {
  if (selectedMethods.length === 0) {
    return '请选择资源归属方式'
  }

  // 按优先级排序选中的方式
  const sortedMethods = selectedMethods
    .map(key => ASSIGNMENT_OPTIONS.find(opt => opt.key === key))
    .filter(Boolean)
    .sort((a, b) => (a?.priority || 0) - (b?.priority || 0))

  if (sortedMethods.length === 1) {
    // 只勾选一项，返回单项说明
    return sortedMethods[0]?.description || ''
  }

  // 多项勾选，生成优先级说明
  const parts: string[] = []

  // 第一项：优先归属
  const first = sortedMethods[0]
  if (first) {
    parts.push(`资源会优先${first.label.replace('根据', '')}`)
  }

  // 中间项：依次兜底
  for (let i = 1; i < sortedMethods.length - 1; i++) {
    const current = sortedMethods[i]
    const prev = sortedMethods[i - 1]
    if (current && prev) {
      parts.push(`${prev.fallbackDescription}会${current.label.replace('根据', '')}`)
    }
  }

  // 最后项：最终兜底
  const last = sortedMethods[sortedMethods.length - 1]
  if (last && sortedMethods.length > 1) {
    if (last.key === 'specify_project') {
      parts.push(`最终不匹配的资源归属到指定项目`)
    } else {
      parts.push(`若仍不匹配则${last.label.replace('根据', '')}`)
    }
  }

  return parts.join('；') + '。'
}

/**
 * 获取当前选中方式需要显示的控件
 */
export function getVisibleControls(selectedMethods: string[]): {
  showSyncPolicySelector: boolean
  showSyncScopeSelector: boolean
  showSpecifyProjectSelector: boolean
  showDefaultProjectSelector: boolean
} {
  return {
    // 勾选"根据同步策略归属"时显示同步策略选择器
    showSyncPolicySelector: selectedMethods.includes('tag_mapping'),
    // 勾选"根据同步策略归属"时显示同步策略生效范围
    showSyncScopeSelector: selectedMethods.includes('tag_mapping'),
    // 勾选"指定项目"时显示指定项目选择器（必填）
    showSpecifyProjectSelector: selectedMethods.includes('specify_project'),
    // 只要存在兜底逻辑就显示缺省项目选择器
    showDefaultProjectSelector: selectedMethods.length > 1 ||
      !selectedMethods.includes('specify_project')
  }
}

/**
 * 组合函数：资源归属方式说明生成器
 */
export function useResourceAssignmentDescription(methodsRef: Ref<string[]>) {
  // 动态说明文案
  const description = computed(() => generateAssignmentDescription(methodsRef.value))

  // 控件显示状态
  const visibleControls = computed(() => getVisibleControls(methodsRef.value))

  // 是否需要选择缺省项目（校验提示）
  const needsDefaultProject = computed(() => {
    // 如果勾选了多项且没有"指定项目"，则需要选择缺省项目
    return methodsRef.value.length > 1 && !methodsRef.value.includes('specify_project')
  })

  // 是否需要选择指定项目（校验提示）
  const needsSpecifyProject = computed(() => {
    return methodsRef.value.includes('specify_project')
  })

  // 校验规则提示
  const validationHint = computed(() => {
    if (methodsRef.value.length === 0) {
      return '请至少选择一种资源归属方式'
    }
    if (needsSpecifyProject.value) {
      return '已选择"指定项目"，请选择目标项目'
    }
    if (needsDefaultProject.value) {
      return '已选择多种归属方式，请选择缺省项目作为最终兜底'
    }
    return ''
  })

  return {
    description,
    visibleControls,
    needsDefaultProject,
    needsSpecifyProject,
    validationHint,
    options: ASSIGNMENT_OPTIONS
  }
}

export default useResourceAssignmentDescription