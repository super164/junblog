# JunBlog 前后台分离博客系统 - 完整设计方案

> 日期：2026-06-04
> 方案：C - 模块化单体

---

## 一、推荐方案

**方案 C：模块化单体**（单 Vue 项目 + 目录隔离）

| 维度 | A 渐进式 | B 独立项目 | C 模块化单体 ✅ |
|------|---------|-----------|----------------|
| 开发效率 | 9 | 5 | 8 |
| 代码复用 | 9 | 4 | 8 |
| 维护成本 | 7 | 4 | 8 |
| 架构清晰 | 5 | 8 | 9 |
| 扩展性 | 4 | 9 | 8 |
| **总分** | 34 | 30 | **41** |

---

## 二、目录结构

```
bolg_forntend/src/
├── modules/
│   ├── public/                    # 前台模块
│   │   ├── router.js              # 前台路由
│   │   ├── stores/                # 前台状态管理
│   │   ├── views/                 # 前台页面
│   │   │   ├── HomePage.vue
│   │   │   ├── ArticlesPage.vue
│   │   │   ├── ArticleDetailPage.vue
│   │   │   ├── AboutPage.vue
│   │   │   ├── LoginPage.vue
│   │   │   └── RegisterPage.vue
│   │   └── components/            # 前台专用组件
│   │       ├── CommentSection.vue
│   │       └── InteractionBar.vue
│   └── admin/                     # 后台模块
│       ├── router.js              # 后台路由 + 权限守卫
│       ├── stores/                # 后台状态管理
│       ├── layouts/
│       │   └── AdminLayout.vue    # 后台布局框架
│       ├── views/
│       │   ├── DashboardPage.vue
│       │   ├── ArticleListPage.vue
│       │   ├── ArticleEditPage.vue
│       │   ├── CategoryManagePage.vue
│       │   ├── TagManagePage.vue
│       │   ├── CommentManagePage.vue
│       │   └── UserManagePage.vue
│       └── components/
│           ├── Sidebar.vue
│           ├── Topbar.vue
│           └── TipTapEditor.vue
├── shared/                        # 共享层
│   ├── api/
│   │   ├── client.js              # axios 实例 + 拦截器
│   │   └── modules/
│   │       ├── article.js
│   │       ├── auth.js
│   │       ├── category.js
│   │       ├── tag.js
│   │       └── comment.js
│   ├── stores/
│   │   └── auth.js                # 认证状态
│   ├── utils/
│   │   ├── date.js
│   │   └── text.js
│   ├── composables/
│   │   └── useSEO.js
│   └── components/
│       ├── Loading.vue
│       ├── Pagination.vue
│       ├── Modal.vue
│       └── Toast.vue
├── router/
│   └── index.js                   # 主路由（合并 public + admin）
├── App.vue
├── main.js
└── style.css
```

---

## 三、API 设计

### 公开 API（/api/v1/）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /articles | 文章列表（分页+筛选+搜索） |
| GET | /articles/:slug | 文章详情（浏览量+1） |
| GET | /articles/hot | 热门文章 |
| GET | /articles/recent | 最新文章 |
| GET | /categories | 分类树 |
| GET | /tags | 标签列表 |

### 需登录 API（/api/v1/ + JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /articles/:id/like | 点赞 toggle |
| POST | /articles/:id/favorite | 收藏 toggle |
| POST | /comments | 发表评论 |
| GET/PUT | /user/profile | 个人信息 |

### 认证 API（/api/v1/auth/）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /register | 注册 |
| POST | /login | 登录 |
| POST | /refresh | 刷新 Token |
| POST | /logout | 退出 |

### 管理 API（/api/v1/admin/ + JWT + admin）

| 方法 | 路径 | 说明 |
|------|------|------|
| CRUD | /articles | 文章管理 + 状态切换 |
| CRUD | /categories | 分类管理 |
| CRUD | /tags | 标签管理 |
| GET/PATCH/DELETE | /comments | 评论审核 |
| GET/PATCH | /users | 用户管理 |
| GET | /stats/overview | 数据统计 |
| POST | /upload | 文件上传 |

---

## 四、安全修复清单

| 级别 | 问题 | 修复方案 |
|------|------|---------|
| CRITICAL | 明文密码 | bcrypt(cost=12) 哈希 |
| CRITICAL | 硬编码密钥 | 迁移到 .env + .gitignore |
| HIGH | JWT 弱密钥 | 32 字节随机密钥 |
| HIGH | 无角色校验 | RequireRole 中间件 |
| MEDIUM | 无 Token 刷新 | access 15min + refresh 7d |

---

## 五、工时估算

| 阶段 | P0 | P1 | P2 | 小计 |
|------|----|----|-----|------|
| 0 基础设施 | 6h | 6h | - | **12h** |
| 1 后端 API | 13h | 9h | 7h | **25h** (核心22h) |
| 2 前端架构 | 8.5h | 2h | - | **10.5h** |
| 3 前台页面 | 8.5h | 5h | - | **13.5h** |
| 4 后台页面 | 16h | 5.5h | 3h | **24.5h** |
| 5 增强功能 | - | - | 15h | **15h** |
| **总计** | **52h** | **27.5h** | **25h** | **100.5h** |

- **MVP（仅 P0）**：52h ≈ 6.5 天
- **核心功能（P0+P1）**：79.5h ≈ 10 天
- **完整功能**：100.5h ≈ 13 天

---

## 六、详细任务清单

→ 见 [development-tasks.md](./development-tasks.md)

---

## 七、执行顺序

```
阶段0（安全修复）全部完成
    ↓
阶段1（后端API）←→ 阶段2（前端架构）并行
    ↓
阶段3（前台页面）←→ 阶段4（后台页面）并行
    ↓
阶段5（增强功能）可选
```

### 并行策略

**第一波**：阶段 0 全部（安全优先，不可跳过）
**第二波**：后端 P0 API + 前端架构重构 并行
**第三波**：前台页面改造 + 后台管理页面 并行
**第四波**：交互功能（评论/点赞/收藏）+ 后台补充 并行
**第五波**：增强功能（搜索/图表/SEO）可选

---

## 八、前后端路由对照

### 前台路由（无需登录）

| 路径 | 页面 | API |
|------|------|-----|
| / | 首页 | GET /articles/hot, /articles/recent |
| /articles | 文章列表 | GET /articles |
| /articles/:slug | 文章详情 | GET /articles/:slug |
| /about | 关于 | - |
| /login | 登录 | POST /auth/login |
| /register | 注册 | POST /auth/register |

### 后台路由（需 JWT + admin）

| 路径 | 页面 | API |
|------|------|-----|
| /admin | Dashboard | GET /admin/stats/overview |
| /admin/articles | 文章列表 | GET /admin/articles |
| /admin/articles/new | 新建文章 | POST /admin/articles |
| /admin/articles/:id/edit | 编辑文章 | PUT /admin/articles/:id |
| /admin/categories | 分类管理 | CRUD /admin/categories |
| /admin/tags | 标签管理 | CRUD /admin/tags |
| /admin/comments | 评论管理 | GET/PATCH /admin/comments |
| /admin/users | 用户管理 | GET/PATCH /admin/users |
| /admin/profile | 个人信息 | GET/PUT /user/profile |
