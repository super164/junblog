package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// articleService 文章服务实现
type articleService struct {
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
}

// NewArticleRepository 创建文章服务
func NewArticleService(articleRepo repository.ArticleRepository, tagRepo repository.TagRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
	}
}

// generateSlug 从标题生成 slug（支持中英文）
func generateSlug(title string) string {
	var result []rune
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			result = append(result, '-')
		}
	}
	slug := strings.Trim(string(result), "-")
	// 压缩连续的连字符
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	if slug == "" {
		slug = "article"
	}
	// 添加时间戳后缀保证唯一性
	return fmt.Sprintf("%s-%d", slug, time.Now().Unix())
}

// CreateArticle 创建文章
func (s *articleService) CreateArticle(authorID uint, req *request.CreateArticleRequest) (*response.ArticleResponse, error) {
	// 自动生成 slug（如果为空）
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	// 检查 slug 唯一性
	existing, err := s.articleRepo.FindBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, bizerrors.New(bizerrors.CodeConflict, "slug 已存在")
	}

	article := &entity.Article{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		AuthorID:   authorID,
		Status:     req.Status,
	}
	if article.Status == "" {
		article.Status = "draft"
	}

	// 创建文章
	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	// 保存标签关联
	if len(req.TagIDs) > 0 {
		if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
			return nil, err
		}
	}

	// 重新查询获取完整关联
	created, err := s.articleRepo.FindByID(article.ID)
	if err != nil {
		return nil, err
	}

	return s.toArticleResponse(created), nil
}

// GetArticleDetail 获取文章详情（前台）
func (s *articleService) GetArticleDetail(id uint) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.ErrNotFound
	}

	// 只有已发布文章可访问（admin 用户可以在后台查看所有文章）
	if article.Status != "published" {
		return nil, bizerrors.ErrNotFound
	}

	// 增加浏览量
	_ = s.articleRepo.IncrementViews(article.ID)
	article.ViewsCount++

	return s.toArticleResponse(article), nil
}

// GetArticleList 获取文章列表
func (s *articleService) GetArticleList(page, size int, categoryID uint, tagID uint, keyword string, sort string) ([]*response.ArticleListItem, int64, error) {
	// 参数标准化
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	articles, total, err := s.articleRepo.List(offset, size, categoryID, tagID, keyword, sort)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*response.ArticleListItem, len(articles))
	for i, a := range articles {
		items[i] = s.toArticleListItem(a)
	}

	return items, total, nil
}

// GetHotArticles 获取热门文章
func (s *articleService) GetHotArticles(limit int) ([]*response.ArticleListItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	articles, err := s.articleRepo.GetHotArticles(limit)
	if err != nil {
		return nil, err
	}

	items := make([]*response.ArticleListItem, len(articles))
	for i, a := range articles {
		items[i] = s.toArticleListItem(a)
	}
	return items, nil
}

// GetRecentArticles 获取最新文章
func (s *articleService) GetRecentArticles(limit int) ([]*response.ArticleListItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 5
	}

	articles, err := s.articleRepo.GetRecentArticles(limit)
	if err != nil {
		return nil, err
	}

	items := make([]*response.ArticleListItem, len(articles))
	for i, a := range articles {
		items[i] = s.toArticleListItem(a)
	}
	return items, nil
}

// GetAdminArticleList 后台获取文章列表
func (s *articleService) GetAdminArticleList(page, size int, status string, keyword string) ([]*response.ArticleListItem, int64, error) {
	// 参数标准化
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	articles, total, err := s.articleRepo.ListAdmin(offset, size, status, keyword)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*response.ArticleListItem, len(articles))
	for i, a := range articles {
		items[i] = s.toArticleListItem(a)
	}

	return items, total, nil
}

// GetAdminArticleDetail 后台获取文章详情
func (s *articleService) GetAdminArticleDetail(id uint) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.ErrNotFound
	}

	return s.toArticleResponse(article), nil
}

// UpdateArticle 更新文章
func (s *articleService) UpdateArticle(id uint, req *request.UpdateArticleRequest) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.ErrNotFound
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Cover != "" {
		article.Cover = req.Cover
	}
	if req.CategoryID > 0 {
		article.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		article.Status = req.Status
	}
	if req.Slug != "" {
		// 检查 slug 唯一性（排除自身）
		existing, err := s.articleRepo.FindBySlug(req.Slug)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, bizerrors.New(bizerrors.CodeConflict, "slug 已存在")
		}
		article.Slug = req.Slug
	}

	// 更新标签
	if req.TagIDs != nil {
		if err := s.articleRepo.UpdateTags(id, req.TagIDs); err != nil {
			return nil, err
		}
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	// 重新查询获取完整关联
	updated, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toArticleResponse(updated), nil
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(id uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if article == nil {
		return bizerrors.ErrNotFound
	}

	return s.articleRepo.Delete(id)
}

// UpdateArticleStatus 更新文章状态
func (s *articleService) UpdateArticleStatus(id uint, status string) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.ErrNotFound
	}

	article.Status = status
	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	// 重新查询获取完整关联
	updated, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toArticleResponse(updated), nil
}

// fillCounts 从 likes/favorites 表实时填充文章的点赞和收藏计数
func (s *articleService) fillCounts(a *entity.Article) {
	if likeCount, err := s.articleRepo.CountLikes(a.ID); err == nil {
		a.LikesCount = int(likeCount)
	}
	if favCount, err := s.articleRepo.CountFavorites(a.ID); err == nil {
		a.FavoritesCount = int(favCount)
	}
}

// --- 转换函数 ---

func (s *articleService) toArticleResponse(a *entity.Article) *response.ArticleResponse {
	s.fillCounts(a)
	return &response.ArticleResponse{
		ID:             a.ID,
		Title:          a.Title,
		Slug:           a.Slug,
		Content:        a.Content,
		Cover:          a.Cover,
		Status:         a.Status,
		ViewsCount:     a.ViewsCount,
		LikesCount:     a.LikesCount,
		FavoritesCount: a.FavoritesCount,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Author: response.UserResponse{
			ID:       a.Author.ID,
			Username: a.Author.Username,
			Avatar:   a.Author.Avatar,
		},
		Category: response.CategoryResponse{
			ID:   a.Category.ID,
			Name: a.Category.Name,
			Slug: a.Category.Slug,
		},
		Tags: s.toTagResponses(a.Tags),
	}
}

func (s *articleService) toArticleListItem(a *entity.Article) *response.ArticleListItem {
	s.fillCounts(a)
	return &response.ArticleListItem{
		ID:             a.ID,
		Title:          a.Title,
		Slug:           a.Slug,
		Summary:        makeSummary(a.Content),
		Cover:          a.Cover,
		Status:         a.Status,
		ViewsCount:     a.ViewsCount,
		LikesCount:     a.LikesCount,
		FavoritesCount: a.FavoritesCount,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Author: response.UserResponse{
			ID:       a.Author.ID,
			Username: a.Author.Username,
			Avatar:   a.Author.Avatar,
		},
		Category: response.CategoryResponse{
			ID:   a.Category.ID,
			Name: a.Category.Name,
			Slug: a.Category.Slug,
		},
		Tags: s.toTagResponses(a.Tags),
	}
}

func (s *articleService) toTagResponses(tags []entity.Tag) []response.TagResponse {
	result := make([]response.TagResponse, len(tags))
	for i, t := range tags {
		result[i] = response.TagResponse{ID: t.ID, Name: t.Name, Slug: t.Slug}
	}
	return result
}

// makeSummary 从 HTML/Markdown 内容提取纯文本摘要（前 200 字符）
func makeSummary(content string) string {
	text := html.UnescapeString(content)
	// 去除 HTML 标签
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<p>", "")
	text = strings.ReplaceAll(text, "</p>", "\n")
	// 简单去除标签
	for {
		start := strings.Index(text, "<")
		end := strings.Index(text, ">")
		if start == -1 || end == -1 || end <= start {
			break
		}
		text = text[:start] + text[end+1:]
	}
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	if len([]rune(text)) > 200 {
		return string([]rune(text)[:200]) + "..."
	}
	if text == "" {
		return "暂无摘要"
	}
	return text
}
