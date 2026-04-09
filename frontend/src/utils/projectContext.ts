// utils/projectContext.ts
import { ref, computed } from 'vue'

// 全局项目上下文状态
const selectedProjectId = ref<number | null>(null)
const selectedProjectName = ref<string>('')

/**
 * 初始化项目上下文
 */
export const initializeProjectContext = () => {
  const storedProjectId = localStorage.getItem('selectedProjectId')
  const storedProjectName = localStorage.getItem('selectedProjectName')

  if (storedProjectId) {
    selectedProjectId.value = parseInt(storedProjectId)
  }

  if (storedProjectName) {
    selectedProjectName.value = storedProjectName
  }
}

/**
 * 设置项目上下文
 */
export const setProjectContext = (id: number | null, name: string) => {
  selectedProjectId.value = id
  selectedProjectName.value = name

  if (id) {
    localStorage.setItem('selectedProjectId', id.toString())
    localStorage.setItem('selectedProjectName', name)
  } else {
    localStorage.removeItem('selectedProjectId')
    localStorage.removeItem('selectedProjectName')
  }
}

/**
 * 获取当前项目ID
 */
export const getCurrentProjectId = () => {
  // 优先使用全局状态，如果未设置则尝试从localStorage读取
  if (selectedProjectId.value !== null) {
    return selectedProjectId.value
  }

  const storedId = localStorage.getItem('selectedProjectId')
  return storedId ? parseInt(storedId) : null
}

/**
 * 获取当前项目名称
 */
export const getCurrentProjectName = () => {
  // 优先使用全局状态，如果未设置则尝试从localStorage读取
  if (selectedProjectName.value) {
    return selectedProjectName.value
  }

  return localStorage.getItem('selectedProjectName') || ''
}

/**
 * 检查是否处于项目模式
 */
export const isInProjectMode = () => {
  return getCurrentProjectId() !== null
}

/**
 * 清除项目上下文
 */
export const clearProjectContext = () => {
  selectedProjectId.value = null
  selectedProjectName.value = ''

  localStorage.removeItem('selectedProjectId')
  localStorage.removeItem('selectedProjectName')
}

/**
 * 获取基础查询参数（包含项目ID过滤）
 */
export const getBaseQueryParams = () => {
  const params: Record<string, any> = {}
  if (isInProjectMode()) {
    const projectId = getCurrentProjectId()
    if (projectId) {
      params.project_id = projectId
    }
  }
  return params
}