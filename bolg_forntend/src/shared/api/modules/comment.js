import http from '../../../services/api'

// 获取文章评论
export function getComments(articleId) {
  return http.get('/comments', { params: { article_id: articleId } })
}

// 发表评论
export function createComment(data) {
  return http.post('/comments', data)
}

// 后台评论列表
export function adminGetComments(params = {}) {
  return http.get('/admin/comments', { params })
}

// 审核评论
export function adminUpdateCommentStatus(id, status) {
  return http.patch(`/admin/comments/${id}/status`, { status })
}

// 删除评论
export function adminDeleteComment(id) {
  return http.delete(`/admin/comments/${id}`)
}
