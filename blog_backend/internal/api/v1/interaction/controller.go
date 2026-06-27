package interaction

import (
	"blog_backend/internal/middleware"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 互动控制器
type Controller struct {
	interactionService service.InteractionService
}

// NewController 创建互动控制器
func NewController(interactionService service.InteractionService) *Controller {
	return &Controller{
		interactionService: interactionService,
	}
}

// GetInteractionStatus 获取互动状态
func (ctrl *Controller) GetInteractionStatus(c *gin.Context) {
	articleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	isLiked, err := ctrl.interactionService.IsLiked(userID, uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	isFavorited, err := ctrl.interactionService.IsFavorited(userID, uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	likeCount, err := ctrl.interactionService.GetLikeCount(uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	favoriteCount, err := ctrl.interactionService.GetFavoriteCount(uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"is_liked":      isLiked,
		"is_favorited":  isFavorited,
		"like_count":    likeCount,
		"favorite_count": favoriteCount,
	})
}

// ToggleLike 切换点赞状态
func (ctrl *Controller) ToggleLike(c *gin.Context) {
	articleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	isLiked, err := ctrl.interactionService.IsLiked(userID, uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if isLiked {
		err = ctrl.interactionService.Unlike(userID, uint(articleID))
	} else {
		err = ctrl.interactionService.Like(userID, uint(articleID))
	}

	if err != nil {
		response.HandleError(c, err)
		return
	}

	likeCount, _ := ctrl.interactionService.GetLikeCount(uint(articleID))
	response.Success(c, gin.H{
		"is_liked":    !isLiked,
		"likes_count": likeCount,
	})
}

// ToggleFavorite 切换收藏状态
func (ctrl *Controller) ToggleFavorite(c *gin.Context) {
	articleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	isFavorited, err := ctrl.interactionService.IsFavorited(userID, uint(articleID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if isFavorited {
		err = ctrl.interactionService.Unfavorite(userID, uint(articleID))
	} else {
		err = ctrl.interactionService.Favorite(userID, uint(articleID))
	}

	if err != nil {
		response.HandleError(c, err)
		return
	}

	favCount, _ := ctrl.interactionService.GetFavoriteCount(uint(articleID))
	response.Success(c, gin.H{
		"is_favorited":     !isFavorited,
		"favorites_count": favCount,
	})
}

// GetUserFavorites 获取用户收藏列表
func (ctrl *Controller) GetUserFavorites(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	favorites, total, err := ctrl.interactionService.GetUserFavorites(userID, page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  favorites,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetUserLikes 获取用户点赞列表
func (ctrl *Controller) GetUserLikes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	likes, total, err := ctrl.interactionService.GetUserLikes(userID, page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  likes,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
