package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// SystemSettingService 系统设置服务接口
type SystemSettingService interface {
	// GetAboutPage 获取关于页面数据
	GetAboutPage() (*response.AboutPageResponse, error)
	// UpdateAboutPage 更新关于页面数据
	UpdateAboutPage(req *request.UpdateAboutPageRequest) error
	// GetSetting 获取单个设置
	GetSetting(key string) (string, error)
	// SetSetting 设置单个值
	SetSetting(key, value string) error
	// GetAllSettings 获取所有设置
	GetAllSettings() ([]*response.SettingResponse, error)
}
