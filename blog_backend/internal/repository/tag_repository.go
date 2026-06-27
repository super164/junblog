package repository

import (
	"errors"

	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
)

// tagRepository 标签仓储实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// FindByID 根据 ID 查找标签
func (r *tagRepository) FindByID(id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindByName 根据名称查找标签
func (r *tagRepository) FindByName(name string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindBySlug 根据 slug 查找标签
func (r *tagRepository) FindBySlug(slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签
func (r *tagRepository) Create(tag *entity.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *tagRepository) Update(tag *entity.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *tagRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 清理关联
		if err := tx.Exec("DELETE FROM article_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Tag{}, id).Error
	})
}

// List 获取所有标签
func (r *tagRepository) List() ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.Preload("Category").Order("created_at ASC").Find(&tags).Error
	return tags, err
}

// ListByCategoryID 根据分类ID获取标签列表
func (r *tagRepository) ListByCategoryID(categoryID uint) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.Where("category_id = ?", categoryID).Preload("Category").Order("created_at ASC").Find(&tags).Error
	return tags, err
}

// ExistsByName 检查名称是否存在
func (r *tagRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Tag{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcludeID 检查名称是否存在（排除指定ID）
func (r *tagRepository) ExistsByNameExcludeID(name string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Tag{}).Where("name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// ExistsBySlug 检查 slug 是否存在
func (r *tagRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Tag{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// ExistsBySlugExcludeID 检查 slug 是否存在（排除指定ID）
func (r *tagRepository) ExistsBySlugExcludeID(slug string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Tag{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count > 0, err
}
