package repository

import (
	"errors"

	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
)

// articleRepository 文章仓储实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// FindByID 根据 ID 查找文章
func (r *articleRepository) FindByID(id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").
		First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindBySlug 根据 slug 查找文章
func (r *articleRepository) FindBySlug(slug string) (*entity.Article, error) {
	var article entity.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// Create 创建文章
func (r *articleRepository) Create(article *entity.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *articleRepository) Update(article *entity.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *articleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清理标签关联
		if err := tx.Exec("DELETE FROM article_tags WHERE article_id = ?", id).Error; err != nil {
			return err
		}
		// 清理评论
		if err := tx.Where("article_id = ?", id).Delete(&entity.Comment{}).Error; err != nil {
			return err
		}
		// 清理点赞
		if err := tx.Where("article_id = ?", id).Delete(&entity.Like{}).Error; err != nil {
			return err
		}
		// 清理收藏
		if err := tx.Where("article_id = ?", id).Delete(&entity.Favorite{}).Error; err != nil {
			return err
		}
		// 删除文章
		return tx.Delete(&entity.Article{}, id).Error
	})
}

// List 分页获取文章列表
func (r *articleRepository) List(offset, limit int, categoryID uint, tagID uint, keyword string, sort string) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).Where("status = ?", "published")

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", tagID)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "created_at DESC"
	if sort == "hottest" {
		orderClause = "views_count DESC"
	}

	err := query.Preload("Author").Preload("Category").Preload("Tags").
		Order(orderClause).
		Offset(offset).Limit(limit).
		Find(&articles).Error

	return articles, total, err
}

// ListAdmin 后台分页获取文章列表
func (r *articleRepository) ListAdmin(offset, limit int, status string, keyword string) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Author").Preload("Category").Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error

	return articles, total, err
}

// IncrementViews 增加浏览量
func (r *articleRepository) IncrementViews(id uint) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

// GetHotArticles 获取热门文章
func (r *articleRepository) GetHotArticles(limit int) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.Where("status = ?", "published").
		Preload("Author").Preload("Category").Preload("Tags").
		Order("views_count DESC").
		Limit(limit).Find(&articles).Error
	return articles, err
}

// GetRecentArticles 获取最新文章
func (r *articleRepository) GetRecentArticles(limit int) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.Where("status = ?", "published").
		Preload("Author").Preload("Category").Preload("Tags").
		Order("created_at DESC").
		Limit(limit).Find(&articles).Error
	return articles, err
}

// CountLikes 实时统计文章点赞数
func (r *articleRepository) CountLikes(articleID uint) (int64, error) {
	var count int64
	err := r.db.Table("likes").Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// CountFavorites 实时统计文章收藏数
func (r *articleRepository) CountFavorites(articleID uint) (int64, error) {
	var count int64
	err := r.db.Table("favorites").Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// UpdateTags 更新文章标签关联
func (r *articleRepository) UpdateTags(articleID uint, tagIDs []uint) error {
	article := &entity.Article{BaseEntity: entity.BaseEntity{ID: articleID}}
	// 先清除旧关联
	if err := r.db.Model(article).Association("Tags").Clear(); err != nil {
		return err
	}
	// 再添加新关联
	if len(tagIDs) > 0 {
		var tags []entity.Tag
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
		if err := r.db.Model(article).Association("Tags").Replace(tags); err != nil {
			return err
		}
	}
	return nil
}
