package comment

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes 注册公开路由
func (ctrl *Controller) RegisterPublicRoutes(r *gin.RouterGroup) {
	commentGroup := r.Group("/comments")
	{
		commentGroup.GET("", ctrl.GetComments)
	}
}

// RegisterAuthRoutes 注册需要认证的路由
func (ctrl *Controller) RegisterAuthRoutes(r *gin.RouterGroup) {
	commentGroup := r.Group("/comments")
	{
		commentGroup.POST("", ctrl.CreateComment)
	}
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	adminCommentGroup := r.Group("/comments")
	{
		adminCommentGroup.GET("", ctrl.AdminGetComments)
		adminCommentGroup.PATCH("/:id/status", ctrl.AdminUpdateCommentStatus)
		adminCommentGroup.DELETE("/:id", ctrl.AdminDeleteComment)
	}
}
