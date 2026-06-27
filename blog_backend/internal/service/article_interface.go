package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// ArticleService 文章服务接口
type ArticleService interface {
	// CreateArticle 创建文章
	CreateArticle(authorID uint, req *request.CreateArticleRequest) (*response.ArticleResponse, error)
	// GetArticleDetail 获取文章详情
	GetArticleDetail(id uint) (*response.ArticleResponse, error)
	// GetArticleList 获取文章列表
	GetArticleList(page, size int, categoryID uint, tagID uint, keyword string, sort string) ([]*response.ArticleListItem, int64, error)
	// GetHotArticles 获取热门文章
	GetHotArticles(limit int) ([]*response.ArticleListItem, error)
	// GetRecentArticles 获取最新文章
	GetRecentArticles(limit int) ([]*response.ArticleListItem, error)
	// GetAdminArticleList 后台获取文章列表
	GetAdminArticleList(page, size int, status string, keyword string) ([]*response.ArticleListItem, int64, error)
	// GetAdminArticleDetail 后台获取文章详情
	GetAdminArticleDetail(id uint) (*response.ArticleResponse, error)
	// UpdateArticle 更新文章
	UpdateArticle(id uint, req *request.UpdateArticleRequest) (*response.ArticleResponse, error)
	// DeleteArticle 删除文章
	DeleteArticle(id uint) error
	// UpdateArticleStatus 更新文章状态
	UpdateArticleStatus(id uint, status string) (*response.ArticleResponse, error)
}
