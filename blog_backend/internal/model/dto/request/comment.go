package request

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"omitempty,oneof=approved pending rejected"`
}
