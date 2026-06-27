<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="error-page">
    <div class="error-content">
      <div class="error-icon">📭</div>
      <h1>404</h1>
      <p class="error-desc">你访问的页面不存在</p>
      <p class="error-hint">可能页面已被移除，或者链接输入有误</p>
      <div class="error-actions">
        <RouterLink class="btn-primary" to="/">返回首页</RouterLink>
        <RouterLink class="btn-ghost" to="/articles">浏览文章</RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
.error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f8f7f4;
}

.error-content {
  text-align: center;
  padding: 48px;
}

.error-icon {
  font-size: 56px;
  margin-bottom: 16px;
  opacity: 0.8;
}

.error-content h1 {
  font-size: 80px;
  font-weight: 900;
  color: #1a1a1a;
  margin: 0;
  line-height: 1;
  letter-spacing: -0.04em;
  opacity: 0.12;
}

.error-desc {
  font-size: 20px;
  color: #1a1a1a;
  margin: 12px 0 8px;
  font-weight: 600;
}

.error-hint {
  font-size: 14px;
  color: #8a8578;
  margin: 0 0 32px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  padding: 12px 28px;
  background: #1a1a1a;
  color: #f5f0e8;
  border-radius: 10px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.25s;
}

.btn-primary:hover {
  background: #333;
  transform: translateY(-1px);
}

.btn-ghost {
  display: inline-flex;
  align-items: center;
  padding: 12px 28px;
  border: 1.5px solid rgba(26, 26, 26, 0.15);
  color: #4a4a4a;
  border-radius: 10px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.25s;
}

.btn-ghost:hover {
  border-color: #1a1a1a;
  color: #1a1a1a;
}
</style>
