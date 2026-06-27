# JunBlog 项目设计 - 阶段1: 调研分析报告

## 一、现状差距分析

### 前端页面改造清单

| 页面 | 当前状态 | 需要改造 |
|------|---------|---------|
| HomePage.vue | mock 数据（articles, stats, profile） | 调用 GET /api/v1/articles + GET /api/v1/stats |
| ArticlesPage.vue | mock 数据，前端 computed 过滤 | 调用后端分页+筛选 API |
| ArticleDetailPage.vue | mock 数据，本地 find 匹配 | 调用详情 API，含评论+点赞状态 |
| DashboardPage.vue | mock 统计数据 | 调用 GET /api/v1/admin/stats |
| LoginPage.vue | ✅ 已对接真实 API | 补充注册入口 + Token 自动刷新 |
| ProfilePage.vue | ✅ 已对接真实 API | 无需改造 |
| **缺失** AdminArticleList | 不存在 | 新建：文章管理列表页 |
| **缺失** AdminArticleEdit | 不存在 | 新建：文章编辑页（含富文本编辑器） |
| **缺失** AdminCategoryManage | 不存在 | 新建：分类管理页 |
| **缺失** AdminTagManage | 不存在 | 新建：标签管理页 |
| **缺失** AdminCommentManage | 不存在 | 新建：评论审核页 |
| **缺失** AdminUserManage | 不存在 | 新建：用户管理页 |

### 后端 API 待实现清单

**P0（核心功能）：**
- 文章管理：前台列表/详情 + 后台 CRUD + 状态切换
- 分类管理：前台树形列表 + 后台 CRUD
- 标签管理：前台列表 + 后台 CRUD

**P1（交互功能）：**
- 评论系统：前台发表/列表 + 后台审核
- 点赞/收藏：toggle 操作
- 用户管理：后台列表、启用禁用

**P2（增强功能）：**
- 数据统计、文件上传 MinIO、站点信息

### 中间件缺失清单

1. CORS 中间件（文件为空）
2. 管理员角色鉴权中间件（缺 AdminAuth）
3. 统一分页参数处理
4. 请求限流中间件
5. JWT Token 刷新机制

### 安全问题（严重）

| 级别 | 问题 | 位置 |
|------|------|------|
| CRITICAL | 明文密码存储，CheckPassword 用 == 比较 | utils/jwt.go:20 |
| CRITICAL | config.yaml 硬编码 MySQL 密码、JWT 密钥、MinIO 凭据 | config/config.yaml |
| HIGH | JWT 密钥强度不足（blog_super_secret_key） | config/config.yaml |
| HIGH | 密码修改直接写明文 | service/user_service.go:108 |
| HIGH | 缺少 .gitignore | 项目根目录 |
| MEDIUM | UpdateArticleRequest.Status 类型不一致 | dto/request/article_request.go |
| MEDIUM | 缺少角色校验中间件 | router/router.go |

---

## 二、最佳实践调研

### 路由方案对比

| 方案 | 开发效率 | 代码复用 | 维护成本 | 架构清晰 | 扩展性 | 总分 |
|------|---------|---------|---------|---------|--------|------|
| A: 渐进式改造 | 9 | 9 | 7 | 5 | 4 | **34** |
| B: 前后台独立 | 5 | 4 | 4 | 8 | 9 | **30** |
| C: 模块化单体 | 8 | 8 | 8 | 9 | 8 | **41** ✅ |

**推荐方案 C：模块化单体**

### 后台管理功能建议

**Dashboard 指标：**
- 文章总数/月新增/待审核
- PV/UV 趋势 7d/30d
- 热门 Top10 文章
- 评论统计/待审核
- 用户增长趋势

**文章管理：**
- 编辑器推荐：TipTap（Vue3 原生，扩展性强）
- 状态流：draft → published + 定时发布
- Markdown + WYSIWYG + 代码高亮

### 后端分层规范

- **Controller**：HTTP 请求绑定 + 调用 service + 统一响应，不含业务逻辑
- **Service**：业务逻辑 + 调用 dao + 事务 + DTO 转换
- **Dao**：纯 CRUD + GORM，不含逻辑

### 安全修复方案

- **密码**：bcrypt.GenerateFromPassword(cost=12) + bcrypt.CompareHashAndPassword
- **JWT**：双 token — access 15min + refresh 7天，refresh 存 httpOnly cookie
- **RBAC**：实现 RequireRole 中间件，路由分组 admin 用 RequireRole(admin)

---

## 三、完整 API 规范

### 公开 API（/api/v1/，无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /articles | 文章列表（分页+分类/标签筛选+搜索） |
| GET | /articles/:slug | 文章详情（浏览量+1） |
| GET | /articles/hot | 热门文章（按浏览量） |
| GET | /articles/recent | 最新文章 |
| GET | /categories | 分类树 |
| GET | /tags | 标签列表 |

### 需登录 API（/api/v1/，JWT 认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /articles/:id/like | 点赞/取消（toggle） |
| POST | /articles/:id/favorite | 收藏/取消（toggle） |
| POST | /comments | 发表评论 |
| GET | /user/profile | 获取个人信息 |
| PUT | /user/profile | 更新个人信息 |

### 认证 API（/api/v1/auth/）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /register | 用户注册 |
| POST | /login | 用户登录 |
| POST | /refresh | 刷新 Token |
| POST | /logout | 退出登录 |

### 管理 API（/api/v1/admin/，JWT + admin 角色）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /articles | 文章列表（含所有状态） |
| POST | /articles | 创建文章 |
| GET | /articles/:id | 文章详情 |
| PUT | /articles/:id | 更新文章 |
| DELETE | /articles/:id | 删除文章 |
| PATCH | /articles/:id/status | 切换状态 |
| GET | /categories | 分类列表 |
| POST | /categories | 创建分类 |
| PUT | /categories/:id | 更新分类 |
| DELETE | /categories/:id | 删除分类 |
| GET | /tags | 标签列表 |
| POST | /tags | 创建标签 |
| PUT | /tags/:id | 更新标签 |
| DELETE | /tags/:id | 删除标签 |
| GET | /comments | 评论列表 |
| PATCH | /comments/:id/status | 审核评论 |
| DELETE | /comments/:id | 删除评论 |
| GET | /users | 用户列表 |
| PATCH | /users/:id/status | 启用/禁用用户 |
| GET | /stats/overview | 数据统计总览 |
