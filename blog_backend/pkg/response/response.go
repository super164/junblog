package response

import (
	"net/http"

	"blog_backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response[any]{
		Code: errors.CodeSuccess,
		Data: data,
		Msg:  "success",
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response[any]{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

// ErrorWithBiz 业务错误响应
func ErrorWithBiz(c *gin.Context, bizErr *errors.BizError) {
	c.JSON(http.StatusOK, Response[any]{
		Code: bizErr.Code,
		Data: nil,
		Msg:  bizErr.Message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, msg)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, msg)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, msg)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, msg)
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, msg)
}

// HandleError 处理错误，自动判断是否为业务错误
func HandleError(c *gin.Context, err error) {
	if bizErr, ok := err.(*errors.BizError); ok {
		ErrorWithBiz(c, bizErr)
	} else {
		InternalError(c, err.Error())
	}
}
