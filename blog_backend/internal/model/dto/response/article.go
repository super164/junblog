package response

import "time"

// ArticleResponse 文章详情响应
type ArticleResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Content        string    `json:"content"`
	Cover          string    `json:"cover"`
	Status         string    `json:"status"`
	ViewsCount     uint      `json:"views_count"`
	LikesCount     int       `json:"likes_count"`
	FavoritesCount int       `json:"favorites_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Author   UserResponse     `json:"author"`
	Category CategoryResponse `json:"category"`
	Tags     []TagResponse    `json:"tags"`
}

// ArticleListItem 文章列表项（不含正文内容）
type ArticleListItem struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Summary        string    `json:"summary"`
	Cover          string    `json:"cover"`
	Status         string    `json:"status"`
	ViewsCount     uint      `json:"views_count"`
	LikesCount     int       `json:"likes_count"`
	FavoritesCount int       `json:"favorites_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Author   UserResponse     `json:"author"`
	Category CategoryResponse `json:"category"`
	Tags     []TagResponse    `json:"tags"`
}
