package response

import "time"

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User       UserResponse       `json:"user"`
	Article    ArticleListItem    `json:"article"`
	Parent     *CommentResponse   `json:"parent,omitempty"`
	Replies    []CommentResponse  `json:"replies,omitempty"`
}
