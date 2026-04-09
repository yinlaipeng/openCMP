import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

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

// 响应拦截器：统一处理错误，401 跳转登录
service.interceptors.response.use(
  (response: AxiosResponse) => response.data,
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      const requestUrl: string = error.config?.url || ''

      // 登录接口的错误由页面自己处理，不走全局拦截
      const isLoginRequest = requestUrl.includes('/auth/login')

      if (status === 401 && !isLoginRequest) {
        // 清除本地凭证
        localStorage.removeItem('token')
        localStorage.removeItem('user')

        const code = data?.code
        if (code === 'TOKEN_EXPIRED') {
          ElMessage.warning('登录已过期，请重新登录')
        } else {
          ElMessage.error('未登录或登录已失效，请重新登录')
        }

        // 避免在登录页重复跳转
        if (window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
      } else if (!isLoginRequest) {
        ElMessage.error(data?.error || `请求失败（${status}）`)
      }
    } else if (error.request) {
      ElMessage.error('网络异常，请检查网络连接')
    } else {
      ElMessage.error('请求发送失败')
    }
    return Promise.reject(error)
  }
)

export default service
