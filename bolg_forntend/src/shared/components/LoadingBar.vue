<script setup>
import { ref } from 'vue'

const visible = ref(false)
const progress = ref(0)
let timer = null

function start() {
  visible.value = true
  progress.value = 0
  clearInterval(timer)
  timer = setInterval(() => {
    if (progress.value < 90) {
      progress.value += Math.random() * 10
    }
  }, 100)
}

function finish() {
  progress.value = 100
  clearInterval(timer)
  setTimeout(() => {
    visible.value = false
    progress.value = 0
  }, 300)
}

defineExpose({ start, finish })
</script>

<template>
  <div v-if="visible" class="loading-bar">
    <div class="loading-bar-inner" :style="{ width: progress + '%' }"></div>
  </div>
</template>

<style scoped>
.loading-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  z-index: 9999;
  background: transparent;
}
.loading-bar-inner {
  height: 100%;
  background: linear-gradient(90deg, #e94560, #ff6b81);
  transition: width 0.2s ease;
  border-radius: 0 2px 2px 0;
  box-shadow: 0 0 8px rgba(233, 69, 96, 0.5);
}
</style>
