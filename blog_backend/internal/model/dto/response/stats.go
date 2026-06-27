package response

// AdminStatsResponse 管理后台统计数据
type AdminStatsResponse struct {
	TotalArticles    int64            `json:"total_articles"`
	PublishedArticles int64           `json:"published_articles"`
	DraftArticles    int64            `json:"draft_articles"`
	TotalViews       int64            `json:"total_views"`
	TotalComments    int64            `json:"total_comments"`
	PendingComments  int64            `json:"pending_comments"`
	TotalUsers       int64            `json:"total_users"`
	TotalCategories  int64            `json:"total_categories"`
	TotalTags        int64            `json:"total_tags"`
	RecentArticles   []ArticleListItem `json:"recent_articles"`
}
