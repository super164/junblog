<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import http from '../../../services/api'
import { useAuth, saveAuthState } from '../../../shared/stores/auth'

const router = useRouter()
const auth = useAuth()
const form = reactive({ username: '', email: '', password: '', confirm_password: '' })
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showPassword = ref(false)
const fieldErrors = reactive({ username: '', email: '', password: '', confirm_password: '' })

// 已登录则跳转首页
onMounted(() => {
  if (auth.token) {
    router.replace('/')
  }
})

function validate() {
  Object.keys(fieldErrors).forEach(k => fieldErrors[k] = '')
  if (!form.username) { fieldErrors.username = '请输入用户名'; return false }
  if (form.username.length < 4 || form.username.length > 10) { fieldErrors.username = '用户名需要 4-10 位'; return false }
  if (!form.email) { fieldErrors.email = '请输入邮箱'; return false }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) { fieldErrors.email = '邮箱格式不正确'; return false }
  if (!form.password) { fieldErrors.password = '请输入密码'; return false }
  if (form.password.length < 6) { fieldErrors.password = '密码至少6位'; return false }
  if (form.password !== form.confirm_password) { fieldErrors.confirm_password = '两次密码不一致'; return false }
  return true
}

async function handleSubmit() {
  errorMessage.value = ''
  successMessage.value = ''
  if (!validate()) return

  loading.value = true
  try {
    const { data } = await http.post('/auth/register', form)
    if (data.code !== 0) throw new Error(data.msg || '注册失败')

    saveAuthState(data.data.token, data.data.refresh_token, data.data.user)
    successMessage.value = '注册成功，正在跳转...'
    setTimeout(() => router.push('/'), 800)
  } catch (error) {
    errorMessage.value = error.response?.data?.msg || error.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-layout">
    <RouterLink to="/" class="back-home">← 返回首页</RouterLink>

    <div class="login-card">
      <p class="section-kicker">Sign Up</p>
      <h2>注册</h2>
      <p class="muted-text">注册一个账号，开始你的博客之旅</p>

      <form class="login-form" @submit.prevent="handleSubmit">
        <transition name="fade">
          <div v-if="errorMessage" class="form-alert error">
            <span class="alert-icon">⚠</span>
            <span>{{ errorMessage }}</span>
          </div>
        </transition>

        <div class="field" :class="{ error: fieldErrors.username }">
          <label>用户名</label>
          <input v-model.trim="form.username" type="text" placeholder="4-10位字符" @input="fieldErrors.username=''" />
          <span class="field-error" v-if="fieldErrors.username">{{ fieldErrors.username }}</span>
        </div>

        <div class="field" :class="{ error: fieldErrors.email }">
          <label>邮箱</label>
          <input v-model.trim="form.email" type="email" placeholder="your@email.com" @input="fieldErrors.email=''" />
          <span class="field-error" v-if="fieldErrors.email">{{ fieldErrors.email }}</span>
        </div>

        <div class="field" :class="{ error: fieldErrors.password }">
          <label>密码</label>
          <div class="password-wrap">
            <input v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="至少6位" @input="fieldErrors.password=''" />
            <button type="button" class="eye-btn" @click="showPassword = !showPassword">
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>
          <span class="field-error" v-if="fieldErrors.password">{{ fieldErrors.password }}</span>
        </div>

        <div class="field" :class="{ error: fieldErrors.confirm_password }">
          <label>确认密码</label>
          <div class="password-wrap">
            <input v-model="form.confirm_password" :type="showPassword ? 'text' : 'password'" placeholder="再次输入密码" @input="fieldErrors.confirm_password=''" />
            <button type="button" class="eye-btn" @click="showPassword = !showPassword">
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>
          <span class="field-error" v-if="fieldErrors.confirm_password">{{ fieldErrors.confirm_password }}</span>
        </div>

        <button class="submit-btn" :disabled="loading" type="submit">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? '注册中...' : '立即注册' }}
        </button>
      </form>

      <div class="form-footer">
        <span>已有账号？</span>
        <RouterLink to="/login" class="link">去登录</RouterLink>
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
