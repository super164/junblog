package interaction

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册互动路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	// 使用 /interactions 路径避免与 /articles/:slug 冲突
	interactionGroup := r.Group("/interactions")
	{
		interactionGroup.GET("/articles/:id", ctrl.GetInteractionStatus)
		interactionGroup.POST("/articles/:id/like", ctrl.ToggleLike)
		interactionGroup.POST("/articles/:id/favorite", ctrl.ToggleFavorite)
	}

	userGroup := r.Group("/user")
	{
		userGroup.GET("/favorites", ctrl.GetUserFavorites)
		userGroup.GET("/likes", ctrl.GetUserLikes)
	}
}
