<script setup>
import { ref, onMounted, reactive } from 'vue'
import { getAboutPage, adminUpdateAboutPage } from '../../../shared/api/modules/setting'
import { useMessage } from '../../../composables/useMessage'

const message = useMessage()

// 当前激活的标签页
const activeTab = ref('profile')

// 加载状态
const loading = ref(false)
const saving = ref(false)

// 表单数据
const form = reactive({
  profile: {
    name: '',
    title: '',
    bio: '',
    social: []
  },
  highlights: [],
  skills: [],
  timeline: [],
  philosophy: {
    title: '',
    content: ''
  }
})

// 新增项的临时数据
const newHighlight = ref('')
const newSkill = reactive({ name: '', category: '' })
const newTimeline = reactive({ year: '', title: '', description: '' })
const newSocial = reactive({ name: '', url: '', icon: '' })

// 加载数据
async function loadData() {
  loading.value = true
  try {
    const { data } = await getAboutPage()
    if (data.data) {
      const d = data.data
      if (d.profile) {
        form.profile.name = d.profile.name || ''
        form.profile.title = d.profile.title || ''
        form.profile.bio = d.profile.bio || ''
        // 兼容旧格式：如果是对象则转换为数组
        if (d.profile.social && !Array.isArray(d.profile.social)) {
          const socialArr = []
          if (d.profile.social.github) {
            socialArr.push({ name: 'GitHub', url: d.profile.social.github, icon: 'github' })
          }
          if (d.profile.social.email) {
            socialArr.push({ name: 'Email', url: 'mailto:' + d.profile.social.email, icon: 'email' })
          }
          form.profile.social = socialArr
        } else {
          form.profile.social = d.profile.social || []
        }
      }
      form.highlights = d.highlights || []
      form.skills = d.skills || []
      form.timeline = d.timeline || []
      if (d.philosophy) {
        form.philosophy.title = d.philosophy.title || ''
        form.philosophy.content = d.philosophy.content || ''
      }
    }
  } catch (e) {
    message.error('加载数据失败: ' + (e.response?.data?.msg || e.message))
  } finally {
    loading.value = false
  }
}

// 保存数据
async function handleSave() {
  saving.value = true
  try {
    await adminUpdateAboutPage({
      profile: form.profile,
      highlights: form.highlights,
      skills: form.skills,
      timeline: form.timeline,
      philosophy: form.philosophy
    })
    message.success('保存成功')
  } catch (e) {
    message.error('保存失败: ' + (e.response?.data?.msg || e.message))
  } finally {
    saving.value = false
  }
}

// 添加社交链接
function addSocial() {
  if (newSocial.name.trim() && newSocial.url.trim()) {
    form.profile.social.push({
      name: newSocial.name.trim(),
      url: newSocial.url.trim(),
      icon: newSocial.icon.trim() || newSocial.name.trim().toLowerCase()
    })
    newSocial.name = ''
    newSocial.url = ''
    newSocial.icon = ''
  }
}

// 删除社交链接
function removeSocial(index) {
  form.profile.social.splice(index, 1)
}

// 添加特点
function addHighlight() {
  if (newHighlight.value.trim()) {
    form.highlights.push(newHighlight.value.trim())
    newHighlight.value = ''
  }
}

// 删除特点
function removeHighlight(index) {
  form.highlights.splice(index, 1)
}

// 添加技术栈
function addSkill() {
  if (newSkill.name.trim() && newSkill.category.trim()) {
    form.skills.push({ name: newSkill.name.trim(), category: newSkill.category.trim() })
    newSkill.name = ''
    newSkill.category = ''
  }
}

// 删除技术栈
function removeSkill(index) {
  form.skills.splice(index, 1)
}

// 添加经历
function addTimeline() {
  if (newTimeline.year.trim() && newTimeline.title.trim()) {
    form.timeline.push({
      year: newTimeline.year.trim(),
      title: newTimeline.title.trim(),
      description: newTimeline.description.trim()
    })
    newTimeline.year = ''
    newTimeline.title = ''
    newTimeline.description = ''
  }
}

// 删除经历
function removeTimeline(index) {
  form.timeline.splice(index, 1)
}

// 上移经历
function moveTimelineUp(index) {
  if (index > 0) {
    const temp = form.timeline[index]
    form.timeline[index] = form.timeline[index - 1]
    form.timeline[index - 1] = temp
  }
}

// 下移经历
function moveTimelineDown(index) {
  if (index < form.timeline.length - 1) {
    const temp = form.timeline[index]
    form.timeline[index] = form.timeline[index + 1]
    form.timeline[index + 1] = temp
  }
}

onMounted(loadData)
</script>

<template>
  <div class="settings-page">
    <div class="page-header">
      <h3>站点设置</h3>
      <button class="primary-btn" @click="handleSave" :disabled="saving">
        {{ saving ? '保存中...' : '保存所有更改' }}
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <p>加载中...</p>
    </div>

    <!-- 设置内容 -->
    <div v-else class="settings-content">
      <!-- 标签页导航 -->
      <div class="tabs">
        <button
          :class="['tab-btn', { active: activeTab === 'profile' }]"
          @click="activeTab = 'profile'"
        >
          个人信息
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'highlights' }]"
          @click="activeTab = 'highlights'"
        >
          特点
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'skills' }]"
          @click="activeTab = 'skills'"
        >
          技术栈
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'timeline' }]"
          @click="activeTab = 'timeline'"
        >
          经历
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'philosophy' }]"
          @click="activeTab = 'philosophy'"
        >
          写作理念
        </button>
      </div>

      <!-- 个人信息 -->
      <div v-show="activeTab === 'profile'" class="tab-content">
        <div class="form-section">
          <h4>基本信息</h4>
          <label class="form-label">
            <span>姓名</span>
            <input v-model="form.profile.name" placeholder="请输入姓名" />
          </label>
          <label class="form-label">
            <span>头衔</span>
            <input v-model="form.profile.title" placeholder="例如：后端开发者 / 博客作者" />
          </label>
          <label class="form-label">
            <span>个人简介</span>
            <textarea v-model="form.profile.bio" rows="3" placeholder="请输入个人简介"></textarea>
          </label>
        </div>

        <div class="form-section">
          <h4>社交链接</h4>
          <div class="item-list">
            <div v-for="(item, index) in form.profile.social" :key="index" class="item-row">
              <div class="social-info">
                <span class="social-name">{{ item.name }}</span>
                <span class="social-url">{{ item.url }}</span>
              </div>
              <button class="danger-btn" @click="removeSocial(index)">删除</button>
            </div>
          </div>
          <div class="add-form social-form">
            <input v-model="newSocial.name" placeholder="名称（如：GitHub）" />
            <input v-model="newSocial.url" placeholder="链接地址" />
            <input v-model="newSocial.icon" placeholder="图标名（可选）" />
            <button class="primary-btn" @click="addSocial">添加</button>
          </div>
        </div>
      </div>

      <!-- 特点 -->
      <div v-show="activeTab === 'highlights'" class="tab-content">
        <div class="form-section">
          <h4>特点列表</h4>
          <div class="item-list">
            <div v-for="(item, index) in form.highlights" :key="index" class="item-row">
              <span class="item-text">{{ item }}</span>
              <button class="danger-btn" @click="removeHighlight(index)">删除</button>
            </div>
          </div>
          <div class="add-form">
            <input v-model="newHighlight" placeholder="输入新的特点" @keyup.enter="addHighlight" />
            <button class="primary-btn" @click="addHighlight">添加</button>
          </div>
        </div>
      </div>

      <!-- 技术栈 -->
      <div v-show="activeTab === 'skills'" class="tab-content">
        <div class="form-section">
          <h4>技术栈列表</h4>
          <div class="item-list">
            <div v-for="(item, index) in form.skills" :key="index" class="item-row">
              <span class="item-text">{{ item.name }}</span>
              <span class="item-category">{{ item.category }}</span>
              <button class="danger-btn" @click="removeSkill(index)">删除</button>
            </div>
          </div>
          <div class="add-form skill-form">
            <input v-model="newSkill.name" placeholder="技术名称" />
            <input v-model="newSkill.category" placeholder="分类（如：语言、前端、工具）" />
            <button class="primary-btn" @click="addSkill">添加</button>
          </div>
        </div>
      </div>

      <!-- 经历 -->
      <div v-show="activeTab === 'timeline'" class="tab-content">
        <div class="form-section">
          <h4>经历时间轴</h4>
          <div class="item-list">
            <div v-for="(item, index) in form.timeline" :key="index" class="item-row timeline-item">
              <div class="item-info">
                <span class="item-year">{{ item.year }}</span>
                <span class="item-title">{{ item.title }}</span>
                <span class="item-desc" v-if="item.description">{{ item.description }}</span>
              </div>
              <div class="item-actions">
                <button @click="moveTimelineUp(index)" :disabled="index === 0">↑</button>
                <button @click="moveTimelineDown(index)" :disabled="index === form.timeline.length - 1">↓</button>
                <button class="danger-btn" @click="removeTimeline(index)">删除</button>
              </div>
            </div>
          </div>
          <div class="add-form timeline-form">
            <input v-model="newTimeline.year" placeholder="年份" />
            <input v-model="newTimeline.title" placeholder="标题" />
            <input v-model="newTimeline.description" placeholder="描述（可选）" />
            <button class="primary-btn" @click="addTimeline">添加</button>
          </div>
        </div>
      </div>

      <!-- 写作理念 -->
      <div v-show="activeTab === 'philosophy'" class="tab-content">
        <div class="form-section">
          <h4>写作理念</h4>
          <label class="form-label">
            <span>标题</span>
            <input v-model="form.philosophy.title" placeholder="例如：Philosophy" />
          </label>
          <label class="form-label">
            <span>内容</span>
            <textarea v-model="form.philosophy.content" rows="4" placeholder="请输入写作理念内容"></textarea>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  width: 100%;
}

.page-header h3 {
  margin: 0;
}

.primary-btn {
  background: #e94560;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.primary-btn:hover:not(:disabled) {
  background: #d13a54;
}

.primary-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state {
  text-align: center;
  padding: 40px;
  color: #666;
}

/* 标签页 */
.tabs {
  display: flex;
  gap: 4px;
  border-bottom: 1px solid #e0e0e0;
  margin-bottom: 24px;
  width: 100%;
}

.tab-btn {
  padding: 12px 20px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #333;
}

.tab-btn.active {
  color: #e94560;
  border-bottom-color: #e94560;
}

/* 表单 */
.tab-content {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  width: 100%;
  box-sizing: border-box;
}

.form-section {
  margin-bottom: 24px;
  width: 100%;
}

.form-section:last-child {
  margin-bottom: 0;
}

.form-section h4 {
  margin: 0 0 16px;
  font-size: 16px;
  color: #333;
}

.form-label {
  display: block;
  margin-bottom: 16px;
  width: 100%;
}

.form-label span {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 6px;
}

.form-label input,
.form-label textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
  transition: border-color 0.2s;
  max-width: 100%;
}

.form-label input:focus,
.form-label textarea:focus {
  outline: none;
  border-color: #e94560;
}

.form-label textarea {
  resize: vertical;
  min-height: 80px;
}

/* 列表项 */
.item-list {
  margin-bottom: 16px;
  width: 100%;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: #f9f9f9;
  border-radius: 6px;
  margin-bottom: 8px;
  width: 100%;
  box-sizing: border-box;
}

.item-row.timeline-item {
  flex-direction: column;
  align-items: flex-start;
}

/* 社交链接 */
.social-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.social-name {
  font-weight: 500;
  font-size: 14px;
}

.social-url {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.social-form {
  flex-wrap: wrap;
}

.social-form input:nth-child(1) {
  flex: 0 0 120px;
}

.social-form input:nth-child(2) {
  flex: 1;
  min-width: 200px;
}

.social-form input:nth-child(3) {
  flex: 0 0 120px;
}

.item-info {
  flex: 1;
}

.item-text {
  flex: 1;
  font-size: 14px;
}

.item-category {
  font-size: 12px;
  color: #888;
  background: #eee;
  padding: 2px 8px;
  border-radius: 4px;
}

.item-year {
  font-weight: 600;
  color: #e94560;
  margin-right: 8px;
}

.item-title {
  font-weight: 500;
  margin-right: 8px;
}

.item-desc {
  font-size: 13px;
  color: #666;
}

.item-actions {
  display: flex;
  gap: 4px;
}

.item-actions button {
  padding: 4px 8px;
  border: 1px solid #ddd;
  background: #fff;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.item-actions button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.danger-btn {
  color: #ff4d4f;
  border-color: #ff4d4f !important;
}

.danger-btn:hover {
  background: #fff1f0;
}

/* 添加表单 */
.add-form {
  display: flex;
  gap: 8px;
  width: 100%;
}

.add-form input {
  flex: 1;
  min-width: 0;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.add-form input:focus {
  outline: none;
  border-color: #e94560;
}

.skill-form input:nth-child(2) {
  flex: 0.6;
}

.timeline-form {
  flex-wrap: wrap;
}

.timeline-form input:nth-child(1) {
  flex: 0 0 100px;
}

.timeline-form input:nth-child(2) {
  flex: 1;
  min-width: 150px;
}

.timeline-form input:nth-child(3) {
  flex: 1;
  min-width: 200px;
}

/* 响应式 */
@media (max-width: 768px) {
  .tabs {
    flex-wrap: wrap;
  }

  .tab-btn {
    flex: 1;
    min-width: 80px;
    text-align: center;
  }

  .add-form {
    flex-wrap: wrap;
  }

  .add-form input {
    flex: 1 1 100%;
  }

  .skill-form input,
  .timeline-form input,
  .social-form input {
    flex: 1 1 100%;
  }

  .item-row {
    flex-wrap: wrap;
  }

  .item-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .social-info {
    width: 100%;
  }
}
</style>
