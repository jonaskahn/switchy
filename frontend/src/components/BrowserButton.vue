<script setup lang="ts">
import { onMounted, computed, toRef } from 'vue'
import type { Browser } from '@/types'
import { useBrowserIcon } from '@/composables/useBrowserIcon'

const props = defineProps<{
  browser: Browser
  shortcut?: string
}>()

defineEmits<{ (e: 'open', remember: boolean): void }>()

const browserRef = toRef(props, 'browser')
const { iconDataUrl, loadIcon } = useBrowserIcon(browserRef)
console.log(iconDataUrl)

const subLabel = computed(() => {
  const parts: string[] = []
  if (props.browser.IsUwp) parts.push('uwp')
  if (props.browser.AlternateLaunches?.length)
    parts.push(...props.browser.AlternateLaunches.map((a) => a.Name.toLowerCase()))
  return parts.join(' · ') || ''
})

onMounted(() => void loadIcon())
</script>

<template>
  <button class="browser-row" role="option" @click="(e: MouseEvent) => $emit('open', e.ctrlKey)">
    <div class="browser-icon">
      <img v-if="iconDataUrl" :src="iconDataUrl" decoding="async" :alt="browser.Name" width="12" />
      <svg
        v-else
        viewBox="0 0 24 24"
        width="24"
        height="24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        style="color: var(--ink-60)"
      >
        <circle cx="12" cy="12" r="10" />
        <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
        <path d="M2 12h20" />
      </svg>
      <img
        v-if="browser.ProfileIcon"
        class="profile-avatar"
        :src="browser.ProfileIcon"
        decoding="async"
        alt="Profile"
        width="10"
        height="10"
      />
    </div>

    <div class="browser-meta">
      <div class="browser-name">{{ browser.Name }}</div>
      <div v-if="subLabel" class="browser-sub">{{ subLabel }}</div>
    </div>

    <div v-if="shortcut" class="keycap">{{ shortcut }}</div>
  </button>
</template>

<style scoped>
/* Stagger animation — parent passes animation-delay via :style */
.browser-row {
  margin: 0 4px;
  animation: rowIn 320ms cubic-bezier(0.2, 0.8, 0.2, 1) both;
}
@keyframes rowIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
