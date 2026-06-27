package entity

// User 用户实体
type User struct {
	BaseEntity
	Username string `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password string `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Phone    string `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Role     string `gorm:"type:varchar(20);default:'user';comment:角色" json:"role"`
	Avatar   string `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	Status      bool   `gorm:"default:true;comment:状态" json:"status"`
	GitHubID    string `gorm:"type:varchar(50);uniqueIndex;comment:GitHub用户ID" json:"github_id,omitempty"`
	GitHubLogin string `gorm:"type:varchar(100);comment:GitHub用户名" json:"github_login,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
