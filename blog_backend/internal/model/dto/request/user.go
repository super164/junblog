package request

// UpdateProfileRequest 修改用户信息基本信息
type UpdateProfileRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Avatar string `json:"avatar" binding:"omitempty,max=255"`
}

// UpdatePasswordRequest 修改用户密码
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// AdminResetPasswordRequest 管理员重置用户密码
type AdminResetPasswordRequest struct {
	Password string `json:"new_password" binding:"required,min=6"`
}

// AdminUpdateUserStatusRequest 管理员更新用户状态
type AdminUpdateUserStatusRequest struct {
	Status bool `json:"status"`
}

// AdminUpdateUserRequest 管理员更新用户信息
type AdminUpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone" binding:"omitempty,max=20"`
	Avatar string `json:"avatar" binding:"omitempty,max=255"`
	Role   string `json:"role" binding:"omitempty,oneof=user admin"`
}
