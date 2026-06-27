package user

import (
	"blog_backend/internal/middleware"
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	"blog_backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 用户控制器
type Controller struct {
	userService service.UserService
}

// NewController 创建用户控制器
func NewController(userService service.UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}

// GetProfile 获取当前用户
func (ctrl *Controller) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, ctrl.userService.GetUserResponse(user))
}

// UpdateProfile 修改用户信息
func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	if err := ctrl.userService.UpdateUser(userID, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdatePassword 修改用户密码
func (ctrl *Controller) UpdatePassword(c *gin.Context) {
	var req request.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	if err := ctrl.userService.ChangePassword(userID, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}

// AdminGetUsers 管理员获取用户列表
func (ctrl *Controller) AdminGetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("page_size", "10")
	keyword := c.DefaultQuery("keyword", "")

	// 解析 page 和 size 参数
	pageNum := 1
	sizeNum := 10
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}
	if s, err := strconv.Atoi(size); err == nil && s > 0 {
		sizeNum = s
	}

	users, total, err := ctrl.userService.ListUsersWithSearch(pageNum, sizeNum, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

// AdminResetPassword 管理员重置用户密码
func (ctrl *Controller) AdminResetPassword(c *gin.Context) {
	var req request.AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.AdminResetPassword(uint(userID), req.Password); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}

// AdminUpdateUserStatus 管理员更新用户状态
func (ctrl *Controller) AdminUpdateUserStatus(c *gin.Context) {
	var req request.AdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	// 检查是否是 admin 用户
	user, err := ctrl.userService.GetUserByID(uint(userID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if user.Role == "admin" {
		response.BadRequest(c, "不能禁用管理员账号")
		return
	}

	if err := ctrl.userService.AdminUpdateUserStatus(uint(userID), req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}

// AdminUpdateUser 管理员更新用户信息
func (ctrl *Controller) AdminUpdateUser(c *gin.Context) {
	var req request.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.AdminUpdateUser(uint(userID), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}
