package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// TagService 标签服务接口
type TagService interface {
	// CreateTag 创建标签
	CreateTag(req *request.CreateTagRequest) (*response.TagResponse, error)
	// GetAllTags 获取所有标签
	GetAllTags() ([]*response.TagResponse, error)
	// GetTagsByCategoryID 根据分类ID获取标签列表
	GetTagsByCategoryID(categoryID uint) ([]*response.TagResponse, error)
	// UpdateTag 更新标签
	UpdateTag(id uint, req *request.UpdateTagRequest) (*response.TagResponse, error)
	// DeleteTag 删除标签
	DeleteTag(id uint) error
}
