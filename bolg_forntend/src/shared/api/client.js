// 统一 API 客户端 - 从 services/api.js 导出
// 后续可在此扩展更多 API 模块
export { default as http } from '../../services/api'
export {
  getStoredAuth,
  saveAuth,
  clearAuth,
  login,
  refreshToken,
  getProfile,
  updateProfile,
  updatePassword,
} from '../../services/api'
