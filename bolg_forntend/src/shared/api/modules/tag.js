import http from '../../../services/api'

// 标签列表（前台）
export function getTags() {
  return http.get('/tags')
}

// 根据分类ID获取标签列表（前台）
export function getTagsByCategory(categoryId) {
  return http.get(`/tags/category/${categoryId}`)
}

// 后台标签列表
export function adminGetTags() {
  return http.get('/admin/tags')
}

// 后台根据分类ID获取标签列表
export function adminGetTagsByCategory(categoryId) {
  return http.get(`/admin/tags/category/${categoryId}`)
}

// 创建标签
export function adminCreateTag(data) {
  return http.post('/admin/tags', data)
}

// 更新标签
export function adminUpdateTag(id, data) {
  return http.put(`/admin/tags/${id}`, data)
}

// 删除标签
export function adminDeleteTag(id) {
  return http.delete(`/admin/tags/${id}`)
}
