package entity

// SystemSetting 系统设置实体
type SystemSetting struct {
	BaseEntity
	Key   string `gorm:"type:varchar(100);uniqueIndex;not null;column:key;comment:设置键" json:"key"`
	Value string `gorm:"type:text;not null;comment:设置值(JSON)" json:"value"`
}

// TableName 指定表名
func (SystemSetting) TableName() string {
	return "system_settings"
}
