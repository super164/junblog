package stats

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册统计路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/stats", ctrl.GetStats)
}
