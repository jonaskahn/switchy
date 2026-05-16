<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps<{ seconds: number }>()
const emit = defineEmits<{ timeout: [] }>()

const remaining = ref(props.seconds)
const progress = ref(100)
let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(() => {
    remaining.value--
    progress.value = (remaining.value / props.seconds) * 100
    if (remaining.value <= 0) {
      if (timer) clearInterval(timer)
      emit('timeout')
    }
  }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="timer-bar-track">
    <div class="timer-bar-fill" :style="{ width: progress + '%' }" />
  </div>
</template>

<style scoped>
.timer-bar-track {
  height: 2px;
  width: 100%;
  background: var(--ink-10);
  overflow: hidden;
}
.timer-bar-fill {
  height: 100%;
  background: rgb(from var(--accent) r g b / 0.55);
  transition: width 1000ms linear;
}
</style>
