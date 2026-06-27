package tag

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 标签控制器
type Controller struct {
	tagService service.TagService
}

// NewController 创建标签控制器
func NewController(tagService service.TagService) *Controller {
	return &Controller{
		tagService: tagService,
	}
}

// GetTags 获取标签列表（前台）
func (ctrl *Controller) GetTags(c *gin.Context) {
	tags, err := ctrl.tagService.GetAllTags()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// GetTagsByCategoryID 根据分类ID获取标签列表（前台）
func (ctrl *Controller) GetTagsByCategoryID(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Param("category_id"), 10, 32)

	tags, err := ctrl.tagService.GetTagsByCategoryID(uint(categoryID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// AdminGetTags 管理员获取标签列表
func (ctrl *Controller) AdminGetTags(c *gin.Context) {
	tags, err := ctrl.tagService.GetAllTags()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// AdminGetTagsByCategoryID 管理员根据分类ID获取标签列表
func (ctrl *Controller) AdminGetTagsByCategoryID(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Param("category_id"), 10, 32)

	tags, err := ctrl.tagService.GetTagsByCategoryID(uint(categoryID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// AdminCreateTag 管理员创建标签
func (ctrl *Controller) AdminCreateTag(c *gin.Context) {
	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := ctrl.tagService.CreateTag(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tag)
}

// AdminUpdateTag 管理员更新标签
func (ctrl *Controller) AdminUpdateTag(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req request.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := ctrl.tagService.UpdateTag(uint(id), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tag)
}

// AdminDeleteTag 管理员删除标签
func (ctrl *Controller) AdminDeleteTag(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.tagService.DeleteTag(uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}
