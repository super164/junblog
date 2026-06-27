package repository

import (
	"errors"

	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
)

// categoryRepository 分类仓储实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindByID 根据 ID 查找分类
func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindByName 根据名称查找分类
func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug 根据 slug 查找分类
func (r *categoryRepository) FindBySlug(slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

// List 获取所有分类
func (r *categoryRepository) List() ([]*entity.Category, error) {
	var categories []*entity.Category
	err := r.db.Order("created_at ASC").Find(&categories).Error
	return categories, err
}

// GetTree 获取分类树
func (r *categoryRepository) GetTree() ([]*entity.Category, error) {
	var categories []*entity.Category
	err := r.db.Where("parent_id IS NULL").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at ASC").
		Find(&categories).Error
	return categories, err
}

// HasArticles 检查分类下是否有文章
func (r *categoryRepository) HasArticles(categoryID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查名称是否存在
func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsBySlug 检查 slug 是否存在
func (r *categoryRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}
