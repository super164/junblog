<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { adminGetArticles, adminDeleteArticle, adminUpdateArticleStatus } from '../../../shared/api/modules/article'
import ConfirmModal from '../../../components/ConfirmModal.vue'
import { useMessage } from '../../../composables/useMessage'

const router = useRouter()
const message = useMessage()
const articles = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const statusFilter = ref('')
const loading = ref(false)

// 确认弹窗
const showConfirm = ref(false)
const deleteId = ref(null)

async function loadArticles() {
  loading.value = true
  try {
    const { data } = await adminGetArticles({ page: page.value, page_size: pageSize.value, keyword: keyword.value, status: statusFilter.value })
    articles.value = data.data.list || []
    total.value = data.data.total || 0
  } catch (e) {
    console.error('加载失败', e)
  } finally {
    loading.value = false
  }
}

function handleDelete(id) {
  deleteId.value = id
  showConfirm.value = true
}

async function doDelete() {
  showConfirm.value = false
  try {
    await adminDeleteArticle(deleteId.value)
    message.success('删除成功')
    loadArticles()
  } catch (e) {
    message.error('删除失败')
  }
}

async function toggleStatus(id, currentStatus) {
  const newStatus = currentStatus === 'published' ? 'draft' : 'published'
  try {
    await adminUpdateArticleStatus(id, newStatus)
    message.success('操作成功')
    loadArticles()
  } catch (e) {
    message.error('操作失败')
  }
}

onMounted(loadArticles)
</script>

<template>
  <div class="article-list-page">
    <div class="page-header">
      <h3>文章管理</h3>
      <button class="primary-btn" @click="router.push('/admin/articles/new')">+ 新建文章</button>
    </div>

    <div class="filters">
      <input v-model="keyword" placeholder="搜索文章标题..." @keyup.enter="loadArticles" />
      <select v-model="statusFilter" @change="loadArticles">
        <option value="">全部状态</option>
        <option value="published">已发布</option>
        <option value="draft">草稿</option>
        <option value="privacy">私密</option>
      </select>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>标题</th>
          <th>分类</th>
          <th>状态</th>
          <th>浏览</th>
          <th>创建时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="a in articles" :key="a.id">
          <td>{{ a.title }}</td>
          <td>{{ a.category?.name || '-' }}</td>
          <td>
            <span class="status-tag" :class="a.status">{{ a.status }}</span>
          </td>
          <td>{{ a.views_count }}</td>
          <td>{{ new Date(a.created_at).toLocaleDateString('zh-CN') }}</td>
          <td class="actions">
            <button @click="router.push(`/admin/articles/${a.id}/edit`)">编辑</button>
            <button @click="toggleStatus(a.id, a.status)">{{ a.status === 'published' ? '撤回' : '发布' }}</button>
            <button class="danger" @click="handleDelete(a.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-if="!loading && articles.length === 0" class="empty">暂无文章</div>

    <ConfirmModal
      :show="showConfirm"
      title="确认删除"
      message="确认删除此文章？删除后不可恢复。"
      :danger="true"
      confirm-text="删除"
      @confirm="doDelete"
      @cancel="showConfirm = false"
    />
  </div>
</template>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h3 { margin: 0; font-size: 18px; }
.primary-btn { background: #e94560; color: #fff; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.primary-btn:hover { background: #d63851; }
.filters { display: flex; gap: 12px; margin-bottom: 16px; }
.filters input, .filters select { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; }
.filters input { flex: 1; max-width: 300px; }
.data-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
.data-table th, .data-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #f0f0f0; font-size: 14px; }
.data-table th { background: #fafafa; font-weight: 600; color: #666; }
.status-tag { padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.status-tag.published { background: #e6f7e6; color: #52c41a; }
.status-tag.draft { background: #fff7e6; color: #faad14; }
.status-tag.privacy { background: #f0f0f0; color: #999; }
.actions { display: flex; gap: 8px; }
.actions button { padding: 4px 10px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; font-size: 13px; }
.actions button:hover { border-color: #e94560; color: #e94560; }
.actions button.danger { color: #ff4d4f; border-color: #ff4d4f; }
.loading, .empty { text-align: center; padding: 40px; color: #999; }
</style>
