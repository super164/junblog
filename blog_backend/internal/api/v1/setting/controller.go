package setting

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 设置控制器
type Controller struct {
	settingService service.SystemSettingService
}

// NewController 创建设置控制器
func NewController(settingService service.SystemSettingService) *Controller {
	return &Controller{
		settingService: settingService,
	}
}

// GetAboutPage 获取关于页面数据（公开接口）
func (ctrl *Controller) GetAboutPage(c *gin.Context) {
	data, err := ctrl.settingService.GetAboutPage()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, data)
}

// AdminGetSettings 管理员获取所有设置
func (ctrl *Controller) AdminGetSettings(c *gin.Context) {
	settings, err := ctrl.settingService.GetAllSettings()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, settings)
}

// AdminGetSetting 管理员获取单个设置
func (ctrl *Controller) AdminGetSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "设置键不能为空")
		return
	}

	value, err := ctrl.settingService.GetSetting(key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"key":   key,
		"value": value,
	})
}

// AdminUpdateSetting 管理员更新单个设置
func (ctrl *Controller) AdminUpdateSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "设置键不能为空")
		return
	}

	var req request.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.settingService.SetSetting(key, req.Value); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// AdminUpdateAboutPage 管理员更新关于页面
func (ctrl *Controller) AdminUpdateAboutPage(c *gin.Context) {
	var req request.UpdateAboutPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.settingService.UpdateAboutPage(&req); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
