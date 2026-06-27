package admin

import (
	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.POST("/upload", ctrl.UploadFile)
}
