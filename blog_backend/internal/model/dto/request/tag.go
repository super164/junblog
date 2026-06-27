package request

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name       string `json:"name" binding:"required,max=50"`
	Slug       string `json:"slug" binding:"required,max=50"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name       string `json:"name" binding:"omitempty,max=50"`
	Slug       string `json:"slug" binding:"omitempty,max=50"`
	CategoryID uint   `json:"category_id" binding:"omitempty"`
}
