package request

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"required,max=50"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,max=50"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}
