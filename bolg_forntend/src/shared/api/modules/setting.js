import http from '../../../services/api'

// 获取关于页面数据（公开接口）
export const getAboutPage = () => {
  return http.get('/about')
}

// 管理员获取所有设置
export const adminGetSettings = () => {
  return http.get('/admin/settings')
}

// 管理员获取单个设置
export const adminGetSetting = (key) => {
  return http.get(`/admin/settings/${key}`)
}

// 管理员更新单个设置
export const adminUpdateSetting = (key, value) => {
  return http.put(`/admin/settings/${key}`, { value })
}

// 管理员更新关于页面
export const adminUpdateAboutPage = (data) => {
  return http.put('/admin/settings/about', data)
}
