package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", ctrl.Register)
		authGroup.POST("/login", ctrl.Login)
		authGroup.POST("/refresh", ctrl.RefreshToken)
		authGroup.POST("/github/callback", ctrl.GitHubCallback)
		// GitHub OAuth 回调（GET 请求，用于接收 GitHub 授权后重定向）
		authGroup.GET("/github/callback", ctrl.GitHubOAuthCallback)
		// 获取 GitHub 授权地址（前端调用，动态获取正确的回调地址）
		authGroup.GET("/github/url", ctrl.GetGitHubAuthURL)
	}
}
