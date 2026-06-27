package category

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 分类控制器
type Controller struct {
	categoryService service.CategoryService
}

// NewController 创建分类控制器
func NewController(categoryService service.CategoryService) *Controller {
	return &Controller{
		categoryService: categoryService,
	}
}

// GetCategories 获取分类列表（前台）
func (ctrl *Controller) GetCategories(c *gin.Context) {
	categories, err := ctrl.categoryService.GetAllCategories()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// GetCategoryTree 获取分类树（前台）
func (ctrl *Controller) GetCategoryTree(c *gin.Context) {
	tree, err := ctrl.categoryService.GetCategoryTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tree)
}

// AdminGetCategories 管理员获取分类列表
func (ctrl *Controller) AdminGetCategories(c *gin.Context) {
	categories, err := ctrl.categoryService.GetAllCategories()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// AdminCreateCategory 管理员创建分类
func (ctrl *Controller) AdminCreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := ctrl.categoryService.CreateCategory(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, category)
}

// AdminUpdateCategory 管理员更新分类
func (ctrl *Controller) AdminUpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := ctrl.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, category)
}

// AdminDeleteCategory 管理员删除分类
func (ctrl *Controller) AdminDeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.categoryService.DeleteCategory(uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}
