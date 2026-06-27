package repository

import (
	"blog_backend/internal/model/entity"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	// FindByID 根据 ID 查找分类
	FindByID(id uint) (*entity.Category, error)
	// FindByName 根据名称查找分类
	FindByName(name string) (*entity.Category, error)
	// FindBySlug 根据 slug 查找分类
	FindBySlug(slug string) (*entity.Category, error)
	// Create 创建分类
	Create(category *entity.Category) error
	// Update 更新分类
	Update(category *entity.Category) error
	// Delete 删除分类
	Delete(id uint) error
	// List 获取所有分类
	List() ([]*entity.Category, error)
	// GetTree 获取分类树
	GetTree() ([]*entity.Category, error)
	// HasArticles 检查分类下是否有文章
	HasArticles(categoryID uint) (bool, error)
	// ExistsByName 检查名称是否存在
	ExistsByName(name string) (bool, error)
	// ExistsBySlug 检查 slug 是否存在
	ExistsBySlug(slug string) (bool, error)
}
