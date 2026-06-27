<script setup>
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuth, clearAuthState } from '../shared/stores/auth'

const route = useRoute()
const router = useRouter()

const navItems = [
  { label: '首页', to: '/' },
  { label: '文章', to: '/articles' },
  { label: '关于', to: '/about' },
]

const auth = useAuth()

function logout() {
  clearAuthState()
  router.push('/')
}
</script>

<template>
  <div class="app-shell">
    <header class="topbar">
      <div class="topbar-inner">
        <RouterLink class="brand" to="/">
          <img src="/logo.png" alt="JunBlog" class="brand-logo" />
        </RouterLink>
        <nav class="nav-links">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            class="nav-link"
            :class="{ active: route.path === item.to }"
            :to="item.to"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
        <div class="topbar-actions">
          <template v-if="auth.user">
            <RouterLink v-if="auth.user.role === 'admin'" class="ghost-button small" to="/admin">管理后台</RouterLink>
            <div class="user-chip">
              <span class="user-avatar">{{ auth.user.username?.slice(0, 1) || 'U' }}</span>
              <div>
                <strong>{{ auth.user.username }}</strong>
                <small>{{ auth.user.role || 'member' }}</small>
              </div>
            </div>
            <button class="ghost-button small" @click="logout">退出</button>
          </template>
          <RouterLink v-else class="primary-button small" to="/login">登录</RouterLink>
        </div>
      </div>
    </header>

    <main class="page-container">
      <slot />
    </main>

    <footer class="site-footer">
      <div class="footer-inner">
        <div class="footer-brand">
          <span class="brand-badge">J</span>
          <span>JunBlog</span>
        </div>
        <span class="footer-copy">© 2026 Jun. 用简洁的设计，承载持续表达的内容。</span>
        <div class="footer-links">
          <RouterLink to="/articles">文章</RouterLink>
          <RouterLink to="/about">关于</RouterLink>
          <RouterLink v-if="auth.user?.role === 'admin'" to="/admin">管理后台</RouterLink>
          <RouterLink v-else to="/login">登录</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>