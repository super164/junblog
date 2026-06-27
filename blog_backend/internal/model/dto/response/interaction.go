package response

import "time"

// LikeResponse 点赞响应
type LikeResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}
