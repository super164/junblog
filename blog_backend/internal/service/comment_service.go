package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
)

// commentService 评论服务实现
type commentService struct {
	commentRepo  repository.CommentRepository
	articleRepo  repository.ArticleRepository
	userRepo     repository.UserRepository
}

// NewCommentService 创建评论服务
func NewCommentService(commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository, userRepo repository.UserRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(userID uint, req *request.CreateCommentRequest) (*response.CommentResponse, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(req.ArticleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.ErrNotFound
	}

	// 检查父评论是否存在
	if req.ParentID != nil {
		parentComment, err := s.commentRepo.FindByID(*req.ParentID)
		if err != nil {
			return nil, err
		}
		if parentComment == nil {
			return nil, bizerrors.ErrNotFound
		}
	}

	comment := &entity.Comment{
		Content:   req.Content,
		UserID:    userID,
		ArticleID: req.ArticleID,
		ParentID:  req.ParentID,
		Status:    "approved",
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 重新查询获取完整关联
	created, err := s.commentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	return s.toCommentResponse(created), nil
}

// GetCommentsByArticleID 根据文章 ID 获取评论列表（树形结构）
func (s *commentService) GetCommentsByArticleID(articleID uint, page, size int) ([]*response.CommentResponse, int64, error) {
	// 参数标准化
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	// 获取该文章的所有评论（包括子评论）
	allComments, total, err := s.commentRepo.ListByArticleID(articleID, offset, size)
	if err != nil {
		return nil, 0, err
	}

	// 构建树形结构：先将所有评论转为 response，并用 map 索引
	commentMap := make(map[uint]*response.CommentResponse)
	for _, c := range allComments {
		commentMap[c.ID] = s.toCommentResponse(c)
	}

	// 分离顶级评论和子评论，构建树
	var topComments []*response.CommentResponse
	for _, c := range allComments {
		resp := commentMap[c.ID]
		if c.ParentID == nil {
			// 顶级评论
			topComments = append(topComments, resp)
		} else {
			// 子评论，挂到父评论下
			if parent, ok := commentMap[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, *resp)
			}
		}
	}

	return topComments, total, nil
}

// UpdateComment 更新评论
func (s *commentService) UpdateComment(id uint, req *request.UpdateCommentRequest) (*response.CommentResponse, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, bizerrors.ErrNotFound
	}

	if req.Content != "" {
		comment.Content = req.Content
	}
	if req.Status != "" {
		comment.Status = req.Status
	}

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment), nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(id uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return bizerrors.ErrNotFound
	}

	return s.commentRepo.Delete(id)
}

// GetAdminCommentList 后台获取评论列表
func (s *commentService) GetAdminCommentList(page, size int, status string) ([]*response.CommentResponse, int64, error) {
	// 参数标准化
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	comments, total, err := s.commentRepo.ListAdmin(offset, size, status)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*response.CommentResponse, len(comments))
	for i, c := range comments {
		result[i] = s.toCommentResponse(c)
	}

	return result, total, nil
}

// toCommentResponse 转换为评论响应
func (s *commentService) toCommentResponse(c *entity.Comment) *response.CommentResponse {
	resp := &response.CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		User: response.UserResponse{
			ID:       c.User.ID,
			Username: c.User.Username,
			Avatar:   c.User.Avatar,
		},
		Article: response.ArticleListItem{
			ID:    c.Article.ID,
			Title: c.Article.Title,
			Slug:  c.Article.Slug,
		},
	}

	if c.Parent != nil {
		resp.Parent = s.toCommentResponse(c.Parent)
	}

	return resp
}
