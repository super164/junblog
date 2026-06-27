package service

import "blog_backend/internal/model/dto/response"

// StatsService 统计服务接口
type StatsService interface {
	// GetStats 获取后台统计数据
	GetStats() (*response.AdminStatsResponse, error)
}
