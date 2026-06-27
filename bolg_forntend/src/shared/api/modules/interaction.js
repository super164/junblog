import http from '../../../services/api'

// 获取互动状态
export function getInteractionStatus(articleId) {
  return http.get(`/interactions/articles/${articleId}`)
}

// 点赞/取消点赞
export function toggleLike(articleId) {
  return http.post(`/interactions/articles/${articleId}/like`)
}

// 收藏/取消收藏
export function toggleFavorite(articleId) {
  return http.post(`/interactions/articles/${articleId}/favorite`)
}
