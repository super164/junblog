<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { marked } from 'marked'
import SiteShell from '../components/SiteShell.vue'
import ShareCardModal from '../components/ShareCardModal.vue'
import { getArticleByID } from '../shared/api/modules/article'
import { getComments, createComment } from '../shared/api/modules/comment'
import { getInteractionStatus, toggleLike } from '../shared/api/modules/interaction'
import { useAuth } from '../shared/stores/auth'
import { timeAgo } from '../shared/utils/date'
import { useMessage } from '../composables/useMessage'

const route = useRoute()
const message = useMessage()

const article = ref(null)
const comments = ref([])
const loading = ref(true)
const commentText = ref('')
const submitting = ref(false)
const readProgress = ref(0)

// 互动状态
const isLiked = ref(false)
const likeCount = ref(0)

// 分享功能
const showShareModal = ref(false)

const auth = useAuth()

// 阅读进度
function handleScroll() {
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  readProgress.value = docHeight > 0 ? Math.min((window.scrollY / docHeight) * 100, 100) : 0
}

async function loadArticle() {
  loading.value = true
  try {
    console.log('正在加载文章，ID:', route.params.id)
    const { data } = await getArticleByID(route.params.id)
    console.log('文章数据:', data)
    article.value = data.data
    likeCount.value = data.data.likes_count || 0
    loadComments()
    if (auth.token) {
      try {
        const interRes = await getInteractionStatus(article.value.id)
        isLiked.value = interRes.data.data?.is_liked || false
      } catch (e) { /* 未登录忽略 */ }
    }
  } catch (e) {
    console.error('加载文章失败', e)
    console.error('错误详情:', e.response?.data || e.message)
  } finally {
    loading.value = false
  }
}

async function loadComments() {
  try {
    const { data } = await getComments(article.value.id)
    // 确保评论数据格式正确，过滤掉无效数据
    const commentList = data.data?.list || data.data || []
    comments.value = Array.isArray(commentList) ? commentList.filter(c => c && c.id) : []
    console.log('评论数据:', comments.value)
  } catch (e) {
    console.error('加载评论失败:', e)
    comments.value = []
  }
}

async function handleLike() {
  if (!auth.token) { message.warning('请先登录'); return }
  try {
    const { data } = await toggleLike(article.value.id)
    isLiked.value = data.data.is_liked
    likeCount.value = data.data.likes_count
  } catch (e) { message.error('操作失败') }
}

async function handleComment() {
  if (!auth.token) { message.warning('请先登录'); return }
  if (!commentText.value.trim()) return
  submitting.value = true
  try {
    await createComment({ article_id: article.value.id, content: commentText.value })
    commentText.value = ''
    await loadComments()
    message.success('评论成功')
  } catch (e) { message.error('评论失败') }
  finally { submitting.value = false }
}

const htmlContent = computed(() => {
  if (!article.value?.content) return ''
  // 修复 Windows 路径中的反斜杠，并确保图片 URL 正确
  let content = article.value.content
  // 将反斜杠替换为正斜杠
  content = content.replace(/\\/g, '/')
  // 修复相对路径的图片 URL（添加后端地址）
  content = content.replace(/!\[\]\((\/uploads\/[^)]+)\)/g, (match, path) => {
    return `![](http://localhost:8080${path})`
  })
  return marked(content)
})

const wordCount = computed(() => {
  if (!article.value?.content) return 0
  return article.value.content.length
})

const readTime = computed(() => {
  return Math.max(1, Math.ceil(wordCount.value / 400))
})

onMounted(() => {
  loadArticle()
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <SiteShell>
    <!-- 阅读进度条 -->
    <div class="read-progress" :style="{ width: readProgress + '%' }"></div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <span>加载中...</span>
    </div>

    <article v-else-if="article" class="article-detail">
      <!-- 文章头部 -->
      <header class="article-hero anim-fade-up">
        <div class="hero-meta-row">
          <RouterLink v-if="article.category" :to="`/articles?category_id=${article.category.id}`" class="hero-category">
            {{ article.category.name }}
          </RouterLink>
          <span class="hero-dot">·</span>
          <time>{{ timeAgo(article.created_at) }}</time>
          <span class="hero-dot">·</span>
          <span>{{ readTime }} 分钟阅读</span>
          <span class="hero-dot">·</span>
          <button class="hero-share-btn" @click="showShareModal = true">
            📤 分享
          </button>
        </div>

        <h1 class="hero-title">{{ article.title }}</h1>

        <div class="hero-tags" v-if="article.tags?.length">
          <RouterLink
            v-for="tag in article.tags"
            :key="tag.id"
            :to="`/articles?tag_id=${tag.id}`"
            class="hero-tag"
          >
            {{ tag.name }}
          </RouterLink>
        </div>

        <div class="hero-divider">
          <span class="divider-ornament">◆</span>
        </div>
      </header>

      <!-- 正文区域 -->
      <div class="article-body-wrap">
        <aside class="article-sidebar anim-fade-up anim-delay-1">
          <div class="sidebar-sticky">
            <button class="sidebar-btn" :class="{ active: isLiked }" @click="handleLike" :title="isLiked ? '取消点赞' : '点赞'">
              <span class="sidebar-icon">{{ isLiked ? '♥' : '♡' }}</span>
              <span class="sidebar-count">{{ likeCount }}</span>
            </button>
            <div class="sidebar-divider"></div>
            <span class="sidebar-views">👁 {{ article.views_count }}</span>
          </div>
        </aside>

        <div class="article-content md-content anim-fade-up anim-delay-2" v-html="htmlContent"></div>
      </div>

      <!-- 文章底部互动 -->
      <footer class="article-footer anim-fade-up anim-delay-3">
        <div class="footer-divider">
          <span class="divider-ornament">◆</span>
        </div>

        <div class="footer-actions">
          <button class="footer-btn" :class="{ liked: isLiked }" @click="handleLike">
            <span>{{ isLiked ? '❤️' : '🤍' }}</span>
            <span>{{ likeCount }} 点赞</span>
          </button>
          <button class="footer-btn share-btn" @click="showShareModal = true">
            <span>📤</span>
            <span>分享</span>
          </button>
        </div>
      </footer>

      <!-- 分享弹窗 -->
      <ShareCardModal
        :visible="showShareModal"
        :article="article"
        @close="showShareModal = false"
      />

      <!-- 评论区 -->
      <section class="comment-section anim-fade-up anim-delay-4">
        <div class="comment-header-row">
          <h3>评论</h3>
          <span class="comment-count">{{ comments.length }}</span>
        </div>

        <div v-if="auth.token" class="comment-form">
          <div class="comment-form-avatar">{{ auth.user?.username?.slice(0, 1) || 'U' }}</div>
          <div class="comment-form-body">
            <textarea
              v-model="commentText"
              placeholder="写下你的想法..."
              rows="3"
            ></textarea>
            <div class="comment-form-footer">
              <span class="comment-form-hint">支持 Markdown</span>
              <button @click="handleComment" :disabled="submitting || !commentText.trim()">
                {{ submitting ? '发送中...' : '发表评论' }}
              </button>
            </div>
          </div>
        </div>
        <div v-else class="login-prompt">
          <RouterLink to="/login">登录</RouterLink> 后参与讨论
        </div>

        <div v-if="comments.length === 0" class="empty-comments">
          <span class="empty-icon">💬</span>
          <p>暂无评论，来说两句吧~</p>
        </div>
        <div v-else class="comment-list">
          <template v-for="c in comments" :key="c?.id || Math.random()">
            <div class="comment-item">
              <div class="comment-avatar">{{ c.user?.username?.slice(0, 1) || 'U' }}</div>
              <div class="comment-body">
                <div class="comment-meta">
                  <strong>{{ c.user?.username }}</strong>
                  <time>{{ timeAgo(c.created_at) }}</time>
                </div>
                <p>{{ c.content }}</p>
              </div>
            </div>
            <template v-if="c.replies?.length">
              <div v-for="child in c.replies" :key="child?.id || Math.random()" class="comment-item child">
                <div class="comment-avatar small">{{ child.user?.username?.slice(0, 1) || 'U' }}</div>
                <div class="comment-body">
                  <div class="comment-meta">
                    <strong>{{ child.user?.username }}</strong>
                    <time>{{ timeAgo(child.created_at) }}</time>
                  </div>
                  <p>{{ child.content }}</p>
                </div>
              </div>
            </template>
          </template>
        </div>
      </section>

      <!-- 返回链接 -->
      <div class="back-link anim-fade-up">
        <RouterLink to="/articles">
          <span class="back-arrow">←</span>
          返回文章列表
        </RouterLink>
      </div>
    </article>

    <div v-else class="empty-state">文章不存在</div>
  </SiteShell>
</template>

<style scoped>
/* 阅读进度条 */
.read-progress {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  background: linear-gradient(90deg, #c23b22, #e94560);
  z-index: 1000;
  transition: width 0.1s linear;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 120px 0;
  color: var(--ink-muted, #8a8578);
  font-size: 14px;
}

.loading-spinner {
  width: 28px;
  height: 28px;
  border: 2.5px solid rgba(194, 59, 34, 0.15);
  border-top-color: #c23b22;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* 文章容器 */
.article-detail {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 24px;
}

/* === 文章头部 === */
.article-hero {
  padding: 48px 0 0;
  text-align: center;
}

.hero-meta-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 12px;
  color: var(--ink-muted, #8a8578);
  letter-spacing: 0.06em;
  text-transform: uppercase;
  margin-bottom: 20px;
}

.hero-category {
  color: #c23b22;
  text-decoration: none;
  font-weight: 600;
  letter-spacing: 0.08em;
  transition: opacity 0.2s;
}

.hero-category:hover {
  opacity: 0.7;
}

.hero-dot {
  color: var(--ink-faint, #b5ae9e);
  font-size: 8px;
}

.hero-share-btn {
  background: none;
  border: none;
  color: var(--ink-muted, #8a8578);
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
  font-family: inherit;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.hero-share-btn:hover {
  color: #c23b22;
}

.hero-title {
  font-family: var(--font-heading, 'Noto Serif SC', serif);
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 700;
  line-height: 1.3;
  letter-spacing: -0.02em;
  color: var(--ink-primary, #1a1a1a);
  margin: 0 0 20px;
  max-width: 640px;
  margin-left: auto;
  margin-right: auto;
}

.hero-tags {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 28px;
}

.hero-tag {
  font-size: 11px;
  color: #c23b22;
  background: rgba(194, 59, 34, 0.06);
  padding: 5px 14px;
  border-radius: 20px;
  text-decoration: none;
  letter-spacing: 0.03em;
  transition: all 0.2s;
  border: 1px solid rgba(194, 59, 34, 0.1);
}

.hero-tag:hover {
  background: #c23b22;
  color: #fff;
  border-color: #c23b22;
}

.hero-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.divider-ornament {
  color: var(--ink-faint, #b5ae9e);
  font-size: 10px;
  letter-spacing: 12px;
}

/* === 正文 + 侧边栏 === */
.article-body-wrap {
  display: flex;
  gap: 32px;
  padding: 0 0 40px;
}

.article-sidebar {
  flex-shrink: 0;
  width: 48px;
}

.sidebar-sticky {
  position: sticky;
  top: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.sidebar-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 10px 0;
  border: none;
  background: none;
  cursor: pointer;
  color: var(--ink-muted, #8a8578);
  transition: all 0.25s;
  border-radius: 8px;
  width: 100%;
}

.sidebar-btn:hover {
  color: #c23b22;
  background: rgba(194, 59, 34, 0.04);
}

.sidebar-btn.active {
  color: #c23b22;
}

.sidebar-icon {
  font-size: 18px;
  line-height: 1;
}

.sidebar-count {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.sidebar-divider {
  width: 20px;
  height: 1px;
  background: rgba(26, 26, 26, 0.08);
  margin: 4px 0;
}

.sidebar-views {
  font-size: 10px;
  color: var(--ink-faint, #b5ae9e);
  text-align: center;
  line-height: 1.4;
}

/* === 文章正文 === */
.article-content {
  flex: 1;
  min-width: 0;
  padding: 0 0 0 0;
  font-size: 15.5px;
  line-height: 1.9;
  color: var(--ink-secondary, #4a4a4a);
}

/* Markdown 渲染内容样式 */
.article-content :deep(h1) {
  font-family: var(--font-heading, 'Noto Serif SC', serif);
  font-size: 1.9em;
  margin: 1.6em 0 0.6em;
  font-weight: 700;
  color: var(--ink-primary, #1a1a1a);
  line-height: 1.3;
}

.article-content :deep(h2) {
  font-family: var(--font-heading, 'Noto Serif SC', serif);
  font-size: 1.45em;
  margin: 1.4em 0 0.5em;
  font-weight: 700;
  color: var(--ink-primary, #1a1a1a);
  padding-bottom: 0.35em;
  border-bottom: 2px solid rgba(194, 59, 34, 0.12);
  line-height: 1.35;
}

.article-content :deep(h3) {
  font-family: var(--font-heading, 'Noto Serif SC', serif);
  font-size: 1.2em;
  margin: 1.2em 0 0.4em;
  font-weight: 700;
  color: var(--ink-primary, #1a1a1a);
}

.article-content :deep(p) {
  margin: 0.9em 0;
}

.article-content :deep(ul), .article-content :deep(ol) {
  padding-left: 1.6em;
  margin: 0.9em 0;
}

.article-content :deep(li) {
  margin: 0.35em 0;
}

.article-content :deep(blockquote) {
  margin: 1.4em 0;
  padding: 16px 24px;
  border-left: 3px solid #c23b22;
  background: linear-gradient(135deg, rgba(194, 59, 34, 0.03), rgba(194, 59, 34, 0.01));
  color: var(--ink-muted, #8a8578);
  border-radius: 0 10px 10px 0;
  font-style: italic;
  line-height: 1.8;
}

.article-content :deep(code) {
  background: rgba(194, 59, 34, 0.06);
  padding: 2px 7px;
  border-radius: 4px;
  font-size: 0.88em;
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  color: #c23b22;
  font-weight: 500;
}

.article-content :deep(pre) {
  background: #1a1a1e;
  color: #e0ddd5;
  padding: 20px 24px;
  border-radius: 12px;
  overflow-x: auto;
  margin: 1.4em 0;
  line-height: 1.65;
  border: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.article-content :deep(pre code) {
  background: none;
  color: inherit;
  padding: 0;
  font-size: 0.87em;
  font-weight: 400;
}

.article-content :deep(img) {
  max-width: 100%;
  border-radius: 10px;
  margin: 1.2em 0;
  box-shadow: 0 2px 12px rgba(80, 60, 30, 0.08);
}

.article-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1.2em 0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(26, 26, 26, 0.08);
}

.article-content :deep(th), .article-content :deep(td) {
  padding: 11px 16px;
  border: 1px solid rgba(26, 26, 26, 0.06);
  text-align: left;
}

.article-content :deep(th) {
  background: var(--bg-warm, #efe9dd);
  font-weight: 600;
  font-size: 0.9em;
  color: var(--ink-primary, #1a1a1a);
}

.article-content :deep(a) {
  color: #c23b22;
  text-decoration: none;
  border-bottom: 1px solid rgba(194, 59, 34, 0.3);
  transition: all 0.2s;
}

.article-content :deep(a:hover) {
  border-bottom-color: #c23b22;
  background: rgba(194, 59, 34, 0.04);
}

.article-content :deep(hr) {
  border: none;
  text-align: center;
  margin: 2.5em 0;
}

.article-content :deep(hr)::after {
  content: '◆ ◆ ◆';
  color: var(--ink-faint, #b5ae9e);
  font-size: 8px;
  letter-spacing: 8px;
}

.article-content :deep(strong) {
  color: var(--ink-primary, #1a1a1a);
  font-weight: 600;
}

/* === 文章底部 === */
.article-footer {
  padding: 0 0 32px;
}

.footer-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 0 28px;
}

.footer-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.footer-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  background: var(--bg-card, #faf7f2);
  border-radius: 28px;
  cursor: pointer;
  font-size: 14px;
  color: var(--ink-secondary, #4a4a4a);
  transition: all 0.3s cubic-bezier(0.22, 1, 0.36, 1);
  font-family: inherit;
}

.footer-btn:hover {
  border-color: #c23b22;
  color: #c23b22;
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(194, 59, 34, 0.1);
}

.footer-btn.liked {
  border-color: #c23b22;
  background: rgba(194, 59, 34, 0.04);
  color: #c23b22;
}

.footer-btn.share-btn:hover {
  border-color: #1da1f2;
  color: #1da1f2;
}

/* === 评论区 === */
.comment-section {
  border-top: 1px solid rgba(26, 26, 26, 0.06);
  padding-top: 36px;
  margin-bottom: 40px;
}

.comment-header-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
}

.comment-header-row h3 {
  margin: 0;
  font-family: var(--font-heading, 'Noto Serif SC', serif);
  font-size: 20px;
  font-weight: 700;
  color: var(--ink-primary, #1a1a1a);
}

.comment-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 8px;
  background: var(--accent-soft, rgba(194, 59, 34, 0.08));
  color: #c23b22;
  font-size: 12px;
  font-weight: 700;
  border-radius: 12px;
}

/* 评论表单 */
.comment-form {
  display: flex;
  gap: 14px;
  margin-bottom: 32px;
}

.comment-form-avatar {
  width: 38px;
  height: 38px;
  background: linear-gradient(135deg, #c23b22 0%, #e94560 100%);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.comment-form-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.comment-form textarea {
  padding: 14px 18px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  border-radius: 12px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  background: var(--bg-card, #faf7f2);
  transition: all 0.25s;
  outline: none;
  line-height: 1.7;
}

.comment-form textarea:focus {
  border-color: #c23b22;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(194, 59, 34, 0.06);
}

.comment-form-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.comment-form-hint {
  font-size: 12px;
  color: var(--ink-faint, #b5ae9e);
}

.comment-form button {
  padding: 10px 24px;
  background: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 13px;
  transition: all 0.2s;
  font-family: inherit;
}

.comment-form button:hover:not(:disabled) {
  background: #333;
  transform: translateY(-1px);
}

.comment-form button:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.login-prompt {
  padding: 24px;
  background: var(--bg-card, #faf7f2);
  border-radius: 12px;
  text-align: center;
  font-size: 14px;
  color: var(--ink-muted, #8a8578);
  margin-bottom: 28px;
  border: 1px dashed rgba(26, 26, 26, 0.1);
}

.login-prompt a {
  color: #c23b22;
  font-weight: 600;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.login-prompt a:hover {
  border-bottom-color: #c23b22;
}

.empty-comments {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-faint, #b5ae9e);
}

.empty-icon {
  font-size: 32px;
  display: block;
  margin-bottom: 8px;
}

.empty-comments p {
  margin: 0;
  font-size: 14px;
}

/* 评论列表 */
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.comment-item {
  display: flex;
  gap: 14px;
  padding: 20px 0;
  border-bottom: 1px solid rgba(26, 26, 26, 0.04);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-item.child {
  margin-left: 52px;
  padding-top: 16px;
  border-top: none;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #c23b22 0%, #e94560 100%);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.comment-avatar.small {
  width: 28px;
  height: 28px;
  font-size: 11px;
  background: var(--ink-muted, #8a8578);
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-meta {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 6px;
}

.comment-meta strong {
  font-size: 13px;
  color: var(--ink-primary, #1a1a1a);
  font-weight: 600;
}

.comment-meta time {
  font-size: 12px;
  color: var(--ink-faint, #b5ae9e);
}

.comment-body p {
  margin: 0;
  font-size: 14px;
  color: var(--ink-secondary, #4a4a4a);
  line-height: 1.7;
}

/* 返回链接 */
.back-link {
  text-align: center;
  padding: 20px 0 60px;
}

.back-link a {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--ink-muted, #8a8578);
  font-size: 13px;
  text-decoration: none;
  transition: all 0.2s;
  padding: 8px 16px;
  border-radius: 8px;
}

.back-link a:hover {
  color: #c23b22;
  background: rgba(194, 59, 34, 0.04);
}

.back-arrow {
  font-size: 16px;
  transition: transform 0.2s;
}

.back-link a:hover .back-arrow {
  transform: translateX(-3px);
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 120px 20px;
  color: var(--ink-faint, #b5ae9e);
  font-size: 15px;
}

/* 响应式 */
@media (max-width: 768px) {
  .article-detail {
    padding: 0 16px;
  }

  .article-body-wrap {
    gap: 0;
  }

  .article-sidebar {
    display: none;
  }

  .hero-title {
    font-size: 1.6rem;
  }

  .footer-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .footer-btn {
    justify-content: center;
  }

  .comment-item.child {
    margin-left: 24px;
  }
}
</style>
