<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import SiteShell from '../components/SiteShell.vue'
import { getArticles } from '../shared/api/modules/article'
import { getCategories } from '../shared/api/modules/category'
import { getTags, getTagsByCategory } from '../shared/api/modules/tag'
import { timeAgo } from '../shared/utils/date'
import { truncate } from '../shared/utils/text'

const route = useRoute()
const router = useRouter()

const articles = ref([])
const categories = ref([])
const tags = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const activeCategoryId = ref(0)
const activeTagId = ref(0)
const keyword = ref('')
const sort = ref('latest')
const loading = ref(false)

// 用于竞态条件防护：记录每次请求的标识，过期请求的结果直接丢弃
let fetchId = 0

async function loadData() {
  const currentFetchId = ++fetchId

  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (activeCategoryId.value) params.category_id = activeCategoryId.value
    if (activeTagId.value) params.tag_id = activeTagId.value
    if (keyword.value) params.keyword = keyword.value
    params.sort = sort.value

    const [artRes, catRes] = await Promise.all([
      getArticles(params),
      getCategories(),
    ])

    // 如果已经有更新的请求在进行中，丢弃本次结果
    if (currentFetchId !== fetchId) return

    articles.value = artRes.data.data?.list || []
    total.value = artRes.data.data?.total || 0
    categories.value = catRes.data.data || []

    // 根据分类加载标签
    await loadTags()
  } catch (e) {
    console.error('加载失败', e)
  } finally {
    if (currentFetchId === fetchId) {
      loading.value = false
    }
  }
}

async function loadTags() {
  try {
    let tagRes
    if (activeCategoryId.value) {
      // 根据分类获取标签
      tagRes = await getTagsByCategory(activeCategoryId.value)
    } else {
      // 获取所有标签
      tagRes = await getTags()
    }
    tags.value = tagRes.data.data || []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

function filterByCategory(id) {
  activeCategoryId.value = activeCategoryId.value === id ? 0 : id
  activeTagId.value = 0  // 切换分类时清空标签选择
  page.value = 1
  loadData()
}

function filterByTag(id) {
  activeTagId.value = activeTagId.value === id ? 0 : id
  page.value = 1
  loadData()
}

function handleSearch() {
  page.value = 1
  loadData()
}

function filterBySort(value) {
  sort.value = value
  page.value = 1
  loadData()
}

function goPage(p) {
  page.value = p
  loadData()
}

onMounted(loadData)
</script>

<template>
  <SiteShell>
    <section class="page-hero anim-fade-up">
      <p class="section-kicker">Article Library</p>
      <h1>文章列表</h1>
      <p>共 {{ total }} 篇文章，分类筛选或搜索找到你想读的内容。</p>
    </section>

    <!-- 搜索栏 -->
    <div class="search-bar anim-fade-up anim-delay-1">
      <input v-model="keyword" placeholder="搜索文章..." @keyup.enter="handleSearch" />
      <button @click="handleSearch">搜索</button>
    </div>

    <!-- 排序筛选 -->
    <section class="filter-row anim-fade-up anim-delay-2">
      <button class="filter-chip" :class="{ active: sort === 'latest' }" @click="filterBySort('latest')">最新</button>
      <button class="filter-chip" :class="{ active: sort === 'hottest' }" @click="filterBySort('hottest')">最热</button>
    </section>

    <!-- 分类筛选 -->
    <section class="filter-row anim-fade-up anim-delay-2">
      <button class="filter-chip" :class="{ active: activeCategoryId === 0 }" @click="filterByCategory(0)">全部</button>
      <button
        v-for="cat in categories"
        :key="cat.id"
        class="filter-chip"
        :class="{ active: activeCategoryId === cat.id }"
        @click="filterByCategory(cat.id)"
      >
        {{ cat.name }}
      </button>
    </section>

    <!-- 标签筛选 -->
    <section class="filter-row anim-fade-up anim-delay-2" v-if="tags.length">
      <span class="filter-label">标签：</span>
      <button
        v-for="tag in tags"
        :key="tag.id"
        class="filter-chip small"
        :class="{ active: activeTagId === tag.id }"
        @click="filterByTag(tag.id)"
      >
        # {{ tag.name }}
      </button>
    </section>

    <!-- 文章列表 -->
    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="articles.length === 0" class="empty-state">暂无文章</div>
    <section v-else class="articles-list">
      <RouterLink
        v-for="(article, index) in articles"
        :key="article.id"
        :to="`/articles/${article.id}`"
        class="list-card anim-fade-up"
      >
        <div v-if="article.cover" class="list-cover">
          <img :src="article.cover" :alt="article.title" />
        </div>
        <div class="list-body">
          <div class="article-meta">
            <span>{{ article.category?.name || '未分类' }}</span>
            <span>{{ timeAgo(article.created_at) }}</span>
            <span>👁 {{ article.views_count }}</span>
          </div>
          <h3>{{ article.title }}</h3>
          <p>{{ truncate(article.summary || article.content, 50) }}</p>
          <div class="tags-row" v-if="article.tags?.length">
            <span v-for="tag in article.tags" :key="tag.id" class="tag-chip"># {{ tag.name }}</span>
          </div>
        </div>
        <span class="text-link">阅读详情</span>
      </RouterLink>
    </section>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination anim-fade-up">
      <button :disabled="page <= 1" @click="goPage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ Math.ceil(total / pageSize) }} 页</span>
      <button :disabled="page >= Math.ceil(total / pageSize)" @click="goPage(page + 1)">下一页</button>
    </div>
  </SiteShell>
</template>

<style scoped>
.search-bar {
  display: flex;
  gap: 10px;
  max-width: 420px;
  margin-bottom: 24px;
}

.search-bar input {
  flex: 1;
  padding: 12px 16px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  border-radius: 10px;
  font-size: 14px;
  background: #faf7f2;
  transition: all 0.25s;
  outline: none;
}

.search-bar input:focus {
  border-color: #1a1a1a;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.04);
}

.search-bar button {
  padding: 12px 24px;
  background: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.25s;
}

.search-bar button:hover {
  background: #333;
  transform: translateY(-1px);
}

.filter-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
  align-items: center;
}

.filter-label {
  font-size: 13px;
  color: #8a8578;
}

.filter-chip {
  padding: 7px 16px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  background: transparent;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
  color: #4a4a4a;
  transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);
}

.filter-chip:hover {
  border-color: #c23b22;
  color: #c23b22;
  background: rgba(194, 59, 34, 0.04);
}

.filter-chip.active {
  background: #c23b22;
  color: #fff;
  border-color: #c23b22;
}

.filter-chip.small {
  padding: 5px 12px;
  font-size: 12px;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.list-card {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 22px 24px;
  background: #fff;
  border-radius: 14px;
  text-decoration: none;
  color: inherit;
  border: 1px solid rgba(26, 26, 26, 0.05);
  transition: all 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}

.list-card:hover {
  box-shadow: 0 8px 28px rgba(80, 60, 30, 0.08);
  transform: translateY(-2px);
  border-color: transparent;
}

.list-cover {
  width: 130px;
  height: 88px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
}

.list-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.list-card:hover .list-cover img {
  transform: scale(1.05);
}

.list-body {
  flex: 1;
}

.list-body h3 {
  margin: 6px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  transition: color 0.2s;
}

.list-card:hover .list-body h3 {
  color: #c23b22;
}

.list-body p {
  margin: 0;
  font-size: 13px;
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.6;
}

.article-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #8a8578;
}

.tags-row {
  display: flex;
  gap: 6px;
  margin-top: 8px;
}

.tag-chip {
  font-size: 11px;
  color: #c23b22;
  background: rgba(194, 59, 34, 0.06);
  padding: 3px 10px;
  border-radius: 4px;
}

.text-link {
  font-size: 13px;
  color: #c23b22;
  white-space: nowrap;
  font-weight: 500;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-top: 40px;
  padding: 20px;
}

.pagination button {
  padding: 10px 20px;
  border: 1.5px solid rgba(26, 26, 26, 0.12);
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  border-color: #1a1a1a;
  color: #1a1a1a;
}

.pagination button:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.pagination span {
  font-size: 13px;
  color: #8a8578;
}

.loading-state, .empty-state {
  text-align: center;
  padding: 60px;
  color: #b5ae9e;
  font-size: 15px;
}
</style>
