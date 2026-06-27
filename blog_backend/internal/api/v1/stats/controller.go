package stats

import (
	"blog_backend/internal/service"
	"blog_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 统计控制器
type Controller struct {
	statsService service.StatsService
}

// NewController 创建统计控制器
func NewController(statsService service.StatsService) *Controller {
	return &Controller{
		statsService: statsService,
	}
}

// GetStats 获取后台统计数据
func (ctrl *Controller) GetStats(c *gin.Context) {
	stats, err := ctrl.statsService.GetStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}
