package entity

// Category 分类实体
type Category struct {
	BaseEntity
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null;comment:名称" json:"name"`
	Slug        string `gorm:"type:varchar(50);uniqueIndex;not null;comment:别名" json:"slug"`
	Description string `gorm:"type:text;comment:描述" json:"description"`
	ParentID    *uint  `gorm:"index;comment:父分类ID" json:"parent_id"`

	Articles []Article  `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
