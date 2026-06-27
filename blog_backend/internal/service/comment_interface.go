package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// CommentService 评论服务接口
type CommentService interface {
	// CreateComment 创建评论
	CreateComment(userID uint, req *request.CreateCommentRequest) (*response.CommentResponse, error)
	// GetCommentsByArticleID 根据文章 ID 获取评论列表
	GetCommentsByArticleID(articleID uint, page, size int) ([]*response.CommentResponse, int64, error)
	// UpdateComment 更新评论
	UpdateComment(id uint, req *request.UpdateCommentRequest) (*response.CommentResponse, error)
	// DeleteComment 删除评论
	DeleteComment(id uint) error
	// GetAdminCommentList 后台获取评论列表
	GetAdminCommentList(page, size int, status string) ([]*response.CommentResponse, int64, error)
}
