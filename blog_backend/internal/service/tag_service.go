package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
)

// tagService 标签服务实现
type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

// CreateTag 创建标签
func (s *tagService) CreateTag(req *request.CreateTagRequest) (*response.TagResponse, error) {
	// 检查名称是否存在
	exists, err := s.tagRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, bizerrors.New(bizerrors.CodeConflict, "标签名称已存在")
	}

	// 检查 slug 是否存在
	exists, err = s.tagRepo.ExistsBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, bizerrors.New(bizerrors.CodeConflict, "标签别名已存在")
	}

	tag := &entity.Tag{
		Name:       req.Name,
		Slug:       req.Slug,
		CategoryID: req.CategoryID,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return &response.TagResponse{
		ID:         tag.ID,
		Name:       tag.Name,
		Slug:       tag.Slug,
		CategoryID: tag.CategoryID,
	}, nil
}

// GetAllTags 获取所有标签
func (s *tagService) GetAllTags() ([]*response.TagResponse, error) {
	tags, err := s.tagRepo.List()
	if err != nil {
		return nil, err
	}

	result := make([]*response.TagResponse, len(tags))
	for i, t := range tags {
		resp := &response.TagResponse{
			ID:         t.ID,
			Name:       t.Name,
			Slug:       t.Slug,
			CategoryID: t.CategoryID,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
		if t.Category != nil {
			resp.Category = &response.CategoryResponse{
				ID:   t.Category.ID,
				Name: t.Category.Name,
				Slug: t.Category.Slug,
			}
		}
		result[i] = resp
	}
	return result, nil
}

// GetTagsByCategoryID 根据分类ID获取标签列表
func (s *tagService) GetTagsByCategoryID(categoryID uint) ([]*response.TagResponse, error) {
	tags, err := s.tagRepo.ListByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	result := make([]*response.TagResponse, len(tags))
	for i, t := range tags {
		resp := &response.TagResponse{
			ID:         t.ID,
			Name:       t.Name,
			Slug:       t.Slug,
			CategoryID: t.CategoryID,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
		if t.Category != nil {
			resp.Category = &response.CategoryResponse{
				ID:   t.Category.ID,
				Name: t.Category.Name,
				Slug: t.Category.Slug,
			}
		}
		result[i] = resp
	}
	return result, nil
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(id uint, req *request.UpdateTagRequest) (*response.TagResponse, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, bizerrors.ErrNotFound
	}

	if req.Name != "" {
		// 检查名称是否存在（排除自身）
		exists, err := s.tagRepo.ExistsByNameExcludeID(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, bizerrors.New(bizerrors.CodeConflict, "标签名称已存在")
		}
		tag.Name = req.Name
	}

	if req.Slug != "" {
		// 检查 slug 是否存在（排除自身）
		exists, err := s.tagRepo.ExistsBySlugExcludeID(req.Slug, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, bizerrors.New(bizerrors.CodeConflict, "标签别名已存在")
		}
		tag.Slug = req.Slug
	}

	// 直接更新 category_id，不管值是多少
	tag.CategoryID = req.CategoryID

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}

	return &response.TagResponse{
		ID:         tag.ID,
		Name:       tag.Name,
		Slug:       tag.Slug,
		CategoryID: tag.CategoryID,
	}, nil
}

// DeleteTag 删除标签
func (s *tagService) DeleteTag(id uint) error {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return err
	}
	if tag == nil {
		return bizerrors.ErrNotFound
	}

	return s.tagRepo.Delete(id)
}
