package service

import (
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/pkg/database"
)

// statsService 统计服务实现
type statsService struct{}

// NewStatsService 创建统计服务
func NewStatsService() StatsService {
	return &statsService{}
}

// GetStats 获取后台统计数据
func (s *statsService) GetStats() (*response.AdminStatsResponse, error) {
	db := database.GetMySQL()

	stats := &response.AdminStatsResponse{}

	// 文章统计
	db.Model(&entity.Article{}).Count(&stats.TotalArticles)
	db.Model(&entity.Article{}).Where("status = ?", "published").Count(&stats.PublishedArticles)
	db.Model(&entity.Article{}).Where("status = ?", "draft").Count(&stats.DraftArticles)

	// 浏览量统计
	db.Model(&entity.Article{}).Select("COALESCE(SUM(views_count), 0)").Scan(&stats.TotalViews)

	// 评论统计
	db.Model(&entity.Comment{}).Count(&stats.TotalComments)
	db.Model(&entity.Comment{}).Where("status = ?", "pending").Count(&stats.PendingComments)

	// 用户统计
	db.Model(&entity.User{}).Count(&stats.TotalUsers)

	// 分类统计
	db.Model(&entity.Category{}).Count(&stats.TotalCategories)

	// 标签统计
	db.Model(&entity.Tag{}).Count(&stats.TotalTags)

	// 最近文章
	var articles []entity.Article
	db.Preload("Author").Order("created_at DESC").Limit(5).Find(&articles)

	stats.RecentArticles = make([]response.ArticleListItem, len(articles))
	for i, a := range articles {
		stats.RecentArticles[i] = response.ArticleListItem{
			ID:         a.ID,
			Title:      a.Title,
			Slug:       a.Slug,
			Status:     a.Status,
			ViewsCount: a.ViewsCount,
			CreatedAt:  a.CreatedAt,
			Author: response.UserResponse{
				ID:       a.Author.ID,
				Username: a.Author.Username,
			},
		}
	}

	return stats, nil
}
