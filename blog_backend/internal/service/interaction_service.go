package service

import (
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
)

// interactionService 互动服务实现
type interactionService struct {
	interactionRepo repository.InteractionRepository
	articleRepo     repository.ArticleRepository
}

// NewInteractionService 创建互动服务
func NewInteractionService(interactionRepo repository.InteractionRepository, articleRepo repository.ArticleRepository) InteractionService {
	return &interactionService{
		interactionRepo: interactionRepo,
		articleRepo:     articleRepo,
	}
}

// Like 点赞
func (s *interactionService) Like(userID, articleID uint) error {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return bizerrors.ErrNotFound
	}

	// 检查是否已点赞
	exists, err := s.interactionRepo.IsLiked(userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.New(bizerrors.CodeConflict, "已经点赞过了")
	}

	return s.interactionRepo.Like(userID, articleID)
}

// Unlike 取消点赞
func (s *interactionService) Unlike(userID, articleID uint) error {
	return s.interactionRepo.Unlike(userID, articleID)
}

// IsLiked 检查是否已点赞
func (s *interactionService) IsLiked(userID, articleID uint) (bool, error) {
	return s.interactionRepo.IsLiked(userID, articleID)
}

// GetLikeCount 获取文章点赞数
func (s *interactionService) GetLikeCount(articleID uint) (int64, error) {
	return s.interactionRepo.GetLikeCount(articleID)
}

// Favorite 收藏
func (s *interactionService) Favorite(userID, articleID uint) error {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return bizerrors.ErrNotFound
	}

	// 检查是否已收藏
	exists, err := s.interactionRepo.IsFavorited(userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.New(bizerrors.CodeConflict, "已经收藏过了")
	}

	return s.interactionRepo.Favorite(userID, articleID)
}

// Unfavorite 取消收藏
func (s *interactionService) Unfavorite(userID, articleID uint) error {
	return s.interactionRepo.Unfavorite(userID, articleID)
}

// IsFavorited 检查是否已收藏
func (s *interactionService) IsFavorited(userID, articleID uint) (bool, error) {
	return s.interactionRepo.IsFavorited(userID, articleID)
}

// GetFavoriteCount 获取文章收藏数
func (s *interactionService) GetFavoriteCount(articleID uint) (int64, error) {
	return s.interactionRepo.GetFavoriteCount(articleID)
}

// GetUserFavorites 获取用户收藏列表
func (s *interactionService) GetUserFavorites(userID uint, page, size int) ([]*response.FavoriteResponse, int64, error) {
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

	favorites, total, err := s.interactionRepo.GetUserFavorites(userID, offset, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*response.FavoriteResponse, len(favorites))
	for i, f := range favorites {
		result[i] = &response.FavoriteResponse{
			ID:        f.ID,
			UserID:    f.UserID,
			ArticleID: f.ArticleID,
			CreatedAt: f.CreatedAt,
		}
	}

	return result, total, nil
}

// GetUserLikes 获取用户点赞列表
func (s *interactionService) GetUserLikes(userID uint, page, size int) ([]*response.LikeResponse, int64, error) {
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

	likes, total, err := s.interactionRepo.GetUserLikes(userID, offset, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*response.LikeResponse, len(likes))
	for i, l := range likes {
		result[i] = &response.LikeResponse{
			ID:        l.ID,
			UserID:    l.UserID,
			ArticleID: l.ArticleID,
			CreatedAt: l.CreatedAt,
		}
	}

	return result, total, nil
}
