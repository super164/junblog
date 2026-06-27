<script setup>
import { useMessage } from '../composables/useMessage'

const { state, close } = useMessage()
</script>

<template>
  <Teleport to="body">
    <Transition name="toast">
      <div v-if="state.show" class="toast-container" :class="state.type">
        <div class="toast-icon">
          <span v-if="state.type === 'success'">✓</span>
          <span v-else-if="state.type === 'error'">✕</span>
          <span v-else-if="state.type === 'warning'">!</span>
          <span v-else>i</span>
        </div>
        <span class="toast-message">{{ state.message }}</span>
        <button class="toast-close" @click="close">✕</button>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06);
  z-index: 3000;
  min-width: 280px;
  max-width: 480px;
  border-left: 4px solid #8a8578;
}

.toast-container.success {
  border-left-color: #52c41a;
}

.toast-container.error {
  border-left-color: #dc3545;
}

.toast-container.warning {
  border-left-color: #faad14;
}

.toast-container.info {
  border-left-color: #1890ff;
}

.toast-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}

.success .toast-icon {
  background: rgba(82, 196, 26, 0.12);
  color: #52c41a;
}

.error .toast-icon {
  background: rgba(220, 53, 69, 0.12);
  color: #dc3545;
}

.warning .toast-icon {
  background: rgba(250, 173, 20, 0.12);
  color: #d48806;
}

.info .toast-icon {
  background: rgba(24, 144, 255, 0.12);
  color: #1890ff;
}

.toast-message {
  flex: 1;
  font-size: 14px;
  color: #2c2825;
  line-height: 1.5;
}

.toast-close {
  background: none;
  border: none;
  color: #b8afa6;
  cursor: pointer;
  padding: 4px;
  font-size: 12px;
  transition: color 0.2s;
  flex-shrink: 0;
}

.toast-close:hover {
  color: #4a4540;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-16px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-16px);
}
</style>
