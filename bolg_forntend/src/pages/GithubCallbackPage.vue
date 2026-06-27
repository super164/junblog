<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { githubLogin } from '../services/api'
import { saveAuthState } from '../shared/stores/auth'

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const errorMessage = ref('')

onMounted(async () => {
  const code = route.query.code
  const token = route.query.token
  const refreshToken = route.query.refresh_token
  const username = route.query.username
  const role = route.query.role
  const error = route.query.error

  // 处理错误情况
  if (error) {
    errorMessage.value = error === 'github_login_failed' ? 'GitHub 登录失败，请重试' : '授权失败'
    loading.value = false
    return
  }

  // 方式1：从 URL 参数直接获取 token（后端重定向方式）
  if (token && refreshToken) {
    try {
      const user = { username: username || '', role: role || 'user' }
      saveAuthState(token, refreshToken, user)

      // 根据角色跳转
      if (role === 'admin') {
        router.replace('/admin')
      } else {
        router.replace('/')
      }
    } catch (e) {
      errorMessage.value = '登录处理失败'
      loading.value = false
    }
    return
  }

  // 方式2：使用 code 调用后端 API（兼容旧方式）
  if (!code) {
    errorMessage.value = '未收到授权码，请重试'
    loading.value = false
    return
  }

  try {
    const result = await githubLogin(code)
    if (result.code !== 0) throw new Error(result.msg || 'GitHub 登录失败')

    saveAuthState(result.data.token, result.data.refresh_token, result.data.user)

    // 根据角色跳转
    if (result.data.user.role === 'admin') {
      router.replace('/admin')
    } else {
      router.replace('/')
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.msg || error.message || 'GitHub 登录失败'
    loading.value = false
  }
})
</script>

<template>
  <div class="callback-layout">
    <div class="callback-card">
      <div v-if="loading" class="loading-state">
        <div class="spinner-large"></div>
        <p>正在完成 GitHub 登录...</p>
      </div>

      <div v-else class="error-state">
        <span class="error-icon">⚠</span>
        <p>{{ errorMessage }}</p>
        <div class="actions">
          <RouterLink to="/login" class="back-link">返回登录页</RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.callback-layout {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
}

.callback-card {
  width: 100%;
  max-width: 360px;
  background: #fff;
  padding: 40px 36px;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06), 0 1px 4px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(26, 26, 26, 0.05);
  text-align: center;
}

.loading-state p {
  margin-top: 16px;
  font-size: 14px;
  color: #8a8578;
}

.spinner-large {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(26, 26, 26, 0.08);
  border-top-color: #c23b22;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.error-icon {
  font-size: 32px;
}

.error-state p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.actions {
  margin-top: 8px;
}

.back-link {
  display: inline-block;
  padding: 10px 24px;
  background: linear-gradient(135deg, #e94560 0%, #c23b22 100%);
  color: #fff;
  border-radius: 10px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
}

.back-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(233, 69, 96, 0.3);
}
</style>
