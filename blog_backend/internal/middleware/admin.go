package middleware

import (
	"blog_backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole 要求用户具有指定角色
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if role == "" {
			response.Error(c, http.StatusForbidden, "无法获取用户角色")
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}
