package service

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/model/entity"
	"blog_backend/internal/repository"
	bizerrors "blog_backend/pkg/errors"
	"blog_backend/pkg/config"
	"blog_backend/pkg/jwt"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// _newGitHubHTTPClient 创建支持代理的 HTTP 客户端
func _newGitHubHTTPClient(timeout time.Duration) *http.Client {
	// 优先使用环境变量中的代理设置
	proxyURL := os.Getenv("HTTPS_PROXY")
	if proxyURL == "" {
		proxyURL = os.Getenv("https_proxy")
	}
	if proxyURL == "" {
		proxyURL = os.Getenv("HTTP_PROXY")
	}
	if proxyURL == "" {
		proxyURL = os.Getenv("http_proxy")
	}
	if proxyURL == "" {
		proxyURL = os.Getenv("ALL_PROXY")
	}
	if proxyURL == "" {
		proxyURL = os.Getenv("all_proxy")
	}

	// 如果没有环境变量，使用默认代理
	if proxyURL == "" {
		proxyURL = "http://127.0.0.1:7890"
	}

	proxy, _ := url.Parse(proxyURL)
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	fmt.Printf("[GitHub OAuth] 使用代理: %s\n", proxyURL)
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (s *authService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, bizerrors.ErrInvalidCredentials
	}

	// 检查用户状态
	if !user.Status {
		return nil, bizerrors.ErrUserDisabled
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, bizerrors.ErrInvalidCredentials
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 构建响应
	return &response.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Avatar:   user.Avatar,
			Status:   user.Status,
		},
	}, nil
}

// RefreshToken 刷新 Token
func (s *authService) RefreshToken(token string) (*response.LoginResponse, error) {
	// 解析并校验 refresh token
	claims, err := jwt.ParseRefreshToken(token)
	if err != nil {
		return nil, bizerrors.ErrInvalidToken
	}

	// 验证用户是否存在
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, bizerrors.ErrUserNotFound
	}

	// 生成新的 access token
	newToken, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 生成新的 refresh token（轮换）
	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Avatar:   user.Avatar,
			Status:   user.Status,
		},
	}, nil
}

// Register 用户注册
func (s *authService) Register(req *request.RegisterRequest) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.ErrUserAlreadyExists
	}

	// 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.New(bizerrors.CodeUserAlreadyExists, "邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Status:   true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}

// githubTokenResponse GitHub access_token 响应
type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

// githubUserResponse GitHub 用户信息响应
type githubUserResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

// githubEmailResponse GitHub 用户邮箱响应
type githubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// GitHubLogin GitHub OAuth 登录
func (s *authService) GitHubLogin(code string, redirectURI string) (*response.LoginResponse, error) {
	cfg := config.Get()

	// 使用前端传来的 redirect_uri，如果没有则使用配置文件中的
	if redirectURI == "" {
		redirectURI = cfg.GitHub.RedirectURI
	}

	// 1. 用 code 换取 access_token
	accessToken, err := _exchangeGitHubCode(code, cfg, redirectURI)
	if err != nil {
		return nil, bizerrors.New(bizerrors.CodeInternalError, "GitHub 授权失败: "+err.Error())
	}

	// 2. 获取 GitHub 用户信息
	ghUser, err := _getGitHubUser(accessToken)
	if err != nil {
		return nil, bizerrors.New(bizerrors.CodeInternalError, "获取 GitHub 用户信息失败: "+err.Error())
	}

	// 3. 获取 GitHub 用户邮箱
	ghEmail, err := _getGitHubEmail(accessToken)
	if err != nil {
		ghEmail = ""
	}

	githubID := strconv.Itoa(ghUser.ID)

	// 4. 按 GitHub ID 查找已有用户
	user, err := s.userRepo.FindByGitHubID(githubID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return _generateLoginResponse(user)
	}

	// 5. 按邮箱查找（可能已通过注册创建）
	if ghEmail != "" {
		user, err = s.userRepo.FindByEmail(ghEmail)
		if err != nil {
			return nil, err
		}
		if user != nil {
			user.GitHubID = githubID
			user.GitHubLogin = ghUser.Login
			if user.Avatar == "" && ghUser.AvatarURL != "" {
				user.Avatar = ghUser.AvatarURL
			}
			if err := s.userRepo.Update(user); err != nil {
				return nil, err
			}
			return _generateLoginResponse(user)
		}
	}

	// 6. 全新用户，创建账号
	username := _generateUniqueUsername(ghUser.Login, s)
	password := _randomPassword()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Username:    username,
		Password:    string(hashedPassword),
		Email:       ghEmail,
		Avatar:      ghUser.AvatarURL,
		Status:      true,
		Role:        "user",
		GitHubID:    githubID,
		GitHubLogin: ghUser.Login,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return _generateLoginResponse(newUser)
}

// _exchangeGitHubCode 用授权码换取 access_token（带重试）
func _exchangeGitHubCode(code string, cfg *config.Config, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.GitHub.ClientID)
	data.Set("client_secret", cfg.GitHub.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	fmt.Printf("[GitHub OAuth] 用 code 换取 token, client_id=%s, redirect_uri=%s\n", cfg.GitHub.ClientID, redirectURI)

	// 重试最多 3 次，每次超时 15 秒
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			fmt.Printf("[GitHub OAuth] 第 %d 次重试...\n", attempt)
			time.Sleep(time.Duration(attempt-1) * 2 * time.Second) // 递增等待
		}

		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := _newGitHubHTTPClient(15 * time.Second)
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求 GitHub 失败: %w", err)
			fmt.Printf("[GitHub OAuth] 请求失败: %v\n", lastErr)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("[GitHub OAuth] GitHub 响应: %s\n", string(body))

		var tokenResp githubTokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			lastErr = fmt.Errorf("解析 GitHub 响应失败: %w", err)
			continue
		}

		if tokenResp.Error != "" {
			return "", fmt.Errorf("%s: %s", tokenResp.Error, tokenResp.ErrorDesc)
		}

		fmt.Printf("[GitHub OAuth] 获取 token 成功\n")
		return tokenResp.AccessToken, nil
	}

	return "", lastErr
}

// _getGitHubUser 获取 GitHub 用户信息
func _getGitHubUser(accessToken string) (*githubUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := _newGitHubHTTPClient(30 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var user githubUserResponse
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("解析 GitHub 用户信息失败: %w", err)
	}

	return &user, nil
}

// _getGitHubEmail 获取 GitHub 用户邮箱
func _getGitHubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := _newGitHubHTTPClient(30 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var emails []githubEmailResponse
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// 优先返回已验证的主邮箱
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	// 其次返回任意已验证邮箱
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", nil
}

// _generateUniqueUsername 生成唯一用户名（基于 GitHub login，冲突时加数字后缀）
func _generateUniqueUsername(ghLogin string, s *authService) string {
	username := ghLogin
	exists, _ := s.userRepo.ExistsByUsername(username)
	if !exists {
		return username
	}
	for i := 1; i <= 999; i++ {
		candidate := fmt.Sprintf("%s%d", ghLogin, i)
		exists, _ = s.userRepo.ExistsByUsername(candidate)
		if !exists {
			return candidate
		}
	}
	// 极端情况用时间戳
	return fmt.Sprintf("%s%d", ghLogin, time.Now().UnixMilli()%100000)
}

// _randomPassword 生成随机密码（OAuth 用户不需要记住密码）
func _randomPassword() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("gh_%08x", r.Uint32())
}

// _generateLoginResponse 生成登录响应
func _generateLoginResponse(user *entity.User) (*response.LoginResponse, error) {
	// 检查用户状态
	if !user.Status {
		return nil, bizerrors.ErrUserDisabled
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Avatar:   user.Avatar,
			Status:   user.Status,
		},
	}, nil
}
