package repository

import (
	"blog_backend/internal/model/entity"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	// FindByID 根据 ID 查找文章
	FindByID(id uint) (*entity.Article, error)
	// FindBySlug 根据 slug 查找文章
	FindBySlug(slug string) (*entity.Article, error)
	// Create 创建文章
	Create(article *entity.Article) error
	// Update 更新文章
	Update(article *entity.Article) error
	// Delete 删除文章
	Delete(id uint) error
	// List 分页获取文章列表
	List(offset, limit int, categoryID uint, tagID uint, keyword string, sort string) ([]*entity.Article, int64, error)
	// ListAdmin 后台分页获取文章列表
	ListAdmin(offset, limit int, status string, keyword string) ([]*entity.Article, int64, error)
	// IncrementViews 增加浏览量
	IncrementViews(id uint) error
	// GetHotArticles 获取热门文章
	GetHotArticles(limit int) ([]*entity.Article, error)
	// GetRecentArticles 获取最新文章
	GetRecentArticles(limit int) ([]*entity.Article, error)
	// CountLikes 实时统计文章点赞数
	CountLikes(articleID uint) (int64, error)
	// CountFavorites 实时统计文章收藏数
	CountFavorites(articleID uint) (int64, error)
	// UpdateTags 更新文章标签关联
	UpdateTags(articleID uint, tagIDs []uint) error
}
