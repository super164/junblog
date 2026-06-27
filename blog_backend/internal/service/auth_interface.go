package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
	// RefreshToken 刷新 Token
	RefreshToken(token string) (*response.LoginResponse, error)
	// Register 用户注册
	Register(req *request.RegisterRequest) error
	// GitHubLogin GitHub OAuth 登录
	GitHubLogin(code string, redirectURI string) (*response.LoginResponse, error)
}
