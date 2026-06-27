package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// CategoryService 分类服务接口
type CategoryService interface {
	// CreateCategory 创建分类
	CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error)
	// GetCategoryByID 根据 ID 获取分类
	GetCategoryByID(id uint) (*response.CategoryResponse, error)
	// GetAllCategories 获取所有分类
	GetAllCategories() ([]*response.CategoryResponse, error)
	// GetCategoryTree 获取分类树
	GetCategoryTree() ([]*response.CategoryTreeResponse, error)
	// UpdateCategory 更新分类
	UpdateCategory(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error)
	// DeleteCategory 删除分类
	DeleteCategory(id uint) error
}
