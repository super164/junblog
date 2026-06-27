package repository

import (
	"blog_backend/internal/model/entity"
)

// TagRepository 标签仓储接口
type TagRepository interface {
	// FindByID 根据 ID 查找标签
	FindByID(id uint) (*entity.Tag, error)
	// FindByName 根据名称查找标签
	FindByName(name string) (*entity.Tag, error)
	// FindBySlug 根据 slug 查找标签
	FindBySlug(slug string) (*entity.Tag, error)
	// Create 创建标签
	Create(tag *entity.Tag) error
	// Update 更新标签
	Update(tag *entity.Tag) error
	// Delete 删除标签
	Delete(id uint) error
	// List 获取所有标签
	List() ([]*entity.Tag, error)
	// ListByCategoryID 根据分类ID获取标签列表
	ListByCategoryID(categoryID uint) ([]*entity.Tag, error)
	// ExistsByName 检查名称是否存在
	ExistsByName(name string) (bool, error)
	// ExistsByNameExcludeID 检查名称是否存在（排除指定ID）
	ExistsByNameExcludeID(name string, excludeID uint) (bool, error)
	// ExistsBySlug 检查 slug 是否存在
	ExistsBySlug(slug string) (bool, error)
	// ExistsBySlugExcludeID 检查 slug 是否存在（排除指定ID）
	ExistsBySlugExcludeID(slug string, excludeID uint) (bool, error)
}
