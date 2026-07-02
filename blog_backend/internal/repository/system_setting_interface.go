package repository

import "blog_backend/internal/model/entity"

// SystemSettingRepository 系统设置仓储接口
type SystemSettingRepository interface {
	FindByKey(key string) (*entity.SystemSetting, error)
	FindByKeys(keys []string) ([]*entity.SystemSetting, error)
	FindAll() ([]*entity.SystemSetting, error)
	Set(key, value string) error
	Delete(key string) error
}
