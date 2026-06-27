package auth

import (
	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/service"
	bizerrors "blog_backend/pkg/errors"
	"blog_backend/pkg/config"
	"blog_backend/pkg/response"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller 认证控制器
type Controller struct {
	authService service.AuthService
}

// NewController 创建认证控制器
func NewController(authService service.AuthService) *Controller {
	return &Controller{
		authService: authService,
	}
}

// GetGitHubAuthURL 获取 GitHub 授权地址（前端调用，动态获取正确的回调地址）
func (ctrl *Controller) GetGitHubAuthURL(c *gin.Context) {
	cfg := config.Get()
	origin := c.Query("origin") // 前端传来的当前访问地址

	// 确定回调地址：优先使用公网地址，否则使用配置的地址
	redirectURI := cfg.GitHub.RedirectURI
	if cfg.GitHub.PublicURL != "" {
		redirectURI = cfg.GitHub.PublicURL + "/api/v1/auth/github/callback"
	}

	// 构建 state，格式: origin|redirect_uri（这样回调时能拿到正确的 redirect_uri）
	frontendOrigin := origin
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}
	state := frontendOrigin + "|" + redirectURI

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		cfg.GitHub.ClientID,
		redirectURI,
		url.QueryEscape(state),
	)

	fmt.Printf("[GitHub OAuth] 生成授权URL: redirect_uri=%s, state=%s\n", redirectURI, state)
	response.Success(c, gin.H{
		"url":            authURL,
		"redirect_uri":   redirectURI,
		"state":          state,
		"public_url":     cfg.GitHub.PublicURL,
	})
}

// Register 用户注册
func (ctrl *Controller) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.authService.Register(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, nil)
}

// Login 用户登录
func (ctrl *Controller) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := ctrl.authService.Login(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// RefreshToken 刷新 Token
func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := ctrl.authService.RefreshToken(req.Token)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GitHubCallback GitHub OAuth 回调（POST，供前端调用）
func (ctrl *Controller) GitHubCallback(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`
		RedirectURI string `json:"redirect_uri"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "缺少授权码参数")
		return
	}

	resp, err := ctrl.authService.GitHubLogin(req.Code, req.RedirectURI)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GitHubOAuthCallback GitHub OAuth 回调（GET，用于接收 GitHub 授权后重定向）
func (ctrl *Controller) GitHubOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state") // state 中保存了前端的 origin

	fmt.Printf("[GitHub OAuth Callback] 收到回调: code=%s, state=%s\n", code, state)

	// 解码 state 获取前端地址和 redirect_uri
	// state 格式: origin|redirect_uri
	frontendOrigin := "http://localhost:5173"
	redirectURI := ""
	if decodedState, err := url.QueryUnescape(state); err == nil && decodedState != "" {
		parts := strings.SplitN(decodedState, "|", 2)
		if len(parts) == 2 {
			frontendOrigin = parts[0]
			redirectURI = parts[1]
		} else {
			frontendOrigin = decodedState
		}
	}
	fmt.Printf("[GitHub OAuth Callback] 前端地址: %s, redirect_uri: %s\n", frontendOrigin, redirectURI)

	if code == "" {
		fmt.Printf("[GitHub OAuth Callback] 错误: 没有 code 参数\n")
		c.Redirect(302, frontendOrigin+"/login?error=missing_code")
		return
	}

	// 使用与授权请求相同的 redirect_uri
	resp, err := ctrl.authService.GitHubLogin(code, redirectURI)
	if err != nil {
		fmt.Printf("[GitHub OAuth Callback] 错误: %v\n", err)
		// 判断是否是用户被禁用
		errorCode := "github_login_failed"
		if bizErr, ok := err.(*bizerrors.BizError); ok && bizErr.Code == bizerrors.CodeUserDisabled {
			errorCode = "user_disabled"
		}
		// 重定向到前端登录页并带上错误信息
		c.Redirect(302, frontendOrigin+"/login?error="+errorCode)
		return
	}

	fmt.Printf("[GitHub OAuth Callback] 成功: username=%s, role=%s\n", resp.User.Username, resp.User.Role)
	// 重定向到前端的 callback 页面，带上 token 信息
	c.Redirect(302, frontendOrigin+"/auth/github/callback?token="+resp.Token+"&refresh_token="+resp.RefreshToken+"&username="+resp.User.Username+"&role="+resp.User.Role)
}
