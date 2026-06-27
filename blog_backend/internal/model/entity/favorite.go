package entity

// Favorite 收藏实体
type Favorite struct {
	BaseEntity
	UserID    uint `gorm:"index;comment:用户ID" json:"user_id"`
	ArticleID uint `gorm:"index;comment:文章ID" json:"article_id"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}
