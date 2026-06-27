<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'

const props = defineProps({
  visible: Boolean,
  article: Object
})

const emit = defineEmits(['close'])

const copied = ref(false)
const downloading = ref(false)
const canvasRef = ref(null)
const cardImage = ref(null)
const showWechatQR = ref(false)
const showQQQR = ref(false)

// 生成文章摘要
const summary = computed(() => {
  if (!props.article?.content) return ''
  // 去除 Markdown 语法，取前 120 个字符
  const plain = props.article.content
    .replace(/#{1,6}\s/g, '')
    .replace(/\*{1,2}([^*]+)\*{1,2}/g, '$1')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/!\[[^\]]*\]\([^)]+\)/g, '')
    .replace(/`{1,3}[^`]*`{1,3}/g, '')
    .replace(/\n/g, ' ')
    .trim()
  return plain.length > 120 ? plain.slice(0, 120) + '...' : plain
})

// 网站 URL
const siteUrl = computed(() => {
  return window.location.origin
})

// 文章 URL（注意路由是 /articles/:id 复数形式）
const articleUrl = computed(() => {
  return `${siteUrl.value}/articles/${props.article?.id}`
})

// 生成二维码（使用简单的 QR API）
const qrUrl = computed(() => {
  const url = encodeURIComponent(articleUrl.value)
  return `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${url}&bgcolor=ffffff&color=c23b22`
})

// 生成卡片
async function generateCard() {
  if (!canvasRef.value || !props.article) return

  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  const width = 800
  const height = 420
  canvas.width = width
  canvas.height = height

  // 背景渐变
  const gradient = ctx.createLinearGradient(0, 0, width, height)
  gradient.addColorStop(0, '#faf7f2')
  gradient.addColorStop(1, '#f5f0e8')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, height)

  // 装饰元素 - 左上角圆弧
  ctx.beginPath()
  ctx.arc(-50, -50, 150, 0, Math.PI * 2)
  ctx.fillStyle = 'rgba(194, 59, 34, 0.06)'
  ctx.fill()

  // 装饰元素 - 右下角圆弧
  ctx.beginPath()
  ctx.arc(width + 50, height + 50, 200, 0, Math.PI * 2)
  ctx.fillStyle = 'rgba(194, 59, 34, 0.04)'
  ctx.fill()

  // 左侧红色装饰条
  ctx.fillStyle = '#c23b22'
  ctx.fillRect(40, 60, 4, 120)

  // 分类标签
  if (props.article.category?.name) {
    ctx.font = '600 13px "Noto Sans SC", sans-serif'
    ctx.fillStyle = '#c23b22'
    ctx.fillText(props.article.category.name.toUpperCase(), 56, 80)
  }

  // 标题
  ctx.font = 'bold 32px "Noto Serif SC", serif'
  ctx.fillStyle = '#1a1a1a'
  const title = props.article.title || 'Untitled'
  // 简单的标题换行
  const titleLines = wrapText(ctx, title, 480, 2)
  titleLines.forEach((line, i) => {
    ctx.fillText(line, 56, 120 + i * 44)
  })

  // 摘要
  ctx.font = '15px "Noto Sans SC", sans-serif'
  ctx.fillStyle = '#666666'
  const summaryLines = wrapText(ctx, summary.value, 460, 2)
  summaryLines.forEach((line, i) => {
    ctx.fillText(line, 56, 230 + i * 24)
  })

  // 底部信息栏
  ctx.fillStyle = 'rgba(26, 26, 26, 0.04)'
  ctx.fillRect(0, height - 80, width, 80)

  // 作者信息
  ctx.font = '500 14px "Noto Sans SC", sans-serif'
  ctx.fillStyle = '#1a1a1a'
  const authorName = props.article.author?.username || 'Jun'
  ctx.fillText(authorName, 56, height - 45)

  // 网站域名
  ctx.font = '13px "Noto Sans SC", sans-serif'
  ctx.fillStyle = '#999999'
  ctx.fillText(siteUrl.value, 56, height - 25)

  // 日期
  const date = new Date(props.article.created_at).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
  ctx.textAlign = 'right'
  ctx.font = '13px "Noto Sans SC", sans-serif'
  ctx.fillStyle = '#999999'
  ctx.fillText(date, width - 180, height - 45)

  // 右侧二维码区域
  ctx.textAlign = 'left'

  // 二维码背景
  ctx.fillStyle = '#ffffff'
  roundRect(ctx, width - 170, height - 160, 140, 140, 12)
  ctx.fill()
  ctx.strokeStyle = 'rgba(194, 59, 34, 0.15)'
  ctx.lineWidth = 1
  roundRect(ctx, width - 170, height - 160, 140, 140, 12)
  ctx.stroke()

  // 加载二维码图片
  const qrImg = new Image()
  qrImg.crossOrigin = 'anonymous'
  await new Promise((resolve) => {
    qrImg.onload = resolve
    qrImg.onerror = resolve
    qrImg.src = qrUrl.value
  })

  if (qrImg.complete && qrImg.naturalWidth > 0) {
    ctx.drawImage(qrImg, width - 158, height - 148, 116, 116)
  } else {
    // 二维码加载失败时的占位
    ctx.font = '12px "Noto Sans SC", sans-serif'
    ctx.fillStyle = '#999'
    ctx.textAlign = 'center'
    ctx.fillText('扫码阅读', width - 100, height - 90)
    ctx.textAlign = 'left'
  }

  // 扫码提示
  ctx.font = '11px "Noto Sans SC", sans-serif'
  ctx.fillStyle = '#999999'
  ctx.textAlign = 'center'
  ctx.fillText('扫码阅读原文', width - 100, height - 15)
  ctx.textAlign = 'left'

  // 生成图片
  cardImage.value = canvas.toDataURL('image/png', 1.0)
}

// 文本换行辅助函数
function wrapText(ctx, text, maxWidth, maxLines) {
  const lines = []
  let currentLine = ''

  for (let i = 0; i < text.length; i++) {
    const testLine = currentLine + text[i]
    const metrics = ctx.measureText(testLine)

    if (metrics.width > maxWidth && currentLine) {
      lines.push(currentLine)
      currentLine = text[i]
      if (lines.length >= maxLines) {
        lines[lines.length - 1] += '...'
        return lines
      }
    } else {
      currentLine = testLine
    }
  }

  if (currentLine) {
    lines.push(currentLine)
  }

  return lines
}

// 圆角矩形辅助函数
function roundRect(ctx, x, y, width, height, radius) {
  ctx.beginPath()
  ctx.moveTo(x + radius, y)
  ctx.lineTo(x + width - radius, y)
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius)
  ctx.lineTo(x + width, y + height - radius)
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height)
  ctx.lineTo(x + radius, y + height)
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius)
  ctx.lineTo(x, y + radius)
  ctx.quadraticCurveTo(x, y, x + radius, y)
  ctx.closePath()
}

// 复制链接
async function copyLink() {
  try {
    await navigator.clipboard.writeText(articleUrl.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // 降级方案
    const input = document.createElement('input')
    input.value = articleUrl.value
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}

// 下载卡片
function downloadCard() {
  if (!cardImage.value) return
  downloading.value = true

  const link = document.createElement('a')
  link.download = `share-${props.article?.id || 'card'}.png`
  link.href = cardImage.value
  link.click()

  setTimeout(() => { downloading.value = false }, 1000)
}

// 监听 visible 变化，生成卡片
watch(() => props.visible, async (val) => {
  if (val && props.article) {
    await nextTick()
    await generateCard()
  }
})

onMounted(() => {
  if (props.visible && props.article) {
    nextTick(() => generateCard())
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="share-modal-overlay" @click.self="emit('close')">
        <div class="share-modal">
          <div class="share-modal-header">
            <h3>分享文章</h3>
            <button class="close-btn" @click="emit('close')">×</button>
          </div>

          <div class="share-modal-body">
            <!-- 卡片预览 -->
            <div class="card-preview">
              <canvas ref="canvasRef" class="hidden-canvas"></canvas>
              <img v-if="cardImage" :src="cardImage" alt="Share Card" class="card-image" />
              <div v-else class="card-loading">
                <div class="loading-spinner"></div>
                <span>生成卡片中...</span>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="share-actions">
              <button class="action-btn copy-btn" @click="copyLink">
                <span class="action-icon">{{ copied ? '✓' : '🔗' }}</span>
                <span>{{ copied ? '已复制!' : '复制链接' }}</span>
              </button>

              <button class="action-btn download-btn" @click="downloadCard" :disabled="!cardImage || downloading">
                <span class="action-icon">{{ downloading ? '⏳' : '📥' }}</span>
                <span>{{ downloading ? '下载中...' : '下载卡片' }}</span>
              </button>
            </div>

            <!-- 分享到社交媒体 -->
            <div class="social-share">
              <span class="social-label">分享到</span>
              <div class="social-buttons">
                <button class="social-btn qq" @click="showQQQR = true" title="QQ">
                  QQ
                </button>
                <button class="social-btn wechat" @click="showWechatQR = true" title="微信">
                  微信
                </button>
              </div>
            </div>

            <!-- 扫码弹窗（QQ 和微信共用） -->
            <Teleport to="body">
              <Transition name="modal">
                <div v-if="showWechatQR || showQQQR" class="scan-modal-overlay" @click.self="showWechatQR = false; showQQQR = false">
                  <div class="scan-modal">
                    <div class="scan-modal-header">
                      <h4>{{ showQQQR ? 'QQ扫码分享' : '微信扫码分享' }}</h4>
                      <button class="close-btn" @click="showWechatQR = false; showQQQR = false">×</button>
                    </div>
                    <div class="scan-modal-body">
                      <img :src="qrUrl" alt="二维码" class="scan-qr" />
                      <p>{{ showQQQR ? '打开QQ，扫一扫分享给好友' : '打开微信，扫一扫分享给好友' }}</p>
                    </div>
                  </div>
                </div>
              </Transition>
            </Teleport>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.hidden-canvas {
  position: absolute;
  left: -9999px;
  top: -9999px;
}

.share-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.share-modal {
  background: #fff;
  border-radius: 20px;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.share-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(26, 26, 26, 0.06);
}

.share-modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(26, 26, 26, 0.05);
  border-radius: 8px;
  cursor: pointer;
  font-size: 20px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(26, 26, 26, 0.1);
  color: #1a1a1a;
}

.share-modal-body {
  padding: 24px;
}

.card-preview {
  background: #f5f0e8;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
  aspect-ratio: 800 / 420;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.card-loading {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #999;
  font-size: 14px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2.5px solid rgba(194, 59, 34, 0.15);
  border-top-color: #c23b22;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.share-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 20px;
  border: 1.5px solid rgba(26, 26, 26, 0.1);
  background: #fff;
  border-radius: 12px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);
  font-family: inherit;
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.copy-btn:hover:not(:disabled) {
  border-color: #c23b22;
  color: #c23b22;
}

.download-btn:hover:not(:disabled) {
  border-color: #c23b22;
  background: #c23b22;
  color: #fff;
}

.action-icon {
  font-size: 16px;
}

.social-share {
  text-align: center;
}

.social-label {
  font-size: 12px;
  color: #999;
  display: block;
  margin-bottom: 12px;
}

.social-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.social-btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.social-btn.qq {
  background: #12b7f5;
  color: #fff;
}

.social-btn.qq {
  background: #12b7f5;
  color: #fff;
}

.social-btn.qq:hover {
  background: #0ea5e0;
}

.social-btn.wechat {
  background: #07c160;
  color: #fff;
}

.social-btn.wechat:hover {
  background: #06ae56;
}

/* 扫码弹窗（QQ 和微信共用） */
.scan-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.scan-modal {
  background: #fff;
  border-radius: 16px;
  width: 280px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.scan-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(26, 26, 26, 0.06);
}

.scan-modal-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.close-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: rgba(26, 26, 26, 0.05);
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.close-btn:hover {
  background: rgba(26, 26, 26, 0.1);
}

.scan-modal-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 24px;
}

.scan-qr {
  width: 180px;
  height: 180px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.scan-modal-body p {
  margin: 0;
  font-size: 13px;
  color: #999;
}

/* 过渡动画 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-active .share-modal,
.modal-leave-active .share-modal {
  transition: transform 0.3s cubic-bezier(0.22, 1, 0.36, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .share-modal,
.modal-leave-to .share-modal {
  transform: scale(0.95) translateY(10px);
}

/* 移动端适配 */
@media (max-width: 560px) {
  .share-modal {
    max-width: 100%;
    margin: 10px;
    border-radius: 16px;
  }

  .share-actions {
    flex-direction: column;
  }
}
</style>
