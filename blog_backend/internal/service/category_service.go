package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
)

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	// 检查名称是否存在
	exists, err := s.categoryRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, bizerrors.New(bizerrors.CodeConflict, "分类名称已存在")
	}

	// 检查 slug 是否存在
	exists, err = s.categoryRepo.ExistsBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, bizerrors.New(bizerrors.CodeConflict, "分类别名已存在")
	}

	// 处理 parent_id：0 视为无父分类
	var parentID *uint
	if req.ParentID != nil && *req.ParentID != 0 {
		parentID = req.ParentID
	}

	category := &entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    parentID,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetCategoryByID 根据 ID 获取分类
func (s *categoryService) GetCategoryByID(id uint) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, bizerrors.ErrNotFound
	}

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetAllCategories 获取所有分类
func (s *categoryService) GetAllCategories() ([]*response.CategoryResponse, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	result := make([]*response.CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = &response.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			ParentID:    c.ParentID,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
	}
	return result, nil
}

// GetCategoryTree 获取分类树
func (s *categoryService) GetCategoryTree() ([]*response.CategoryTreeResponse, error) {
	categories, err := s.categoryRepo.GetTree()
	if err != nil {
		return nil, err
	}

	return s.toCategoryTreeResponse(categories), nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, bizerrors.ErrNotFound
	}

	// 直接更新字段，不需要检查名称/Slug 是否存在（因为唯一索引会自动校验）
	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Slug != "" {
		category.Slug = req.Slug
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	if req.ParentID != nil {
		// 处理 parent_id：0 视为无父分类
		if *req.ParentID == 0 {
			category.ParentID = nil
		} else {
			// 不能将父分类设为自己
			if *req.ParentID == id {
				return nil, bizerrors.New(bizerrors.CodeBadRequest, "不能将分类设为自己的子分类")
			}
			category.ParentID = req.ParentID
		}
	} else {
		// 前端未传 parent_id 时保持不变
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return bizerrors.ErrNotFound
	}

	// 检查分类下是否有文章
	hasArticles, err := s.categoryRepo.HasArticles(id)
	if err != nil {
		return err
	}
	if hasArticles {
		return bizerrors.New(bizerrors.CodeConflict, "分类下还有文章，无法删除")
	}

	return s.categoryRepo.Delete(id)
}

// toCategoryTreeResponse 转换为分类树响应
func (s *categoryService) toCategoryTreeResponse(categories []*entity.Category) []*response.CategoryTreeResponse {
	result := make([]*response.CategoryTreeResponse, len(categories))
	for i, c := range categories {
		// 将 []entity.Category 转换为 []*entity.Category
		children := make([]*entity.Category, len(c.Children))
		for j := range c.Children {
			children[j] = &c.Children[j]
		}

		result[i] = &response.CategoryTreeResponse{
			CategoryResponse: response.CategoryResponse{
				ID:          c.ID,
				Name:        c.Name,
				Slug:        c.Slug,
				Description: c.Description,
				ParentID:    c.ParentID,
				CreatedAt:   c.CreatedAt,
				UpdatedAt:   c.UpdatedAt,
			},
			Children: s.toCategoryTreeResponse(children),
		}
	}
	return result
}
