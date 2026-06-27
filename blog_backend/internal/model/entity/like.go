package entity

// Like 点赞实体
type Like struct {
	BaseEntity
	UserID    uint `gorm:"index;comment:用户ID" json:"user_id"`
	ArticleID uint `gorm:"index;comment:文章ID" json:"article_id"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
