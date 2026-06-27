package response

import "time"

// TagResponse 标签响应
type TagResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	CategoryID uint      `json:"category_id"`
	Category   *CategoryResponse `json:"category,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
