<script setup>
import { onMounted, reactive, ref } from 'vue'
import SiteShell from '../components/SiteShell.vue'
import { getProfile, updatePassword, updateProfile } from '../services/api'
import { useAuth, saveAuthState } from '../shared/stores/auth'

const auth = useAuth()

const profileForm = reactive({
  email: '',
  avatar: '',
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
})

const currentUser = ref(null)
const loadingProfile = ref(false)
const savingProfile = ref(false)
const savingPassword = ref(false)
const profileMessage = ref('')
const profileError = ref('')
const passwordMessage = ref('')
const passwordError = ref('')

async function loadProfile() {
  loadingProfile.value = true
  profileError.value = ''

  try {
    const result = await getProfile()
    if (result.code !== 0) {
      throw new Error(result.msg || '获取个人信息失败')
    }

    currentUser.value = result.data
    profileForm.email = result.data.email || ''
    profileForm.avatar = result.data.avatar || ''
  } catch (error) {
    profileError.value =
      error.response?.data?.msg || error.message || '获取个人信息失败，请检查登录状态'
  } finally {
    loadingProfile.value = false
  }
}

async function handleProfileSubmit() {
  savingProfile.value = true
  profileMessage.value = ''
  profileError.value = ''

  try {
    const result = await updateProfile(profileForm)
    if (result.code !== 0) {
      throw new Error(result.msg || '保存失败')
    }

    currentUser.value = result.data
    saveAuthState(auth.token, localStorage.getItem('junblog_refresh_token') || '', result.data)
    profileMessage.value = '个人资料已更新'
  } catch (error) {
    profileError.value = error.response?.data?.msg || error.message || '保存失败'
  } finally {
    savingProfile.value = false
  }
}

async function handlePasswordSubmit() {
  savingPassword.value = true
  passwordMessage.value = ''
  passwordError.value = ''

  try {
    const result = await updatePassword(passwordForm)
    if (result.code !== 0) {
      throw new Error(result.msg || '密码修改失败')
    }

    passwordMessage.value = result.data?.message || '密码修改成功'
    passwordForm.old_password = ''
    passwordForm.new_password = ''
  } catch (error) {
    passwordError.value = error.response?.data?.msg || error.message || '密码修改失败'
  } finally {
    savingPassword.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <SiteShell>
    <section class="page-hero anim-fade-up">
      <p class="section-kicker">Profile Settings</p>
      <h1>个人设置</h1>
      <p>这个页面专门用于联调你后端的个人资料和密码修改接口。先登录，再进入这里即可直接测试。</p>
    </section>

    <section class="profile-grid">
      <article class="content-card anim-fade-up anim-delay-2">
        <div class="section-header compact">
          <div>
            <p class="section-kicker">Current User</p>
            <h2>当前资料</h2>
          </div>
          <button class="ghost-button" :disabled="loadingProfile" @click="loadProfile">
            {{ loadingProfile ? '加载中...' : '刷新资料' }}
          </button>
        </div>

        <div class="profile-summary">
          <div class="profile-avatar">
            {{ currentUser?.username?.slice(0, 1) || 'U' }}
          </div>
          <div>
            <strong>{{ currentUser?.username || '未获取到用户名' }}</strong>
            <p class="muted-text">角色：{{ currentUser?.role || '未知' }}</p>
            <p class="muted-text">邮箱：{{ currentUser?.email || '暂无' }}</p>
            <p class="muted-text">头像：{{ currentUser?.avatar || '暂无' }}</p>
          </div>
        </div>

        <p v-if="profileError && loadingProfile" class="error-text">{{ profileError }}</p>
      </article>

      <article class="content-card anim-fade-up anim-delay-3">
        <p class="section-kicker">Update Profile</p>
        <h2>修改资料</h2>
        <form class="settings-form" @submit.prevent="handleProfileSubmit">
          <label>
            邮箱
            <input v-model.trim="profileForm.email" type="email" placeholder="请输入新邮箱" />
          </label>
          <label>
            头像地址
            <input
              v-model.trim="profileForm.avatar"
              type="text"
              placeholder="/uploads/avatar/a.png"
            />
          </label>
          <button class="primary-button" :disabled="savingProfile" type="submit">
            {{ savingProfile ? '保存中...' : '保存资料' }}
          </button>
        </form>

        <p v-if="profileMessage" class="success-text">{{ profileMessage }}</p>
        <p v-if="profileError && !loadingProfile" class="error-text">{{ profileError }}</p>
      </article>
    </section>

    <section class="profile-grid single">
      <article class="content-card anim-fade-up anim-delay-4">
        <p class="section-kicker">Update Password</p>
        <h2>修改密码</h2>
        <form class="settings-form" @submit.prevent="handlePasswordSubmit">
          <label>
            旧密码
            <input
              v-model.trim="passwordForm.old_password"
              type="password"
              placeholder="请输入旧密码"
            />
          </label>
          <label>
            新密码
            <input
              v-model.trim="passwordForm.new_password"
              type="password"
              placeholder="请输入新密码"
            />
          </label>
          <button class="primary-button" :disabled="savingPassword" type="submit">
            {{ savingPassword ? '提交中...' : '修改密码' }}
          </button>
        </form>

        <p v-if="passwordMessage" class="success-text">{{ passwordMessage }}</p>
        <p v-if="passwordError" class="error-text">{{ passwordError }}</p>
      </article>
    </section>
  </SiteShell>
</template>
