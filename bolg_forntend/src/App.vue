<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import LoadingBar from './shared/components/LoadingBar.vue'
import MessageToast from './components/MessageToast.vue'

const loadingBar = ref(null)
const router = useRouter()

router.beforeEach(() => {
  loadingBar.value?.start()
})

router.afterEach(() => {
  setTimeout(() => {
    loadingBar.value?.finish()
  }, 200)
})
</script>

<template>
  <LoadingBar ref="loadingBar" />
  <MessageToast />
  <router-view />
</template>
