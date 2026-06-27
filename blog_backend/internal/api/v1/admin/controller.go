package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"blog_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 管理控制器
type Controller struct{}

// NewController 创建管理控制器
func NewController() *Controller {
	return &Controller{}
}

// UploadFile 上传文件
func (ctrl *Controller) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 限制文件大小 (10MB)
	if file.Size > 10<<20 {
		response.BadRequest(c, "文件大小不能超过 10MB")
		return
	}

	// 校验文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".svg": true,
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "不支持的文件类型，仅支持 jpg/png/gif/webp/svg")
		return
	}

	// 按年月生成存储路径: uploads/2026/06/
	now := time.Now()
	uploadDir := filepath.Join("uploads", now.Format("2006"), now.Format("01"))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名: 时间戳+随机数.ext
	filename := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.InternalError(c, "保存文件失败")
		return
	}

	// 返回可访问的 URL（确保使用正斜杠）
	url := "/" + filepath.ToSlash(savePath)
	response.Success(c, gin.H{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
