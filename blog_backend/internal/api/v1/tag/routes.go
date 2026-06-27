package tag

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes 注册公开路由
func (ctrl *Controller) RegisterPublicRoutes(r *gin.RouterGroup) {
	tagGroup := r.Group("/tags")
	{
		tagGroup.GET("", ctrl.GetTags)
		tagGroup.GET("/category/:category_id", ctrl.GetTagsByCategoryID)
	}
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	adminTagGroup := r.Group("/tags")
	{
		adminTagGroup.GET("", ctrl.AdminGetTags)
		adminTagGroup.GET("/category/:category_id", ctrl.AdminGetTagsByCategoryID)
		adminTagGroup.POST("", ctrl.AdminCreateTag)
		adminTagGroup.PUT("/:id", ctrl.AdminUpdateTag)
		adminTagGroup.DELETE("/:id", ctrl.AdminDeleteTag)
	}
}
