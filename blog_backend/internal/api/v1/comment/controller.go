package comment

import (
	"blog_backend/internal/middleware"
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 评论控制器
type Controller struct {
	commentService service.CommentService
}

// NewController 创建评论控制器
func NewController(commentService service.CommentService) *Controller {
	return &Controller{
		commentService: commentService,
	}
}

// GetComments 获取评论列表（前台）
func (ctrl *Controller) GetComments(c *gin.Context) {
	articleID, _ := strconv.ParseUint(c.Query("article_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	comments, total, err := ctrl.commentService.GetCommentsByArticleID(uint(articleID), page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  comments,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// CreateComment 创建评论（需登录）
func (ctrl *Controller) CreateComment(c *gin.Context) {
	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	comment, err := ctrl.commentService.CreateComment(userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, comment)
}

// AdminGetComments 管理员获取评论列表
func (ctrl *Controller) AdminGetComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")

	comments, total, err := ctrl.commentService.GetAdminCommentList(page, size, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  comments,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// AdminUpdateCommentStatus 管理员更新评论状态
func (ctrl *Controller) AdminUpdateCommentStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := ctrl.commentService.UpdateComment(uint(id), &request.UpdateCommentRequest{Status: req.Status})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, comment)
}

// AdminDeleteComment 管理员删除评论
func (ctrl *Controller) AdminDeleteComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.commentService.DeleteComment(uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}
