import http from '../../../services/api'

// ==================== 文件上传 ====================

export function adminUploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return http.post('/admin/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// ==================== 前台文章 API ====================

// 文章列表（分页+筛选）
export function getArticles(params = {}) {
  return http.get('/articles', { params })
}

// 文章详情
export function getArticleByID(id) {
  return http.get(`/articles/${id}`)
}

// 热门文章
export function getHotArticles(limit = 10) {
  return http.get('/articles/hot', { params: { limit } })
}

// 最新文章
export function getRecentArticles(limit = 5) {
  return http.get('/articles/recent', { params: { limit } })
}

// ==================== 后台文章 API ====================

// 后台文章列表
export function adminGetArticles(params = {}) {
  return http.get('/admin/articles', { params })
}

// 后台文章详情
export function adminGetArticle(id) {
  return http.get(`/admin/articles/${id}`)
}

// 创建文章
export function adminCreateArticle(data) {
  return http.post('/admin/articles', data)
}

// 更新文章
export function adminUpdateArticle(id, data) {
  return http.put(`/admin/articles/${id}`, data)
}

// 删除文章
export function adminDeleteArticle(id) {
  return http.delete(`/admin/articles/${id}`)
}

// 切换文章状态
export function adminUpdateArticleStatus(id, status) {
  return http.patch(`/admin/articles/${id}/status`, { status })
}
