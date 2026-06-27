<script setup>
import { ref, onMounted, reactive } from 'vue'
import http from '../../../services/api'
import { useMessage } from '../../../composables/useMessage'

const message = useMessage()
const users = ref([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const loading = ref(false)

// 编辑弹窗
const showEditModal = ref(false)
const editForm = reactive({ id: 0, username: '', email: '', phone: '', avatar: '', role: '' })
const editLoading = ref(false)

// 重置密码弹窗
const showResetModal = ref(false)
const resetForm = reactive({ id: 0, username: '', new_password: '' })
const resetLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await http.get('/admin/users', { params: { page: page.value, page_size: 10, keyword: keyword.value } })
    users.value = data.data.list || []
    total.value = data.data.total || 0
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function toggleStatus(id, currentStatus) {
  try {
    await http.patch(`/admin/users/${id}/status`, { status: !currentStatus })
    message.success('操作成功')
    load()
  } catch (e) { message.error('操作失败') }
}

// 打开编辑弹窗
async function openEdit(user) {
  editForm.id = user.id
  editForm.username = user.username
  editForm.email = user.email || ''
  editForm.phone = user.phone || ''
  editForm.avatar = user.avatar || ''
  editForm.role = user.role
  showEditModal.value = true
}

// 保存编辑
async function handleEdit() {
  editLoading.value = true
  try {
    await http.put(`/admin/users/${editForm.id}`, {
      email: editForm.email,
      phone: editForm.phone,
      avatar: editForm.avatar,
      role: editForm.role,
    })
    message.success('更新成功')
    showEditModal.value = false
    load()
  } catch (e) { message.error('更新失败: ' + (e.response?.data?.msg || e.message)) }
  finally { editLoading.value = false }
}

// 打开重置密码弹窗
function openReset(user) {
  resetForm.id = user.id
  resetForm.username = user.username
  resetForm.new_password = ''
  showResetModal.value = true
}

// 重置密码
async function handleReset() {
  if (resetForm.new_password.length < 6) { message.warning('密码至少6位'); return }
  resetLoading.value = true
  try {
    await http.patch(`/admin/users/${resetForm.id}/password`, { new_password: resetForm.new_password })
    message.success('密码重置成功')
    showResetModal.value = false
  } catch (e) { message.error('重置失败: ' + (e.response?.data?.msg || e.message)) }
  finally { resetLoading.value = false }
}

onMounted(load)
</script>

<template>
  <div class="user-manage">
    <div class="page-header">
      <h3>用户管理</h3>
      <div class="header-actions">
        <input v-model="keyword" placeholder="搜索用户名或邮箱..." @keyup.enter="load" class="search-input" />
        <button class="search-btn" @click="load">搜索</button>
      </div>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>用户名</th>
          <th>邮箱</th>
          <th>角色</th>
          <th>状态</th>
          <th>注册时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="u in users" :key="u.id">
          <td>
            <div class="user-cell">
              <span class="avatar-sm">{{ u.username?.slice(0, 1) || 'U' }}</span>
              <span>{{ u.username }}</span>
            </div>
          </td>
          <td>{{ u.email || '-' }}</td>
          <td><span class="role-tag" :class="u.role">{{ u.role }}</span></td>
          <td>
            <span class="status-dot" :class="{ active: u.status }"></span>
            {{ u.status ? '正常' : '禁用' }}
          </td>
          <td>{{ new Date(u.created_at).toLocaleDateString('zh-CN') }}</td>
          <td class="actions">
            <button @click="openEdit(u)">编辑</button>
            <button @click="openReset(u)">重置密码</button>
            <button v-if="u.role !== 'admin'" :class="u.status ? 'warning' : 'success'" @click="toggleStatus(u.id, u.status)">
              {{ u.status ? '禁用' : '启用' }}
            </button>
            <span v-else class="admin-badge">管理员</span>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-if="!loading && users.length === 0" class="empty">暂无用户</div>

    <!-- 分页 -->
    <div v-if="total > 10" class="pagination">
      <button :disabled="page <= 1" @click="page--; load()">上一页</button>
      <span>第 {{ page }} / {{ Math.ceil(total / 10) }} 页</span>
      <button :disabled="page >= Math.ceil(total / 10)" @click="page++; load()">下一页</button>
    </div>

    <!-- 编辑用户弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal=false">
      <div class="modal">
        <h4>编辑用户：{{ editForm.username }}</h4>
        <form @submit.prevent="handleEdit">
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="editForm.email" type="email" placeholder="user@example.com" />
          </div>
          <div class="form-group">
            <label>手机号</label>
            <input v-model="editForm.phone" placeholder="可选" />
          </div>
          <div class="form-group">
            <label>头像 URL</label>
            <input v-model="editForm.avatar" placeholder="可选" />
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="editForm.role">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="submit" class="primary-btn" :disabled="editLoading">
              {{ editLoading ? '保存中...' : '保存' }}
            </button>
            <button type="button" class="cancel-btn" @click="showEditModal=false">取消</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 重置密码弹窗 -->
    <div v-if="showResetModal" class="modal-overlay" @click.self="showResetModal=false">
      <div class="modal modal-sm">
        <h4>重置密码：{{ resetForm.username }}</h4>
        <form @submit.prevent="handleReset">
          <div class="form-group">
            <label>新密码</label>
            <input v-model="resetForm.new_password" type="password" placeholder="至少6位" required />
          </div>
          <div class="modal-actions">
            <button type="submit" class="primary-btn" :disabled="resetLoading">
              {{ resetLoading ? '重置中...' : '确认重置' }}
            </button>
            <button type="button" class="cancel-btn" @click="showResetModal=false">取消</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h3 { margin: 0; font-size: 18px; }
.header-actions { display: flex; gap: 8px; }
.search-input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; width: 220px; font-size: 14px; }
.search-btn { padding: 8px 16px; background: #e94560; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; }

.data-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
.data-table th, .data-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #f0f0f0; font-size: 14px; }
.data-table th { background: #fafafa; font-weight: 600; color: #666; font-size: 13px; }

.user-cell { display: flex; align-items: center; gap: 8px; }
.avatar-sm { width: 28px; height: 28px; background: #e94560; color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; flex-shrink: 0; }

.role-tag { padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.role-tag.admin { background: #fff0f0; color: #e94560; }
.role-tag.user { background: #f0f0f0; color: #666; }

.status-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #ccc; margin-right: 4px; }
.status-dot.active { background: #52c41a; }

.actions { display: flex; gap: 6px; flex-wrap: wrap; }
.actions button { padding: 4px 10px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; font-size: 12px; transition: all 0.2s; }
.actions button:hover { border-color: #e94560; color: #e94560; }
.actions button.warning { border-color: #faad14; color: #faad14; }
.actions button.warning:hover { background: #fff7e6; }
.actions button.success { border-color: #52c41a; color: #52c41a; }
.actions button.success:hover { background: #f6ffed; }
.admin-badge { padding: 4px 10px; background: #f0f0f0; color: #999; border-radius: 4px; font-size: 12px; }

.pagination { display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 20px; }
.pagination button { padding: 6px 14px; border: 1px solid #ddd; background: #fff; border-radius: 6px; cursor: pointer; font-size: 13px; }
.pagination button:disabled { opacity: 0.4; cursor: not-allowed; }
.pagination span { font-size: 13px; color: #666; }

.loading, .empty { text-align: center; padding: 32px; color: #999; }

/* 弹窗 */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: #fff; padding: 28px; border-radius: 12px; width: 460px; max-width: 90vw; box-shadow: 0 8px 32px rgba(0,0,0,0.15); }
.modal.modal-sm { width: 360px; }
.modal h4 { margin: 0 0 20px; font-size: 16px; color: #333; }

.form-group { margin-bottom: 16px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: #333; margin-bottom: 4px; }
.form-group input, .form-group select { width: 100%; padding: 9px 12px; border: 1.5px solid #e0e0e0; border-radius: 8px; font-size: 14px; box-sizing: border-box; transition: border-color 0.2s; }
.form-group input:focus, .form-group select:focus { outline: none; border-color: #e94560; }

.modal-actions { display: flex; gap: 12px; margin-top: 20px; }
.primary-btn { flex: 1; padding: 10px; background: #e94560; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
.primary-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.cancel-btn { padding: 10px 20px; border: 1px solid #ddd; background: #fff; border-radius: 8px; font-size: 14px; cursor: pointer; }
.cancel-btn:hover { border-color: #999; }
</style>
