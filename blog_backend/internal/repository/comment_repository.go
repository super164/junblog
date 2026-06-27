package repository

import (
	"errors"

	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
)

// commentRepository 评论仓储实现
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// FindByID 根据 ID 查找评论
func (r *commentRepository) FindByID(id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Preload("User").Preload("Article").First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// Create 创建评论
func (r *commentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

// Update 更新评论
func (r *commentRepository) Update(comment *entity.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Comment{}, id).Error
}

// ListByArticleID 根据文章 ID 获取评论列表
func (r *commentRepository) ListByArticleID(articleID uint, offset, limit int) ([]*entity.Comment, int64, error) {
	var comments []*entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).Where("article_id = ?", articleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").Preload("Article").Preload("Parent.User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments).Error

	return comments, total, err
}

// ListAdmin 后台分页获取评论列表
func (r *commentRepository) ListAdmin(offset, limit int, status string) ([]*entity.Comment, int64, error) {
	var comments []*entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments).Error

	return comments, total, err
}

// CountByArticleID 统计文章评论数
func (r *commentRepository) CountByArticleID(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Comment{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}
