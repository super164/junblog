<script setup>
import SiteShell from '../components/SiteShell.vue'
import { profile } from '../data/mock'

const highlights = [
  '以 Go 后端开发为主，关注工程化与 API 设计',
  '习惯把开发过程沉淀成文章，强调长期输出',
  '喜欢简洁、克制、有层次的界面设计风格',
]
</script>

<template>
  <SiteShell>
    <!-- 作者卡片 -->
    <section class="author-card anim-fade-up">
      <div class="author-main">
        <div class="author-avatar">{{ profile.name.charAt(0) }}</div>
        <div class="author-info">
          <h1>{{ profile.name }}</h1>
          <p class="author-title">{{ profile.title }}</p>
          <p class="author-bio">{{ profile.bio }}</p>
          <div class="author-social">
            <a :href="profile.social.github" target="_blank" rel="noopener" class="social-link">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
              </svg>
              <span>GitHub</span>
            </a>
            <a :href="'mailto:' + profile.social.email" class="social-link">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="4" width="20" height="16" rx="2"/>
                <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
              </svg>
              <span>{{ profile.social.email }}</span>
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- 特点 -->
    <section class="content-card anim-fade-up anim-delay-2">
      <p class="section-kicker">Highlights</p>
      <h2>特点</h2>
      <ul class="plain-list">
        <li v-for="item in highlights" :key="item">{{ item }}</li>
      </ul>
    </section>

    <!-- 技术栈 -->
    <section class="section-block anim-fade-up anim-delay-3">
      <div class="section-header">
        <h2>技术栈</h2>
      </div>
      <div class="skills-grid">
        <div v-for="skill in profile.skills" :key="skill.name" class="skill-tag">
          <span class="skill-name">{{ skill.name }}</span>
          <span class="skill-cat">{{ skill.category }}</span>
        </div>
      </div>
    </section>

    <!-- 经历时间轴 -->
    <section class="section-block anim-fade-up anim-delay-4">
      <div class="section-header">
        <h2>经历</h2>
      </div>
      <div class="timeline">
        <div v-for="(item, index) in profile.timeline" :key="index" class="timeline-item">
          <div class="timeline-dot"></div>
          <div class="timeline-content">
            <span class="timeline-year">{{ item.year }}</span>
            <h3>{{ item.title }}</h3>
            <p>{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 写作理念 -->
    <section class="intro-banner anim-fade-up anim-delay-5">
      <div>
        <p class="section-kicker">Philosophy</p>
        <h2>写作理念</h2>
        <p class="lead-text">技术写作不是文档，而是思考的痕迹。</p>
      </div>
      <p class="intro-copy">
        每篇文章都是一次对技术问题的深度拆解。从问题发现到方案设计，从代码实现到复盘总结，
        这个博客记录的不只是结论，更是思考的过程。长期输出，持续精进。
      </p>
    </section>
  </SiteShell>
</template>

<style scoped>
.author-card {
  background: var(--bg-card);
  border: 1px solid var(--rule);
  border-radius: var(--radius-md);
  padding: var(--space-xl);
  margin-bottom: var(--space-xl);
}

.author-main {
  display: flex;
  gap: var(--space-xl);
  align-items: flex-start;
}

.author-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 88px;
  height: 88px;
  background: var(--ink-primary);
  color: var(--bg-page);
  font-family: var(--font-display);
  font-size: 2.2rem;
  font-weight: 900;
  flex-shrink: 0;
  position: relative;
}

.author-avatar::after {
  content: '';
  position: absolute;
  bottom: -3px;
  right: -3px;
  width: 10px;
  height: 10px;
  background: var(--accent);
  border-radius: 50%;
}

.author-info {
  flex: 1;
}

.author-info h1 {
  font-size: 2rem;
  margin-bottom: var(--space-xs);
}

.author-title {
  color: var(--accent);
  font-family: var(--font-display);
  font-size: 0.9rem;
  font-style: italic;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-md);
}

.author-bio {
  color: var(--ink-secondary);
  font-size: 0.95rem;
  line-height: 1.8;
  margin-bottom: var(--space-lg);
  max-width: 600px;
}

.author-social {
  display: flex;
  gap: var(--space-md);
  flex-wrap: wrap;
}

.social-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-warm);
  border: 1px solid var(--rule);
  color: var(--ink-secondary);
  font-size: 0.85rem;
  letter-spacing: 0.02em;
  transition: all 0.25s var(--ease-out);
}

.social-link:hover {
  color: var(--accent);
  border-color: rgba(194, 59, 34, 0.2);
  background: var(--accent-soft);
}

/* 技术栈 */
.skills-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.skill-tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--rule);
  transition: all 0.25s var(--ease-out);
}

.skill-tag:hover {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.skill-name {
  font-weight: 700;
  font-size: 0.9rem;
  color: var(--ink-primary);
}

.skill-cat {
  font-size: 0.72rem;
  color: var(--ink-muted);
  letter-spacing: 0.03em;
}

/* 时间轴 */
.timeline {
  position: relative;
  padding-left: 28px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 8px;
  bottom: 8px;
  width: 1px;
  background: var(--border-strong);
}

.timeline-item {
  position: relative;
  padding-bottom: var(--space-xl);
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: -28px;
  top: 6px;
  width: 11px;
  height: 11px;
  background: var(--bg-page);
  border: 2px solid var(--accent);
  border-radius: 50%;
  z-index: 1;
}

.timeline-year {
  display: inline-block;
  font-family: var(--font-display);
  font-size: 0.75rem;
  font-style: italic;
  color: var(--accent);
  letter-spacing: 0.08em;
  margin-bottom: var(--space-xs);
}

.timeline-content h3 {
  font-size: 1.05rem;
  margin-bottom: var(--space-xs);
}

.timeline-content p {
  color: var(--ink-secondary);
  font-size: 0.9rem;
  line-height: 1.7;
}

/* 响应式 */
@media (max-width: 640px) {
  .author-main {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .author-bio {
    max-width: 100%;
  }

  .author-social {
    justify-content: center;
  }

  .skills-grid {
    justify-content: center;
  }
}
</style>