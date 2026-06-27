// 全局响应式认证状态
import { reactive, readonly } from 'vue'

const state = reactive({
  token: localStorage.getItem('junblog_token') || '',
  user: JSON.parse(localStorage.getItem('junblog_user') || 'null'),
})

// 获取当前认证状态（只读）
export function useAuth() {
  return readonly(state)
}

// 保存登录信息
export function saveAuthState(token, refreshToken, user) {
  state.token = token
  state.user = user
  localStorage.setItem('junblog_token', token)
  localStorage.setItem('junblog_refresh_token', refreshToken)
  localStorage.setItem('junblog_user', JSON.stringify(user))
}

// 清除登录状态
export function clearAuthState() {
  state.token = ''
  state.user = null
  localStorage.removeItem('junblog_token')
  localStorage.removeItem('junblog_refresh_token')
  localStorage.removeItem('junblog_user')
}

// 获取存储的认证信息（兼容旧代码）
export function getStoredAuth() {
  return { token: state.token, user: state.user }
}
