// 前台路由
import HomePage from '../../pages/HomePage.vue'
import ArticlesPage from '../../pages/ArticlesPage.vue'
import ArticleDetailPage from '../../pages/ArticleDetailPage.vue'
import AboutPage from '../../pages/AboutPage.vue'
import LoginPage from '../../pages/LoginPage.vue'
import RegisterPage from './views/RegisterPage.vue'

export const publicRoutes = [
  { path: '/', component: HomePage, meta: { title: 'JunBlog | 首页' } },
  { path: '/articles', component: ArticlesPage, meta: { title: 'JunBlog | 文章' } },
  { path: '/articles/:id', component: ArticleDetailPage, meta: { title: 'JunBlog | 文章详情' } },
  { path: '/about', component: AboutPage, meta: { title: 'JunBlog | 关于' } },
  { path: '/login', component: LoginPage, meta: { title: 'JunBlog | 登录' } },
  { path: '/register', component: RegisterPage, meta: { title: 'JunBlog | 注册' } },
]
