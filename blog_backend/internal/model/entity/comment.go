package entity

// Comment 评论实体
type Comment struct {
	BaseEntity
	Content   string `gorm:"type:text;not null;comment:内容" json:"content"`
	UserID    uint   `gorm:"index;comment:用户ID" json:"user_id"`
	ArticleID uint   `gorm:"index;comment:文章ID" json:"article_id"`
	ParentID  *uint  `gorm:"index;comment:父评论ID" json:"parent_id"`
	Status    string `gorm:"type:varchar(20);default:'approved';comment:状态" json:"status"`

	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article Article  `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Parent  *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
