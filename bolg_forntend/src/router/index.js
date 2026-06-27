import { createRouter, createWebHistory } from 'vue-router'
import { publicRoutes } from '../modules/public/router'
import { adminRoutes, adminGuard } from '../modules/admin/router'
import { getStoredAuth } from '../shared/stores/auth'
import NotFoundPage from '../pages/NotFoundPage.vue'
import ForbiddenPage from '../pages/ForbiddenPage.vue'
import GithubCallbackPage from '../pages/GithubCallbackPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    ...publicRoutes,
    ...adminRoutes,
    { path: '/auth/github/callback', component: GithubCallbackPage, meta: { title: 'JunBlog | GitHub 登录' } },
    { path: '/403', component: ForbiddenPage, meta: { title: 'JunBlog | 无权限' } },
    { path: '/:pathMatch(.*)*', component: NotFoundPage, meta: { title: 'JunBlog | 404' } },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || 'JunBlog'

  // 后台路由权限检查
  if (to.path.startsWith('/admin')) {
    return adminGuard(to, from, next)
  }

  // 前台需登录页面检查
  if (to.meta.requiresAuth && !getStoredAuth().token) {
    next('/login')
    return
  }

  next()
})

export default router
