package response

// AboutPageResponse 关于页面响应
type AboutPageResponse struct {
	Profile    *AboutProfile    `json:"profile"`
	Highlights []string         `json:"highlights"`
	Skills     []*AboutSkill    `json:"skills"`
	Timeline   []*AboutTimeline `json:"timeline"`
	Philosophy *AboutPhilosophy `json:"philosophy"`
}

// AboutProfile 个人信息
type AboutProfile struct {
	Name   string         `json:"name"`
	Title  string         `json:"title"`
	Bio    string         `json:"bio"`
	Social []*AboutSocial `json:"social"`
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

// SettingResponse 设置响应
type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
