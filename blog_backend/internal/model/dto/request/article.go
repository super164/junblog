package request

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	Cover      string `json:"cover"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published privacy"`
	Slug       string `json:"slug" binding:"omitempty,max=255"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string `json:"title" binding:"omitempty,max=255"`
	Content    string `json:"content" binding:"omitempty"`
	Cover      string `json:"cover" binding:"omitempty"`
	CategoryID uint   `json:"category_id" binding:"omitempty"`
	TagIDs     []uint `json:"tag_ids" binding:"omitempty"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published privacy"`
	Slug       string `json:"slug" binding:"omitempty,max=255"`
}
