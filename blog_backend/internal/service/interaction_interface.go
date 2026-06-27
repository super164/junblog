package service

import (
	"blog_backend/internal/model/dto/response"
)

// InteractionService 互动服务接口
type InteractionService interface {
	// Like 点赞
	Like(userID, articleID uint) error
	// Unlike 取消点赞
	Unlike(userID, articleID uint) error
	// IsLiked 检查是否已点赞
	IsLiked(userID, articleID uint) (bool, error)
	// GetLikeCount 获取文章点赞数
	GetLikeCount(articleID uint) (int64, error)

	// Favorite 收藏
	Favorite(userID, articleID uint) error
	// Unfavorite 取消收藏
	Unfavorite(userID, articleID uint) error
	// IsFavorited 检查是否已收藏
	IsFavorited(userID, articleID uint) (bool, error)
	// GetFavoriteCount 获取文章收藏数
	GetFavoriteCount(articleID uint) (int64, error)

	// GetUserFavorites 获取用户收藏列表
	GetUserFavorites(userID uint, page, size int) ([]*response.FavoriteResponse, int64, error)
	// GetUserLikes 获取用户点赞列表
	GetUserLikes(userID uint, page, size int) ([]*response.LikeResponse, int64, error)
}
