# JunBlog 开发任务清单

> 基于方案 C（模块化单体）拆解，后端 Go + Gin 三层架构，前端 Vue3 模块化目录隔离。
> 日期：2026-06-04

---

## 阶段 0：基础设施修复

**目标**：修复安全隐患，补齐中间件基础能力，为后续功能开发扫清障碍。

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T0.1 | 密码 bcrypt 加密 | 将明文密码存储改为 bcrypt(cost=12) 哈希存储；重写 `CheckPassword` 使用 `bcrypt.CompareHashAndPassword`；注册/创建用户时自动哈希 | `blog_backend/utils/password.go`（新建）, `blog_backend/service/user_service.go`, `blog_backend/dao/user_dao.go` | P0 | 2 | - |
| T0.2 | 配置外置化到 .env | 将 `config.yaml` 中硬编码的 MySQL 密码、JWT 密钥、MinIO 凭据等敏感信息迁移到 `.env` 文件；使用 `godotenv` 或 `viper` 读取环境变量；提供 `.env.example` 模板 | `blog_backend/config/config.go`, `blog_backend/config/config.yaml`, `blog_backend/.env.example`（新建）, `blog_backend/.env`（新建，gitignore） | P0 | 1.5 | - |
| T0.3 | JWT 密钥强化 + refresh token | 替换弱密钥 `blog_super_secret_key` 为至少 32 字节随机密钥（从 .env 读取）；实现双 token 机制：access token 15min + refresh token 7 天；refresh token 存 httpOnly cookie | `blog_backend/utils/jwt.go`, `blog_backend/middleware/jwt.go`, `blog_backend/config/jwt.go` | P0 | 2 | - |
| T0.4 | 密码修改接口哈希修复 | 修改密码更新接口，确保新密码经过 bcrypt 哈希后再写入数据库 | `blog_backend/service/user_service.go` | P0 | 0.5 | T0.1 |
| T0.5 | 添加 .gitignore | 在项目根目录和 `blog_backend/` 下添加 `.gitignore`，排除 `.env`、`node_modules/`、`dist/`、`*.exe`、IDE 配置等 | `F:\junblog\.gitignore`（新建） | P0 | 0.5 | - |
| T0.6 | CORS 中间件 | 当前 `middleware/cors.go` 文件为空；实现 CORS 中间件，允许前端开发服务器跨域访问，生产环境限制来源白名单 | `blog_backend/middleware/cors.go` | P1 | 1 | - |
| T0.7 | 管理员角色鉴权中间件 | 实现 `RequireRole` 中间件，从 JWT claims 中提取用户角色，非 admin 角色返回 403；在 admin 路由组上挂载 | `blog_backend/middleware/admin.go`（新建）, `blog_backend/router/router.go` | P1 | 1.5 | T0.3 |
| T0.8 | 统一分页处理 | 封装分页参数解析工具函数（page, page_size, sort_by, order），提供默认值和边界校验；统一分页响应结构 `{list, total, page, page_size}` | `blog_backend/utils/pagination.go`（新建） | P1 | 1 | - |
| T0.9 | Token 刷新机制 | 实现 `/api/v1/auth/refresh` 接口；前端在 access token 过期前自动用 refresh token 换取新 token；拦截 401 响应触发刷新流程 | `blog_backend/service/auth_service.go`（新建）, `blog_backend/controller/auth_controller.go`（新建）, `bolg_forntend/src/services/api.js` | P1 | 2 | T0.3 |

**阶段 0 小计**：12h

---

## 阶段 1：后端核心 API

**目标**：实现文章、分类、标签、评论、互动等全部后端业务 API。

### P0（核心功能）

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T1.1 | 文章 CRUD 服务 | 实现文章的创建、查询、更新、删除、状态切换；支持 slug 唯一性校验；文章关联分类和标签（多对多）；发布时自动生成摘要 | `blog_backend/dao/article_dao.go`, `blog_backend/service/article_service.go`, `blog_backend/controller/article_controller.go`, `blog_backend/dto/request/article_request.go`, `blog_backend/dto/response/article_response.go` | P0 | 3 | - |
| T1.2 | 文章公开 API | 实现前台文章接口：列表（分页+分类/标签筛选+关键词搜索）、详情（浏览量+1）、热门文章（按浏览量排序）、最新文章 | `blog_backend/controller/article_controller.go`, `blog_backend/router/router.go` | P0 | 2 | T1.1, T0.8 |
| T1.3 | 文章管理 API | 实现后台文章管理：列表（含所有状态/筛选）、创建、编辑、删除、状态切换（draft/published）；支持定时发布字段 | `blog_backend/controller/article_controller.go`, `blog_backend/router/router.go` | P0 | 2 | T1.1, T0.7 |
| T1.4 | 分类服务 + 树形结构 | 实现分类 CRUD；支持父子关系构建树形结构；分类删除时校验是否有关联文章 | `blog_backend/dao/category_dao.go`（新建）, `blog_backend/service/category_service.go`（新建）, `blog_backend/controller/category_controller.go`（新建） | P0 | 2 | - |
| T1.5 | 分类 API | 前台返回树形分类列表；后台提供完整 CRUD 接口 | `blog_backend/controller/category_controller.go`, `blog_backend/router/router.go` | P0 | 1.5 | T1.4 |
| T1.6 | 标签服务 | 实现标签 CRUD；支持按使用频次排序；标签删除时解除文章关联 | `blog_backend/dao/tag_dao.go`（新建）, `blog_backend/service/tag_service.go`（新建）, `blog_backend/controller/tag_controller.go`（新建） | P0 | 1.5 | - |
| T1.7 | 标签 API | 前台返回标签列表（含文章数量）；后台提供完整 CRUD 接口 | `blog_backend/controller/tag_controller.go`, `blog_backend/router/router.go` | P0 | 1 | T1.6 |

### P1（交互功能）

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T1.8 | 评论服务 + 审核 | 实现评论创建（关联文章+用户）、列表查询（支持嵌套回复）、审核状态管理（pending/approved/rejected）、敏感词基础过滤 | `blog_backend/dao/comment_dao.go`（新建）, `blog_backend/service/comment_service.go`（新建）, `blog_backend/controller/comment_controller.go`（新建） | P1 | 2.5 | - |
| T1.9 | 评论 API | 前台：发表评论（需登录）、按文章获取评论列表；后台：评论列表（含筛选）、审核/删除评论 | `blog_backend/controller/comment_controller.go`, `blog_backend/router/router.go` | P1 | 2 | T1.8 |
| T1.10 | 点赞收藏 toggle 服务 | 实现点赞/取消点赞、收藏/取消收藏的 toggle 操作；幂等性保证（重复操作返回当前状态）；文章详情返回当前用户的点赞/收藏状态 | `blog_backend/dao/like_dao.go`（新建）, `blog_backend/dao/favorite_dao.go`（新建）, `blog_backend/service/interaction_service.go`（新建）, `blog_backend/controller/interaction_controller.go`（新建） | P1 | 2 | - |
| T1.11 | 点赞收藏 API | 前台接口：POST 点赞 toggle、POST 收藏 toggle；文章详情附带当前用户互动状态 | `blog_backend/controller/interaction_controller.go`, `blog_backend/router/router.go` | P1 | 1 | T1.10 |
| T1.12 | 用户管理 API | 后台用户列表（分页+搜索）、启用/禁用用户状态；不含敏感信息（密码哈希等） | `blog_backend/controller/user_controller.go`, `blog_backend/service/user_service.go`, `blog_backend/router/router.go` | P1 | 1.5 | T0.7 |

### P2（增强功能）

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T1.13 | 数据统计 API | 实现 `/api/v1/admin/stats/overview`：文章总数/月新增、PV/UV 趋势（7d/30d）、热门 Top10、评论统计、用户增长 | `blog_backend/service/stats_service.go`（新建）, `blog_backend/controller/stats_controller.go`（新建）, `blog_backend/router/router.go` | P2 | 2 | T1.1 |
| T1.14 | MinIO 文件上传服务 | 封装 MinIO 客户端初始化；实现文件上传（图片/附件）、文件名唯一化、返回访问 URL；支持文件类型和大小校验 | `blog_backend/service/upload_service.go`（新建）, `blog_backend/config/minio.go` | P2 | 3 | - |
| T1.15 | 文件上传 API | 实现 POST `/api/v1/admin/upload` 接口，支持 multipart/form-data，返回文件 URL；支持多文件上传 | `blog_backend/controller/upload_controller.go`（新建）, `blog_backend/router/router.go` | P2 | 1 | T1.14 |
| T1.16 | 站点信息 API | 实现站点基本信息接口（站点名称、描述、公告、社交链接等），支持后台修改 | `blog_backend/service/site_service.go`（新建）, `blog_backend/controller/site_controller.go`（新建） | P2 | 1 | - |

**阶段 1 小计**：25h

---

## 阶段 2：前端架构重构

**目标**：从扁平目录结构迁移到模块化单体架构，实现前后台代码隔离。

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T2.1 | 创建模块化目录结构 | 创建 `src/modules/public/`（前台模块）和 `src/modules/admin/`（后台模块）以及 `src/shared/`（共享层）目录；迁移现有页面到对应模块；配置 Vite 别名 | `bolg_forntend/src/modules/public/`（新建目录）, `bolg_forntend/src/modules/admin/`（新建目录）, `bolg_forntend/src/shared/`（新建目录）, `bolg_forntend/vite.config.js` | P0 | 2 | - |
| T2.2 | 共享层封装 | 在 `src/shared/` 下封装：API 客户端（axios 实例 + 拦截器）、工具函数（日期格式化、文本截断等）、通用组件（Loading、Pagination、Modal、Toast） | `bolg_forntend/src/shared/api/client.js`（新建）, `bolg_forntend/src/shared/utils/`（新建）, `bolg_forntend/src/shared/components/`（新建） | P0 | 3 | T2.1 |
| T2.3 | 前台路由迁移 | 将现有前台页面（Home、Articles、ArticleDetail、About、Login）迁移到 `modules/public/` 下；配置前台子路由；更新导航组件引用 | `bolg_forntend/src/modules/public/router.js`（新建）, `bolg_forntend/src/modules/public/views/`（迁移页面）, `bolg_forntend/src/router/index.js` | P0 | 1.5 | T2.1 |
| T2.4 | 后台路由 + 权限守卫 | 在 `modules/admin/` 下配置后台路由；实现路由守卫：未登录跳转登录页、非 admin 角色跳转 403 页面；配置 AdminLayout 布局 | `bolg_forntend/src/modules/admin/router.js`（新建）, `bolg_forntend/src/modules/admin/views/`（新建目录）, `bolg_forntend/src/modules/admin/layouts/AdminLayout.vue`（新建） | P0 | 2 | T2.1, T0.7 |
| T2.5 | Pinia 状态管理按模块拆分 | 将全局状态按模块拆分：`useAuthStore`（认证）、`useArticleStore`（文章）、`useAdminStore`（后台管理）；迁移现有 mock 数据逻辑到对应 store | `bolg_forntend/src/shared/stores/auth.js`（新建）, `bolg_forntend/src/modules/public/stores/article.js`（新建）, `bolg_forntend/src/modules/admin/stores/admin.js`（新建） | P1 | 2 | T2.1 |

**阶段 2 小计**：10.5h

---

## 阶段 3：前台页面改造

**目标**：将前台所有 mock 数据页面切换为真实 API 调用。

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T3.1 | API 客户端封装 | 基于 axios 创建统一 API 客户端：请求拦截器自动附加 JWT token、响应拦截器处理 401 自动刷新 token、统一错误提示（Toast）、请求取消（页面切换时） | `bolg_forntend/src/shared/api/client.js`, `bolg_forntend/src/shared/api/modules/article.js`（新建）, `bolg_forntend/src/shared/api/modules/auth.js`（新建） | P0 | 2 | T0.9 |
| T3.2 | 首页对接真实 API | 替换 `HomePage.vue` 中的 mock 数据：调用 `GET /api/v1/articles/hot` 展示热门文章、`GET /api/v1/articles/recent` 展示最新文章、`GET /api/v1/stats` 展示站点统计 | `bolg_forntend/src/modules/public/views/HomePage.vue` | P0 | 2 | T1.2 |
| T3.3 | 文章列表页对接分页筛选 | 替换 `ArticlesPage.vue` 中的 computed 过滤逻辑：调用后端分页 API、支持分类/标签筛选、关键词搜索、排序切换；实现分页组件交互 | `bolg_forntend/src/modules/public/views/ArticlesPage.vue` | P0 | 2.5 | T1.2 |
| T3.4 | 文章详情页对接 | 替换 `ArticleDetailPage.vue` 中的本地 find 匹配：调用 `GET /api/v1/articles/:slug` 获取详情、展示评论列表、显示当前用户点赞/收藏状态 | `bolg_forntend/src/modules/public/views/ArticleDetailPage.vue` | P0 | 2 | T1.2 |
| T3.5 | 评论组件 | 实现评论组件：评论列表展示（支持嵌套）、评论输入框、提交评论（需登录提示）、评论分页加载 | `bolg_forntend/src/modules/public/components/CommentSection.vue`（新建）, `bolg_forntend/src/modules/public/components/CommentItem.vue`（新建） | P1 | 3 | T1.9 |
| T3.6 | 点赞收藏组件 | 实现互动组件：点赞按钮（toggle + 计数动画）、收藏按钮（toggle）、未登录时引导登录 | `bolg_forntend/src/modules/public/components/InteractionBar.vue`（新建） | P1 | 2 | T1.11 |
| T3.7 | 登录注册页 | 在现有 `LoginPage.vue` 基础上补充注册功能：注册表单、表单校验、注册成功自动跳转登录、登录成功后 token 存储 + 路由跳转 | `bolg_forntend/src/modules/public/views/LoginPage.vue`, `bolg_forntend/src/modules/public/views/RegisterPage.vue`（新建） | P0 | 2.5 | - |

**阶段 3 小计**：16h

---

## 阶段 4：后台管理页面

**目标**：实现完整的后台管理系统界面。

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T4.1 | AdminLayout 布局 | 实现后台管理布局：左侧可折叠侧边栏（导航菜单）、顶部栏（用户信息+退出）、内容区域；响应式适配移动端 | `bolg_forntend/src/modules/admin/layouts/AdminLayout.vue`, `bolg_forntend/src/modules/admin/components/Sidebar.vue`（新建）, `bolg_forntend/src/modules/admin/components/Topbar.vue`（新建） | P0 | 3 | T2.4 |
| T4.2 | Dashboard 统计看板 | 实现后台首页 Dashboard：文章总数/月新增/待审核卡片、PV/UV 趋势折线图（7d/30d 切换）、热门文章排行、评论统计；使用 ECharts 或 Chart.js | `bolg_forntend/src/modules/admin/views/DashboardPage.vue` | P1 | 3 | T1.13 |
| T4.3 | 文章列表管理页 | 实现后台文章管理：表格列表（标题/分类/标签/状态/发布时间/操作）、状态筛选、关键词搜索、批量操作（删除/状态切换）、分页 | `bolg_forntend/src/modules/admin/views/ArticleListPage.vue`（新建） | P0 | 3 | T1.3 |
| T4.4 | 文章编辑页 + 富文本编辑器 | 实现文章编辑页：标题/Slug/摘要/正文（TipTap 编辑器）、分类选择（下拉树）、标签选择（多选+新建）、封面图上传、状态切换（保存草稿/发布）；支持 Markdown 快捷键 | `bolg_forntend/src/modules/admin/views/ArticleEditPage.vue`（新建）, `bolg_forntend/src/modules/admin/components/TipTapEditor.vue`（新建） | P0 | 5 | T1.3 |
| T4.5 | 分类管理页 | 实现分类管理：树形表格展示、新增/编辑分类（名称/别名/父级/描述）、删除确认弹窗、拖拽排序（可选） | `bolg_forntend/src/modules/admin/views/CategoryManagePage.vue`（新建） | P0 | 3 | T1.5 |
| T4.6 | 标签管理页 | 实现标签管理：标签列表表格、新增/编辑标签（名称/别名/颜色）、删除确认、搜索过滤 | `bolg_forntend/src/modules/admin/views/TagManagePage.vue`（新建） | P0 | 2 | T1.7 |
| T4.7 | 评论审核页 | 实现评论管理：评论列表（含文章标题/评论者/内容/时间/状态）、状态筛选（全部/待审核/已通过/已拒绝）、审核操作（通过/拒绝）、删除评论 | `bolg_forntend/src/modules/admin/views/CommentManagePage.vue`（新建） | P1 | 3 | T1.9 |
| T4.8 | 用户管理页 | 实现用户管理：用户列表表格（用户名/邮箱/角色/状态/注册时间）、搜索过滤、启用/禁用用户操作、角色标识展示 | `bolg_forntend/src/modules/admin/views/UserManagePage.vue`（新建） | P1 | 2.5 | T1.12 |

**阶段 4 小计**：24.5h

---

## 阶段 5：增强功能

**目标**：提升搜索、数据可视化、SEO 等体验。

| ID | 标题 | 描述 | 涉及文件 | 优先级 | 工时(h) | 依赖 |
|----|------|------|----------|--------|---------|------|
| T5.1 | 全文搜索 | 基于 MySQL FULLTEXT 或引入 MeiliSearch/Elasticsearch，实现文章标题+内容全文搜索；前台搜索框支持实时搜索建议 | `blog_backend/service/search_service.go`（新建）, `blog_backend/controller/search_controller.go`（新建）, `bolg_forntend/src/modules/public/components/SearchBar.vue`（新建） | P2 | 5 | T1.1 |
| T5.2 | 统计图表 | 丰富 Dashboard 图表：文章发布趋势柱状图、分类分布饼图、标签词云、月度 PV/UV 对比图；支持时间范围选择器 | `bolg_forntend/src/modules/admin/components/charts/`（新建目录）, `bolg_forntend/src/modules/admin/views/DashboardPage.vue` | P2 | 4 | T1.13 |
| T5.3 | 图片上传组件 | 实现通用图片上传组件：拖拽上传、粘贴上传、进度条、图片预览、裁剪（可选）；集成 MinIO 上传 API；编辑器内嵌图片上传 | `bolg_forntend/src/shared/components/ImageUploader.vue`（新建）, `bolg_forntend/src/modules/admin/components/TipTapEditor.vue` | P2 | 3 | T1.15 |
| T5.4 | SEO 优化 | 前台页面添加 meta 信息（title/description/keywords）；文章页 Open Graph 标签；使用 vue-meta 或 @vueuse/head；生成 sitemap.xml；SSR 预渲染关键页面（可选） | `bolg_forntend/src/modules/public/views/*.vue`, `bolg_forntend/src/shared/composables/useSEO.js`（新建） | P2 | 3 | T3.2 |

**阶段 5 小计**：15h

---

## 关键路径分析

关键路径（决定项目最短完成时间的最长依赖链）：

```
T0.1 → T0.4
T0.3 → T0.7 → T1.3 → T4.3 / T4.4
T0.3 → T0.9 → T3.1
T1.1 → T1.2 → T3.2 / T3.3 / T3.4
T2.1 → T2.4 → T4.1 → T4.3 / T4.4
```

**最长路径**：

```
T0.3 (2h) → T0.7 (1.5h) → T1.3 (2h) → T4.4 (5h) = 10.5h
T2.1 (2h) → T2.4 (2h) → T4.1 (3h) → T4.4 (5h) = 12h（前端关键路径）
```

前后端关键路径在 T4.4（文章编辑页）汇合，总关键路径约 **12h**（前端侧），需前后端并行推进才能压缩。

---

## 总工时汇总表

| 阶段 | P0 | P1 | P2 | 小计 |
|------|----|----|----|------|
| 阶段 0：基础设施修复 | 6h | 6h | - | **12h** |
| 阶段 1：后端核心 API | 13h | 9h | 7h | **25h**（P0+P1 核心 = 22h） |
| 阶段 2：前端架构重构 | 8.5h | 2h | - | **10.5h** |
| 阶段 3：前台页面改造 | 8.5h | 5h | - | **13.5h**（P0 = 8.5h） |
| 阶段 4：后台管理页面 | 16h | 5.5h | 3h | **24.5h** |
| 阶段 5：增强功能 | - | - | 15h | **15h** |
| **总计** | **52h** | **27.5h** | **25h** | **100.5h** |

**MVP 范围（仅 P0）**：52h
**核心功能范围（P0 + P1）**：79.5h
**完整功能范围（P0 + P1 + P2）**：100.5h

---

## 并行执行策略

### 第一波：基础设施（阶段 0）—— 2 人并行

| 线程 A（安全） | 线程 B（基础设施） |
|----------------|-------------------|
| T0.1 密码 bcrypt (2h) | T0.2 配置外置化 (1.5h) |
| T0.3 JWT 强化 (2h) | T0.5 .gitignore (0.5h) |
| T0.4 密码修改修复 (0.5h) | T0.6 CORS 中间件 (1h) |
| T0.7 管理员鉴权 (1.5h) | T0.8 统一分页 (1h) |
| T0.9 Token 刷新 (2h) | |

完成后进入下一波。

### 第二波：后端 API + 前端架构（阶段 1 + 2）—— 2 人并行

| 线程 A（后端） | 线程 B（前端架构） |
|---------------|-------------------|
| T1.1 文章 CRUD (3h) | T2.1 目录结构 (2h) |
| T1.4 分类服务 (2h) | T2.2 共享层封装 (3h) |
| T1.6 标签服务 (1.5h) | T2.3 前台路由迁移 (1.5h) |
| T1.2 文章公开 API (2h) | T2.4 后台路由+守卫 (2h) |
| T1.5 分类 API (1.5h) | T2.5 状态管理拆分 (2h) |
| T1.7 标签 API (1h) | |
| T1.3 文章管理 API (2h) | |

### 第三波：页面对接（阶段 3 + 4 部分）—— 2 人并行

| 线程 A（前台页面） | 线程 B（后台页面） |
|-------------------|-------------------|
| T3.7 登录注册页 (2.5h) | T4.1 AdminLayout (3h) |
| T3.1 API 客户端封装 (2h) | T4.5 分类管理页 (3h) |
| T3.2 首页对接 (2h) | T4.6 标签管理页 (2h) |
| T3.3 文章列表页 (2.5h) | T4.3 文章列表管理 (3h) |
| T3.4 文章详情页 (2h) | T4.4 文章编辑页 (5h) |

### 第四波：交互功能 + 后台补充（阶段 1 P1 + 阶段 4 剩余）

| 线程 A（后端交互 API） | 线程 B（前端交互 + 后台） |
|------------------------|--------------------------|
| T1.8 评论服务 (2.5h) | T3.5 评论组件 (3h) |
| T1.9 评论 API (2h) | T3.6 点赞收藏组件 (2h) |
| T1.10 点赞收藏服务 (2h) | T4.7 评论审核页 (3h) |
| T1.11 点赞收藏 API (1h) | T4.8 用户管理页 (2.5h) |
| T1.12 用户管理 API (1.5h) | T4.2 Dashboard (3h) |

### 第五波：增强功能（阶段 5）—— 可选并行

| 线程 A | 线程 B |
|--------|--------|
| T5.1 全文搜索 (5h) | T5.3 图片上传组件 (3h) |
| T5.2 统计图表 (4h) | T5.4 SEO 优化 (3h) |

### 依赖关系图（简化）

```
阶段0 全部完成
  ├── 阶段1 (后端API)  ──┐
  └── 阶段2 (前端架构) ──┤
                          ├── 阶段3 (前台页面)
                          └── 阶段4 (后台页面)
                                └── 阶段5 (增强功能)
```

**建议**：优先完成阶段 0 全部任务（安全修复不可跳过），然后阶段 1 和阶段 2 可由不同开发者并行推进。阶段 3 和阶段 4 的前台/后台页面也可以并行开发，只要对应后端 API 已就绪。
