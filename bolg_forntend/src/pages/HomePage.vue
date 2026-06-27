<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import SiteShell from '../components/SiteShell.vue'
import { getRecentArticles } from '../shared/api/modules/article'
import { timeAgo } from '../shared/utils/date'
import { truncate } from '../shared/utils/text'

const recentArticles = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await getRecentArticles(5)
    recentArticles.value = data.data || []
  } catch (e) {
    console.error('加载首页数据失败', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <SiteShell>
    <!-- Hero -->
    <section class="hero-section">
      <div class="hero-content anim-fade-up">
        <span class="hero-badge">JunBlog</span>
        <h1>技术不设限，分享即成长</h1>
        <p>保持好奇，持续探索。在这里记录技术实践与思考。</p>
        <div class="hero-actions">
          <RouterLink class="primary-button" to="/articles">浏览文章</RouterLink>
          <RouterLink class="ghost-button" to="/about">了解作者</RouterLink>
        </div>
      </div>
    </section>

    <!-- 最新文章 -->
    <section class="section-block">
      <div class="section-header anim-fade-up anim-delay-4">
        <div>
          <p class="section-kicker">Recent Posts</p>
          <h2>最新文章</h2>
        </div>
        <RouterLink class="text-link" to="/articles">查看全部</RouterLink>
      </div>

      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="recentArticles.length === 0" class="empty-state">暂无文章</div>
      <div v-else class="articles-grid">
        <RouterLink
          v-for="(article, index) in recentArticles"
          :key="article.id"
          :to="`/articles/${article.id}`"
          class="article-card anim-fade-up"
        >
          <div v-if="article.cover" class="article-card-cover">
            <img :src="article.cover" :alt="article.title" />
          </div>
          <div class="article-card-body">
            <span class="article-number">No.{{ String(index + 1).padStart(2, '0') }}</span>
            <div class="article-meta">
              <span>{{ article.category?.name || '未分类' }}</span>
              <span>{{ timeAgo(article.created_at) }}</span>
            </div>
            <h3>{{ article.title }}</h3>
            <p>{{ truncate(article.summary || article.content, 50) }}</p>
            <div class="tags-row" v-if="article.tags?.length">
              <span v-for="tag in article.tags" :key="tag.id" class="tag-chip"># {{ tag.name }}</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </section>

    <!-- Author Intro -->
    <section class="section-block intro-banner anim-fade-up anim-delay-8">
      <div>
        <p class="section-kicker">About The Author</p>
        <h2>Jun</h2>
        <p class="lead-text">全栈开发者 / 写作者</p>
      </div>
      <p class="intro-copy">热爱技术与设计，用代码构建有温度的产品。</p>
    </section>
  </SiteShell>
</template>

<style scoped>
.loading-state, .empty-state { text-align: center; padding: 40px; color: #999; }
</style>
