package article

import (
	bizerrors "blog_backend/pkg/errors"
	"blog_backend/internal/middleware"
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 文章控制器
type Controller struct {
	articleService service.ArticleService
}

// NewController 创建文章控制器
func NewController(articleService service.ArticleService) *Controller {
	return &Controller{
		articleService: articleService,
	}
}

// GetArticles 获取文章列表（前台）
func (ctrl *Controller) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	tagID, _ := strconv.ParseUint(c.Query("tag_id"), 10, 32)
	keyword := c.Query("keyword")
	sort := c.DefaultQuery("sort", "latest")

	articles, total, err := ctrl.articleService.GetArticleList(page, size, uint(categoryID), uint(tagID), keyword, sort)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  articles,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetArticleByID 根据 ID 获取文章详情
func (ctrl *Controller) GetArticleByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	article, err := ctrl.articleService.GetArticleDetail(uint(id))
	if err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	response.Success(c, article)
}

// GetHotArticles 获取热门文章
func (ctrl *Controller) GetHotArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := ctrl.articleService.GetHotArticles(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, articles)
}

// GetRecentArticles 获取最新文章
func (ctrl *Controller) GetRecentArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	articles, err := ctrl.articleService.GetRecentArticles(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, articles)
}

// AdminCreateArticle 管理员创建文章
func (ctrl *Controller) AdminCreateArticle(c *gin.Context) {
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	article, err := ctrl.articleService.CreateArticle(userID, &req)
	if err != nil {
		// 安全的类型断言
		if bizErr, ok := err.(*bizerrors.BizError); ok {
			response.ErrorWithBiz(c, bizErr)
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, article)
}

// AdminGetArticles 管理员获取文章列表
func (ctrl *Controller) AdminGetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	articles, total, err := ctrl.articleService.GetAdminArticleList(page, size, status, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  articles,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// AdminGetArticleByID 管理员根据 ID 获取文章
func (ctrl *Controller) AdminGetArticleByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	article, err := ctrl.articleService.GetAdminArticleDetail(uint(id))
	if err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	response.Success(c, article)
}

// AdminUpdateArticle 管理员更新文章
func (ctrl *Controller) AdminUpdateArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	article, err := ctrl.articleService.UpdateArticle(uint(id), &req)
	if err != nil {
		// 安全的类型断言
		if bizErr, ok := err.(*bizerrors.BizError); ok {
			response.ErrorWithBiz(c, bizErr)
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, article)
}

// AdminDeleteArticle 管理员删除文章
func (ctrl *Controller) AdminDeleteArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.articleService.DeleteArticle(uint(id)); err != nil {
		// 安全的类型断言
		if bizErr, ok := err.(*bizerrors.BizError); ok {
			response.ErrorWithBiz(c, bizErr)
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// AdminUpdateArticleStatus 管理员更新文章状态
func (ctrl *Controller) AdminUpdateArticleStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	article, err := ctrl.articleService.UpdateArticleStatus(uint(id), req.Status)
	if err != nil {
		// 安全的类型断言
		if bizErr, ok := err.(*bizerrors.BizError); ok {
			response.ErrorWithBiz(c, bizErr)
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, article)
}
