package entity

// Article 文章实体
type Article struct {
	BaseEntity
	Title          string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Slug           string `gorm:"type:varchar(255);uniqueIndex;not null;comment:别名" json:"slug"`
	Content        string `gorm:"type:text;comment:内容" json:"content"`
	CategoryID     uint   `gorm:"index;comment:分类ID" json:"category_id"`
	Status         string `gorm:"type:varchar(20);default:'published';comment:状态" json:"status"`
	Cover          string `gorm:"type:varchar(255);comment:封面" json:"cover"`
	AuthorID       uint   `gorm:"index;comment:作者ID" json:"author_id"`
	ViewsCount     uint   `gorm:"default:0;comment:浏览次数" json:"views_count"`
	LikesCount     int    `gorm:"default:0;comment:点赞总数" json:"likes_count"`
	FavoritesCount int    `gorm:"default:0;comment:收藏总数" json:"favorites_count"`

	Author   User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag    `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
