package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID 根据 ID 获取用户
func (s *userService) GetUserByID(id uint) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, bizerrors.ErrUserNotFound
	}
	return user, nil
}

// GetUserResponse 获取用户响应
func (s *userService) GetUserResponse(user *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id uint, req *request.UpdateProfileRequest) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id uint, req *request.UpdatePasswordRequest) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return bizerrors.ErrInvalidCredentials
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// ListUsers 分页获取用户列表
func (s *userService) ListUsers(page, size int) ([]*response.UserResponse, int64, error) {
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

	// 计算偏移量
	offset := (page - 1) * size

	// 查询数据
	users, total, err := s.userRepo.List(offset, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应 DTO
	list := make([]*response.UserResponse, 0, len(users))
	for _, user := range users {
		list = append(list, s.GetUserResponse(user))
	}

	return list, total, nil
}

// ListUsersWithSearch 分页获取用户列表（支持搜索）
func (s *userService) ListUsersWithSearch(page, size int, keyword string) ([]*response.UserResponse, int64, error) {
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

	// 计算偏移量
	offset := (page - 1) * size

	// 查询数据
	users, total, err := s.userRepo.ListWithSearch(offset, size, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应 DTO
	list := make([]*response.UserResponse, 0, len(users))
	for _, user := range users {
		list = append(list, s.GetUserResponse(user))
	}

	return list, total, nil
}

// AdminResetPassword 管理员重置用户密码
func (s *userService) AdminResetPassword(userID uint, password string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// AdminUpdateUserStatus 管理员更新用户状态
func (s *userService) AdminUpdateUserStatus(userID uint, status bool) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Status = status
	return s.userRepo.Update(user)
}

// AdminUpdateUser 管理员更新用户信息
func (s *userService) AdminUpdateUser(userID uint, req *request.AdminUpdateUserRequest) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	return s.userRepo.Update(user)
}
