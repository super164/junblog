package category

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes 注册公开路由
func (ctrl *Controller) RegisterPublicRoutes(r *gin.RouterGroup) {
	categoryGroup := r.Group("/categories")
	{
		categoryGroup.GET("", ctrl.GetCategories)
		categoryGroup.GET("/tree", ctrl.GetCategoryTree)
	}
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	adminCategoryGroup := r.Group("/categories")
	{
		adminCategoryGroup.GET("", ctrl.AdminGetCategories)
		adminCategoryGroup.POST("", ctrl.AdminCreateCategory)
		adminCategoryGroup.PUT("/:id", ctrl.AdminUpdateCategory)
		adminCategoryGroup.DELETE("/:id", ctrl.AdminDeleteCategory)
	}
}
