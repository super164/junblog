<script setup>
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuth, clearAuthState } from '../../../shared/stores/auth'

const route = useRoute()
const router = useRouter()
const sidebarCollapsed = ref(false)

const auth = useAuth()

const menuItems = [
  { label: '控制台', icon: '📊', to: '/admin/dashboard' },
  { label: '文章管理', icon: '📝', to: '/admin/articles' },
  { label: '分类管理', icon: '📁', to: '/admin/categories' },
  { label: '标签管理', icon: '🏷️', to: '/admin/tags' },
  { label: '评论管理', icon: '💬', to: '/admin/comments' },
  { label: '用户管理', icon: '👥', to: '/admin/users' },
]

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function logout() {
  clearAuthState()
  router.push('/login')
}

function goToPublic() {
  router.push('/')
}
</script>

<template>
  <div class="admin-layout" :class="{ collapsed: sidebarCollapsed }">
    <!-- 侧边栏 -->
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <img class="brand-logo" src="/logo.png" alt="JunBlog" />
        <span v-if="!sidebarCollapsed" class="brand-text">JunBlog 管理</span>
        <button class="collapse-btn" @click="toggleSidebar">
          {{ sidebarCollapsed ? '▶' : '◀' }}
        </button>
      </div>

      <nav class="sidebar-nav">
        <RouterLink
          v-for="item in menuItems"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          :class="{ active: route.path.startsWith(item.to) }"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span v-if="!sidebarCollapsed" class="nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <button class="nav-item" @click="goToPublic">
          <span class="nav-icon">🌐</span>
          <span v-if="!sidebarCollapsed" class="nav-label">返回前台</span>
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="admin-main">
      <!-- 顶部栏 -->
      <header class="admin-topbar">
        <div class="topbar-left">
          <h2 class="page-title">{{ route.meta.title?.replace('JunBlog | ', '') || '管理后台' }}</h2>
        </div>
        <div class="topbar-right">
          <div class="user-info">
            <span class="user-avatar">{{ auth.user?.username?.slice(0, 1) || 'A' }}</span>
            <span class="user-name">{{ auth.user?.username }}</span>
          </div>
          <RouterLink to="/admin/profile" class="topbar-link">设置</RouterLink>
          <button class="topbar-link logout" @click="logout">退出</button>
        </div>
      </header>

      <!-- 页面内容 -->
      <main class="admin-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: #f8f7f4;
}

/* 侧边栏 */
.admin-sidebar {
  width: 240px;
  background: linear-gradient(180deg, #2c2825 0%, #1f1c1a 100%);
  color: #f5f0eb;
  display: flex;
  flex-direction: column;
  transition: width 0.35s cubic-bezier(0.22, 1, 0.36, 1);
  flex-shrink: 0;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.12);
}

.admin-layout.collapsed .admin-sidebar {
  width: 68px;
}

.sidebar-header {
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.brand-logo {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  object-fit: cover;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.brand-text {
  font-weight: 600;
  font-size: 14px;
  white-space: nowrap;
  letter-spacing: -0.01em;
}

.collapse-btn {
  margin-left: auto;
  background: none;
  border: none;
  color: #8a8578;
  cursor: pointer;
  font-size: 11px;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #f5f0eb;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 20px;
  color: #b8afa6;
  text-decoration: none;
  transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);
  cursor: pointer;
  border: none;
  background: none;
  width: 100%;
  font-size: 14px;
  margin: 2px 8px;
  border-radius: 8px;
  position: relative;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #f5f0eb;
}

.nav-item.active {
  background: rgba(194, 59, 34, 0.18);
  color: #e8734a;
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: #c23b22;
  border-radius: 0 3px 3px 0;
}

.nav-icon {
  font-size: 17px;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.nav-label {
  white-space: nowrap;
}

.sidebar-footer {
  padding: 16px 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

/* 主内容区 */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.admin-topbar {
  height: 60px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(138, 133, 120, 0.12);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  flex-shrink: 0;
}

.page-title {
  font-size: 15px;
  font-weight: 600;
  color: #2c2825;
  margin: 0;
  letter-spacing: -0.01em;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #c23b22 0%, #a52d18 100%);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(194, 59, 34, 0.25);
}

.user-name {
  font-size: 14px;
  color: #4a4540;
  font-weight: 500;
}

.topbar-link {
  font-size: 13px;
  color: #666;
  text-decoration: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
}

.topbar-link:hover {
  color: #c23b22;
  background: rgba(194, 59, 34, 0.06);
}

.topbar-link.logout {
  color: #999;
}

.admin-content {
  flex: 1;
  padding: 28px 32px;
  overflow-y: auto;
  background: #f8f7f4;
}
</style>
