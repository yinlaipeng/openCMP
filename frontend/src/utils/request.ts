import axios from 'axios'
import type { AxiosInstance, AxiosResponse, AxiosError } from 'axios'
import { ElMessage } from 'element-plus'

// 错误码对应的中文提示
const ERROR_MESSAGES: Record<string, string> = {
  // 认证相关
  'TOKEN_EXPIRED': '登录已过期，请重新登录',
  'TOKEN_INVALID': '登录凭证无效，请重新登录',
  'UNAUTHORIZED': '未登录或无权限访问',
  'PERMISSION_DENIED': '权限不足，无法执行此操作',
  'PROJECT_ACCESS_DENIED': '无权访问该项目',

  // 资源相关
  'RESOURCE_NOT_FOUND': '资源不存在',
  'ACCOUNT_NOT_FOUND': '云账号不存在',
  'VM_NOT_FOUND': '虚拟机不存在',
  'VPC_NOT_FOUND': 'VPC不存在',

  // 操作相关
  'OPERATION_FAILED': '操作失败',
  'VALIDATION_ERROR': '参数验证失败',
  'INVALID_PARAMETER': '参数格式错误',

  // 云厂商相关
  'CLOUD_API_ERROR': '云厂商API调用失败',
  'CREDENTIAL_INVALID': '云账号凭证无效',
  'REGION_NOT_FOUND': '区域不存在',
  'QUOTA_EXCEEDED': '配额不足',

  // 网络相关
  'NETWORK_ERROR': '网络连接失败',
  'TIMEOUT_ERROR': '请求超时，请稍后重试'
}

// 最大重试次数
const MAX_RETRY_COUNT = 3
// 重试间隔（毫秒）
const RETRY_DELAY = 2000
// 需要重试的错误状态码
const RETRY_STATUS_CODES = [408, 429, 500, 502, 503, 504]

const service: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

// 请求拦截器：自动附加 token
service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 获取友好的错误消息
function getFriendlyErrorMessage(error: AxiosError): string {
  const data = error.response?.data as any
  const code = data?.code
  const status = error.response?.status

  // 优先使用后端返回的code对应的中文消息
  if (code && ERROR_MESSAGES[code]) {
    return ERROR_MESSAGES[code]
  }

  // 根据状态码返回消息
  if (status === 403) {
    return '权限不足，无法执行此操作'
  }
  if (status === 404) {
    return data?.error || '请求的资源不存在'
  }
  if (status === 400) {
    return data?.error || '请求参数错误'
  }
  if (status === 422) {
    return '参数验证失败，请检查输入'
  }
  if (status === 500) {
    return '服务器内部错误，请稍后重试'
  }
  if (status === 502 || status === 503 || status === 504) {
    return '服务暂时不可用，请稍后重试'
  }
  if (status === 408) {
    return '请求超时，请稍后重试'
  }
  if (status === 429) {
    return '请求过于频繁，请稍后重试'
  }

  // 使用后端返回的error字段
  if (data?.error) {
    return data.error
  }

  return '请求失败'
}

// 响应拦截器：统一处理错误，支持自动重试
service.interceptors.response.use(
  (response: AxiosResponse) => response.data,
  async (error: AxiosError) => {
    const config = error.config as any

    // 初始化重试计数
    if (!config.__retryCount) {
      config.__retryCount = 0
    }

    // 判断是否需要重试
    const shouldRetry =
      config.__retryCount < MAX_RETRY_COUNT &&
      (error.response?.status && RETRY_STATUS_CODES.includes(error.response?.status)) ||
      (error.code === 'ECONNABORTED' || error.code === 'ETIMEDOUT') ||
      !error.response

    if (shouldRetry) {
      config.__retryCount++

      // 延迟后重试
      await new Promise(resolve => setTimeout(resolve, RETRY_DELAY))

      // 显示重试提示
      ElMessage.warning(`请求失败，正在重试（第${config.__retryCount}次）`)

      return service.request(config)
    }

    if (error.response) {
      const { status, data } = error.response
      const requestUrl: string = error.config?.url || ''

      // 登录接口的错误由页面自己处理，不走全局拦截
      // 检查 URL 是否包含 login（使用更宽松的匹配）
      const isLoginRequest = requestUrl.includes('/auth/login') || requestUrl === '/auth/login'

      // 权限和用户信息接口的错误也不清除 token（登录后调用）
      const isAuthInfoRequest = requestUrl.includes('/auth/permissions') ||
                                requestUrl.includes('/auth/user') ||
                                requestUrl.includes('/auth/me')

      if (status === 401 && !isLoginRequest && !isAuthInfoRequest) {
        // 清除本地凭证
        localStorage.removeItem('token')
        localStorage.removeItem('user')

        const code = (data as any)?.code
        if (code === 'TOKEN_EXPIRED') {
          ElMessage.warning(ERROR_MESSAGES['TOKEN_EXPIRED'])
        } else {
          ElMessage.error(ERROR_MESSAGES['UNAUTHORIZED'])
        }

        // 避免在登录页重复跳转
        if (window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
      } else if (!isLoginRequest) {
        ElMessage.error(getFriendlyErrorMessage(error))
      }
    } else if (error.code === 'ECONNABUTED' || error.code === 'ETIMEDOUT') {
      ElMessage.error(ERROR_MESSAGES['TIMEOUT_ERROR'])
    } else if (!error.request) {
      ElMessage.error(ERROR_MESSAGES['NETWORK_ERROR'])
    } else {
      ElMessage.error('请求发送失败')
    }
    return Promise.reject(error)
  }
)

export default service
