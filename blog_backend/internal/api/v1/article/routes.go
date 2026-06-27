package article

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes 注册公开路由
func (ctrl *Controller) RegisterPublicRoutes(r *gin.RouterGroup) {
	articleGroup := r.Group("/articles")
	{
		articleGroup.GET("", ctrl.GetArticles)
		articleGroup.GET("/hot", ctrl.GetHotArticles)
		articleGroup.GET("/recent", ctrl.GetRecentArticles)
		articleGroup.GET("/:id", ctrl.GetArticleByID)
	}
}

// RegisterAuthRoutes 注册需要认证的路由
func (ctrl *Controller) RegisterAuthRoutes(r *gin.RouterGroup) {
	// 目前不需要认证的文章路由
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	adminArticleGroup := r.Group("/articles")
	{
		adminArticleGroup.GET("", ctrl.AdminGetArticles)
		adminArticleGroup.POST("", ctrl.AdminCreateArticle)
		adminArticleGroup.GET("/:id", ctrl.AdminGetArticleByID)
		adminArticleGroup.PUT("/:id", ctrl.AdminUpdateArticle)
		adminArticleGroup.DELETE("/:id", ctrl.AdminDeleteArticle)
		adminArticleGroup.PATCH("/:id/status", ctrl.AdminUpdateArticleStatus)
	}
}
