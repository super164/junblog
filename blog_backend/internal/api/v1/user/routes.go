package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/profile", ctrl.GetProfile)
		userGroup.PUT("/profile", ctrl.UpdateProfile)
		userGroup.PUT("/profile/password", ctrl.UpdatePassword)
	}
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	adminGroup := r.Group("/users")
	{
		adminGroup.GET("", ctrl.AdminGetUsers)
		adminGroup.PUT("/:id", ctrl.AdminUpdateUser)
		adminGroup.PATCH("/:id/password", ctrl.AdminResetPassword)
		adminGroup.PATCH("/:id/status", ctrl.AdminUpdateUserStatus)
	}
}
