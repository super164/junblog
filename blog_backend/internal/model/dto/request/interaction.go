package request

// LikeRequest 点赞请求
type LikeRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
}

// FavoriteRequest 收藏请求
type FavoriteRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
}
