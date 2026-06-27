package entity

// Tag 标签实体
type Tag struct {
	BaseEntity
	Name       string `gorm:"type:varchar(50);uniqueIndex;not null;comment:名称" json:"name"`
	Slug       string `gorm:"type:varchar(50);uniqueIndex;not null;comment:别名" json:"slug"`
	CategoryID uint   `gorm:"index;comment:所属分类ID" json:"category_id"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Articles []Article  `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
