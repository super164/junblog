package api

import (
	"blog_backend/internal/api/v1/auth"
	"blog_backend/internal/api/v1/user"
	"blog_backend/internal/api/v1/article"
	"blog_backend/internal/api/v1/category"
	"blog_backend/internal/api/v1/tag"
	"blog_backend/internal/api/v1/comment"
	"blog_backend/internal/api/v1/interaction"
	"blog_backend/internal/api/v1/stats"
	"blog_backend/internal/api/v1/admin"
	"blog_backend/internal/middleware"
	"blog_backend/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	userCtrl        *user.Controller
	authCtrl        *auth.Controller
	articleCtrl     *article.Controller
	categoryCtrl    *category.Controller
	tagCtrl         *tag.Controller
	commentCtrl     *comment.Controller
	interactionCtrl *interaction.Controller
	statsCtrl       *stats.Controller
	adminCtrl       *admin.Controller
}

// NewRouter 创建路由
func NewRouter(
	userService service.UserService,
	authService service.AuthService,
	articleService service.ArticleService,
	categoryService service.CategoryService,
	tagService service.TagService,
	commentService service.CommentService,
	interactionService service.InteractionService,
	statsService service.StatsService,
) *Router {
	return &Router{
		userCtrl:        user.NewController(userService),
		authCtrl:        auth.NewController(authService),
		articleCtrl:     article.NewController(articleService),
		categoryCtrl:    category.NewController(categoryService),
		tagCtrl:         tag.NewController(tagService),
		commentCtrl:     comment.NewController(commentService),
		interactionCtrl: interaction.NewController(interactionService),
		statsCtrl:       stats.NewController(statsService),
		adminCtrl:       admin.NewController(),
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 静态文件服务：上传的图片
	engine.Static("/uploads", "./uploads")

	// API v1 路由组
	v1 := engine.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "blog API is running",
			})
		})

		// 认证路由（无需认证）
		r.authCtrl.RegisterRoutes(v1)

		// 公开路由（无需认证）
		r.articleCtrl.RegisterPublicRoutes(v1)
		r.categoryCtrl.RegisterPublicRoutes(v1)
		r.tagCtrl.RegisterPublicRoutes(v1)
		r.commentCtrl.RegisterPublicRoutes(v1)

		// 需要认证的路由
		authorized := v1.Group("/")
		authorized.Use(middleware.Auth())
		{
			r.userCtrl.RegisterRoutes(authorized)
			r.articleCtrl.RegisterAuthRoutes(authorized)
			r.commentCtrl.RegisterAuthRoutes(authorized)
			r.interactionCtrl.RegisterRoutes(authorized)
		}

		// 管理路由（需要 admin 角色）
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth())
		admin.Use(middleware.RequireRole("admin"))
		{
			r.adminCtrl.RegisterAdminRoutes(admin)
			r.statsCtrl.RegisterRoutes(admin)
			r.articleCtrl.RegisterAdminRoutes(admin)
			r.categoryCtrl.RegisterAdminRoutes(admin)
			r.tagCtrl.RegisterAdminRoutes(admin)
			r.commentCtrl.RegisterAdminRoutes(admin)
			r.userCtrl.RegisterAdminRoutes(admin)
		}
	}
}
