<script setup>
import { ref, onMounted } from 'vue'
import http from '../services/api'

const stats = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await http.get('/admin/stats')
    stats.value = data.data
  } catch (e) {
    console.error('获取统计数据失败', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard">
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <span>加载中...</span>
    </div>

    <template v-else-if="stats">
      <!-- 欢迎语 -->
      <div class="welcome-banner">
        <h2>控制台</h2>
        <p>欢迎回来，这是你的博客数据概览</p>
      </div>

      <!-- 数据卡片 -->
      <div class="stats-grid">
        <div class="stat-card" style="--card-accent: #e94560">
          <div class="stat-icon">📝</div>
          <div class="stat-body">
            <strong>{{ stats.total_articles }}</strong>
            <span>文章总数</span>
          </div>
          <div class="stat-detail">已发布 {{ stats.published_articles }} · 草稿 {{ stats.draft_articles }}</div>
        </div>

        <div class="stat-card" style="--card-accent: #6366f1">
          <div class="stat-icon">👁</div>
          <div class="stat-body">
            <strong>{{ stats.total_views }}</strong>
            <span>总浏览量</span>
          </div>
        </div>

        <div class="stat-card" style="--card-accent: #f59e0b">
          <div class="stat-icon">💬</div>
          <div class="stat-body">
            <strong>{{ stats.total_comments }}</strong>
            <span>评论总数</span>
          </div>
          <div class="stat-detail" v-if="stats.pending_comments > 0">待审核 {{ stats.pending_comments }}</div>
        </div>

        <div class="stat-card" style="--card-accent: #10b981">
          <div class="stat-icon">👥</div>
          <div class="stat-body">
            <strong>{{ stats.total_users }}</strong>
            <span>注册用户</span>
          </div>
        </div>

        <div class="stat-card" style="--card-accent: #8b5cf6">
          <div class="stat-icon">📁</div>
          <div class="stat-body">
            <strong>{{ stats.total_categories }}</strong>
            <span>分类数</span>
          </div>
        </div>

        <div class="stat-card" style="--card-accent: #ec4899">
          <div class="stat-icon">🏷️</div>
          <div class="stat-body">
            <strong>{{ stats.total_tags }}</strong>
            <span>标签数</span>
          </div>
        </div>
      </div>

      <!-- 最近文章 -->
      <div class="section-card">
        <div class="section-header">
          <h3>最近发布</h3>
        </div>
        <table class="data-table" v-if="stats.recent_articles?.length">
          <thead>
            <tr>
              <th>标题</th>
              <th>状态</th>
              <th>浏览</th>
              <th>发布时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in stats.recent_articles" :key="a.id">
              <td class="title-cell">{{ a.title }}</td>
              <td>
                <span class="status-badge" :class="a.status">
                  {{ a.status === 'published' ? '已发布' : a.status === 'draft' ? '草稿' : '私密' }}
                </span>
              </td>
              <td class="views-cell">{{ a.views_count }}</td>
              <td class="time-cell">{{ new Date(a.created_at).toLocaleDateString('zh-CN') }}</td>
            </tr>
          </tbody>
        </table>
        <div v-else class="empty-state">
          <span class="empty-icon">📭</span>
          <p>暂无文章</p>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 欢迎横幅 */
.welcome-banner {
  margin-bottom: 8px;
}

.welcome-banner h2 {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.02em;
}

.welcome-banner p {
  margin: 0;
  font-size: 14px;
  color: #8a8578;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  --card-accent: #e94560;
  background: #fff;
  padding: 22px;
  border-radius: 14px;
  border: 1px solid rgba(26, 26, 26, 0.05);
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.3s cubic-bezier(0.22, 1, 0.36, 1);
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--card-accent);
  opacity: 0;
  transition: opacity 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  border-color: transparent;
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-icon {
  font-size: 26px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(26, 26, 26, 0.03);
  border-radius: 12px;
}

.stat-body {
  display: flex;
  flex-direction: column;
}

.stat-body strong {
  font-size: 30px;
  font-weight: 800;
  color: #1a1a1a;
  line-height: 1;
  letter-spacing: -0.03em;
}

.stat-body span {
  font-size: 13px;
  color: #8a8578;
  margin-top: 4px;
}

.stat-detail {
  font-size: 12px;
  color: var(--card-accent);
  background: rgba(233, 69, 96, 0.06);
  padding: 5px 10px;
  border-radius: 6px;
  width: fit-content;
  font-weight: 500;
}

/* 内容卡片 */
.section-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid rgba(26, 26, 26, 0.05);
  overflow: hidden;
}

.section-header {
  padding: 20px 24px 0;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #1a1a1a;
}

/* 数据表格 */
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 14px 24px;
  text-align: left;
  font-size: 13px;
}

.data-table thead th {
  font-weight: 600;
  color: #8a8578;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid rgba(26, 26, 26, 0.06);
  background: rgba(26, 26, 26, 0.01);
}

.data-table tbody tr {
  transition: background 0.2s;
}

.data-table tbody tr:hover {
  background: rgba(26, 26, 26, 0.02);
}

.data-table tbody td {
  border-bottom: 1px solid rgba(26, 26, 26, 0.04);
  color: #4a4a4a;
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}

.title-cell {
  font-weight: 500;
  color: #1a1a1a;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.views-cell {
  font-variant-numeric: tabular-nums;
}

.time-cell {
  color: #8a8578;
}

/* 状态标签 */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.published {
  background: rgba(16, 185, 129, 0.1);
  color: #059669;
}

.status-badge.draft {
  background: rgba(245, 158, 11, 0.1);
  color: #d97706;
}

.status-badge.privacy {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #b5ae9e;
}

.empty-icon {
  font-size: 36px;
  display: block;
  margin-bottom: 8px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px;
  color: #8a8578;
  font-size: 14px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2.5px solid rgba(26, 26, 26, 0.08);
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 响应式 */
@media (max-width: 900px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
