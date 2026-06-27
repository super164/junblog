package repository

import (
	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
)

// interactionRepository 互动仓储实现
type interactionRepository struct {
	db *gorm.DB
}

// NewInteractionRepository 创建互动仓储
func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

// Like 点赞
func (r *interactionRepository) Like(userID, articleID uint) error {
	like := &entity.Like{
		UserID:    userID,
		ArticleID: articleID,
	}
	return r.db.Create(like).Error
}

// Unlike 取消点赞
func (r *interactionRepository) Unlike(userID, articleID uint) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&entity.Like{}).Error
}

// IsLiked 检查是否已点赞
func (r *interactionRepository) IsLiked(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

// GetLikeCount 获取文章点赞数
func (r *interactionRepository) GetLikeCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}

// Favorite 收藏
func (r *interactionRepository) Favorite(userID, articleID uint) error {
	favorite := &entity.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}
	return r.db.Create(favorite).Error
}

// Unfavorite 取消收藏
func (r *interactionRepository) Unfavorite(userID, articleID uint) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&entity.Favorite{}).Error
}

// IsFavorited 检查是否已收藏
func (r *interactionRepository) IsFavorited(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

// GetFavoriteCount 获取文章收藏数
func (r *interactionRepository) GetFavoriteCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Favorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}

// GetUserFavorites 获取用户收藏列表
func (r *interactionRepository) GetUserFavorites(userID uint, offset, limit int) ([]*entity.Favorite, int64, error) {
	var favorites []*entity.Favorite
	var total int64

	query := r.db.Model(&entity.Favorite{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Article").Preload("Article.Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&favorites).Error

	return favorites, total, err
}

// GetUserLikes 获取用户点赞列表
func (r *interactionRepository) GetUserLikes(userID uint, offset, limit int) ([]*entity.Like, int64, error) {
	var likes []*entity.Like
	var total int64

	query := r.db.Model(&entity.Like{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Article").Preload("Article.Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&likes).Error

	return likes, total, err
}
