package setting

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes 注册公开路由
func (ctrl *Controller) RegisterPublicRoutes(rg *gin.RouterGroup) {
	about := rg.Group("/about")
	{
		about.GET("", ctrl.GetAboutPage)
	}
}

// RegisterAdminRoutes 注册管理员路由
func (ctrl *Controller) RegisterAdminRoutes(rg *gin.RouterGroup) {
	settings := rg.Group("/settings")
	{
		settings.GET("", ctrl.AdminGetSettings)
		settings.GET("/:key", ctrl.AdminGetSetting)
		settings.PUT("/:key", ctrl.AdminUpdateSetting)

		// 关于页面专用接口
		settings.PUT("/about", ctrl.AdminUpdateAboutPage)
	}
}
