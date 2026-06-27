package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
)

// UserService 用户服务接口
type UserService interface {
	// GetUserByID 根据 ID 获取用户
	GetUserByID(id uint) (*entity.User, error)
	// GetUserResponse 获取用户响应
	GetUserResponse(user *entity.User) *response.UserResponse
	// UpdateUser 更新用户信息
	UpdateUser(id uint, req *request.UpdateProfileRequest) error
	// ChangePassword 修改密码
	ChangePassword(id uint, req *request.UpdatePasswordRequest) error
	// ListUsers 分页获取用户列表
	ListUsers(page, size int) ([]*response.UserResponse, int64, error)
	// ListUsersWithSearch 分页获取用户列表（支持搜索）
	ListUsersWithSearch(page, size int, keyword string) ([]*response.UserResponse, int64, error)
	// AdminResetPassword 管理员重置用户密码
	AdminResetPassword(userID uint, password string) error
	// AdminUpdateUserStatus 管理员更新用户状态
	AdminUpdateUserStatus(userID uint, status bool) error
	// AdminUpdateUser 管理员更新用户信息
	AdminUpdateUser(userID uint, req *request.AdminUpdateUserRequest) error
}
