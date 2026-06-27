<script setup>
import { ref, onMounted } from 'vue'
import { adminGetComments, adminUpdateCommentStatus, adminDeleteComment } from '../../../shared/api/modules/comment'
import ConfirmModal from '../../../components/ConfirmModal.vue'
import { useMessage } from '../../../composables/useMessage'

const message = useMessage()
const comments = ref([])
const total = ref(0)
const page = ref(1)
const statusFilter = ref('')

// 确认弹窗
const showConfirm = ref(false)
const deleteId = ref(null)

async function load() {
  const { data } = await adminGetComments({ page: page.value, page_size: 10, status: statusFilter.value })
  comments.value = data.data.list || []
  total.value = data.data.total || 0
}

async function handleStatus(id, status) {
  try {
    await adminUpdateCommentStatus(id, status)
    message.success('操作成功')
    load()
  } catch (e) { message.error('操作失败') }
}

function handleDelete(id) {
  deleteId.value = id
  showConfirm.value = true
}

async function doDelete() {
  showConfirm.value = false
  try {
    await adminDeleteComment(deleteId.value)
    message.success('删除成功')
    load()
  } catch (e) { message.error('删除失败') }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-header"><h3>评论管理</h3>
      <select v-model="statusFilter" @change="load"><option value="">全部</option><option value="approved">已通过</option><option value="pending">待审核</option><option value="deleted">已删除</option></select>
    </div>
    <table class="data-table">
      <thead><tr><th>用户</th><th>内容</th><th>文章</th><th>状态</th><th>时间</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="c in comments" :key="c.id">
          <td>{{ c.user?.username }}</td>
          <td class="content-cell">{{ c.content }}</td>
          <td>{{ c.article?.title || '-' }}</td>
          <td><span class="status-tag" :class="c.status">{{ c.status }}</span></td>
          <td>{{ new Date(c.created_at).toLocaleDateString('zh-CN') }}</td>
          <td class="actions">
            <button v-if="c.status !== 'approved'" @click="handleStatus(c.id, 'approved')">通过</button>
            <button v-if="c.status !== 'pending'" @click="handleStatus(c.id, 'pending')">待审</button>
            <button class="danger" @click="handleDelete(c.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <ConfirmModal
      :show="showConfirm"
      title="确认删除"
      message="确认删除此评论？删除后不可恢复。"
      :danger="true"
      confirm-text="删除"
      @confirm="doDelete"
      @cancel="showConfirm = false"
    />
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h3 { margin: 0; }
.page-header select { padding: 6px 12px; border: 1px solid #ddd; border-radius: 6px; }
.data-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
.data-table th, .data-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #f0f0f0; font-size: 14px; }
.data-table th { background: #fafafa; font-weight: 600; }
.content-cell { max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.status-tag { padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.status-tag.approved { background: #e6f7e6; color: #52c41a; }
.status-tag.pending { background: #fff7e6; color: #faad14; }
.status-tag.deleted { background: #f0f0f0; color: #999; }
.actions { display: flex; gap: 8px; }
.actions button { padding: 4px 10px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; font-size: 13px; }
.actions button.danger { color: #ff4d4f; border-color: #ff4d4f; }
</style>
