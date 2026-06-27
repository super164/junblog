<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'

const props = defineProps({
  show: Boolean,
  imageSrc: String,
})
const emit = defineEmits(['confirm', 'cancel'])

const canvasRef = ref(null)
const containerRef = ref(null)

const img = ref(null)
const imgWidth = ref(0)
const imgHeight = ref(0)

// 裁剪框在图片坐标系中的位置
const cropX = ref(0)
const cropY = ref(0)
const cropW = ref(0)
const cropH = ref(0)

// Canvas 显示尺寸
const canvasW = ref(0)
const canvasH = ref(0)
const scale = ref(1)

// 拖拽状态
const dragging = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)
const dragOrigX = ref(0)
const dragOrigY = ref(0)

const RATIO = 16 / 9

watch(() => props.show, (v) => {
  if (v && props.imageSrc) loadImage()
})

function loadImage() {
  const image = new Image()
  image.crossOrigin = 'anonymous'
  image.onload = () => {
    img.value = image
    imgWidth.value = image.naturalWidth
    imgHeight.value = image.naturalHeight
    nextTick(() => initCrop())
  }
  image.src = props.imageSrc
}

function initCrop() {
  const container = containerRef.value
  if (!container) return

  const maxW = Math.min(container.clientWidth - 48, 720)
  const maxH = window.innerHeight * 0.6

  // 计算图片在 canvas 上的显示尺寸
  const ratio = imgWidth.value / imgHeight.value
  let dw, dh
  if (ratio > maxW / maxH) {
    dw = maxW
    dh = maxW / ratio
  } else {
    dh = maxH
    dw = maxH * ratio
  }

  canvasW.value = dw
  canvasH.value = dh
  scale.value = imgWidth.value / dw

  // 初始裁剪框：尽可能大的 16:9 区域
  let cw, ch
  if (dw / dh > RATIO) {
    ch = dh
    cw = dh * RATIO
  } else {
    cw = dw
    ch = dw / RATIO
  }

  cropW.value = cw
  cropH.value = ch
  cropX.value = (dw - cw) / 2
  cropY.value = (dh - ch) / 2

  draw()
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas || !img.value) return
  const ctx = canvas.getContext('2d')

  canvas.width = canvasW.value * 2
  canvas.height = canvasH.value * 2
  ctx.scale(2, 2)

  // 画图片
  ctx.drawImage(img.value, 0, 0, canvasW.value, canvasH.value)

  // 半透明遮罩
  ctx.fillStyle = 'rgba(0, 0, 0, 0.55)'
  ctx.fillRect(0, 0, canvasW.value, canvasH.value)

  // 裁剪区域透出原图
  ctx.clearRect(cropX.value, cropY.value, cropW.value, cropH.value)
  ctx.drawImage(
    img.value,
    cropX.value * scale.value, cropY.value * scale.value,
    cropW.value * scale.value, cropH.value * scale.value,
    cropX.value, cropY.value, cropW.value, cropH.value
  )

  // 裁剪框边框
  ctx.strokeStyle = '#fff'
  ctx.lineWidth = 2
  ctx.setLineDash([])
  ctx.strokeRect(cropX.value, cropY.value, cropW.value, cropH.value)

  // 四角手柄
  const hs = 8
  ctx.fillStyle = '#fff'
  const corners = [
    [cropX.value, cropY.value],
    [cropX.value + cropW.value - hs, cropY.value],
    [cropX.value, cropY.value + cropH.value - hs],
    [cropX.value + cropW.value - hs, cropY.value + cropH.value - hs],
  ]
  corners.forEach(([x, y]) => ctx.fillRect(x, y, hs, hs))

  // 三分线
  ctx.strokeStyle = 'rgba(255,255,255,0.3)'
  ctx.lineWidth = 1
  for (let i = 1; i <= 2; i++) {
    const lx = cropX.value + (cropW.value * i) / 3
    ctx.beginPath()
    ctx.moveTo(lx, cropY.value)
    ctx.lineTo(lx, cropY.value + cropH.value)
    ctx.stroke()

    const ly = cropY.value + (cropH.value * i) / 3
    ctx.beginPath()
    ctx.moveTo(cropX.value, ly)
    ctx.lineTo(cropX.value + cropW.value, ly)
    ctx.stroke()
  }
}

function onPointerDown(e) {
  dragging.value = true
  dragStartX.value = e.clientX
  dragStartY.value = e.clientY
  dragOrigX.value = cropX.value
  dragOrigY.value = cropY.value
  e.preventDefault()
}

function onPointerMove(e) {
  if (!dragging.value) return
  const dx = e.clientX - dragStartX.value
  const dy = e.clientY - dragStartY.value

  let nx = dragOrigX.value + dx
  let ny = dragOrigY.value + dy

  // 边界限制
  nx = Math.max(0, Math.min(canvasW.value - cropW.value, nx))
  ny = Math.max(0, Math.min(canvasH.value - cropH.value, ny))

  cropX.value = nx
  cropY.value = ny
  draw()
}

function onPointerUp() {
  dragging.value = false
}

function confirmCrop() {
  const canvas = document.createElement('canvas')
  const outW = Math.round(cropW.value * scale.value)
  const outH = Math.round(cropH.value * scale.value)
  canvas.width = outW
  canvas.height = outH
  const ctx = canvas.getContext('2d')
  ctx.drawImage(
    img.value,
    cropX.value * scale.value, cropY.value * scale.value,
    outW, outH,
    0, 0, outW, outH
  )
  canvas.toBlob((blob) => {
    emit('confirm', blob)
  }, 'image/jpeg', 0.92)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="crop-modal">
      <div v-if="show" class="crop-overlay" @click.self="emit('cancel')">
        <div class="crop-dialog" ref="containerRef">
          <div class="crop-header">
            <h3>裁剪封面图</h3>
            <span class="crop-hint">拖拽选择区域 · 固定 16:9 比例</span>
          </div>

          <div class="crop-canvas-wrap">
            <canvas
              ref="canvasRef"
              :style="{ width: canvasW + 'px', height: canvasH + 'px', cursor: dragging ? 'grabbing' : 'grab' }"
              @pointerdown="onPointerDown"
              @pointermove="onPointerMove"
              @pointerup="onPointerUp"
              @pointerleave="onPointerUp"
            />
          </div>

          <div class="crop-actions">
            <button class="crop-btn-cancel" @click="emit('cancel')">取消</button>
            <button class="crop-btn-confirm" @click="confirmCrop">确认裁剪</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.crop-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(2px);
}

.crop-dialog {
  background: #faf7f2;
  border-radius: 12px;
  padding: 24px;
  max-width: 800px;
  width: 90vw;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.crop-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 16px;
}

.crop-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}

.crop-hint {
  font-size: 13px;
  color: #8a8578;
}

.crop-canvas-wrap {
  display: flex;
  justify-content: center;
  overflow: hidden;
  border-radius: 8px;
  background: #1a1a1a;
}

.crop-canvas-wrap canvas {
  display: block;
  touch-action: none;
}

.crop-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 16px;
}

.crop-btn-cancel,
.crop-btn-confirm {
  padding: 9px 22px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.2s;
}

.crop-btn-cancel {
  background: transparent;
  border: 1.5px solid rgba(26, 26, 26, 0.15);
  color: #4a4a4a;
}

.crop-btn-cancel:hover {
  border-color: #1a1a1a;
  color: #1a1a1a;
}

.crop-btn-confirm {
  background: #c23b22;
  border: none;
  color: #fff;
}

.crop-btn-confirm:hover {
  background: #a83219;
}

/* Transition */
.crop-modal-enter-active,
.crop-modal-leave-active {
  transition: opacity 0.25s ease;
}
.crop-modal-enter-active .crop-dialog,
.crop-modal-leave-active .crop-dialog {
  transition: transform 0.25s ease, opacity 0.25s ease;
}
.crop-modal-enter-from,
.crop-modal-leave-to {
  opacity: 0;
}
.crop-modal-enter-from .crop-dialog {
  transform: scale(0.95);
  opacity: 0;
}
.crop-modal-leave-to .crop-dialog {
  transform: scale(0.95);
  opacity: 0;
}
</style>
