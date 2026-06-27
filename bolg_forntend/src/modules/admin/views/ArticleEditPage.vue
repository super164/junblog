<script setup>
import { ref, onMounted, computed, defineAsyncComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminGetArticle, adminCreateArticle, adminUpdateArticle, adminUploadFile } from '../../../shared/api/modules/article'
import { adminGetCategories } from '../../../shared/api/modules/category'
import { adminGetTags, adminCreateTag } from '../../../shared/api/modules/tag'
import 'md-editor-v3/lib/style.css'
const MdEditor = defineAsyncComponent(() => import('md-editor-v3').then(m => m.MdEditor))
import ImageCropModal from '../../../components/ImageCropModal.vue'
import { useMessage } from '../../../composables/useMessage'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const isEdit = ref(false)
const articleId = ref(null)

const form = ref({
  title: '', content: '', cover: '', category_id: '', tag_ids: [], status: 'published'
})
const categories = ref([])
const tags = ref([])
const loading = ref(false)
const uploading = ref(false)
const coverPreview = ref('')
const showTagInput = ref(false)
const newTagName = ref('')
const tagCreating = ref(false)
const mdEditorRef = ref(null)

// 裁剪相关
const showCropModal = ref(false)
const cropImageSrc = ref('')

async function loadData() {
  const [catRes, tagRes] = await Promise.all([adminGetCategories(), adminGetTags()])
  categories.value = catRes.data.data || []
  tags.value = tagRes.data.data || []

  if (route.params.id) {
    isEdit.value = true
    articleId.value = route.params.id
    const { data } = await adminGetArticle(articleId.value)
    const a = data.data
    // 修复内容中的路径（反斜杠转正斜杠，相对路径转完整URL）
    let content = a.content || ''
    content = content.replace(/\\/g, '/')
    content = content.replace(/!\[\]\((\/uploads\/[^)]+)\)/g, (match, path) => {
      return `![](http://localhost:8080${path})`
    })
    form.value = {
      title: a.title,
      content: content,
      cover: a.cover || '',
      category_id: a.category?.id || '',
      tag_ids: (a.tags || []).map(t => t.id),
      status: a.status
    }
    // 修复封面图路径
    let cover = a.cover || ''
    if (cover.startsWith('/uploads')) {
      cover = `http://localhost:8080${cover}`
    }
    coverPreview.value = cover
    form.value.cover = cover
  }
}

async function handleUploadCover(e) {
  const file = e.target.files?.[0]
  if (!file) return
  // 读取图片并打开裁剪弹窗
  const reader = new FileReader()
  reader.onload = (ev) => {
    cropImageSrc.value = ev.target.result
    showCropModal.value = true
  }
  reader.readAsDataURL(file)
  // 清空 input 允许重复选择同一文件
  e.target.value = ''
}

async function handleCropConfirm(blob) {
  showCropModal.value = false
  uploading.value = true
  try {
    const file = new File([blob], 'cover.jpg', { type: 'image/jpeg' })
    const { data } = await adminUploadFile(file)
    // 将相对路径转换为完整的后端 URL
    let url = data.data.url
    if (url.startsWith('/uploads')) {
      url = `http://localhost:8080${url}`
    }
    form.value.cover = url
    coverPreview.value = url
  } catch (err) {
    message.error('上传失败: ' + (err.response?.data?.msg || err.message))
  } finally {
    uploading.value = false
    cropImageSrc.value = ''
  }
}

function handleCropCancel() {
  showCropModal.value = false
  cropImageSrc.value = ''
}

function removeCover() {
  form.value.cover = ''
  coverPreview.value = ''
}

// Markdown 编辑器图片上传
async function handleMdEditorUploadImage(files, callback) {
  try {
    const results = await Promise.all(
      files.map(async (file) => {
        const { data } = await adminUploadFile(file)
        // 将相对路径转换为完整的后端 URL
        const url = data.data.url
        if (url.startsWith('/uploads')) {
          return `http://localhost:8080${url}`
        }
        return url
      })
    )
    callback(results)
  } catch (err) {
    message.error('图片上传失败: ' + (err.response?.data?.msg || err.message))
    callback([])
  }
}

async function handleCreateTag() {
  const name = newTagName.value.trim()
  if (!name) return
  tagCreating.value = true
  try {
    const slug = name.toLowerCase().replace(/[^a-z0-9一-龥]+/g, '-').replace(/^-+|-+$/g, '')
    const { data } = await adminCreateTag({ name, slug: slug || name.toLowerCase() })
    const newTag = data.data
    tags.value.push(newTag)
    form.value.tag_ids.push(newTag.id)
    newTagName.value = ''
    showTagInput.value = false
  } catch (err) {
    message.error('创建标签失败: ' + (err.response?.data?.msg || err.message))
  } finally {
    tagCreating.value = false
  }
}

async function handleSubmit() {
  if (!form.value.cover) {
    message.warning('请上传封面图')
    return
  }
  if (form.value.tag_ids.length === 0) {
    message.warning('请至少选择一个标签')
    return
  }
  loading.value = true
  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await adminUpdateArticle(articleId.value, payload)
    } else {
      await adminCreateArticle(payload)
    }
    message.success('保存成功')
    router.push('/admin/articles')
  } catch (e) {
    message.error('保存失败: ' + (e.response?.data?.msg || e.message))
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <div class="article-edit">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑文章' : '新建文章' }}</h2>
      <p class="page-desc">{{ isEdit ? '修改文章内容和设置' : '创建一篇新的博客文章' }}</p>
    </div>

    <form @submit.prevent="handleSubmit" class="edit-form">
      <!-- 标题 -->
      <div class="form-group">
        <label class="form-label">标题</label>
        <input v-model="form.title" class="form-input" placeholder="输入文章标题" required />
      </div>

      <!-- 分类 + 状态 -->
      <div class="form-row">
        <div class="form-group flex-1">
          <label class="form-label">分类</label>
          <select v-model="form.category_id" class="form-input" required>
            <option value="">请选择分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="form-group flex-1">
          <label class="form-label">状态</label>
          <select v-model="form.status" class="form-input">
            <option value="draft">草稿</option>
            <option value="published">发布</option>
            <option value="privacy">私密</option>
          </select>
        </div>
      </div>

      <!-- 封面图上传 -->
      <div class="form-group">
        <label class="form-label">封面图 <span class="label-required">*</span></label>
        <div v-if="coverPreview" class="cover-preview">
          <img :src="coverPreview" alt="封面预览" />
          <button type="button" class="cover-remove" @click="removeCover">✕</button>
        </div>
        <label v-else class="cover-upload" :class="{ uploading }">
          <input type="file" accept="image/*" @change="handleUploadCover" hidden />
          <span v-if="uploading" class="upload-spinner"></span>
          <span v-else class="upload-icon">📷</span>
          <span>{{ uploading ? '上传中...' : '点击上传封面图' }}</span>
          <span class="upload-hint">支持 jpg/png/gif/webp，最大 10MB</span>
        </label>
      </div>

      <!-- 标签 -->
      <div class="form-group">
        <label class="form-label">标签 <span class="label-required">*</span></label>
        <div class="tag-grid">
          <label
            v-for="t in tags"
            :key="t.id"
            class="tag-chip"
            :class="{ selected: form.tag_ids.includes(t.id) }"
          >
            <input type="checkbox" :value="t.id" v-model="form.tag_ids" hidden />
            <span class="tag-dot"></span>
            {{ t.name }}
          </label>
          <span v-if="tags.length === 0 && !showTagInput" class="empty-hint">暂无标签，请先创建</span>

          <!-- 新建标签 -->
          <div v-if="showTagInput" class="tag-new-inline">
            <input
              v-model="newTagName"
              class="tag-new-input"
              placeholder="输入标签名"
              @keyup.enter="handleCreateTag"
              @keyup.escape="showTagInput = false; newTagName = ''"
              ref="tagInputRef"
            />
            <button type="button" class="tag-new-btn" @click="handleCreateTag" :disabled="tagCreating || !newTagName.trim()">
              {{ tagCreating ? '...' : '✓' }}
            </button>
            <button type="button" class="tag-new-cancel" @click="showTagInput = false; newTagName = ''">✕</button>
          </div>
          <button v-else type="button" class="tag-add-btn" @click="showTagInput = true; $nextTick(() => $refs.tagInputRef?.focus())">
            + 新建标签
          </button>
        </div>
      </div>

      <!-- 内容 (Markdown) -->
      <div class="form-group">
        <label class="form-label">内容 <span class="label-hint">（Markdown 格式）</span></label>
        <div class="md-editor-wrapper">
          <MdEditor
            ref="mdEditorRef"
            v-model="form.content"
            :theme="'light'"
            language="zh-CN"
            :preview="true"
            :toolbarsExclude="['github', 'htmlPreview']"
            placeholder="输入文章内容..."
            style="height: 500px"
            @onUploadImg="handleMdEditorUploadImage"
          />
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button type="submit" class="btn-primary" :disabled="loading || uploading">
          <span v-if="loading" class="btn-spinner"></span>
          {{ loading ? '保存中...' : '保存文章' }}
        </button>
        <button type="button" class="btn-ghost" @click="router.back()">取消</button>
      </div>
    </form>

    <!-- 裁剪弹窗 -->
    <ImageCropModal
      :show="showCropModal"
      :imageSrc="cropImageSrc"
      @confirm="handleCropConfirm"
      @cancel="handleCropCancel"
    />
  </div>
</template>

<style scoped>
.article-edit {
  max-width: 100%;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h2 {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: #8a8578;
}

/* 表单样式 */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 {
  flex: 1;
}

.form-label {
  font-size: 13px;
  font-weight: 700;
  color: #333;
  letter-spacing: 0.02em;
}

.label-hint {
  font-weight: 400;
  color: #b5ae9e;
  font-size: 12px;
}

.label-required {
  color: #c23b22;
  margin-left: 2px;
}

.form-input {
  width: 100%;
  padding: 11px 14px;
  border: 1.5px solid rgba(26, 26, 26, 0.12);
  background: #faf7f2;
  color: #1a1a1a;
  font-size: 14px;
  font-family: inherit;
  border-radius: 8px;
  outline: none;
  transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);
  box-sizing: border-box;
}

.form-input:focus {
  border-color: #1a1a1a;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.04);
}

.form-input::placeholder {
  color: #b5ae9e;
}

.form-textarea {
  resize: vertical;
  line-height: 1.7;
  min-height: 300px;
}

.form-help {
  margin: 0;
  font-size: 12px;
  color: #8a8578;
}

.form-help code {
  background: rgba(26, 26, 26, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #c23b22;
}

/* 封面图上传 */
.cover-preview {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(26, 26, 26, 0.08);
}

.cover-preview img {
  width: 100%;
  max-height: 240px;
  object-fit: cover;
  display: block;
}

.cover-remove {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.cover-remove:hover {
  background: rgba(0, 0, 0, 0.8);
}

.cover-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 20px;
  border: 2px dashed rgba(26, 26, 26, 0.12);
  border-radius: 8px;
  background: #faf7f2;
  cursor: pointer;
  transition: all 0.25s;
  color: #8a8578;
  font-size: 14px;
}

.cover-upload:hover {
  border-color: #1a1a1a;
  background: #fff;
}

.cover-upload.uploading {
  opacity: 0.7;
  cursor: not-allowed;
}

.upload-icon {
  font-size: 28px;
}

.upload-hint {
  font-size: 12px;
  color: #b5ae9e;
}

.upload-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(26, 26, 26, 0.1);
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 标签选择 */
.tag-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  border-radius: 20px;
  font-size: 13px;
  color: #4a4a4a;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
  background: #faf7f2;
}

.tag-chip:hover {
  border-color: #c23b22;
  color: #c23b22;
}

.tag-chip.selected {
  background: #c23b22;
  border-color: #c23b22;
  color: #fff;
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.4;
}

.tag-chip.selected .tag-dot {
  opacity: 1;
}

.empty-hint {
  font-size: 13px;
  color: #b5ae9e;
}

/* 新建标签 */
.tag-add-btn {
  display: inline-flex;
  align-items: center;
  padding: 7px 16px;
  border: 1.5px dashed rgba(26, 26, 26, 0.2);
  border-radius: 20px;
  font-size: 13px;
  color: #8a8578;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.tag-add-btn:hover {
  border-color: #c23b22;
  color: #c23b22;
}

.tag-new-inline {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tag-new-input {
  width: 120px;
  padding: 6px 12px;
  border: 1.5px solid #c23b22;
  border-radius: 20px;
  font-size: 13px;
  color: #1a1a1a;
  background: #fff;
  outline: none;
  font-family: inherit;
}

.tag-new-btn {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: #c23b22;
  color: #fff;
  font-size: 13px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s;
}

.tag-new-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.tag-new-cancel {
  width: 28px;
  height: 28px;
  border: 1.5px solid rgba(26, 26, 26, 0.15);
  border-radius: 50%;
  background: transparent;
  color: #8a8578;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.tag-new-cancel:hover {
  border-color: #c23b22;
  color: #c23b22;
}

/* Markdown 编辑器 */
.md-editor-wrapper {
  border: 1.5px solid rgba(26, 26, 26, 0.12);
  border-radius: 8px;
  overflow: hidden;
  background: #faf7f2;
}

.md-editor-wrapper:focus-within {
  border-color: #1a1a1a;
  box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.04);
}

/* 操作按钮 */
.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
  border-top: 1px solid rgba(26, 26, 26, 0.06);
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 11px 28px;
  background: #1a1a1a;
  color: #f5f0e8;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.btn-primary:hover {
  background: #333;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.btn-ghost {
  padding: 11px 28px;
  background: transparent;
  color: #4a4a4a;
  border: 1.5px solid rgba(26, 26, 26, 0.15);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.btn-ghost:hover {
  border-color: #1a1a1a;
  color: #1a1a1a;
}

/* 响应式 */
@media (max-width: 640px) {
  .form-row {
    flex-direction: column;
  }
}
</style>
