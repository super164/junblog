// 后台路由
import { getStoredAuth } from '../../shared/stores/auth'

// 后台页面懒加载
const AdminLayout = () => import('./layouts/AdminLayout.vue')
const DashboardPage = () => import('../../pages/DashboardPage.vue')
const ProfilePage = () => import('../../pages/ProfilePage.vue')
const ArticleListPage = () => import('./views/ArticleListPage.vue')
const ArticleEditPage = () => import('./views/ArticleEditPage.vue')
const CategoryManagePage = () => import('./views/CategoryManagePage.vue')
const TagManagePage = () => import('./views/TagManagePage.vue')
const CommentManagePage = () => import('./views/CommentManagePage.vue')
const UserManagePage = () => import('./views/UserManagePage.vue')
const SiteSettingPage = () => import('./views/SiteSettingPage.vue')

export const adminRoutes = [
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    redirect: '/admin/dashboard',
    children: [
      { path: 'dashboard', component: DashboardPage, meta: { title: 'JunBlog | 控制台' } },
      { path: 'profile', component: ProfilePage, meta: { title: 'JunBlog | 个人设置' } },
      { path: 'articles', component: ArticleListPage, meta: { title: 'JunBlog | 文章管理' } },
      { path: 'articles/new', component: ArticleEditPage, meta: { title: 'JunBlog | 新建文章' } },
      { path: 'articles/:id/edit', component: ArticleEditPage, meta: { title: 'JunBlog | 编辑文章' } },
      { path: 'categories', component: CategoryManagePage, meta: { title: 'JunBlog | 分类管理' } },
      { path: 'tags', component: TagManagePage, meta: { title: 'JunBlog | 标签管理' } },
      { path: 'comments', component: CommentManagePage, meta: { title: 'JunBlog | 评论管理' } },
      { path: 'users', component: UserManagePage, meta: { title: 'JunBlog | 用户管理' } },
      { path: 'settings', component: SiteSettingPage, meta: { title: 'JunBlog | 站点设置' } },
    ],
  },
]

// 后台路由守卫
export function adminGuard(to, from, next) {
  if (to.meta.requiresAuth || to.meta.requiresAdmin) {
    const auth = getStoredAuth()
    if (!auth.token) {
      next('/login')
      return
    }
    if (to.meta.requiresAdmin && auth.user?.role !== 'admin') {
      next('/403')
      return
    }
  }
  next()
}
