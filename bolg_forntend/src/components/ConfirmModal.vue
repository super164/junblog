<script setup>
const props = defineProps({
  show: Boolean,
  title: { type: String, default: '确认操作' },
  message: { type: String, default: '确认执行此操作？' },
  confirmText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' },
  danger: { type: Boolean, default: false },
})

const emit = defineEmits(['confirm', 'cancel'])
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="modal-overlay" @click.self="emit('cancel')">
        <div class="confirm-modal">
          <div class="modal-icon" :class="{ danger }">
            <span v-if="danger">!</span>
            <span v-else>?</span>
          </div>
          <h4>{{ title }}</h4>
          <p>{{ message }}</p>
          <div class="modal-actions">
            <button class="btn-cancel" @click="emit('cancel')">{{ cancelText }}</button>
            <button class="btn-confirm" :class="{ danger }" @click="emit('confirm')">{{ confirmText }}</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(2px);
}

.confirm-modal {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  width: 380px;
  max-width: 90vw;
  text-align: center;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  animation: modalIn 0.25s ease-out;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  background: rgba(194, 59, 34, 0.1);
  color: #c23b22;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
}

.modal-icon.danger {
  background: rgba(220, 53, 69, 0.1);
  color: #dc3545;
}

.confirm-modal h4 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #2c2825;
}

.confirm-modal p {
  margin: 0 0 24px;
  font-size: 14px;
  color: #8a8578;
  line-height: 1.5;
}

.modal-actions {
  display: flex;
  gap: 12px;
}

.btn-cancel {
  flex: 1;
  padding: 12px;
  background: #f5f0eb;
  color: #4a4540;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background: #ebe5de;
}

.btn-confirm {
  flex: 1;
  padding: 12px;
  background: linear-gradient(135deg, #c23b22 0%, #a52d18 100%);
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-confirm:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(194, 59, 34, 0.3);
}

.btn-confirm.danger {
  background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
}

.btn-confirm.danger:hover {
  box-shadow: 0 4px 12px rgba(220, 53, 69, 0.3);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
