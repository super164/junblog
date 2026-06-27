import { reactive } from 'vue'

const state = reactive({
  show: false,
  message: '',
  type: 'info',
  duration: 3000,
})

let timer = null

export function useMessage() {
  function showMessage(msg, type = 'info', duration = 3000) {
    if (timer) clearTimeout(timer)
    state.message = msg
    state.type = type
    state.duration = duration
    state.show = true

    if (duration > 0) {
      timer = setTimeout(() => {
        state.show = false
      }, duration)
    }
  }

  function success(msg, duration = 3000) {
    showMessage(msg, 'success', duration)
  }

  function error(msg, duration = 3000) {
    showMessage(msg, 'error', duration)
  }

  function warning(msg, duration = 3000) {
    showMessage(msg, 'warning', duration)
  }

  function info(msg, duration = 3000) {
    showMessage(msg, 'info', duration)
  }

  function close() {
    if (timer) clearTimeout(timer)
    state.show = false
  }

  return {
    state,
    showMessage,
    success,
    error,
    warning,
    info,
    close,
  }
}
