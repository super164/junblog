<script setup>
import { ref, onMounted } from 'vue'
import { adminGetTags, adminCreateTag, adminUpdateTag, adminDeleteTag } from '../../../shared/api/modules/tag'
import { adminGetCategories } from '../../../shared/api/modules/category'
import ConfirmModal from '../../../components/ConfirmModal.vue'
import { useMessage } from '../../../composables/useMessage'

const message = useMessage()
const tags = ref([])
const categories = ref([])
const form = ref({ name: '', slug: '', category_id: null })
const editId = ref(null)
const showForm = ref(false)

// 确认弹窗
const showConfirm = ref(false)
const deleteId = ref(null)
const confirmTitle = ref('确认删除')
const confirmMessage = ref('确认删除？')

async function load() {
  const [tagRes, catRes] = await Promise.all([
    adminGetTags(),
    adminGetCategories()
  ])
  console.log('Tags:', tagRes.data.data)
  console.log('Categories:', catRes.data.data)
  tags.value = tagRes.data.data || []
  categories.value = catRes.data.data || []
}

function openCreate() {
  form.value = { name: '', slug: '', category_id: null }
  editId.value = null
  showForm.value = true
}

function openEdit(t) {
  form.value = { name: t.name, slug: t.slug, category_id: t.category_id || null }
  editId.value = t.id
  showForm.value = true
}

function getCategoryName(id) {
  if (!id) return '未分类'
  const numId = Number(id)
  const cat = categories.value.find(c => c.id === numId || c.id === id)
  return cat ? cat.name : '未分类'
}

async function handleSubmit() {
  try {
    if (editId.value) await adminUpdateTag(editId.value, form.value)
    else await adminCreateTag(form.value)
    message.success('保存成功')
    showForm.value = false
    load()
  } catch (e) { message.error('操作失败') }
}

function handleDelete(id) {
  deleteId.value = id
  confirmTitle.value = '确认删除'
  confirmMessage.value = '确认删除此标签？删除后不可恢复。'
  showConfirm.value = true
}

async function doDelete() {
  showConfirm.value = false
  try {
    await adminDeleteTag(deleteId.value)
    message.success('删除成功')
    load()
  } catch (e) { message.error('删除失败') }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-header"><h3>标签管理</h3><button class="primary-btn" @click="openCreate">+ 新建标签</button></div>
    <table class="data-table">
      <thead><tr><th>名称</th><th>Slug</th><th>所属分类</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="t in tags" :key="t.id">
          <td>{{ t.name }}</td>
          <td>{{ t.slug }}</td>
          <td>{{ getCategoryName(t.category_id) }}</td>
          <td class="actions">
            <button @click="openEdit(t)">编辑</button>
            <button class="danger" @click="handleDelete(t.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="showForm" class="modal-overlay" @click.self="showForm=false">
      <div class="modal">
        <h4>{{ editId ? '编辑标签' : '新建标签' }}</h4>
        <form @submit.prevent="handleSubmit">
          <label>名称 <input v-model="form.name" required /></label>
          <label>Slug <input v-model="form.slug" required /></label>
          <label>所属分类
            <select v-model="form.category_id" required>
              <option :value="null" disabled>请选择分类</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </label>
          <div class="form-actions">
            <button type="submit" class="primary-btn">保存</button>
            <button type="button" @click="showForm=false">取消</button>
          </div>
        </form>
      </div>
    </div>

    <ConfirmModal
      :show="showConfirm"
      :title="confirmTitle"
      :message="confirmMessage"
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
.primary-btn { background: #e94560; color: #fff; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.data-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
.data-table th, .data-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #f0f0f0; font-size: 14px; }
.data-table th { background: #fafafa; font-weight: 600; }
.actions { display: flex; gap: 8px; }
.actions button { padding: 4px 10px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; }
.actions button.danger { color: #ff4d4f; border-color: #ff4d4f; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: #fff; padding: 24px; border-radius: 12px; width: 400px; }
.modal h4 { margin: 0 0 16px; }
.modal label { display: block; margin-bottom: 12px; font-size: 14px; font-weight: 600; }
.modal input, .modal select { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; margin-top: 4px; box-sizing: border-box; }
.modal select { background: #fff; cursor: pointer; }
.form-actions { display: flex; gap: 12px; margin-top: 16px; }
.form-actions button[type="button"] { padding: 8px 16px; border: 1px solid #ddd; background: #fff; border-radius: 6px; cursor: pointer; }
</style>
