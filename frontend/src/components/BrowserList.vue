<script setup lang="ts">
import type { Browser } from '@/types'
import BrowserButton from '@/components/BrowserButton.vue'
import { browserShortcut, browserAnimationStyle } from '@/utils/browser'
defineProps<{
  browsers: Browser[]
  loading?: boolean
}>()
defineEmits<{
  (e: 'open', name: string): void
  (e: 'open-settings'): void
}>()
</script>

<template>
  <div class="browser-list no-drag" role="listbox">
    <div v-if="loading" class="list-loading">Loading…</div>
    <div v-else-if="browsers.length === 0" class="list-empty">
      <svg
        width="32"
        height="32"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="10" />
        <line x1="12" y1="8" x2="12" y2="12" />
        <line x1="12" y1="16" x2="12.01" y2="16" />
      </svg>
      <p>No browsers configured</p>
      <button class="btn-ghost" style="width: auto" @click="$emit('open-settings')">
        Open settings →
      </button>
    </div>
    <template v-else>
      <BrowserButton
        v-for="(browser, idx) in browsers"
        :key="browser.Name"
        :browser="browser"
        :shortcut="browserShortcut(idx)"
        :style="browserAnimationStyle(idx)"
        @open="$emit('open', browser.Name)"
      />
    </template>
  </div>
</template>

<style scoped>
.browser-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px 56px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-height: 0;
}

.list-loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-60);
  font-family: var(--font-mono);
  font-size: 12px;
  letter-spacing: 0.02em;
}

.list-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--ink-60);
  font-size: 13px;
  padding: 32px 24px;
  margin: 8px;
  border: 1px dashed var(--hairline);
  border-radius: var(--r-row);
  text-align: center;
}
</style>
