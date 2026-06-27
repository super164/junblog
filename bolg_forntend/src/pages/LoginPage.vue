<script setup>
import { computed, reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../services/api'
import { useAuth, saveAuthState, clearAuthState } from '../shared/stores/auth'

const router = useRouter()
const auth = useAuth()
const form = reactive({ username: '', password: '' })
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showPassword = ref(false)
const fieldErrors = reactive({ username: '', password: '' })

// GitHub OAuth 跳转地址（从后端动态获取，支持 ngrok/公网地址）
const githubAuthUrl = ref('#')

// 页面加载时获取正确的 GitHub 授权地址
onMounted(async () => {
  try {
    const origin = window.location.origin
    const { data } = await fetch(`/api/v1/auth/github/url?origin=${encodeURIComponent(origin)}`).then(r => r.json())
    if (data?.url) {
      githubAuthUrl.value = data.url
    }
  } catch (e) {
    console.error('获取 GitHub 授权地址失败', e)
  }
})

// 已登录则跳转首页
onMounted(() => {
  if (auth.token) {
    router.replace('/')
  }

  // 检查 URL 中的错误参数
  const urlParams = new URLSearchParams(window.location.search)
  const error = urlParams.get('error')
  if (error) {
    const errorMessages = {
      'github_login_failed': 'GitHub 登录失败，请重试',
      'user_disabled': '账号已被封禁，请联系管理员',
      'missing_code': '授权码丢失，请重试',
    }
    errorMessage.value = errorMessages[error] || '登录失败，请重试'
    // 清除 URL 中的 error 参数
    window.history.replaceState({}, document.title, window.location.pathname)
  }
})

const demoHint = computed(() => {
  return auth.user ? `当前已登录：${auth.user.username}` : '输入你的账号密码登录'
})

function validate() {
  fieldErrors.username = ''
  fieldErrors.password = ''
  if (!form.username) { fieldErrors.username = '请输入用户名'; return false }
  if (form.username.length < 2) { fieldErrors.username = '用户名至少2位'; return false }
  if (!form.password) { fieldErrors.password = '请输入密码'; return false }
  if (form.password.length < 6) { fieldErrors.password = '密码至少6位'; return false }
  return true
}

async function handleSubmit() {
  errorMessage.value = ''
  successMessage.value = ''
  if (!validate()) return

  loading.value = true
  try {
    const result = await login(form)
    if (result.code !== 0) throw new Error(result.msg || '登录失败')

    saveAuthState(result.data.token, result.data.refresh_token, result.data.user)
    successMessage.value = '登录成功，正在跳转...'
    setTimeout(() => {
      router.push(result.data.user.role === 'admin' ? '/admin' : '/')
    }, 600)
  } catch (error) {
    clearAuthState()
    errorMessage.value = error.response?.data?.msg || error.message || '请求失败，请检查后端服务是否已启动'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-layout">
    <RouterLink to="/" class="back-home">← 返回首页</RouterLink>

    <div class="login-card">
      <p class="section-kicker">Sign In</p>
      <h2>登录</h2>
      <p class="muted-text">{{ demoHint }}</p>

      <form class="login-form" @submit.prevent="handleSubmit">
        <transition name="fade">
          <div v-if="errorMessage" class="form-alert error">
            <span class="alert-icon">⚠</span>
            <span>{{ errorMessage }}</span>
          </div>
        </transition>

        <div class="field" :class="{ error: fieldErrors.username }">
          <label>用户名</label>
          <input v-model.trim="form.username" type="text" placeholder="请输入用户名" @input="fieldErrors.username=''" />
          <span class="field-error" v-if="fieldErrors.username">{{ fieldErrors.username }}</span>
        </div>

        <div class="field" :class="{ error: fieldErrors.password }">
          <label>密码</label>
          <div class="password-wrap">
            <input v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="请输入密码" @input="fieldErrors.password=''" />
            <button type="button" class="eye-btn" @click="showPassword = !showPassword">
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>
          <span class="field-error" v-if="fieldErrors.password">{{ fieldErrors.password }}</span>
        </div>

        <button class="submit-btn" :disabled="loading" type="submit">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? '登录中...' : '立即登录' }}
        </button>
      </form>

      <div class="divider">
        <span>或者</span>
      </div>

      <a class="github-btn" :href="githubAuthUrl">
        <svg viewBox="0 0 16 16" width="18" height="18" fill="currentColor">
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
        </svg>
        使用 GitHub 登录
      </a>

      <div class="form-footer">
        <span>还没有账号？</span>
        <RouterLink to="/register" class="link">立即注册</RouterLink>
      </div>

      <transition name="fade">
        <div v-if="successMessage" class="toast success">{{ successMessage }}</div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
.login-layout {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 48px 24px;
  position: relative;
}

.back-home {
  position: absolute;
  top: 24px;
  left: 24px;
  color: #8a8578;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.back-home:hover {
  color: #1a1a1a;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  padding: 36px;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06), 0 1px 4px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(26, 26, 26, 0.05);
}

.section-kicker {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 2.5px;
  color: #c23b22;
  margin: 0 0 4px;
  font-weight: 500;
}

.login-card h2 {
  margin: 0 0 4px;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.muted-text {
  font-size: 13px;
  color: #8a8578;
  margin: 0 0 24px;
}

.field {
  margin-bottom: 18px;
}

.field label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 6px;
}

.field input {
  width: 100%;
  padding: 12px 14px;
  border: 1.5px solid rgba(26, 26, 26, 0.12);
  border-radius: 10px;
  font-size: 14px;
  box-sizing: border-box;
  background: #faf7f2;
  transition: all 0.25s;
}

.field input:focus {
  outline: none;
  border-color: #1a1a1a;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.04);
}

.field.error input {
  border-color: #ff4d4f;
}

.field-error {
  display: block;
  font-size: 12px;
  color: #ff4d4f;
  margin-top: 5px;
}

.password-wrap {
  position: relative;
}

.password-wrap input {
  padding-right: 44px;
}

.eye-btn {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  opacity: 0.5;
  transition: opacity 0.2s;
}

.eye-btn:hover {
  opacity: 1;
}

.submit-btn {
  width: 100%;
  padding: 13px;
  background: linear-gradient(135deg, #e94560 0%, #c23b22 100%);
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s;
  margin-top: 4px;
}

.submit-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(233, 69, 96, 0.3);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.form-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #8a8578;
}

.form-footer .link {
  color: #c23b22;
  text-decoration: none;
  margin-left: 4px;
  font-weight: 500;
}

.form-footer .link:hover {
  text-decoration: underline;
}

.divider {
  display: flex;
  align-items: center;
  margin: 20px 0;
  gap: 12px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: rgba(26, 26, 26, 0.08);
}

.divider span {
  font-size: 12px;
  color: #8a8578;
  white-space: nowrap;
}

.github-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  border: 1.5px solid rgba(26, 26, 26, 0.12);
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  background: #fff;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.25s;
}

.github-btn:hover {
  border-color: #333;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.toast {
  text-align: center;
  font-size: 13px;
  margin-top: 12px;
  padding: 10px 14px;
  border-radius: 8px;
}

.toast.success {
  background: rgba(16, 185, 129, 0.08);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.form-alert {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 13px;
  margin-bottom: 4px;
}

.form-alert.error {
  background: rgba(255, 77, 79, 0.06);
  color: #ff4d4f;
  border: 1px solid rgba(255, 77, 79, 0.15);
}

.alert-icon {
  font-size: 15px;
  flex-shrink: 0;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

@media (max-width: 768px) {
  .login-layout {
    padding: 24px 16px;
  }

  .back-home {
    position: relative;
    top: auto;
    left: auto;
    margin-bottom: 16px;
    align-self: flex-start;
  }
}
</style>
