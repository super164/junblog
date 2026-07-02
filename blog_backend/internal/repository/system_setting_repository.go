package repository

import (
	"blog_backend/internal/model/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// systemSettingRepository 系统设置仓储实现
type systemSettingRepository struct {
	db *gorm.DB
}

// NewSystemSettingRepository 创建系统设置仓储
func NewSystemSettingRepository(db *gorm.DB) SystemSettingRepository {
	return &systemSettingRepository{db: db}
}

// FindByKey 根据键查找设置
func (r *systemSettingRepository) FindByKey(key string) (*entity.SystemSetting, error) {
	var setting entity.SystemSetting
	err := r.db.Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

// FindByKeys 根据多个键查找设置
func (r *systemSettingRepository) FindByKeys(keys []string) ([]*entity.SystemSetting, error) {
	var settings []*entity.SystemSetting
	err := r.db.Where("`key` IN ?", keys).Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// FindAll 查找所有设置
func (r *systemSettingRepository) FindAll() ([]*entity.SystemSetting, error) {
	var settings []*entity.SystemSetting
	err := r.db.Order("`key` ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// Set 设置键值对（upsert）
func (r *systemSettingRepository) Set(key, value string) error {
	setting := entity.SystemSetting{
		Key:   key,
		Value: value,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "`key`"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&setting).Error
}

// Delete 删除设置
func (r *systemSettingRepository) Delete(key string) error {
	return r.db.Where("`key` = ?", key).Delete(&entity.SystemSetting{}).Error
}
