package request

// AboutProfile 关于页面个人信息
type AboutProfile struct {
	Name   string        `json:"name"`
	Title  string        `json:"title"`
	Bio    string        `json:"bio"`
	Social []AboutSocial `json:"social"`
}

// AboutSocial 社交链接
type AboutSocial struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

// AboutSkill 技术栈
type AboutSkill struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// AboutTimeline 经历时间轴
type AboutTimeline struct {
	Year        string `json:"year"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AboutPhilosophy 写作理念
type AboutPhilosophy struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateAboutPageRequest 更新关于页面请求
type UpdateAboutPageRequest struct {
	Profile    *AboutProfile    `json:"profile"`
	Highlights []string         `json:"highlights"`
	Skills     []AboutSkill     `json:"skills"`
	Timeline   []AboutTimeline  `json:"timeline"`
	Philosophy *AboutPhilosophy `json:"philosophy"`
}

// UpdateSettingRequest 更新单个设置请求
type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}
