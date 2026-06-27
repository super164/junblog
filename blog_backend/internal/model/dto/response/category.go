package response

import "time"

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryTreeResponse 分类树响应
type CategoryTreeResponse struct {
	CategoryResponse
	Children []*CategoryTreeResponse `json:"children"`
}
