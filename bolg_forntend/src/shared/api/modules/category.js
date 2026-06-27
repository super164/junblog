import http from '../../../services/api'

// 分类树（前台）
export function getCategories() {
  return http.get('/categories')
}

// 后台分类列表
export function adminGetCategories() {
  return http.get('/admin/categories')
}

// 创建分类
export function adminCreateCategory(data) {
  return http.post('/admin/categories', data)
}

// 更新分类
export function adminUpdateCategory(id, data) {
  return http.put(`/admin/categories/${id}`, data)
}

// 删除分类
export function adminDeleteCategory(id) {
  return http.delete(`/admin/categories/${id}`)
}
