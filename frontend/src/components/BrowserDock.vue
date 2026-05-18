<script setup lang="ts">
import type { Browser } from '@/types'
import BrowserDockTile from '@/components/BrowserDockTile.vue'
import { browserShortcut, browserAnimationDelay } from '@/utils/browser'
defineProps<{
  browsers: Browser[]
  loading?: boolean
}>()
const emit = defineEmits<{ (e: 'open', name: string, remember: boolean): void; (e: 'open-settings'): void }>()
</script>

<template>
  <div class="dock-root">
    <div v-if="loading" class="dock-loading">Loading…</div>
    <div v-else class="dock no-drag" role="listbox">
      <BrowserDockTile
        v-for="(browser, idx) in browsers"
        :key="browser.Name"
        :browser="browser"
        :shortcut="browserShortcut(idx)"
        :focused="idx === 0"
        :animation-delay-ms="browserAnimationDelay(idx)"
        @open="(remember: boolean) => emit('open', browser.Name, remember)"
      />
    </div>
    <div v-if="!loading && browsers.length > 0" class="dock-floor"></div>
  </div>
</template>

<style scoped>
.dock-root {
  flex-shrink: 0;
  min-width: 0;
}

.dock {
  padding: 22px 22px 16px;
  display: flex;
  align-items: flex-start;
  justify-content: safe center;
  flex-wrap: nowrap;
  gap: 10px;
  overflow-x: auto;
  min-width: 0;
}

.dock-floor {
  height: 1px;
  margin: 0 22px;
  background: linear-gradient(
    to right,
    transparent 0%,
    var(--hairline) 20%,
    var(--hairline) 80%,
    transparent 100%
  );
}

.dock::-webkit-scrollbar {
  height: 3px;
}
.dock::-webkit-scrollbar-track {
  background: transparent;
}
.dock::-webkit-scrollbar-thumb {
  background: var(--ink-20);
  border-radius: 99px;
}
.dock::-webkit-scrollbar-thumb:hover {
  background: var(--ink-40);
}

.dock-loading {
  padding: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-60);
  font-family: var(--font-mono);
  font-size: 12px;
  letter-spacing: 0.02em;
}

.dock-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 32px 24px;
  margin: 16px;
  border: 1px dashed var(--hairline);
  border-radius: var(--r-row);
  color: var(--ink-60);
  font-size: 13px;
  text-align: center;
}
</style>
