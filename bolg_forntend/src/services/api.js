import axios from 'axios'
import { clearAuthState } from '../shared/stores/auth'
import { useMessage } from '../composables/useMessage'

const message = useMessage()

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 30000, // 30秒超时，GitHub OAuth 需要较长时间
})

// 是否正在刷新 token
let isRefreshing = false
let refreshSubscribers = []

// token 刷新后执行等待中的请求
function onRefreshed(newToken) {
  refreshSubscribers.forEach((callback) => callback(newToken))
  refreshSubscribers = []
}

// 将请求加入等待队列
function addRefreshSubscriber(callback) {
  refreshSubscribers.push(callback)
}

// 请求拦截器：自动附加 JWT token
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('junblog_token')
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：处理 401 自动刷新 token，处理 403 账号被封禁
http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // 如果是 403 账号被封禁，清除登录状态并跳转
    if (error.response?.status === 403) {
      clearAuthState()
      message.error('账号已被封禁，请联系管理员')
      setTimeout(() => {
        window.location.href = '/login'
      }, 1500)
      return Promise.reject(error)
    }

    // 如果是 401 且不是刷新/登录请求本身，且未重试过
    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !originalRequest.url.includes('/auth/refresh') &&
      !originalRequest.url.includes('/auth/login')
    ) {
      if (isRefreshing) {
        // 正在刷新中，将请求加入队列等待
        return new Promise((resolve) => {
          addRefreshSubscriber((newToken) => {
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            resolve(http(originalRequest))
          })
        })
      }

      isRefreshing = true
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('junblog_refresh_token')
        if (!refreshToken) {
          throw new Error('无 refresh token')
        }

        const { data } = await axios.post('/api/v1/auth/refresh', {
          token: refreshToken,
        })

        const newToken = data.data.token
        const newRefreshToken = data.data.refresh_token

        localStorage.setItem('junblog_token', newToken)
        localStorage.setItem('junblog_refresh_token', newRefreshToken)

        onRefreshed(newToken)

        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return http(originalRequest)
      } catch (refreshError) {
        // 刷新失败，清除登录状态跳转登录
        clearAuthState()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

// 认证状态已迁移至 shared/stores/auth.js（响应式）
// 此处保留兼容导出，新代码请直接从 shared/stores/auth 引入
export { getStoredAuth, saveAuthState as saveAuth, clearAuthState as clearAuth } from '../shared/stores/auth'

// === 认证 API ===

export async function login(payload) {
  const { data } = await http.post('/auth/login', payload)
  return data
}

export async function refreshToken(refreshTokenStr) {
  const { data } = await http.post('/auth/refresh', {
    token: refreshTokenStr,
  })
  return data
}

// === 用户 API ===

export async function getProfile() {
  const { data } = await http.get('/user/profile')
  return data
}

export async function updateProfile(payload) {
  const { data } = await http.put('/user/profile', payload)
  return data
}

export async function updatePassword(payload) {
  const { data } = await http.put('/user/profile/password', payload)
  return data
}

// === GitHub 登录 ===

export async function githubLogin(code) {
  // 动态获取当前访问地址作为回调地址
  const redirect_uri = `${window.location.origin}/auth/github/callback`
  const { data } = await http.post('/auth/github/callback', { code, redirect_uri })
  return data
}

export default http
