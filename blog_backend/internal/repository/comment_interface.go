package repository

import (
	"blog_backend/internal/model/entity"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	// FindByID 根据 ID 查找评论
	FindByID(id uint) (*entity.Comment, error)
	// Create 创建评论
	Create(comment *entity.Comment) error
	// Update 更新评论
	Update(comment *entity.Comment) error
	// Delete 删除评论
	Delete(id uint) error
	// ListByArticleID 根据文章 ID 获取评论列表
	ListByArticleID(articleID uint, offset, limit int) ([]*entity.Comment, int64, error)
	// ListAdmin 后台分页获取评论列表
	ListAdmin(offset, limit int, status string) ([]*entity.Comment, int64, error)
	// CountByArticleID 统计文章评论数
	CountByArticleID(articleID uint) (int64, error)
}
